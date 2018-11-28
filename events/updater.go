package events

import (
	"log"
	"time"

	"gitlab.com/alpinefresh/tcrpartybot/contracts"
)

// ScheduleUpdates finds new/existing applications/challenge polls and
// schedules tasks to fire their on-chain updater functions once they have
// matured. Example: calling updateStatus on an application once the
// application period has passed.
func ScheduleUpdates(eventChan <-chan *ETHEvent, errChan chan<- error) {
	// First let's instantiate ourselves with existing challenges/applications
	// that need to be scheduled
	listings, err := contracts.GetUnwhitelistedListings()
	if err != nil {
		errChan <- err
		return
	}

	for _, listing := range listings {
		go scheduleApplication(listing, errChan)
	}

	// Listen for any incoming applications and queue them up
	for {
		event := <-eventChan
		var err error

		switch event.EventType {
		case ETHEventNewTCRApplication:
			ethEvent, err := contracts.DecodeApplicationEvent(event.Topics, event.Data)
			if err != nil {
				errChan <- err
				continue
			}

			listing, err := contracts.GetListingFromHash(ethEvent.ListingHash)
			if err != nil {
				errChan <- err
				continue
			}

			go scheduleApplication(listing, errChan)
			break
		}

		if err != nil {
			errChan <- err
		}
	}
}

func scheduleApplication(application *contracts.RegistryListing, errChan chan<- error) {
	expirationTime := time.Unix(application.ApplicationExpiry.Int64(), 0)
	log.Printf("[updater] Watching application 0x%x (exp: %s)", application.ListingHash, expirationTime.String())

	// If we haven't hit the expiration time yet let's sleep until we do
	if expirationTime.After(time.Now()) {
		time.Sleep(time.Until(expirationTime) + (5 * time.Second))
		scheduleApplication(application, errChan)
		return
	}

	tx, err := contracts.UpdateStatus(application.ListingHash)
	if err != nil {
		errChan <- err
		return
	}
	log.Printf("[updater] Application 0x%x's period has expired. Updating tx: %s", application.ListingHash, tx.Hash().Hex())
}
