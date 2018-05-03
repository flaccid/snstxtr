package snstxtr

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

const (
	listenPort = 8080
)

func Serve() {
	log.Printf("initialize snstxtr")
	log.Debug("debug logging enabled")

	http.HandleFunc("/", reqHandler)
	http.HandleFunc("/health", healthCheckHandler)
	http.HandleFunc("/health/", healthCheckHandler)

	log.Info("listening for requests on :" + fmt.Sprintf("%v", listenPort))
	if err := http.ListenAndServe(":"+fmt.Sprintf("%v", listenPort), nil); err != nil {
		log.Fatalf("http; %v", err)
	}
}
