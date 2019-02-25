package main

import (
	"log"

	"github.com/joho/godotenv"
	"gitlab.com/alpinefresh/tcrpartybot/contracts"
	"gitlab.com/alpinefresh/tcrpartybot/errors"
	"gitlab.com/alpinefresh/tcrpartybot/events"
	"gitlab.com/alpinefresh/tcrpartybot/models"
)

func main() {
	// Some pre-boot config
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	godotenv.Load()

	ethEvents := make(chan *events.ETHEvent)
	errChan := make(chan error)

	go events.StartScannerListener(ethEvents, errChan)
	go errors.LogErrors(errChan)
	go events.UpdateBalances(errChan)

	for event := range ethEvents {
		var decoded interface{}
		var err error

		switch event.EventType {
		////
		// Registry events
		//
		case events.ETHEventNewTCRApplication:
			decoded, err = contracts.DecodeApplicationEvent(event.Topics, event.Data)
		case events.ETHEventNewTCRChallenge:
			decoded, err = contracts.DecodeChallengeEvent(event.Topics, event.Data)
		case events.ETHEventTCRDeposit:
			decoded, err = contracts.DecodeDepositEvent(event.Topics, event.Data)
		case events.ETHEventTCRWithdrawal:
			decoded, err = contracts.DecodeWithdrawalEvent(event.Topics, event.Data)
		case events.ETHEventTCRApplicationWhitelisted:
			decoded, err = contracts.DecodeApplicationWhitelistedEvent(event.Topics, event.Data)
		case events.ETHEventTCRApplicationRemoved:
			decoded, err = contracts.DecodeApplicationRemovedEvent(event.Topics, event.Data)
		case events.ETHEventTCRListingRemoved:
			decoded, err = contracts.DecodeListingRemovedEvent(event.Topics, event.Data)
		case events.ETHEventTCRListingWithdrawn:
			decoded, err = contracts.DecodeListingWithdrawnEvent(event.Topics, event.Data)
		case events.ETHEventTCRTouchAndRemoved:
			decoded, err = contracts.DecodeTouchAndRemovedEvent(event.Topics, event.Data)
		case events.ETHEventTCRChallengeFailed:
			decoded, err = contracts.DecodeChallengeFailedEvent(event.Topics, event.Data)
		case events.ETHEventTCRChallengeSucceeded:
			decoded, err = contracts.DecodeChallengeSucceededEvent(event.Topics, event.Data)
		case events.ETHEventTCRRewardClaimed:
			decoded, err = contracts.DecodeRewardClaimedEvent(event.Topics, event.Data)
		case events.ETHEventTCRExitInitialized:
			decoded, err = contracts.DecodeExitInitializedEvent(event.Topics, event.Data)

		////
		// Token events
		//
		case events.ETHEventTokenTransfer:
			decoded, err = contracts.DecodeTransferEvent(event.Topics, event.Data)
		case events.ETHEventTokenApproval:
			decoded, err = contracts.DecodeApprovalEvent(event.Topics, event.Data)
		case events.ETHEventTokenMint:
			decoded, err = contracts.DecodeMintEvent(event.Topics, event.Data)
		case events.ETHEventTokenMintFinished:
			decoded, err = contracts.DecodeMintFinishedEvent(event.Topics, event.Data)

		////
		// PLCR events
		//
		case events.ETHEventPLCRVoteCommitted:
			decoded, err = contracts.DecodeVoteCommittedEvent(event.Topics, event.Data)
		case events.ETHEventPLCRPollCreated:
			decoded, err = contracts.DecodePollCreatedEvent(event.Topics, event.Data)
		case events.ETHEventPLCRVoteRevealed:
			decoded, err = contracts.DecodeVoteRevealedEvent(event.Topics, event.Data)
		case events.ETHEventPLCRVotingRightsGranted:
			decoded, err = contracts.DecodeVotingRightsGrantedEvent(event.Topics, event.Data)
		case events.ETHEventPLCRVotingRightsWithdrawn:
			decoded, err = contracts.DecodeVotingRightsWithdrawnEvent(event.Topics, event.Data)
		case events.ETHEventPLCRTokensRescued:
			decoded, err = contracts.DecodeTokensRescuedEvent(event.Topics, event.Data)

		case events.ETHEventNewMultisigWallet:
			continue

		default:
			log.Printf("Unrecognized event %s", event.EventType)
			continue
		}

		if err != nil {
			errChan <- errors.Wrap(err)
			continue
		}

		event := &models.ETHEvent{
			EventType:   event.EventType,
			BlockNumber: event.BlockNumber,
			CreatedAt:   event.CreatedAt,
			TxHash:      event.TxHash,
			TxIndex:     event.TxIndex,
			LogIndex:    event.LogIndex,
		}

		err = models.CreateETHEvent(event, decoded)
		if err != nil {
			errChan <- errors.Wrap(err)
			continue
		}
	}
}
