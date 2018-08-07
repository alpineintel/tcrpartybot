package main

import (
	"log"
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

type FeedSource interface {
	GetEventFeed(since time.Time) chan Event
	GetError() error
}

func listenToTwitter(apiKey string, eventChan chan<- Event, errorChan chan<- error) {
}

func processEvents(eventChan <-chan Event, errorChan chan<- error) {
	for {
		event := <-eventChan
		switch event.EventType {
		case EventTypeDM:
			break
		}
	}
}

func serveAPI() {
}

func logErrors(errorChan <-chan error) {
	for {
		log.Println(<-errorChan)
	}
}

func main() {
	eventChan := make(chan Event)
	errorChan := make(chan error)
	go listenToTwitter("API_KEY", eventChan, errorChan)
	go processEvents(eventChan, errorChan)
	go serveAPI()
}
