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
	EventTypeFollow        EventType = "EventTypeFollow"
)

type Event struct {
	EventType    EventType // type of event
	Time         time.Time // timestamp
	SourceHandle string    // twitter handle sending
	SourceID     int64     // twitter ID of the handle
	Message      string    // whole message
}

func ProcessEvents(eventChan <-chan *Event, errorChan chan<- error) {
	for {
		event := <-eventChan
		switch event.EventType {
		case EventTypeMention:
			processMention(event, errorChan)
			break

		case EventTypeFollow:
			processFollow(event, errorChan)
			break

		case EventTypeDM:
			processDM(event, errorChan)
			break
		}
	}
}
