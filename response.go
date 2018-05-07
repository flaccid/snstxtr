package snstxtr

import (
	"io"
	"net/http"

	log "github.com/Sirupsen/logrus"
)

func sendResponse(w http.ResponseWriter, status int, body string) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, body)
}

func sendText(recipients []string, msg string, w http.ResponseWriter) {
	var sendFail bool
	var err error

	log.WithFields(log.Fields{
		"recipients": recipients,
		"msg":   msg,
	}).Debug("sending sms")

	// send sms to each recipient
	for _, phone := range recipients {
		err = Send(msg, phone)
		if err != nil {
			sendFail = true
		}
		log.WithFields(log.Fields{
			"info": err,
		}).Debug("send result")
	}

	if sendFail {
		json := "{\"error\": \"" + err.Error() + "\"}"
		sendResponse(w, http.StatusInternalServerError, json)
	} else {
		sendResponse(w, http.StatusOK, `{"status": "all-sent"}`)
	}
}
