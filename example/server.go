package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/devsisters/goquic"
	"github.com/devsisters/gospdyquic"
)

var numOfServers int
var port int
var serveRoot string
var logLevel int

func httpHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("This is an example server.\n"))
}

func init() {
	flag.IntVar(&numOfServers, "n", 1, "Number of concurrent quic dispatchers")
	flag.IntVar(&port, "port", 8080, "TCP/UDP port number to listen")
	flag.StringVar(&serveRoot, "root", "/tmp", "Root of path to serve under https://127.0.0.1/files/")
	flag.IntVar(&logLevel, "loglevel", -1, "Log level")
}

func main() {
	goquic.Initialize()
	goquic.SetLogLevel(logLevel)

	flag.Parse()

	log.Printf("About to listen on %d. Go to https://127.0.0.1:%d/", port, port)
	portStr := fmt.Sprintf(":%d", port)

	http.HandleFunc("/", httpHandler)
	http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir(serveRoot))))
	err := gospdyquic.ListenAndServeSecure(portStr, "cert.pem", "key.pem", numOfServers, nil)
	if err != nil {
		log.Fatal(err)
	}
}
