package snstxtr

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func pingdomHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	switch r.Method {
	case "GET":
		// pingdom only supports webhook POST payloads
		sendResponse(w, http.StatusMethodNotAllowed, `{"error": "method not allowed or supported"}`)
	case "POST":
		recipients := r.URL.Query().Get("recipients")

		// check for required query params first
		if len(recipients) < 1 {
			sendResponse(w, http.StatusBadRequest, `{"error": "insufficient parameters", "reason": "no recipients provided"}`)
			return
		}

		var bodyBytes []byte
		var payload map[string]interface{}

		// read the body content and unmarshal the expected json
		bodyBytes, _ = ioutil.ReadAll(r.Body)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		json.Unmarshal(bodyBytes, &payload)

		// this should be a pingdom payload
		if _, ok := payload["check_id"]; ok {
			msg := "pingdom: " + payload["check_name"].(string) + " is now " + payload["current_state"].(string)
			sendText(recipients, msg, w)
		} else {
			sendResponse(w, http.StatusMethodNotAllowed, `{"error": "method not allowed or supported"}`)
		}
	default:
		sendResponse(w, http.StatusMethodNotAllowed, `{"error": "method not allowed or supported"}`)
	}
}
