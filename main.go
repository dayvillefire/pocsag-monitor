package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/dayvillefire/pocsag-monitor/db"
	"github.com/dayvillefire/pocsag-monitor/obj"
)

var (
	configFile = flag.String("config", "config.yaml", "Configuration file")
)

func main() {
	flag.Parse()

	config, err := LoadConfigWithDefaults(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	_, exists := os.Stat(config.DbFile)

	DB, err := db.OpenDB(config.DbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	if exists != nil {
		err = db.InitDB(DB)
		if err != nil {
			log.Printf("ERR: %s", err.Error())
		}
	}

	_, err = getDiscordClient(config.DiscordToken)
	if err != nil {
		panic(err)
	}

	rtlArg := fmt.Sprintf("-f %s -p %d -s 22050", config.Frequency, config.PPM)
	rtlCmd := exec.Command(config.RtlFmBinary, strings.Split(rtlArg, " ")...)

	rtlStderr, _ := rtlCmd.StderrPipe()

	mmonArg := fmt.Sprintf("-t raw -a POCSAG512 -f alpha -u /dev/stdin")
	mmonCmd := exec.Command(config.MultiMonBinary, strings.Split(mmonArg, " ")...)

	mmonCmd.Stdin, _ = rtlCmd.StdoutPipe()
	//mmonCmd.Stdout = os.Stdout
	stdout, err := mmonCmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	//stderr, err := mmonCmd.StderrPipe()
	//if err != nil {
	//	panic(err)
	//}
	err = mmonCmd.Start()
	if err != nil {
		panic(err)
	}
	err = rtlCmd.Start()
	if err != nil {
		panic(err)
	}
	defer stdout.Close()
	//defer stderr.Close()
	defer rtlStderr.Close()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-sig
		log.Printf("Caught signal %s, terminating", s.String())
		rtlCmd.Process.Kill()
		mmonCmd.Process.Kill()
	}()

	router := Router{config.ChannelMappings}

	scanner := bufio.NewScanner(io.MultiReader(stdout, rtlStderr))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		//log.Printf("DEBUG: Found line '%s'", m)
		ts := time.Now()
		alpha, err := obj.ParseAlphaMessage(ts, m)
		if err != nil {
			log.Printf("ERR: %s", err.Error())
			continue
		}
		if alpha.Valid {
			log.Printf("CAP: %s\tMSG: %s", alpha.CapCode, alpha.Message)
			dest := router.MapMessage(alpha)
			for _, c := range dest {
				sendDiscordMessage(
					config.DiscordChannels[c],
					fmt.Sprintf(
						"%s: %s [%s]",
						alpha.CapCode,
						alpha.Message,
						alpha.Timestamp.Format("2006-01-02 15:04:05"),
					),
				)
			}
			db.Record(DB, alpha)
			continue
		}
	}
	mmonCmd.Wait()
}
