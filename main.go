package main

import (
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
	"time"
)

type EventType string

/*
* polls: [from blockchain]
*	active challenges
* accounts:
*	twitter_handle
*	eth_address
*	private_key [not on blockchain]
*
 */

const (
	EventTypeMention       EventType = "EventTypeMention"
	EventTypeDM            EventType = "EventTypeDM"
	EventTypeVote          EventType = "EventTypeVote"
	EventTypePollCompleted EventType = "EventTypePollCompleted"
)

type TCRContract interface {
	Apply(applicant string, nominee string, amount uint) (string, error)
	Vote(voter string, nominee string, accept bool) error
	Challenge(challengee string, challenger string) error
	GetExpiry(nominee string) (time.Time, error)
}

type Event struct {
	EventType    EventType // type of event
	Time         time.Time // timestamp
	SourceHandle string    // twitter handle sending
	Message      string    // whole message
}

type TwitterCredentials struct {
	ConsumerKey    string
	ConsumerSecret string
}

func listenToTwitter(twitterCredentials *TwitterCredentials, eventChan chan<- *Event, errorChan chan<- error) {
	// listen for DMs
	// listen for @replies
	// listen for votes
	for {
		time.Sleep(1 * time.Second)
		eventChan <- &Event{
			EventType:    EventTypeMention,
			Time:         time.Now(),
			SourceHandle: "stevenleeg",
			Message:      "Let's party, @TCRPartyVIP",
		}
	}
}

func processEvents(eventChan <-chan *Event, errorChan chan<- error) {
	log.Println("processing events")
	for {
		event := <-eventChan
		switch event.EventType {
		case EventTypeMention:
			log.Println("Got a mention!")
			break
		}
	}
}

func serveAPI() {
	log.Println("serving api")
}

func logErrors(errorChan <-chan error) {
	for {
		log.Println(<-errorChan)
	}
}

func main() {
	twitterCredentials := &TwitterCredentials{
		ConsumerKey:    os.Getenv("TWITTER_CONSUMER_KEY"),
		ConsumerSecret: os.Getenv("TWITTER_CONSUMER_SECRET"),
	}

	eventChan := make(chan *Event)
	errorChan := make(chan error)
	go listenToTwitter(twitterCredentials, eventChan, errorChan)
	go processEvents(eventChan, errorChan)
	//go serveAPI()

	for range eventChan {
	}
}
