package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"gitlab.com/alpinefresh/tcrpartybot/events"
	"log"
	"net/http"
	"os"
)

type incomingDM struct {
	Type             string `json:"type"`
	ID               string `json:"id"`
	CreatedTimestamp string `json:"created_timestamp"`
	MessageCreated   struct {
		SenderID    string `json:"sender_id"`
		MessageData struct {
			Text string `json:"text"`
		} `json:"message_data"`
	} `json:"message_create"`
}

type incomingWebhook struct {
	ForUserID           string       `json:"for_user_id"`
	DirectMessageEvents []incomingDM `json:"direct_message_events"`
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	// A GET request signals that Twitter is attempting a CRC request
	if r.Method == "GET" {
		keys, ok := r.URL.Query()["crc_token"]
		if !ok || len(keys) < 1 {
			w.WriteHeader(400)
			w.Write([]byte("Bad request"))
			return
		}

		mac := hmac.New(sha256.New, []byte(os.Getenv("TWITTER_CONSUMER_SECRET")))
		mac.Write([]byte(keys[0]))

		token := "sha256=" + base64.StdEncoding.EncodeToString(mac.Sum(nil))
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("{\"response_token\": \"" + token + "\"}"))
		return
	}

	decoder := json.NewDecoder(r.Body)
	data := &incomingWebhook{}
	err := decoder.Decode(data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(400)
		w.Write([]byte("Bad request"))
		return
	}

	log.Println(data.DirectMessageEvents[0])
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

// StartServer spins up a webserver for the API
func StartServer(eventsChan chan<- *events.Event, errChan chan<- error) {
	http.HandleFunc("/webhooks/direct-message", handleWebhook)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		errChan <- err
	}
}
