package events

import (
	"log"
	"time"
)

// EventType is an alias type for event constants' values
type twitterEventType string

const (
	TwitterEventTypeMention       twitterEventType = "EventTypeMention"
	TwitterEventTypeDM            twitterEventType = "EventTypeDM"
	TwitterEventTypeVote          twitterEventType = "EventTypeVote"
	TwitterEventTypePollCompleted twitterEventType = "EventTypePollCompleted"
	TwitterEventTypeFollow        twitterEventType = "EventTypeFollow"

	// ETHEventNewMultisigWallet is triggered when the multisig wallet factory instantiates a new wallet
	ETHEventNewMultisigWallet     = "ContractInstantiation"
	ETHEventNewTCRApplication     = "_Application"
	ETHEventNewMultisigSubmission = "Submission"
)

// TwitterEvent represents an incoming event from Twitter
type TwitterEvent struct {
	EventType    twitterEventType // type of event
	Time         time.Time        // timestamp
	SourceHandle string           // twitter handle sending
	SourceID     int64            // twitter ID of the handle
	Message      string           // whole message
}

// ETHEvent represents an incoming event from the blockchain
type ETHEvent struct {
	EventType string
	Data      []byte
}

// ProcessTwitterEvents listens for twitter events and fires of a corresponding handler
func ProcessTwitterEvents(eventChan <-chan *TwitterEvent, errorChan chan<- error) {
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

// ProcessETHEvents listens for blockchain events and fires a corresponding handler
func ProcessETHEvents(eventChan <-chan *ETHEvent, errChan chan<- error) {
	for {
		event := <-eventChan
		var err error

		log.Printf("New ETH event: %s", event.EventType)

		switch event.EventType {
		case ETHEventNewMultisigWallet:
			err = processMultisigWalletCreation(event)
			break
		}

		if err != nil {
			errChan <- err
		}
	}
}
