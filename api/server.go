package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"gitlab.com/alpinefresh/tcrpartybot/events"
	"gitlab.com/alpinefresh/tcrpartybot/models"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
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

type user struct {
	ID         int64  `json:"id"`
	ScreenName string `json:"screen_name"`
}

type incomingTweet struct {
	Text string `json:"text"`
	User user   `json:"user"`
}

type incomingFollow struct {
	Source struct {
		ScreenName string `json:"screen_name"`
		ID         string `json:"id"`
	} `json:"source"`
}

type incomingWebhook struct {
	ForUserID           string           `json:"for_user_id"`
	DirectMessageEvents []incomingDM     `json:"direct_message_events"`
	TweetCreateEvents   []incomingTweet  `json:"tweet_create_events"`
	FollowEvents        []incomingFollow `json:"follow_events"`
}

// Server holds relevant data for running an API server
type Server struct {
	errChan    chan<- error
	eventsChan chan<- *events.Event
}

func (server *Server) processDMs(dms []incomingDM) {
	botToken, err := models.FindOAuthTokenByHandle(os.Getenv("VIP_BOT_HANDLE"))
	if err != nil || botToken == nil {
		return
	}

	for _, dm := range dms {
		fromID, err := strconv.ParseInt(dm.MessageCreated.SenderID, 10, 64)
		if err != nil {
			server.errChan <- err
			continue
		}

		// Ignore outgoing DM events
		if fromID == botToken.TwitterID {
			continue
		}

		account, err := models.FindAccountByID(fromID)
		if err != nil {
			server.errChan <- err
			continue
		} else if account == nil {
			server.errChan <- errors.New("Could not find account for incoming DM")
			continue
		}

		server.eventsChan <- &events.Event{
			EventType:    events.EventTypeDM,
			Time:         time.Now(),
			SourceHandle: account.TwitterHandle,
			SourceID:     account.TwitterID,
			Message:      dm.MessageCreated.MessageData.Text,
		}
	}
}

func (server *Server) processMentions(tweets []incomingTweet) {
	for _, tweet := range tweets {
		server.eventsChan <- &events.Event{
			EventType:    events.EventTypeMention,
			Time:         time.Now(),
			SourceHandle: tweet.User.ScreenName,
			SourceID:     tweet.User.ID,
			Message:      tweet.Text,
		}
	}
}

func (server *Server) processFollows(follows []incomingFollow) {
	botToken, err := models.FindOAuthTokenByHandle(os.Getenv("VIP_BOT_HANDLE"))
	if err != nil || botToken == nil {
		return
	}

	for _, follow := range follows {
		id, err := strconv.ParseInt(follow.Source.ID, 10, 64)
		if err != nil {
			server.errChan <- err
			continue
		}

		// Ignore outgoing DM events
		if id == botToken.TwitterID {
			continue
		}

		server.eventsChan <- &events.Event{
			EventType:    events.EventTypeFollow,
			Time:         time.Now(),
			SourceHandle: follow.Source.ScreenName,
			SourceID:     id,
		}
	}
}

func (server *Server) handleTwitterWebhook(w http.ResponseWriter, r *http.Request) {
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

	if len(data.DirectMessageEvents) > 0 {
		go server.processDMs(data.DirectMessageEvents)
	}

	if len(data.TweetCreateEvents) > 0 {
		go server.processMentions(data.TweetCreateEvents)
	}

	if len(data.FollowEvents) > 0 {
		go server.processFollows(data.FollowEvents)
	}

	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

// StartServer spins up a webserver for the API
func StartServer(eventsChan chan<- *events.Event, errChan chan<- error) *Server {
	server := &Server{
		eventsChan: eventsChan,
		errChan:    errChan,
	}

	http.HandleFunc("/webhooks/twitter", server.handleTwitterWebhook)
	err := http.ListenAndServe(os.Getenv("SERVER_HOST"), nil)
	if err != nil {
		errChan <- err
	}

	return server
}