package main

import (
	"fmt"
	"net/http"
	"os"
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
	host := ":9000"
	shareDir := "./"
	argLen := len(os.Args)
	if argLen == 2 {
		if os.Args[1] == "-h" || os.Args[1] == "--help" || os.Args[1] == "help" {
			fmt.Printf("Usage: sfs [host:port] [share_dir]\n")
			fmt.Printf("       sfs [share_dir]\n")
			fmt.Printf("Example: sfs 0.0.0.0:9000 ./\n")
			fmt.Printf("         sfs :9000 ./\n")
			fmt.Printf("         sfs ./\n")
			return
		}

		// 设置dir
		shareDir = os.Args[1]
		s, err := os.Stat(shareDir)
		if err != nil {
			fmt.Printf("Error, %s not a file or dir\n", shareDir)
			return
		}
		if !s.IsDir() {
			fmt.Printf("Error, %s not a dir\n", shareDir)
			return
		}

	} else if argLen == 3 {
		// 设置dir和host
		host = os.Args[1]
		shareDir = os.Args[2]
	}
	handler := http.FileServer(http.Dir(shareDir))
	ips, _ := net.IntranetIP()

	http.Handle("/", handler)

	fmt.Println("Share File Server listing on", host, shareDir)
	go func() {
		for {
			showQrcode(host, ips...)
		}
	}()
	err := http.ListenAndServe(host, nil)
	if err != nil {
		fmt.Printf("Error, %s\n", err)
	}
}
