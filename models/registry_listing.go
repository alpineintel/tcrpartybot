package models

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"time"
)

// RegistryListing represents an application for a listing in the registry
// (which may later be whitelisted)
type RegistryListing struct {
	ID                   string     `db:"id"`
	ListingHash          string     `db:"listing_hash"`
	TwitterHandle        string     `db:"twitter_handle"`
	Deposit              string     `db:"deposit"`
	ApplicationCreatedAt *time.Time `db:"application_created_at"`
	ApplicationEndedAt   *time.Time `db:"application_ended_at"`
	WhitelistedAt        *time.Time `db:"whitelisted_at"`
	RemovedAt            *time.Time `db:"removed_at"`
}

func FindLatestRegistryListingByHash(listingHash [32]byte) (*RegistryListing, error) {
	db := GetDBSession()

	listing := RegistryListing{}
	err := db.Get(&listing, `
		SELECT
			*
		FROM registry_listings
		WHERE listing_hash=$1
		ORDER BY application_created_at DESC
		LIMIT 1
	`, fmt.Sprintf("%x", listingHash))

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return &listing, nil
}

// Create inserts the registry challenge into the database
func (listing *RegistryListing) Create() error {
	db := GetDBSession()

	if listing.ID == "" {
		listing.generateID()
	}

	_, err := db.NamedExec(`
		INSERT INTO registry_listings (
			id,
			listing_hash,
			twitter_handle,
			deposit,
			application_created_at,
			application_ended_at,
			whitelisted_at,
			removed_at
		) VALUES(
			:id,
			:listing_hash,
			:twitter_handle,
			:deposit,
			:application_created_at,
			:application_ended_at,
			:whitelisted_at,
			:removed_at
		)
	`, listing)

	return err
}

func (listing *RegistryListing) Save() error {
	db := GetDBSession()
	_, err := db.NamedExec(`
		UPDATE registry_listings SET
			listing_hash = :listing_hash,
			twitter_handle = :twitter_handle,
			deposit = :deposit,
			application_created_at = :application_created_at,
			application_ended_at = :application_ended_at,
			whitelisted_at = :whitelisted_at,
			removed_at = :removed_at
		WHERE id=:id
	`, listing)

	return err
}

func (listing *RegistryListing) generateID() {
	hash := sha256.New()
	hash.Write([]byte(listing.ListingHash))
	hash.Write([]byte(fmt.Sprintf("%d", listing.ApplicationCreatedAt.Unix())))

	listing.ID = fmt.Sprintf("%x", hash.Sum(nil))
}
