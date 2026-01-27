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
	"syscall"
	"time"

	"github.com/coreos/go-systemd/daemon"
	"github.com/dayvillefire/pocsag-monitor/config"
	"github.com/dayvillefire/pocsag-router/client"
	"github.com/dayvillefire/pocsag-router/obj"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	configFile        = flag.String("config", "config.yaml", "Configuration file")
	dynamicConfigFile = flag.String("dynamic-config", "dynamic.yaml", "Dynamic configuration file")
	testConfig        = flag.Bool("test-config", false, "Test config")
	daemonize         = flag.Bool("daemon", false, "Daemonize")

	Version string
	//logRoute string
	cfg    *config.Config
	router *client.Client
)

func main() {
	flag.Parse()

	var err error

	err = godotenv.Load()
	if err != nil {
		panic(err)
	}

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

	log.Printf("INFO: Connecting to router at %s", cfg.Router.URL)
	router, err = client.NewClient(
		cfg.Router.URL,
		client.ClientTLSConfig{
			ClientCert: os.Getenv("CLIENT_CERT_FILE"),
			ClientKey:  os.Getenv("CLIENT_KEY_FILE"),
			RootCA:     os.Getenv("CA_CERT"),
		})
	if err != nil {
		panic(err)
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
		log.Printf("INFO: Caught signal %s, terminating", s.String())

		log.Print("INFO: Terminating pocsag router " + config.GetConfig().InstanceName + " version " + Version + " at " + time.Now().Local().Format(time.RFC3339))

		rtlCmd.Process.Kill()
		mmonCmd.Process.Kill()
	}(sig, rtlCmd, mmonCmd)
	defer func(rtlCmd *exec.Cmd) {
		// If, for some reason, this doesn't die gracefully, kill it with fire
		log.Printf("INFO: Non-gracefully terminating rtl_fm")
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

	log.Print("INFO: Initialized pocsag router " + config.GetConfig().InstanceName + " version " + Version + " at " + time.Now().Local().Format(time.RFC3339))

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
			log.Printf("ERR: ParseAlphaMessage: %s", err.Error())
			continue
		}
		if alpha.Valid {
			log.Printf("CAP: %s\tMSG: %s", alpha.CapCode, alpha.Message)
			// transmit
			router.Publish(cfg.Router.Topic, alpha)
			continue
		}
	}
	mmonCmd.Wait()
}
