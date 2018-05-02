package snstxtr

import (
	"io"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func logRequest(r *http.Request) {
	// we only log in debug mode due to possible exposure of PI data in request uri
	log.WithFields(log.Fields{
		"method": r.Method,
	}).Debug(r.URL.String())
}

func sendResponse(w http.ResponseWriter, status int, body string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, body)
}

func reqHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	// we only want to process real requests i.e. reject robots, favicon etc.
	if r.URL.Path != "/" {
		sendResponse(w, http.StatusMethodNotAllowed, `{"error": "method not allowed or supported"}`)
		return
	}

	switch r.Method {
	case "GET":
		phone := r.URL.Query().Get("phone")
		msg := r.URL.Query().Get("msg")

		if len(phone) < 1 || len(msg) < 1 {
			sendResponse(w, http.StatusBadRequest, `{"error": "insufficient parameters"}`)
		} else {
			log.Debug("sending sms to ", phone)
			err := Send(msg, phone)
			if err != nil {
				log.Error(err)
				json := "{\"error\": \"" + err.Error() + "\"}"
				sendResponse(w, http.StatusInternalServerError, json)
			} else {
				log.Info(err)
				sendResponse(w, http.StatusOK, `{"status": "sent"}`)
			}
		}
	case "POST":
		// todo: use for handling pingdom webhook payloads
		sendResponse(w, http.StatusMethodNotAllowed, `{"error": "method not allowed or supported"}`)
	default:
		sendResponse(w, http.StatusMethodNotAllowed, `{"error": "method not allowed or supported"}`)
	}
}
