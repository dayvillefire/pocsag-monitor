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
	dbfile = flag.String("db", "scan.db", "Database file")
	rtlfm  = flag.String("rtlfm", "rtl_fm", "Binary path to rtl_fm")
	freq   = flag.String("freq", "152.00750M", "Frequency (append with M, G, etc)")
	mmon   = flag.String("mmon", "multimon-ng", "Binary path to multimon-ng")
	ppm    = flag.String("ppm", "0", "Inaccuracy correction")
)

func main() {
	flag.Parse()

	_, exists := os.Stat(*dbfile)

	DB, err := db.OpenDB(*dbfile)
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

	rtlArg := fmt.Sprintf("-f %s -p %s -s 22050", *freq, *ppm)
	rtlCmd := exec.Command(*rtlfm, strings.Split(rtlArg, " ")...)

	rtlStderr, _ := rtlCmd.StderrPipe()

	mmonArg := fmt.Sprintf("-t raw -a POCSAG512 -f alpha -u /dev/stdin")
	mmonCmd := exec.Command(*mmon, strings.Split(mmonArg, " ")...)

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
			db.Record(DB, alpha)
			continue
		}
	}
	mmonCmd.Wait()
}
