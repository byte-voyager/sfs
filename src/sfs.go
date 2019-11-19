package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
	"github.com/toolkits/net"
)

func showQrcode(host string, ips ...string) {

	for i, v := range ips {
		fmt.Printf("IP %d : %s\n", i+1, v)
	}

	fmt.Printf("Enter the ip number to display the QR code:\n")
	var numStr string
	fmt.Scanln(&numStr)
	num, err := strconv.Atoi(numStr)
	if err != nil {
		println(numStr)
		fmt.Println("Invalid num!")
		return
	}
	if num > len(ips) {
		fmt.Printf("Error, There are at most %d IP\n", len(ips))
		return
	}
	index := strings.LastIndex(host, ":")
	if index == -1 {
		fmt.Printf("Error, the host: %s is invalid\n", host)
		return
	}

	port := host[index:len(host)]

	content := "http://" + ips[num-1] + port

	obj := qrcodeTerminal.New()
	obj.Get(content).Print()
	fmt.Printf("URL: %s\n", content)
}

func main() {
	var host = flag.String("a", "0.0.0.0:9000", "server bind address")
	var shareDir = flag.String("d", "./", "directory to share")
	var disable = flag.Bool("i", false, "ignore input(when use nohup)")
	var maxCount = flag.Uint("m", 10, "input max count")
	flag.Parse()

	handler := http.FileServer(http.Dir(*shareDir))
	ips, _ := net.IntranetIP()

	http.Handle("/", handler)

	fmt.Println("Share File Server listing on", *host, *shareDir, *disable)

	if !*disable {
		go func() {
			for {
				if *maxCount <= 0 {
					return
				}
				*maxCount--
				showQrcode(*host, ips...)
			}
		}()
	}

	err := http.ListenAndServe(*host, nil)
	if err != nil {
		fmt.Printf("Error, %s\n", err)
	}
}
