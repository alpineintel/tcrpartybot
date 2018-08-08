package events

import (
	"time"
)

type EventType string

const (
	EventTypeMention       EventType = "EventTypeMention"
	EventTypeDM            EventType = "EventTypeDM"
	EventTypeVote          EventType = "EventTypeVote"
	EventTypePollCompleted EventType = "EventTypePollCompleted"
)

type Event struct {
	EventType    EventType // type of event
	Time         time.Time // timestamp
	SourceHandle string    // twitter handle sending
	Message      string    // whole message
}

func ProcessEvents(eventChan <-chan *Event, errorChan chan<- error) {
	for {
		event := <-eventChan
		switch event.EventType {
		case EventTypeMention:
			processMention(event, errorChan)
			break
		}
	}
}
