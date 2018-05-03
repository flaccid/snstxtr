package snstxtr

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

const (
	listenPort = 8080
)

var (
	w http.ResponseWriter
	r *http.Request
)

func Serve(allowGet bool) {
	log.Printf("initialize snstxtr")
	log.Debug("debug logging enabled")
	if allowGet {
		log.Debug("allowing get requests")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		reqHandler(w, r, allowGet)
	})
	http.HandleFunc("/health", healthCheckHandler)
	http.HandleFunc("/health/", healthCheckHandler)
	http.HandleFunc("/pingdom-webhook", pingdomHandler)
	http.HandleFunc("/pingdom-webhook/", pingdomHandler)

	log.Info("listening for requests on :" + fmt.Sprintf("%v", listenPort))
	if err := http.ListenAndServe(":"+fmt.Sprintf("%v", listenPort), nil); err != nil {
		log.Fatalf("http; %v", err)
	}
}
