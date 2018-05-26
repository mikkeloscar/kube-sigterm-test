package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var healthy = true

func mainHandler(w http.ResponseWriter, r *http.Request) {
	if !healthy {
		log.Println("request, after sigterm")
	} else {
		log.Println("request")
	}
	w.WriteHeader(200)
	w.Write([]byte("hello"))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	typ := r.URL.Query().Get("type")
	if !healthy {
		log.Printf("not healthy - %s", typ)
		w.WriteHeader(503)
		w.Write([]byte("not ready"))
	} else {
		log.Printf("healthy - %s", typ)
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}
}

var (
	keepAlive      bool
	sigtermTimeout time.Duration
)

func main() {
	flag.BoolVar(&keepAlive, "keep-alive", false, "Don't disable keep-alives after SIGTERM.")
	flag.DurationVar(&sigtermTimeout, "sigterm-timeout", 20*time.Second, "How long to wait after SIGTERM before terminating.")
	flag.Parse()

	sigTerm := make(chan os.Signal, 1)
	signal.Notify(sigTerm, syscall.SIGTERM)

	mux := http.NewServeMux()
	mux.HandleFunc("/", mainHandler)
	mux.HandleFunc("/healthz", healthHandler)

	server := &http.Server{
		Addr:        ":8080",
		Handler:     mux,
		ReadTimeout: 5 * time.Second,
	}

	go server.ListenAndServe()

	<-sigTerm
	log.Println("got SIGTERM, prepare to shut down in 20s")
	healthy = false
	if !keepAlive {
		log.Println("SetKeepAlivesEnabled=false")
		server.SetKeepAlivesEnabled(false)
	}
	time.Sleep(sigtermTimeout)
	log.Println("exiting")
}
