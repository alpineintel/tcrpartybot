package events

import (
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"gitlab.com/alpinefresh/tcrpartybot/contracts"
	"gitlab.com/alpinefresh/tcrpartybot/models"
	"gitlab.com/alpinefresh/tcrpartybot/twitter"
)

const (
	challengeArgErrorMsg            = "Whoops, looks like you forgot something. Try again with something like 'challenge [twitter handle]'. Eg: 'challenge weratedogs'"
	challengeNotFoundMsg            = "Looks like nobody has tried creating a listing for this twitter handle yet."
	challengeAlreadyExistsMsg       = "Looks like somebody has already begun a challenge for this twitter handle. You can support this challenge by voting kick (respond with 'vote %s kick')."
	challengeInsufficientFundsMsg   = "Drat, looks like you don't have enough TCRP to start a challenge with a deposit of %d. Your current balance is %d TCRP."
	challengeInsufficientDepositMsg = "Ack, @%s's listing has a %d TCRP deposit. In order to challenge this listing you need to match their deposit. Try replying with \"challenge @%s %d\""
	challengeBeginMsg               = "Got it! I've just begun submitting your challenge to the registry and will send a message once everything is confirmed"
	challengeSubmissionErrorMsg     = "There was an error trying to submit your challenge: %s. Try tweeting at @stevenleeg for help."
	challengeSubmissionSuccessMsg   = "Done! Your challenge has been submitted to the registry. Your deposit was %d TCRP, leaving you with a balance of %d TCRP. Keep an eye on @TCRPartyVIP for updates.\n\nTX Hash: %s"

	nominateArgErrorMsg          = "Whoops, looks like you forgot something. Try again with something like 'nominate [twitter handle]'. Eg: 'apply weratedogs'"
	nominateAlreadyAppliedMsg    = "Looks like that Twitter handle has already been submitted to the TCR. A twitter handle can only appear on the TCR once, so you'll need to wait for a successful challenge (or a delisting) in order to re-nominate them."
	nominateInsufficientFundsMsg = "Drat! Looks like you don't have enough TCRP to start a nomination. Nominations require a deposit of at least 500 TCRP. Your current balance is %d."
	nominateSubmissionErrorMsg   = "There was an error trying to submit your nomination: %s. Try tweeting at @stevenleeg for help."
	nominateSuccessMsg           = "All done! Your nomination was submitted successfully and your new balance is %d TCRP.\n\nKeep an eye on @TCRPartyVIP for an announcement.\n\nTX Hash: %s"
	nominateBeginMsg             = "Got it! I've just begun submitting your nomination to the registry and will send a message once everything is confirmed."
	invalidHandleMsg             = "Hmm, it looks like @%s isn't a valid Twitter user"
	getOuttaHereMsg              = "Get 'outta here"
	tooManyListingsMsg           = "Ruh roh. Looks like the list has hit its limit of %d slots. In order to nominate a new member you'll need to challenge an existing member's slot."

	defaultDepositAmount = 500
)

func handleNomination(account *models.Account, argv []string, sendDM func(string)) error {
	if len(argv) < 2 {
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

	// Determine their deposit
	depositAmount := int64(defaultDepositAmount)
	if len(argv) == 3 {
		depositAmount, err = strconv.ParseInt(argv[2], 10, 64)
		if err != nil {
			depositAmount = defaultDepositAmount
		}
	}

	// Do they have enough funds to nominate?
	if balance.Cmp(contracts.GetAtomicTokenAmount(depositAmount)) == -1 {
		sendDM(fmt.Sprintf(nominateInsufficientFundsMsg, contracts.GetHumanTokenAmount(balance)))
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

	// Is the handle valid?
	if !verifyHandle(handle) {
		sendDM(fmt.Sprintf(invalidHandleMsg, handle))
		return nil
	}

	// Is this handle real?
	_, err = twitter.GetIDFromHandle(handle)
	if err != nil {
		sendDM(fmt.Sprintf(invalidHandleMsg, handle))
		return nil
	}

	// Is this us?
	if handle == "tcrpartybot" || handle == "tcrpartyvip" {
		sendDM(getOuttaHereMsg)
		return nil
	}

	// Aaand finally, have we hit our limits?
	listings, err := contracts.GetAllListings()
	if len(listings) >= 99 {
		sendDM(fmt.Sprintf(tooManyListingsMsg, len(listings)))
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

	balance, err = contracts.GetTokenBalance(account.MultisigAddress.String)
	if err != nil {
		return err
	}
	humanBalance := contracts.GetHumanTokenAmount(balance).Int64()
	msg := fmt.Sprintf(nominateSuccessMsg, humanBalance, tx.Hash().Hex())
	sendDM(msg)

	return nil
}

func handleChallenge(account *models.Account, argv []string, sendDM func(string)) error {
	if len(argv) < 2 {
		sendDM(challengeArgErrorMsg)
		return nil
	}

	if account.MultisigAddress == nil || !account.MultisigAddress.Valid {
		return errors.New("User attempted to challenge without a multisig address")
	}

	handle := parseHandle(argv[1])

	balance, err := contracts.GetTokenBalance(account.MultisigAddress.String)
	if err != nil {
		return err
	}

	// Determine their deposit
	intDepositAmount := int64(defaultDepositAmount)
	if len(argv) == 3 {
		intDepositAmount, err = strconv.ParseInt(argv[2], 10, 64)
		if err != nil {
			intDepositAmount = defaultDepositAmount
		}
	}
	depositAmount := contracts.GetAtomicTokenAmount(intDepositAmount)

	// Make sure they have enough in their balance or deposit to challenge this
	// listing
	if balance.Cmp(depositAmount) == -1 {
		humanBalance := contracts.GetHumanTokenAmount(balance)
		sendDM(fmt.Sprintf(challengeInsufficientFundsMsg, intDepositAmount, humanBalance))
		return nil
	}

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

	// Make sure their deposit is large enough to challenge this listing
	listingDeposit, err := contracts.GetListingDepositFromHash(listing.ListingHash)
	humanListingDeposit := contracts.GetHumanTokenAmount(listingDeposit)
	if depositAmount.Cmp(listingDeposit) == -1 {
		sendDM(fmt.Sprintf(
			challengeInsufficientDepositMsg,
			handle,
			humanListingDeposit,
			handle,
			humanListingDeposit,
		))
		return nil
	}

	// Send them a message letting them know the gears are in motion
	sendDM(challengeBeginMsg)

	tx, err := contracts.CreateChallenge(account.MultisigAddress.String, listingDeposit, handle)
	if err != nil {
		sendDM(fmt.Sprintf(challengeSubmissionErrorMsg, err.Error()))
		return err
	}

	balance, err = contracts.GetTokenBalance(account.MultisigAddress.String)
	if err != nil {
		return err
	}
	humanBalance := contracts.GetHumanTokenAmount(balance).Int64()

	msg := fmt.Sprintf(challengeSubmissionSuccessMsg, humanListingDeposit, humanBalance, tx.Hash().Hex())
	sendDM(msg)

	return nil
}
