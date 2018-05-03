package snstxtr

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func Serve() {
	log.Printf("start snstxtr")
	log.Debug("debug logging enabled")

	http.HandleFunc("/", reqHandler)
	http.HandleFunc("/health", healthCheckHandler)
	http.HandleFunc("/health/", healthCheckHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("http; %v", err)
	}
}
