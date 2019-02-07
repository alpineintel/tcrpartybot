package models

import (
	"database/sql"
	"time"
)

// Vote represents a vote submitted to the PLCR contract. It is mostly used for
// stats and a temporary store of salts which are used in the reveal phase
type Vote struct {
	PollID          int64      `db:"poll_id"`
	AccountID       int64      `db:"account_id"`
	Salt            int64      `db:"salt"`
	Vote            bool       `db:"vote"`
	Weight          int64      `db:"weight"`
	Revealed        *time.Time `db:"revealed_at"`
	RewardClaimedAt *time.Time `db:"reward_claimed_at"`
	CreatedAt       *time.Time `db:"created_at"`
}

// CreateVote instantiates a new Vote struct and persists it to the database
func CreateVote(account *Account, pollID int64, salt int64, voteVal bool, weight int64) (*Vote, error) {
	db := GetDBSession()

	vote := &Vote{
		PollID:    pollID,
		AccountID: account.ID,
		Salt:      salt,
		Vote:      voteVal,
		Weight:    weight,
	}

	_, err := db.Exec(`
		INSERT INTO votes (
			poll_id,
			account_id,
			salt,
			vote,
			weight
		) VALUES($1, $2, $3, $4, $5)
	`, vote.PollID, vote.AccountID, vote.Salt, vote.Vote, vote.Weight)

	if err != nil {
		return nil, err
	}

	return vote, nil
}

// FindVote returns a vote for a given poll or account ID
func FindVote(pollID, accountID int64) (*Vote, error) {
	db := GetDBSession()

	vote := &Vote{}
	err := db.Get(
		vote,
		"SELECT * FROM votes WHERE poll_id=$1 AND account_id=$2",
		pollID,
		accountID,
	)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return vote, nil
}

// FindUnrevealedVotesFromPoll returns all unrevealed votes associated with a
// given poll ID
func FindUnrevealedVotesFromPoll(pollID int64) ([]*Vote, error) {
	db := GetDBSession()

	votes := []*Vote{}
	err := db.Select(&votes, "SELECT * FROM votes WHERE poll_id=$1 AND revealed_at IS NULL", pollID)

	return votes, err
}

// FindUnrewardedVotesFromPoll returns all votes which have yet to reveal their
// rewards for a given poll and vote direction
func FindUnrewardedVotesFromPoll(pollID int64, voteDirection bool) ([]*Vote, error) {
	db := GetDBSession()

	votes := []*Vote{}
	err := db.Select(&votes, `
		SELECT * FROM votes
		WHERE
			poll_id=$1 AND
			vote=$2 AND
			reward_claimed_at IS NULL
	`, pollID, voteDirection)

	return votes, err
}

// MarkRevealed sets the revealed_at column of the vote to the current timestamp
func (v *Vote) MarkRevealed() error {
	db := GetDBSession()

	now := time.Now().UTC()
	_, err := db.Exec(`
		UPDATE votes
			SET revealed_at = $1
		WHERE account_id=$2 AND poll_id = $3
	`, &now, v.AccountID, v.PollID)

	return err
}

// MarkRewardClaimed sets the reward_claimed_at column of the vote to the
// current timestamp
func (v *Vote) MarkRewardClaimed() error {
	db := GetDBSession()

	now := time.Now().UTC()
	_, err := db.Exec(`
		UPDATE votes
			SET reward_claimed_at = $1
		WHERE account_id=$2 AND poll_id = $3
	`, &now, v.AccountID, v.PollID)

	return err
}
