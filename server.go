package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	staticDir  *string
	listenAddr *string
)

func init() {
	staticDir = flag.String("static_dir", "static", "static files directory")
	listenAddr = flag.String("listen_addr", "localhost:8321", "http serve address")

	flag.Parse()
}

func main() {
	http.Handle("/", http.FileServer(http.Dir(*staticDir)))
	http.HandleFunc("/healthz", health)

	log.Printf("server listening on %s pid: %d", *listenAddr, os.Getpid())
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, os.Interrupt)
	go func() {
		if err := http.ListenAndServe(*listenAddr, nil); err != nil {
			log.Fatalf("error serving: %v", err)
		}
	}()
	<-c // blocks until SIGTERM is received
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}
