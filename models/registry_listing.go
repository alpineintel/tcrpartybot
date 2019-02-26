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
	ID                   string     `db:"id" json:"-"`
	ListingHash          string     `db:"listing_hash" json:"listing_hash"`
	TwitterHandle        string     `db:"twitter_handle" json:"twitter_handle"`
	Deposit              string     `db:"deposit" json:"deposit"`
	ApplicationCreatedAt *time.Time `db:"application_created_at" json:"application_created_at"`
	ApplicationEndedAt   *time.Time `db:"application_ended_at" json:"application_ended_at"`
	WhitelistedAt        *time.Time `db:"whitelisted_at" json:"whitelisted_at"`
	RemovedAt            *time.Time `db:"removed_at" json:"-"`
}

type RegistryListingWithChallenge struct {
	RegistryListing
	RegistryChallenge `db:"registry_challenge" json:"challenge"`
}

// FindLatestRegistryListingByHash finds a registry listing by its listing hash
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

// FindWhitelistedRegistryListings returns a slice of all currently whitelisted
// listings
func FindWhitelistedRegistryListings() ([]*RegistryListing, error) {
	db := GetDBSession()

	listings := []*RegistryListing{}
	err := db.Select(&listings, "SELECT * FROM registry_listings WHERE whitelisted_at IS NOT NULL and removed_at IS NULL")

	return listings, err
}

// FindUnchallengedWhitelistedListings returns a list of registry listings that
// have been whitelisted but do not have an active challenge against them
func FindUnchallengedWhitelistedListings() ([]*RegistryListing, error) {
	db := GetDBSession()

	listings := []*RegistryListing{}
	err := db.Select(&listings, `
		SELECT
			registry_listings.*
		FROM
			registry_listings
		LEFT OUTER JOIN registry_challenges ON
			registry_challenges.listing_id = registry_listings.id
		WHERE
			registry_listings.whitelisted_at IS NOT NULL AND
			registry_listings.removed_at IS NULL AND
			registry_challenges.succeeded_at IS NULL AND
			registry_challenges.failed_at IS NULL
	`)

	return listings, err
}

// FindChallengedRegistryListings returns a list of registry listings (in the
// application and whitelisted stages) that are under active challenge
func FindChallengedRegistryListings() ([]*RegistryListingWithChallenge, error) {
	db := GetDBSession()

	listings := []*RegistryListingWithChallenge{}
	err := db.Select(&listings, `
		SELECT
			registry_listings.*,
			registry_challenges.poll_id AS "registry_challenge.poll_id",
			registry_challenges.created_at AS "registry_challenge.created_at",
			registry_challenges.commit_ends_at AS "registry_challenge.commit_ends_at",
			registry_challenges.reveal_ends_at AS "registry_challenge.reveal_ends_at"
		FROM
			registry_listings
		INNER JOIN registry_challenges ON
			registry_challenges.listing_id = registry_listings.id
		WHERE
			registry_listings.removed_at IS NULL AND
			registry_challenges.succeeded_at IS NULL AND
			registry_challenges.failed_at IS NULL
	`)

	return listings, err
}

// FindUnchallengedApplications returns a list of listings that are currently
// applying to be on the TCR and do not have an active challenge against them.
func FindUnchallengedApplications() ([]*RegistryListing, error) {
	db := GetDBSession()

	listings := []*RegistryListing{}
	err := db.Select(&listings, `
		SELECT
			registry_listings.*
		FROM
			registry_listings
		LEFT OUTER JOIN registry_challenges ON
			registry_challenges.listing_id = registry_listings.id
		WHERE
			registry_listings.whitelisted_at IS NULL AND
			registry_listings.removed_at IS NULL AND
			registry_challenges.poll_id IS NULL
	`)

	return listings, err
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

// Save updates the registry listing in the database based on its current state
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
