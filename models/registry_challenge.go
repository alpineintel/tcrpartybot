package models

import (
	"database/sql"
	"time"
)

// RegistryChallenge represents a challenge on a RegistryListing
type RegistryChallenge struct {
	PollID              int64            `db:"poll_id" json:"poll_id"`
	ListingHash         string           `db:"listing_hash" json:"-"`
	ListingID           string           `db:"listing_id" json:"-"`
	ChallengerAccountID sql.NullInt64    `db:"challenger_account_id" json:"-"`
	ChallengerAddress   string           `db:"challenger_address" json:"-"`
	CreatedAt           *time.Time       `db:"created_at" json:"created_at"`
	CommitEndsAt        *time.Time       `db:"commit_ends_at" json:"commit_ends_at"`
	RevealEndsAt        *time.Time       `db:"reveal_ends_at" json:"reveal_ends_at"`
	SucceededAt         *time.Time       `db:"succeeded_at" json:"-"`
	FailedAt            *time.Time       `db:"failed_at" json:"-"`
	Listing             *RegistryListing `db:"registry_listing" json:"-"`
}

// FindRegistryChallengeByPollID finds a challenge given its poll ID (also
// known as a challenge ID)
func FindRegistryChallengeByPollID(pollID int64) (*RegistryChallenge, error) {
	db := GetDBSession()

	challenge := RegistryChallenge{}
	err := db.Get(&challenge, `
		SELECT
			registry_challenges.*,
			registry_listings.twitter_handle AS "registry_listing.twitter_handle"
		FROM registry_challenges
		LEFT JOIN registry_listings ON
			registry_listings.id = registry_challenges.listing_id
		WHERE poll_id=$1
	`, pollID)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return &challenge, nil
}

// FindActiveChallenges returns a list of challenges current in progress
func FindActiveChallenges() ([]*RegistryChallenge, error) {
	db := GetDBSession()

	challenges := []*RegistryChallenge{}
	err := db.Select(&challenges, `
		SELECT * FROM registry_challenges WHERE succeeded_at IS NULL AND failed_at IS NULL
	`)

	return challenges, err
}

// Create inserts the registry challenge into the database
func (challenge *RegistryChallenge) Create() error {
	db := GetDBSession()

	_, err := db.NamedExec(`
		INSERT INTO registry_challenges (
			poll_id,
			listing_hash,
			listing_id,
			created_at,
			commit_ends_at,
			reveal_ends_at,
			challenger_account_id,
			challenger_address,
			succeeded_at,
			failed_at
		) VALUES(
			:poll_id,
			:listing_hash,
			:listing_id,
			:created_at,
			:commit_ends_at,
			:reveal_ends_at,
			:challenger_account_id,
			:challenger_address,
			:succeeded_at,
			:failed_at
		)
	`, challenge)

	return err
}

// Save updates the registry challenge's row in the database
func (challenge *RegistryChallenge) Save() error {
	db := GetDBSession()

	_, err := db.NamedExec(`
		UPDATE registry_challenges SET
			listing_hash = :listing_hash,
			listing_id = :listing_id,
			created_at = :created_at,
			commit_ends_at = :commit_ends_at,
			reveal_ends_at = :reveal_ends_at,
			challenger_account_id = :challenger_account_id,
			challenger_address = :challenger_address,
			succeeded_at = :succeeded_at,
			failed_at = :failed_at
		WHERE poll_id=:poll_id
	`, challenge)

	return err
}
