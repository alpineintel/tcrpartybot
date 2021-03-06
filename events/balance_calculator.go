package events

import (
	"github.com/jmoiron/sqlx/types"
	"log"
	"math/big"
	"time"

	"gitlab.com/alpinefresh/tcrpartybot/errors"
	"gitlab.com/alpinefresh/tcrpartybot/models"
)

type transferEvent struct {
	From  string   `json:"from"`
	To    string   `json:"to"`
	Value *big.Int `json:"value"`
}

type votingRightsGrantedEvent struct {
	Voter     string   `json:"Voter"`
	NumTokens *big.Int `json:"NumTokens"`
}

type votingRightsWithdrawnEvent struct {
	Voter     string   `json:"Voter"`
	NumTokens *big.Int `json:"NumTokens"`
}

type applicationEvent struct {
	Applicant   string   `json:"Applicant"`
	Data        string   `json:"Data"`
	ListingHash [32]byte `json:"ListingHash"`
	Deposit     *big.Int `json:"Deposit"`
	AppEndDate  int64    `json:"AppEndDate"`
}

type applicationWhitelistedEvent struct {
	ListingHash [32]byte `json:"ListingHash"`
}

type listingRemovedEvent struct {
	ListingHash [32]byte `json:"ListingHash"`
}

type applicationRemovedEvent struct {
	ListingHash [32]byte `json:"ListingHash"`
}

type challengeEvent struct {
	ListingHash   [32]byte `json:"ListingHash"`
	ChallengeID   int64    `json:"ChallengeID"`
	Challenger    string   `json:"Challenger"`
	CommitEndDate int64    `json:"CommitEndDate"`
	RevealEndDate int64    `json:"RevealEndDate"`
}

type challengeSucceededEvent struct {
	ChallengeID int64    `json:"ChallengeID"`
	ListingHash [32]byte `json:"ListingHash"`
}

type challengeFailedEvent struct {
	ChallengeID int64    `json:"ChallengeID"`
	ListingHash [32]byte `json:"ListingHash"`
}

func negate(v *big.Int) *big.Int {
	negOne := big.NewInt(-1)
	return v.Mul(v, negOne)
}

// UpdateBalances polls for new eth events and updates the balances table
// whenever relevant eth events are found.
func UpdateBalances(errChan chan<- error) {
	var moreAvailable = false
	latestEventID := int64(0)

	for {
		// moreAvailable will be true if we are performing an initial sync and
		// are looping through a large quantity of events. We only want to
		// delay if we're performing an incremental sync.
		if !moreAvailable {
			time.Sleep(3 * time.Second)
		} else {
			log.Printf("More available in window, syncing since %d...", latestEventID)
		}

		// Find the last block we've synced
		balance, err := models.FindLatestBalance()
		if err != nil {
			errChan <- errors.Wrap(err)
			continue
		}

		// If we don't have a balance here that means we need to generate
		// everything from the beginning of our event log
		if balance != nil {
			latestEventID = balance.ETHEventID
		}

		// Fetch all events since that block number
		events, available, err := models.FindETHEventsSinceID(latestEventID)
		moreAvailable = available
		if err != nil {
			errChan <- errors.Wrap(err)
			continue
		}

		// Create balance events
		for _, event := range events {
			latestEventID = event.ID

			switch event.EventType {
			case ETHEventTokenTransfer:
				// Remove from someone's wallet balance, add to someone else's
				data, err := unmarshalTransferEvent(event.Data)
				if err != nil {
					errChan <- errors.Wrap(err)
					continue
				}

				// Subtract from the from account
				fromAccount, err := models.FindAccountByMultisigAddress(data.From)
				if err != nil {
					errChan <- errors.Wrap(err)
					continue
				} else if fromAccount != nil {
					if _, err := fromAccount.AddToWalletBalance(event.ID, negate(data.Value)); err != nil {
						errChan <- errors.Wrap(err)
						continue
					}
				}

				// Add to the toAccount
				toAccount, err := models.FindAccountByMultisigAddress(data.To)
				if err != nil {
					errChan <- errors.Wrap(err)
					continue
				} else if toAccount != nil {
					if _, err := toAccount.AddToWalletBalance(event.ID, data.Value); err != nil {
						errChan <- errors.Wrap(err)
						continue
					}
				}

			case ETHEventPLCRVotingRightsGranted:
				data, err := unmarshalVotingRightsGranted(event.Data)
				if err != nil {
					errChan <- errors.Wrap(err)
					continue
				}

				// Add to someone's PLCR balance
				voterAccount, err := models.FindAccountByMultisigAddress(data.Voter)
				if err != nil {
					errChan <- errors.Wrap(err)
					continue
				} else if voterAccount != nil {
					if _, err := voterAccount.AddToPLCRBalance(event.ID, data.NumTokens); err != nil {
						errChan <- errors.Wrap(err)
						continue
					}
				}

			case ETHEventPLCRVotingRightsWithdrawn:
				// Remove from someone's PLCR balance
				data, err := unmarshalVotingRightsWithdrawn(event.Data)
				if err != nil {
					errChan <- errors.Wrap(err)
					continue
				}
				// Remove from someone's PLCR balance
				voterAccount, err := models.FindAccountByMultisigAddress(data.Voter)
				if err != nil {
					errChan <- errors.Wrap(err)
					continue
				} else if voterAccount != nil {
					if _, err := voterAccount.AddToPLCRBalance(event.ID, negate(data.NumTokens)); err != nil {
						errChan <- errors.Wrap(err)
						continue
					}
				}
			default:
				continue
			}
		}
	}
}

