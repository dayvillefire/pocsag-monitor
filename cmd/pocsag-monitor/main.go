package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/coreos/go-systemd/daemon"
	"github.com/dayvillefire/pocsag-monitor/config"
	"github.com/dayvillefire/pocsag-monitor/pocsag"
	"github.com/dayvillefire/pocsag-monitor/sdr"
	"github.com/dayvillefire/pocsag-router/client"
	routerobj "github.com/dayvillefire/pocsag-router/obj"
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
	cfg     *config.Config
	router  *client.Client
	decoder *pocsag.Decoder
)

func parseFrequency(f string) (int, error) {
	f = strings.TrimSuffix(f, "M")
	f = strings.TrimSuffix(f, "m")
	var hz float64
	if _, err := fmt.Sscanf(f, "%f", &hz); err != nil {
		return 0, err
	}
	return int(hz * 1000000), nil
}

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

	// Native SDR + POCSAG pipeline
	log.Printf("INFO: Initializing native SDR source")
	src := sdr.NewRTLSource(cfg.SDR.DeviceIndex)

	freqHz, err := parseFrequency(cfg.Frequency)
	if err != nil {
		log.Fatalf("ERR: invalid frequency %s: %s", cfg.Frequency, err.Error())
	}

	iqCh, err := src.Start(freqHz, cfg.SDR.SampleRate, cfg.PPM, cfg.SDR.Gain)
	if err != nil {
		log.Fatalf("ERR: SDR start failed: %s", err.Error())
	}
	defer src.Stop()

	log.Printf("INFO: Initializing native POCSAG512 decoder")
	decoder = pocsag.NewDecoder()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGQUIT)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGABRT)
	go func(sig chan os.Signal) {
		s := <-sig
		log.Printf("INFO: Caught signal %s, terminating", s.String())

		log.Print("INFO: Terminating pocsag router " + config.GetConfig().InstanceName + " version " + Version + " at " + time.Now().Local().Format(time.RFC3339))

		src.Stop()
		os.Exit(0)
	}(sig)

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

	for alpha := range decoder.Decode(iqCh) {
		if alpha.Valid {
			log.Printf("CAP: %s\tMSG: %s", alpha.CapCode, alpha.Message)
			router.Publish(cfg.Router.Topic, routerobj.AlphaMessage(alpha))
		}
	}
}
