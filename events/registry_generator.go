package events

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"gitlab.com/alpinefresh/tcrpartybot/errors"
	"gitlab.com/alpinefresh/tcrpartybot/models"
)

// GenerateAndWatchRegistry incrementally syncs the registry_* tables with
// events from the eth_events table in order to form an up-to-date state of the
// registry
func GenerateAndWatchRegistry(errChan chan<- error) {
	for {
		time.Sleep(5 * time.Second)
		if err := generateRegistry(); err != nil {
			errChan <- err
		}
	}
}

func generateRegistry() error {
	// Find the last event ID we've seen
	latestEventID, err := models.GetKey(models.LatestRegistrySyncedEventKey)
	if err != nil {
		return err
	}
	intID, err := strconv.ParseInt(latestEventID, 10, 64)
	if err != nil {
		intID = 0
	}

	// Get all of the events we're missing
	events, moreAvailable, err := models.FindETHEventsSinceID(intID)
	if err != nil {
		return err
	}
	if len(events) == 0 {
		return nil
	}

	log.Printf("Updating registry since %d", intID)
	// @TODO Start a transaction
	for _, event := range events {
		switch event.EventType {
		case ETHEventNewTCRApplication:
			decoded, err := unmarshalApplicationEvent(event.Data)
			if err != nil {
				return errors.Wrap(err)
			}

			applicationEndedAt := time.Unix(decoded.AppEndDate, 0)
			listing := &models.RegistryListing{
				ListingHash:          fmt.Sprintf("%x", decoded.ListingHash),
				TwitterHandle:        decoded.Data,
				Deposit:              decoded.Deposit.String(),
				ApplicationCreatedAt: event.CreatedAt,
				ApplicationEndedAt:   &applicationEndedAt,
			}
			if err := listing.Create(); err != nil {
				return errors.Wrap(err)
			}
		case ETHEventTCRApplicationWhitelisted:
			decoded, err := unmarshalApplicationWhitelistedEvent(event.Data)
			if err != nil {
				return errors.Wrap(err)
			}

			listing, err := models.FindLatestRegistryListingByHash(decoded.ListingHash)
			if err != nil {
				return errors.Wrap(err)
			} else if listing == nil {
				return errors.Errorf("could not find listing %x in db to whitelist", decoded.ListingHash)
			}

			listing.WhitelistedAt = event.CreatedAt
			if err := listing.Save(); err != nil {
				return errors.Wrap(err)
			}
		case ETHEventTCRListingRemoved:
			decoded, err := unmarshalListingRemovedEvent(event.Data)
			if err != nil {
				return errors.Wrap(err)
			}

			listing, err := models.FindLatestRegistryListingByHash(decoded.ListingHash)
			if err != nil {
				return errors.Wrap(err)
			} else if listing == nil {
				return errors.Errorf("could not find listing %x in db to remove", decoded.ListingHash)
			}

			listing.RemovedAt = event.CreatedAt
			if err := listing.Save(); err != nil {
				return errors.Wrap(err)
			}
		case ETHEventTCRApplicationRemoved:
			decoded, err := unmarshalApplicationRemovedEvent(event.Data)
			if err != nil {
				return errors.Wrap(err)
			}

			listing, err := models.FindLatestRegistryListingByHash(decoded.ListingHash)
			if err != nil {
				return errors.Wrap(err)
			} else if listing == nil {
				return errors.Errorf("could not find listing %x in db to remove", decoded.ListingHash)
			}

			listing.RemovedAt = event.CreatedAt
			if err := listing.Save(); err != nil {
				return errors.Wrap(err)
			}

		case ETHEventNewTCRChallenge:
			decoded, err := unmarshalChallengeEvent(event.Data)
			if err != nil {
				return errors.Wrap(err)
			}

			// Find the associated listing
			listing, err := models.FindLatestRegistryListingByHash(decoded.ListingHash)
			if err != nil {
				return errors.Wrap(err)
			} else if listing == nil {
				return errors.Errorf("could not find listing %x in db to remove", decoded.ListingHash)
			}

			// Find the challenger account
			accountID := sql.NullInt64{}
			challenger, err := models.FindAccountByMultisigAddress(decoded.Challenger)
			if err != nil {
				return errors.Wrap(err)
			} else if challenger != nil {
				accountID.Int64 = challenger.ID
				accountID.Valid = true
			}

			commitDate := time.Unix(decoded.CommitEndDate, 0)
			revealDate := time.Unix(decoded.CommitEndDate, 0)

			challenge := &models.RegistryChallenge{
				PollID:              decoded.ChallengeID,
				ListingHash:         fmt.Sprintf("%x", decoded.ListingHash),
				ListingID:           listing.ID,
				ChallengerAccountID: accountID,
				ChallengerAddress:   decoded.Challenger,
				CreatedAt:           event.CreatedAt,
				CommitEndsAt:        &commitDate,
				RevealEndsAt:        &revealDate,
			}

			if err := challenge.Create(); err != nil {
				return errors.Wrap(err)
			}

		case ETHEventTCRChallengeSucceeded:
			decoded, err := unmarshalChallengeEvent(event.Data)
			if err != nil {
				return errors.Wrap(err)
			}

			challenge, err := models.FindRegistryChallengeByPollID(decoded.ChallengeID)
			if err != nil {
				return errors.Wrap(err)
			} else if challenge == nil {
				return errors.Errorf("could not find challenge %d in db to mark succeeded", decoded.ChallengeID)
			}

			challenge.SucceededAt = event.CreatedAt
			if err := challenge.Save(); err != nil {
				return errors.Wrap(err)
			}

		case ETHEventTCRChallengeFailed:
			decoded, err := unmarshalChallengeEvent(event.Data)
			if err != nil {
				return errors.Wrap(err)
			}

			challenge, err := models.FindRegistryChallengeByPollID(decoded.ChallengeID)
			if err != nil {
				return errors.Wrap(err)
			} else if challenge == nil {
				return errors.Errorf("could not find challenge %d in db to mark failed", decoded.ChallengeID)
			}

			challenge.FailedAt = event.CreatedAt
			if err := challenge.Save(); err != nil {
				return errors.Wrap(err)
			}
		}
	}

	// Set the last seen event ID and recurse if there's more available
	err = models.SetKey(models.LatestRegistrySyncedEventKey, fmt.Sprintf("%d", events[len(events)-1].ID))
	if moreAvailable {
		return generateRegistry()
	}

	return nil
}
