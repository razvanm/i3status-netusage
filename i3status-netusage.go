package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	iface = flag.String("interface", "eth0", "What interface to use")
)

// stats fetches the cumulative rx/tx bytes for network interface
// iface.
func stats() (rx, tx uint64) {
	b, err := ioutil.ReadFile("/proc/net/dev")
	if err != nil {
		return 0, 0
	}
	buff := bytes.NewBuffer(b)
	for l, err := buff.ReadString('\n'); err == nil; {
		l = strings.Trim(l, " \n")
		if !strings.HasPrefix(l, *iface) {
			l, err = buff.ReadString('\n')
			continue
		}
		re := regexp.MustCompile(" +")
		s := strings.Split(re.ReplaceAllString(l, " "), " ")
		rx, err := strconv.ParseUint(s[1], 10, 64)
		if err != nil {
			return 0, 0
		}
		tx, err := strconv.ParseUint(s[9], 10, 64)
		if err != nil {
			return 0, 0
		}
		return rx, tx
	}
	return 0, 0
}

// humanize converts a number of bytes in KiB or MiB
func humanize(i uint64) string {
	if i < 1024 {
		return fmt.Sprintf("%d bytes", float64(i))
	}
	if i < 1024*1024 {
		return fmt.Sprintf("%.1f KiB", float64(i)/1024)
	}
	return fmt.Sprintf("%.1f MiB", float64(i)/1024/1024)
}

func main() {
	flag.Parse()

	prevRx, prevTx := uint64(0), uint64(0)
	bio := bufio.NewReader(os.Stdin)
	for {
		line, err := bio.ReadString('\n')
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if !(strings.HasPrefix(line, ",[{") && strings.HasSuffix(line, "}]\n")) {
			fmt.Print(line)
			continue
		}
		line = strings.TrimSuffix(line, "]\n")
		rx, tx := stats()
		fmt.Printf("%s,{\"full_text\":\"%10s/s↓ %10s/s↑\"}]\n", line, humanize(rx-prevRx), humanize(tx-prevTx))
		prevRx, prevTx = rx, tx
	}

}
