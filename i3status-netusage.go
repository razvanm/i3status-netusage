// Public domain.
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
	"time"
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

// humanize converts a number of bytes in KiB or MiB.
func humanize(i float64) string {
	if i < 1024 {
		return fmt.Sprintf("%.0f bytes", i)
	}
	if i < 1024*1024 {
		return fmt.Sprintf("%.1f KiB", i/1024)
	}
	return fmt.Sprintf("%.1f MiB", i/1024/1024)
}

func main() {
	flag.Parse()

	prevRx, prevTx := uint64(0), uint64(0)
	bio := bufio.NewReader(os.Stdin)
	prev := time.Now()
	for {
		line, err := bio.ReadString('\n')
		now := time.Now()
		window := now.Sub(prev).Seconds()
		prev = now
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
		rxRate := humanize(float64(rx-prevRx) / window)
		txRate := humanize(float64(tx-prevTx) / window)
		fmt.Printf("%s,{\"full_text\":\"%10s/s↓ %10s/s↑\"}]\n", line, rxRate, txRate)
		prevRx, prevTx = rx, tx
	}

}