func unmarshalTransferEvent(json types.JSONText) (*transferEvent, error) {
	data := &transferEvent{}
	if err := json.Unmarshal(data); err != nil {
		return nil, err
	}
	return data, nil
}

func unmarshalApplicationEvent(json types.JSONText) (*applicationEvent, error) {
	data := &applicationEvent{}
	if err := json.Unmarshal(data); err != nil {
		return nil, err
	}
	return data, nil
}

func unmarshalApplicationWhitelistedEvent(json types.JSONText) (*applicationWhitelistedEvent, error) {
	data := &applicationWhitelistedEvent{}
	if err := json.Unmarshal(data); err != nil {
		return nil, err
	}
	return data, nil
}

func unmarshalListingRemovedEvent(json types.JSONText) (*listingRemovedEvent, error) {
	data := &listingRemovedEvent{}
	if err := json.Unmarshal(data); err != nil {
		return nil, err
	}
	return data, nil
}

func unmarshalApplicationRemovedEvent(json types.JSONText) (*applicationRemovedEvent, error) {
	data := &applicationRemovedEvent{}
	if err := json.Unmarshal(data); err != nil {
		return nil, err
	}
	return data, nil
}

func unmarshalChallengeEvent(json types.JSONText) (*challengeEvent, error) {
	data := &challengeEvent{}
	if err := json.Unmarshal(data); err != nil {
		return nil, err
	}
	return data, nil
}

func unmarshalChallengeSucceededEvent(json types.JSONText) (*challengeSucceededEvent, error) {
	data := &challengeSucceededEvent{}
	if err := json.Unmarshal(data); err != nil {
		return nil, err
	}
	return data, nil
}

func unmarshalChallengeFailedEvent(json types.JSONText) (*challengeFailedEvent, error) {
	data := &challengeFailedEvent{}
	if err := json.Unmarshal(data); err != nil {
		return nil, err
	}
	return data, nil
}

func unmarshalVotingRightsGranted(json types.JSONText) (*votingRightsGrantedEvent, error) {
	data := &votingRightsGrantedEvent{}
	if err := json.Unmarshal(data); err != nil {
		return nil, err
	}
	return data, nil
}

func unmarshalVotingRightsWithdrawn(json types.JSONText) (*votingRightsWithdrawnEvent, error) {
	data := &votingRightsWithdrawnEvent{}
	if err := json.Unmarshal(data); err != nil {
		return nil, err
	}
	return data, nil
}
