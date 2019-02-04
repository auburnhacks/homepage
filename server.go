package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang/glog"
)

var (
	staticDir *string
)

func init() {
	staticDir = flag.String("static_dir", "static", "static files directory")

	flag.Set("v", "0")
	flag.Set("logtostderr", "true")

	flag.Parse()
}

func main() {
	http.Handle("/", http.FileServer(http.Dir(*staticDir)))
	http.HandleFunc("/healthz", health)
	//http.HandleFunc("/metadata", metadataHandler)

	addr := ""
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
		addr = "localhost:" + port
	} else {
		addr = ":" + port
	}

	glog.Infof("starting server on %s pid: %d", addr, os.Getpid())

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGTERM, os.Interrupt)
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalf("error serving: %v", err)
		}
	}()
	<-c // blocks until SIGTERM is received
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("ok"))
}
