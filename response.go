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

func sendText(phone string, msg string, w http.ResponseWriter) {
	log.WithFields(log.Fields{
		"phone": phone,
		"msg":   msg,
	}).Debug("sending sms")

	err := Send(msg, phone)
	if err != nil {
		log.WithFields(log.Fields{
			"info": err,
		}).Error("send result")
		json := "{\"error\": \"" + err.Error() + "\"}"
		sendResponse(w, http.StatusInternalServerError, json)
	} else {
		sendResponse(w, http.StatusOK, `{"status": "sent"}`)
	}
	log.WithFields(log.Fields{
		"info": err,
	}).Debug("send result")
}
