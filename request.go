package snstxtr

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type SendPayload struct {
	Recipients []string `json:"recipients"`
	Message    string   `json:"msg"`
}

func reqHandler(w http.ResponseWriter, r *http.Request, allowGet bool) {
	logRequest(r)

	// we only want to process real requests i.e. reject robots, favicon etc.
	if r.URL.Path != "/" {
		sendResponse(w, http.StatusMethodNotAllowed, `{"error": "method not allowed or supported"}`)
		return
	}

	var jsonResponse string

	switch r.Method {
	case "GET":
		if allowGet {
			msg := r.URL.Query().Get("msg")
			recipients := strings.Split(r.URL.Query().Get("recipients"), ",")

			if len(recipients) < 1 || len(msg) < 1 {
				sendResponse(w, http.StatusBadRequest, `{"error": "insufficient parameters"}`)
			} else {
				sendText(recipients, msg, w)
			}
		} else {
			sendResponse(w, http.StatusMethodNotAllowed, `{"error": "method not allowed or supported"}`)
		}
	case "POST":
		var bodyBytes []byte
		var payload SendPayload

		// read the body content and unmarshal the expected json
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			jsonResponse = `{"error": "internal server error"}`
			sendResponse(w, http.StatusInternalServerError, jsonResponse)
			log.Error("error reading request body: ", err)
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		err = json.Unmarshal(bodyBytes, &payload)
		if err != nil {
			log.Error("error unmarshalling json: ", err)
		}

		if len(payload.Recipients) < 1 {
			sendResponse(w, http.StatusBadRequest, `{"error": "insufficient parameters"}`)
		} else {
			sendText(payload.Recipients, payload.Message, w)
		}
	default:
		sendResponse(w, http.StatusMethodNotAllowed, `{"error": "method not allowed or supported"}`)
	}
}
