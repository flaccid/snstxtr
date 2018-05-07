package snstxtr

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type PingdomPayload struct {
	CheckId               float64                `json:"check_id"`
	CheckName             string                 `json:"check_name"`
	CheckType             string                 `json:"check_type"`
	CheckParams           map[string]interface{} `json:"check_params"`
	Tags                  []string               `json:"tags"`
	PreviousState         string                 `json:"previous_state"`
	CurrentState          string                 `json:"current_state"`
	ImportanceLevel       string                 `json:"importance_level"`
	StateChangedTimestamp float64                `json:"state_changed_timestamp"`
	StateChangedUtcTime   string                 `json:"state_changed_utc_time"`
	LongDescription       string                 `json:"long_description"`
	Description           string                 `json:"description"`
	FirstProbe            map[string]interface{} `json:first_probe`
	SecondProbe           map[string]interface{} `json:second_probe`
}

func pingdomHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	var jsonResponse string

	switch r.Method {
	case "POST":
		// recipients must be specified in the query string
		if len(r.URL.Query().Get("recipients")) < 1 {
			jsonResponse = `{"error": "insufficient parameters", "reason": "no recipients provided"}`
			sendResponse(w, http.StatusBadRequest, jsonResponse)
			return
		}

		// get the recipients
		recipients := strings.Split(r.URL.Query().Get("recipients"), ",")

		var bodyBytes []byte
		var payload PingdomPayload

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

		// this must be a pingdom payload
		if payload.CheckId > 0 {
			msg := "pingdom: " + payload.CheckName + " is now " + payload.CurrentState

			// send sms to each recipient
			sendText(recipients, msg, w)
		} else {
			log.Error("pingdom check_id not found in payload")
			jsonResponse = `{"error": "method not allowed or supported"}`
			sendResponse(w, http.StatusMethodNotAllowed, jsonResponse)
		}
	default:
		// pingdom only supports webhook POST payloads
		jsonResponse = `{"error": "method not allowed or supported"}`
		sendResponse(w, http.StatusMethodNotAllowed, jsonResponse)
	}
}
