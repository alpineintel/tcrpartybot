package twitter

import (
	set "github.com/deckarep/golang-set"
	"log"
	"strings"

	"gitlab.com/alpinefresh/tcrpartybot/contracts"
	"gitlab.com/alpinefresh/tcrpartybot/errors"
)

// SyncList will periodically update the Twitter list to ensure it is in sync
// with our TCR
func SyncList() error {
	log.Println("Syncing Twitter list...")
	// Fetch the current status of the TCR and Twitter list
	members, err := GetList()
	if err != nil {
		return errors.Wrap(err)
	}

	whitelistedHandles, err := contracts.GetWhitelistedHandles()
	if err != nil {
		return errors.Wrap(err)
	}

	// Find any differences
	listHandles := set.NewSet()
	for _, member := range members.Users {
		listHandles.Add(member.ScreenName)
	}

	tcrHandles := set.NewSet()
	for _, handle := range whitelistedHandles {
		tcrHandles.Add(handle)
	}

	toRemove := []string{}
	for handle := range listHandles.Difference(tcrHandles).Iterator().C {
		toRemove = append(toRemove, handle.(string))
	}

	toAdd := []string{}
	for handle := range tcrHandles.Difference(listHandles).Iterator().C {
		toAdd = append(toAdd, handle.(string))
	}

	if len(toAdd) > 0 {
		log.Printf("\tAdding: %s", strings.Join(toAdd, ","))
		err = AddHandlesToList(toAdd)
		if err != nil {
			return errors.Wrap(err)
		}
	}

	if len(toRemove) > 0 {
		log.Printf("\tRemoving: %s", strings.Join(toRemove, ","))
		err = RemoveHandlesToList(toRemove)
		if err != nil {
			return errors.Wrap(err)
		}
	}

	if len(toAdd)+len(toRemove) == 0 {
		log.Println("\tNothing to do")
	}

	return nil
}
