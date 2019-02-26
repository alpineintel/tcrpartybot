package twitter

import (
	"fmt"
	set "github.com/deckarep/golang-set"
	"log"
	"strings"

	"gitlab.com/alpinefresh/tcrpartybot/errors"
	"gitlab.com/alpinefresh/tcrpartybot/models"
)

// SyncList will periodically update the Twitter list to ensure it is in sync
// with our TCR
func SyncList() error {
	log.Println("Syncing Twitter list...")
	toLog := ""

	// Fetch the current status of the TCR and Twitter list
	members, err := GetList()
	if err != nil {
		return errors.Wrap(err)
	}

	whitelistedListings, err := models.FindWhitelistedRegistryListings()
	if err != nil {
		return errors.Wrap(err)
	}

	// Find any differences
	listHandles := set.NewSet()
	toLog += "\nList sync summary:\n\tCurrent members: "
	for _, member := range members.Users {
		handle := strings.TrimSpace(strings.ToLower(member.ScreenName))
		listHandles.Add(handle)
		toLog += fmt.Sprintf("%s,", handle)
	}

	tcrHandles := set.NewSet()
	toLog += "\n\tTCR Members: "
	for _, listing := range whitelistedListings {
		handle := strings.TrimSpace(strings.ToLower(listing.TwitterHandle))
		tcrHandles.Add(handle)
		toLog += fmt.Sprintf("%s,", handle)
	}

	toLog += "\n\tTo remove: "
	toRemove := []string{}
	for handle := range listHandles.Difference(tcrHandles).Iterator().C {
		toRemove = append(toRemove, handle.(string))
		toLog += fmt.Sprintf("%s,", handle.(string))
	}

	toLog += "\n\tTo add: "
	toAdd := []string{}
	for handle := range tcrHandles.Difference(listHandles).Iterator().C {
		toAdd = append(toAdd, handle.(string))
		toLog += fmt.Sprintf("%s,", handle.(string))
	}

	if len(toAdd) > 0 {
		err = AddHandlesToList(toAdd)
		if err != nil {
			return errors.Wrap(err)
		}
	}

	if len(toRemove) > 0 {
		err = RemoveHandlesToList(toRemove)
		if err != nil {
			return errors.Wrap(err)
		}
	}

	log.Println(toLog)

	return nil
}
