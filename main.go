package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/genjidb/genji"
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

	db, err := genji.Open(*dbfile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db = db.WithContext(context.Background())

	err = initDb(db)
	if err != nil {
		log.Printf("ERR: %s", err.Error())
	}

	rtlArg := fmt.Sprintf("-f %s -p %s -s 22050", *freq, *ppm)
	rtlCmd := exec.Command(*rtlfm, strings.Split(rtlArg, " ")...)

	stderr, _ := rtlCmd.StderrPipe()
	//rtlCmd.Start()

	mmonArg := fmt.Sprintf("-t raw -a POCSAG512 -f alpha /dev/stdin")
	mmonCmd := exec.Command(*mmon, strings.Split(mmonArg, " ")...)

	mmonCmd.Stdin, _ = rtlCmd.StdoutPipe()
	mmonCmd.Stdout = os.Stdout
	err = mmonCmd.Start()
	if err != nil {
		panic(err)
	}
	err = rtlCmd.Start()
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		ts := time.Now()
		alpha, err := parse(ts, m)
		if err != nil {
			record(db, alpha)
		}
	}
	mmonCmd.Wait()

	/*
		    precmd := exec.Command("cat", "scan.log")
		    precmdStderr, _ := precmd.StdoutPipe()
			cmd := exec.Command("grep", "POCSAG:")
		    stdin, err := cmd.StdinPipe()
			stdout, err := cmd.StdoutPipe()

			if err != nil {
				log.Fatal(err)
			}

			cmd.Start()

			buf := bufio.NewReader(stdout)
			num := 0

			for {
				line, _, _ := buf.ReadLine()
				if num > 3 {
					os.Exit(0)
				}
				num += 1
				fmt.Println(string(line))
			}
	*/

	/*
		scanner := bufio.NewScanner(&buf)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			m := scanner.Text()
			fmt.Printf("%s: %s\n", time.Now().Format("2006-01-02 15:04:05"), m)
		}
	*/
}

func initDb(db *genji.DB) error {
	return db.Exec(`
    CREATE TABLE entry (
        timestamp       INT     PRIMARY KEY,
        cap             TEXT    NOT NULL,
        message         TEXT    NOT NULL,
        INDEX(cap)
    )
`)
}

func record(db *genji.DB, alpha alphaMessage) error {
	return db.Exec(`INSERT INTO entry ( timestamp, cap, message ) VALUES ( ?, ?, ? )`, alpha.Timestamp.Unix(), alpha.CapCode, alpha.Message)
}
