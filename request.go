package snstxtr

import (
	"bytes"
	"io"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"os"

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

func sendText(phone string, msg string, w http.ResponseWriter) {
	log.WithFields(log.Fields{
		"phone": phone,
		"msg": msg,
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

func reqHandler(w http.ResponseWriter, r *http.Request, allowGet bool) {
	logRequest(r)

	// we only want to process real requests i.e. reject robots, favicon etc.
	if r.URL.Path != "/" {
		sendResponse(w, http.StatusMethodNotAllowed, `{"error": "method not allowed or supported"}`)
		return
	}

	switch r.Method {
	case "GET":
		if allowGet {
			phone := r.URL.Query().Get("phone")
			msg := r.URL.Query().Get("msg")

			if len(phone) < 1 || len(msg) < 1 {
				sendResponse(w, http.StatusBadRequest, `{"error": "insufficient parameters"}`)
			} else {
				sendText(phone, msg, w)
			}
		} else {
			sendResponse(w, http.StatusMethodNotAllowed, `{"error": "method not allowed or supported"}`)
		}
	case "POST":
		var bodyBytes []byte
		var payload map[string]interface{}

		// read the body content and unmarshal the expected json
		bodyBytes, _ = ioutil.ReadAll(r.Body)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		json.Unmarshal(bodyBytes, &payload)

		// assume this is a pingdom payload
		if _, ok := payload["check_id"]; ok {
			// we rely on this being set in the env as its not included in the payload
			phone := os.Getenv("PHONE")
			msg := "pingdom: " + payload["check_name"].(string) + " is now " + payload["current_state"].(string)
			sendText(phone, msg, w)
		} else {
			sendText(payload["phone"].(string), payload["msg"].(string), w)
		}
	default:
		sendResponse(w, http.StatusMethodNotAllowed, `{"error": "method not allowed or supported"}`)
	}
}
