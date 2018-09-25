package main

import (
	"auburn-hacks-landing/metadata"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/glog"
)

var (
	staticDir    *string
	listenAddr   *string
	pollDuration *time.Duration

	meta *metadata.AuburnHacks
)

const (
	MetaFileURL = "https://drive.google.com/uc?id=1AXg6vBbyZ4XR7m8skvQ_062ZEgfGdGvX&export=download"
)

func init() {
	staticDir = flag.String("static_dir", "static", "static files directory")
	listenAddr = flag.String("listen_addr", "localhost:8321", "http serve address")
	pollDuration = flag.Duration("poll_duration", 3*time.Second, "poll duration for metadata watch")

	flag.Set("v", "0")
	flag.Set("logtostderr", "true")

	flag.Parse()
}

func main() {
	http.Handle("/", http.FileServer(http.Dir(*staticDir)))
	http.HandleFunc("/healthz", health)
	http.HandleFunc("/metadata", metadataHandler)

	glog.Infof("starting server on %s pid: %d", *listenAddr, os.Getpid())

	// global metadata object
	meta = metadata.New(MetaFileURL)

	// start watching the file
	go meta.Watch(*pollDuration)

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

func metadataHandler(w http.ResponseWriter, r *http.Request) {
	// allowing parallel reads
	meta.RLock()
	defer meta.RUnlock()
	w.Header().Set("Content-Type", "application/json")

	bb, err := json.Marshal(meta)
	if err != nil {
		http.Error(w, fmt.Sprintf("json marshal error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Write(bb)
}
