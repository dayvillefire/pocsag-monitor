package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/coreos/go-systemd/daemon"
	"github.com/dayvillefire/pocsag-monitor/config"
	"github.com/dayvillefire/pocsag-monitor/obj"
	"github.com/dayvillefire/pocsag-monitor/output"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

var (
	configFile        = flag.String("config", "config.yaml", "Configuration file")
	dynamicConfigFile = flag.String("dynamic-config", "dynamic.yaml", "Dynamic configuration file")
	testConfig        = flag.Bool("test-config", false, "Test config")
	daemonize         = flag.Bool("daemon", false, "Daemonize")

	Version      string
	logRoute     string
	cfg          *config.Config
	router       Router
	outputs      map[string]output.Output
	outputsMutex = &sync.Mutex{}
	routerMutex  = &sync.Mutex{}
)

func logger(s string) {
	log.Print(s)
	if logRoute != "" {
		outputs[logRoute].SendMessage(
			obj.AlphaMessage{},
			cfg.Dynamic.OutputChannels[logRoute].Channel,
			s,
		)
		return
	}
}

func main() {
	flag.Parse()

	var err error
	cfg, err = config.LoadConfigWithDefaults(*configFile, *dynamicConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	if *testConfig {
		log.Printf("%#v", cfg)
		os.Exit(0)
	}

	// Daemon stuff if we're configured for it.
	if *daemonize {
		go func() {
			log.Printf("Daemon: INFO: Spawning systemd integration")

			interval, err := daemon.SdWatchdogEnabled(false)
			if err != nil {
				log.Printf("ERR: %s", err.Error())
				return
			}
			if interval == 0 {
				log.Printf("ERR: interval == 0")
				return
			}
			for {
				daemon.SdNotify(false, daemon.SdNotifyWatchdog)
				time.Sleep(interval / 3)
			}
		}()
	}

	rtlArg := fmt.Sprintf("-f %s -p %d -s 22050", cfg.Frequency, cfg.PPM)
	rtlCmd := exec.Command(cfg.RtlFmBinary, strings.Split(rtlArg, " ")...)

	rtlStderr, _ := rtlCmd.StderrPipe()

	mmonArg := "-t raw -a POCSAG512 -f alpha -u /dev/stdin"
	mmonCmd := exec.Command(cfg.MultiMonBinary, strings.Split(mmonArg, " ")...)

	mmonCmd.Stdin, _ = rtlCmd.StdoutPipe()
	stdout, err := mmonCmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	err = mmonCmd.Start()
	if err != nil {
		panic(err)
	}
	err = rtlCmd.Start()
	if err != nil {
		panic(err)
	}
	defer stdout.Close()
	defer rtlStderr.Close()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGQUIT)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGABRT)
	go func(sig chan os.Signal, rtlCmd *exec.Cmd, mmonCmd *exec.Cmd) {
		s := <-sig
		log.Printf("Caught signal %s, terminating", s.String())
		rtlCmd.Process.Kill()
		mmonCmd.Process.Kill()
	}(sig, rtlCmd, mmonCmd)
	defer func(rtlCmd *exec.Cmd) {
		// If, for some reason, this doesn't die gracefully, kill it with fire
		log.Printf("Non-gracefully terminating rtl_fm")
		rtlCmd.Process.Kill()
	}(rtlCmd)

	go func() {
		log.Printf("INFO: Initializing web services")
		m := gin.New()
		m.Use(gin.Recovery())

		// Enable gzip compression
		m.Use(gzip.Gzip(gzip.DefaultCompression))

		InitApi(m)

		go func() {
			log.Printf("INFO: Initializing on :%d", cfg.ApiPort)
			if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.ApiPort), m); err != nil {
				log.Fatal(err)
			}
		}()
	}()

	// Dynamic channel mapping init
	if cfg.Debug {
		log.Printf("DEBUG: Locking routerMutex")
	}
	routerMutex.Lock()
	router = Router{cfg.Dynamic.ChannelMappings}
	outputs = map[string]output.Output{}
	wg := &sync.WaitGroup{}
	for k, v := range cfg.Dynamic.OutputChannels {
		wg.Add(1)
		go func(k string, v config.OutputMapping) {
			defer wg.Done()
			log.Printf("[%s] Instantiating %s", k, v.Plugin)
			o, err := output.InstantiateOutput(v.Plugin)
			if err != nil {
				log.Printf(k + "| ERR: " + err.Error() + " - output disabled")
				o, _ = output.InstantiateOutput("dummy")
				outputsMutex.Lock()
				outputs[k] = o
				outputsMutex.Unlock()
				return
			}
			err = o.Init(v.Option)
			if err != nil {
				log.Printf(k + "| ERR(Init): " + err.Error() + " - output disabled")
				o, _ = output.InstantiateOutput("dummy")
			}
			outputsMutex.Lock()
			outputs[k] = o
			log.Printf("[%s] Init completed", k)
			outputsMutex.Unlock()
		}(k, v)
	}
	wg.Wait()
	if cfg.Debug {
		log.Printf("DEBUG: Unlocking routerMutex")
	}
	routerMutex.Unlock()

	// Establish log route
	logRoute = router.LogRoute()
	logger("Initialized pocsag router version " + Version + " at " + time.Now().String())

	if cfg.Debug {
		log.Printf("DEBUG: Instantiating scanner")
	}
	scanner := bufio.NewScanner(io.MultiReader(stdout, rtlStderr))
	if cfg.Debug {
		log.Printf("DEBUG: scanner.Split(scan lines)")
	}
	scanner.Split(bufio.ScanLines)
	if cfg.Debug {
		log.Printf("DEBUG: Loop through scan")
	}
	for scanner.Scan() {
		if cfg.Debug {
			log.Printf("DEBUG: scanner.Text()")
		}
		m := scanner.Text()
		if cfg.Debug {
			log.Printf("DEBUG: Found line '%s'", m)
		}
		ts := time.Now()
		alpha, err := obj.ParseAlphaMessage(ts, m)
		if err != nil {
			logger("ParseAlphaMessage: ERR: " + err.Error())
			continue
		}
		if alpha.Valid {
			log.Printf("CAP: %s\tMSG: %s", alpha.CapCode, alpha.Message)
			routerMutex.Lock()
			dest := router.MapMessage(alpha)
			for _, c := range dest {
				msg := fmt.Sprintf(
					"%s: %s [%s]",
					alpha.CapCode,
					alpha.Message,
					alpha.Timestamp.Format("2006-01-02 15:04:05"),
				)
				if outputs[c] == nil {
					logger(fmt.Sprintf("ERROR: dest=%s|outputchannel[%s]|msg=%s",
						dest,
						c,
						msg,
					))
					continue
				}
				if cfg.Debug {
					log.Printf("DEBUG: dest=%s|option=%s|msg=%s",
						dest,
						cfg.Dynamic.OutputChannels[c].Channel,
						msg,
					)
				}
				outputs[c].SendMessage(
					alpha,
					cfg.Dynamic.OutputChannels[c].Channel,
					msg,
				)
			}
			routerMutex.Unlock()
			//db.Record(DB, alpha)
			continue
		}
	}
	mmonCmd.Wait()
}
