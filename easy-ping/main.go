package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/Meepoljdx/easy-tools/utils"
)

var (
	help     = flag.Bool("h", false, "Command usage")
	ip       = flag.String("ip", "127.0.0.1", "If set ip, you can use ip1,ip2,ip3 to specify the server on which the ping test is to be performed.")
	packNum  = flag.Int("packet-num", 10, "The num of packets will be send to remote server.")
	packSize = flag.Int("packet-size", 64, "The size of packets will be send to remote server.")
	file     = flag.String("file", "", "The ip flag will be ignored if file has been set, you can write all ip which the ping test is to be performed in a file.")
	output   = flag.String("output", "stdout", "The output of the result, you can set stdout, csv, json, excel.")
	channl   = flag.Int("channle", 30, "Number of pings performed at the same time")
)

func usage() {
	fmt.Fprintf(os.Stderr, `Use Easy ping to send icmp packet
Author: Meepo

Options:
`)
	flag.PrintDefaults()
}

func init() {
	flag.Usage = usage
}

func ServerPing(list []string, t string) *Result {
	lens := len(list)
	wg := sync.WaitGroup{}
	o := NewResult("test", t)
	if lens < *channl {
		*channl = lens
	}
	ch := make(chan string, *channl)
	// push to ch
	wg.Add(1 + *channl)
	go func(l []string) {
		for _, v := range l {
			ch <- v
		}
		close(ch)
		wg.Done()
	}(list)

	for i := 0; i < *channl; i++ {
		go PingIP(ch, &wg, o)
	}

	wg.Wait()

	return o
}

func PingIP(ch <-chan string, wg *sync.WaitGroup, result *Result) {
	defer wg.Done()
	for {
		ip, ok := <-ch
		if !ok {
			break
		}
		p := NewPing(ip, *packNum, *packSize)

		if err := p.Run(); err != nil {
			continue
		}
		result.Lock.Lock()
		result.Output = append(result.Output, *p)
		result.Lock.Unlock()
	}
}

func main() {
	flag.Parse()

	if *help {
		flag.Usage()
	}
	var ipList []string

	if *file != "" && utils.FileExisted(*file) {
		var err error
		ipList, err = utils.ReadIPFromFile(*file)
		if err != nil {
			fmt.Println("Read ip from file failed.")
		}
	} else {
		// 处理命令行的ip
		ipList = strings.Split(*ip, ",")
	}

	o := ServerPing(ipList, *output)
	o.ResultOutPut()
}
