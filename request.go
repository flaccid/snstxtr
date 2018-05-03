package snstxtr

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

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

		// todo: add check for required params
		sendText(payload["phone"].(string), payload["msg"].(string), w)
	default:
		sendResponse(w, http.StatusMethodNotAllowed, `{"error": "method not allowed or supported"}`)
	}
}
