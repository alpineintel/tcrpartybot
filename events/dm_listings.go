package events

import (
	"errors"
	"fmt"
	"math/big"

	"gitlab.com/alpinefresh/tcrpartybot/contracts"
	"gitlab.com/alpinefresh/tcrpartybot/models"
)

const (
	challengeArgErrorMsg          = "Whoops, looks like you forgot something. Try again with something like 'challenge [twitter handle]'. Eg: 'challenge weratedogs'"
	challengeNotFoundMsg          = "Looks like nobody has tried creating a listing for this twitter handle yet."
	challengeAlreadyExistsMsg     = "Looks like somebody has already begun a challenge for this twitter handle. You can support this challenge by voting kick (respond with 'vote %s kick')."
	challengeInsufficientFundsMsg = "Drat, looks like you don't have enough TCRP to start a challenge. You'll need 500 available in your wallet."
	challengeSubmissionErrorMsg   = "There was an error trying to submit your challenge. The admins have been notified!"
	challengeSubmissionSuccessMsg = "We've submitted your challenge to the registry (tx: %s). Keep an eye on @TCRPartyVIP for updates."

	nominateArgErrorMsg          = "Whoops, looks like you forgot something. Try again with something like 'nominate [twitter handle]'. Eg: 'apply weratedogs'"
	nominateAlreadyAppliedMsg    = "Looks like that Twitter handle has already been submitted to the TCR. A twitter handle can only appear on the TCR once, so you'll need to wait for a successful challenge (or a delisting) in order to re-nominate them."
	nominateInsufficientFundsMsg = "Drat! Looks like you don't have enough TCRP to start a nomination. You'll need 500 available in your wallet."
	nominateSubmissionErrorMsg   = "There was an error trying to submit your nomination: %s. Try tweeting at @stevenleeg for help."
	nominateSuccessMsg           = "All done! Your nomination was submitted in the following transaction: %s.\n\nKeep an eye on @TCRPartyVIP for an announcement."
	nominateBeginMsg             = "Got it! I've just begun submitting your nomination to the registry and will send a message once everything is confirmed."
)

func handleNomination(account *models.Account, argv []string, sendDM func(string)) error {
	if len(argv) != 2 {
		sendDM(nominateArgErrorMsg)
		return nil
	}

	if account.MultisigAddress == nil || !account.MultisigAddress.Valid {
		return errors.New("User attempted to nominate without a multisig address")
	}

	balance, err := contracts.GetTokenBalance(account.MultisigAddress.String)
	if err != nil {
		return err
	}

	// Do they have enough funds to nominate?
	if balance.Cmp(contracts.GetAtomicTokenAmount(depositAmount)) == -1 {
		sendDM(nominateInsufficientFundsMsg)
		return nil
	}

	// Does this handle already have an active listing?
	handle := parseHandle(argv[1])
	alreadyApplied, err := contracts.ApplicationWasMade(handle)
	if err != nil {
		return err
	}

	if alreadyApplied {
		sendDM(nominateAlreadyAppliedMsg)
		return nil
	}

	// Send them a message letting them know the gears are in motion
	sendDM(nominateBeginMsg)

	// Apply
	tokens := contracts.GetAtomicTokenAmount(depositAmount)
	tx, err := contracts.Apply(account.MultisigAddress.String, tokens, handle)
	if err != nil {
		sendDM(fmt.Sprintf(nominateSubmissionErrorMsg, err.Error()))
		return err
	}
	msg := fmt.Sprintf(nominateSuccessMsg, tx.Hash().Hex())
	sendDM(msg)

	return nil
}

func handleChallenge(account *models.Account, argv []string, sendDM func(string)) error {
	if len(argv) != 2 {
		sendDM(challengeArgErrorMsg)
		return nil
	}

	if account.MultisigAddress == nil || !account.MultisigAddress.Valid {
		return errors.New("User attempted to challenge without a multisig address")
	}

	balance, err := contracts.GetTokenBalance(account.MultisigAddress.String)
	if err != nil {
		return err
	}

	if balance.Cmp(contracts.GetAtomicTokenAmount(depositAmount)) == -1 {
		sendDM(challengeInsufficientFundsMsg)
		return nil
	}

	handle := parseHandle(argv[1])
	listing, err := contracts.GetListingFromHandle(handle)
	if err != nil {
		return err
	} else if listing == nil {
		sendDM(challengeNotFoundMsg)
		return nil
	} else if listing.ChallengeID.Cmp(big.NewInt(0)) != 0 {
		challenge, err := contracts.GetChallenge(listing.ChallengeID)
		if err != nil {
			return err
		}

		if !challenge.Resolved {
			sendDM(fmt.Sprintf(challengeAlreadyExistsMsg, handle))
			return nil
		}
	}

	tokens := contracts.GetAtomicTokenAmount(depositAmount)
	tx, err := contracts.CreateChallenge(account.MultisigAddress.String, tokens, handle)
	if err != nil {
		sendDM(challengeSubmissionErrorMsg)
		return err
	}
	msg := fmt.Sprintf(challengeSubmissionSuccessMsg, tx.Hash().Hex())
	sendDM(msg)

	return nil
}
