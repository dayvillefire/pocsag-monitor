package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/dayvillefire/pocsag-monitor/pocsag"
	"github.com/dayvillefire/pocsag-monitor/sdr"
)

func main() {
	freq := flag.String("f", "152.00750M", "Frequency (e.g. 152.00750M)")
	ppm := flag.Int("p", 0, "PPM frequency correction")
	gain := flag.Int("g", 0, "Tuner gain (0 = auto)")
	rate := flag.Int("r", 22050, "Sample rate")
	dev := flag.Int("d", 0, "Device index")
	flag.Parse()

	freqHz, err := parseFrequency(*freq)
	if err != nil {
		log.Fatalf("invalid frequency %q: %v", *freq, err)
	}

	src := sdr.NewRTLSource(*dev)
	iqCh, err := src.Start(freqHz, *rate, *ppm, *gain)
	if err != nil {
		log.Fatalf("SDR: %v", err)
	}
	defer src.Stop()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		src.Stop()
		os.Exit(0)
	}()

	dec := pocsag.NewDecoder()
	for alpha := range dec.Decode(iqCh) {
		if alpha.Valid {
			fmt.Printf("CAP: %s\tMSG: %s\n", alpha.CapCode, alpha.Message)
		}
	}
}

func parseFrequency(f string) (int, error) {
	f = strings.TrimSuffix(f, "M")
	f = strings.TrimSuffix(f, "m")
	var hz float64
	if _, err := fmt.Sscanf(f, "%f", &hz); err != nil {
		return 0, err
	}
	return int(hz * 1_000_000), nil
}
