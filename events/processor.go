package events

import (
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
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
	ETHEventNewMultisigWallet         = "ContractInstantiation"
	ETHEventNewMultisigSubmission     = "Submission"
	ETHEventNewTCRApplication         = "_Application"
	ETHEventNewTCRChallenge           = "_Challenge"
	ETHEventTCRDeposit                = "_Deposit"
	ETHEventTCRWithdrawal             = "_Withdrawal"
	ETHEventTCRApplicationWhitelisted = "_ApplicationWhitelisted"
	ETHEventTCRApplicationRemoved     = "_ApplicationRemoved"
	ETHEventTCRListingRemoved         = "_ListingRemoved"
	ETHEventTCRListingWithdrawn       = "_ListingWithdrawn"
	ETHEventTCRTouchAndRemoved        = "_TouchAndRemoved"
	ETHEventTCRChallengeFailed        = "_ChallengeFailed"
	ETHEventTCRChallengeSucceeded     = "_ChallengeSucceeded"
	ETHEventTCRRewardClaimed          = "_RewardClaimed"
	ETHEventTCRExitInitialized        = "_ExitInitialized"

	ETHEventTokenTransfer     = "Transfer"
	ETHEventTokenApproval     = "Approval"
	ETHEventTokenMint         = "Mint"
	ETHEventTokenMintFinished = "MintFinished"

	ETHEventPLCRVoteCommitted         = "_VoteCommitted"
	ETHEventPLCRPollCreated           = "_PollCreated"
	ETHEventPLCRVoteRevealed          = "_VoteRevealed"
	ETHEventPLCRVotingRightsGranted   = "_VotingRightsGranted"
	ETHEventPLCRVotingRightsWithdrawn = "_VotingRightsWithdrawn"
	ETHEventPLCRTokensRescued         = "_TokensRescued"
)

// TwitterEvent represents an incoming event from Twitter
type TwitterEvent struct {
	EventType    twitterEventType // type of event
	Time         time.Time        // timestamp
	ObjectID     string           // ID of incoming event's object (ie a tweet ID for a mention)
	SourceHandle string           // twitter handle sending
	SourceID     int64            // twitter ID of the handle
	Message      string           // whole message
}

// ETHEvent represents an incoming event from the blockchain
type ETHEvent struct {
	EventType   string
	CreatedAt   *time.Time
	BlockNumber uint64
	TxIndex     uint
	LogIndex    uint
	TxHash      string
	Data        []byte
	Topics      []common.Hash
}

// ProcessTwitterEvents listens for twitter events and fires of a corresponding handler
func ProcessTwitterEvents(eventChan <-chan *TwitterEvent, errorChan chan<- error) {
	for {
		event := <-eventChan
		switch event.EventType {
		case TwitterEventTypeMention:
			go processMention(event, errorChan)
			break

		case TwitterEventTypeFollow:
			go processFollow(event, errorChan)
			break

		case TwitterEventTypeDM:
			go processDM(event, errorChan)
			break
		}
	}
}

// ProcessETHEvents listens for blockchain events and fires a corresponding handler
func ProcessETHEvents(eventChan <-chan *ETHEvent, errChan chan<- error) {
	var err error

	for {
		event := <-eventChan
		go scheduleUpdateForEvent(event, errChan)

		log.Printf("Found event %s", event.EventType)
		switch event.EventType {
		case ETHEventNewMultisigWallet:
			err = processMultisigWalletCreation(event)
			break
		case ETHEventNewTCRApplication:
			err = processNewApplication(event)
			break
		case ETHEventTCRApplicationWhitelisted:
			err = processApplicationWhitelisted(event)
			break
		case ETHEventTCRApplicationRemoved:
			err = processApplicationRemoved(event)
			break
		case ETHEventTCRChallengeSucceeded:
			err = processChallengeSucceeded(event)
			break
		case ETHEventTCRChallengeFailed:
			err = processChallengeFailed(event)
			break
		case ETHEventNewTCRChallenge:
			err = processNewChallenge(event)
			break
		case ETHEventTCRWithdrawal:
			err = processWithdrawal(event)
			break
		case ETHEventTCRRewardClaimed:
			err = processRewardClaimed(event)
			break
		}

		if err != nil {
			errChan <- err
		}
	}
}
