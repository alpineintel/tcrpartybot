package events

import (
	"time"
)

type EventType string

const (
	TwitterEventTypeMention       EventType = "EventTypeMention"
	TwitterEventTypeDM            EventType = "EventTypeDM"
	TwitterEventTypeVote          EventType = "EventTypeVote"
	TwitterEventTypePollCompleted EventType = "EventTypePollCompleted"
	TwitterEventTypeFollow        EventType = "EventTypeFollow"
)

type TwitterEvent struct {
	EventType    EventType // type of event
	Time         time.Time // timestamp
	SourceHandle string    // twitter handle sending
	SourceID     int64     // twitter ID of the handle
	Message      string    // whole message
}

func ProcessEvents(eventChan <-chan *TwitterEvent, errorChan chan<- error) {
	for {
		event := <-eventChan
		switch event.EventType {
		case TwitterEventTypeMention:
			processMention(event, errorChan)
			break

		case TwitterEventTypeFollow:
			processFollow(event, errorChan)
			break

		case TwitterEventTypeDM:
			processDM(event, errorChan)
			break
		}
	}
}
