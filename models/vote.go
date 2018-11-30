package models

import (
	"time"
)

// Vote represents a vote submitted to the PLCR contract. It is mostly used for
// stats and a temporary store of salts which are used in the reveal phase
type Vote struct {
	PollID    int64      `db:"poll_id"`
	AccountID int64      `db:"account_id"`
	Salt      int64      `db:"salt"`
	Vote      bool       `db:"vote"`
	Revealed  *time.Time `db:"revealed_at"`
	CreatedAt *time.Time `db:"created_at"`
}

// CreateVote instantiates a new Vote struct and persists it to the database
func CreateVote(account *Account, pollID int64, salt int64, voteVal bool) (*Vote, error) {
	db := GetDBSession()

	vote := &Vote{
		PollID:    pollID,
		AccountID: account.ID,
		Salt:      salt,
		Vote:      voteVal,
	}

	_, err := db.Exec(`
		INSERT INTO votes (
			poll_id,
			account_id,
			salt,
			vote
		) VALUES($1, $2, $3, $4)
		RETURNING id
	`, vote.PollID, vote.AccountID, vote.Salt, vote.Vote)

	if err != nil {
		return nil, err
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

// MarkRevealed sets the revealed_at column of the vote to the current timestamp
func (v *Vote) MarkRevealed() error {
	db := GetDBSession()

	now := time.Now()
	_, err := db.Exec(`
		UPDATE votes
			SET revealed_at = $1
		WHERE account_id=$2 AND poll_id = $3
	`, &now, v.AccountID, v.PollID)

	return err
}
