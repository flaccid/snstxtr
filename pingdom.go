package snstxtr

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
)

func pingdomHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	switch r.Method {
	case "GET":
		// pingdom only supports webhook POST payloads
		sendResponse(w, http.StatusMethodNotAllowed, `{"error": "method not allowed or supported"}`)
	case "POST":
		var bodyBytes []byte
		var payload map[string]interface{}

		// read the body content and unmarshal the expected json
		bodyBytes, _ = ioutil.ReadAll(r.Body)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		json.Unmarshal(bodyBytes, &payload)

		// this should be a pingdom payload
		if _, ok := payload["check_id"]; ok {
			// we rely on this being set in the env as its not included in the payload
			phone := os.Getenv("PHONE")
			msg := "pingdom: " + payload["check_name"].(string) + " is now " + payload["current_state"].(string)
			sendText(phone, msg, w)
		} else {
			sendResponse(w, http.StatusMethodNotAllowed, `{"error": "method not allowed or supported"}`)
		}
	default:
		sendResponse(w, http.StatusMethodNotAllowed, `{"error": "method not allowed or supported"}`)
	}
}
