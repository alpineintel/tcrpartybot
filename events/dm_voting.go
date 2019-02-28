package events

import (
	"errors"
	"fmt"
	"github.com/dustin/go-humanize"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"

	"gitlab.com/alpinefresh/tcrpartybot/contracts"
	"gitlab.com/alpinefresh/tcrpartybot/models"
)

const (
	votingArgErrorMsg               = "Whoops, looks like you forgot something. Try again with something like 'vote [twitter handle] kick' or 'vote [twitter handle] keep [vote weight, default 50]'"
	votingListingNotFoundMsg        = "Hmm, I couldn't find a registry listing for that twitter handle. Are you sure they've been nominated to the registry?"
	votingChallengeNotFoundMsg      = "Looks like that twitter handle doesn't have a challenge opened on it yet. If you'd like to challenge their place on the registry respond with 'challenge %s'."
	votingChallengeErrorMsg         = "There was an error committing your vote. The admins have been notified!"
	votedAlreadyMsg                 = "Oops! Looks like you've already voted %s on this challenge."
	votingEndedMsg                  = "Ack! Looks like the voting period has ended for this challenge. Hang tight, we'll announce the result on %s."
	votingBeginMsg                  = "We've submitted your vote. Hang tight, we'll notify you once everything is confirmed."
	votingSuccessMsg                = "Your vote to %s %s's listing with a weight of %d has been confirmed!\n\nWe'll announce the results on %s.\n\nTx hash: %s"
	voteInsufficientFundsMsg        = "You don't have enough funds locked up to vote with a weight of %d. You currently have %d available. If you would like to lock up more tokens to increase your voting weight, reply with vote-deposit [amount]. 50 is usually a good starting number.\n\nRemember that any tokens you lock up for voting will be unavailable for use in nominations/challenges."
	voteBalanceMsg                  = "You have %d tokens deposited to vote. This means you can vote with a maximum weight of %d."
	plcrWithdrawArgErrorMsg         = "Whoops, looks like you forgot something. Try vote-withdraw [amount]"
	plcrDepositArgErrorMsg          = "Whoops, looks like you forgot something. Try vote-deposit [amount]"
	plcrDepositInsufficientFundsMsg = "Whoops, looks like you don't have enough tokens to deposit this amount to your maximum voting weight. Your current balance is %d"
	plcrDepositBeginMsg             = "I've submitted your tokens for deposit. Hang tight, I'll let you know when everything clears."
	plcrDepositSuccessMsg           = "Your tokens have been deposited successfully!\n\nYou now have %d tokens locked up to vote and %d tokens in your wallet.\n\nTX hash: %s"
	plcrWithdrawInsufficientFunds   = "Whoops, you only have %d tokens locked up."
	plcrWithdrawBeginMsg            = "I've submitted your request to withdraw. Hang tight, I'll let you know when everything clears."
	plcrWithdrawSuccessMsg          = "Your tokens have been withdrawn successfully!\n\nYou now have %d tokens locked up to vote and %d tokens in your wallet.\n\nTX hash: %s"
	plcrLockedTokensMsg             = "Whoops! You currently have %d tokens locked up in existing challenges. These tokens cannot be withdrawn until challenge for %s's listing has completed (%s)."
	invalidNumberMsg                = "That doesn't look like a valid number to deposit... Get outta here."
)

func handleVote(account *models.Account, argv []string, sendDM func(string)) error {
	if len(argv) < 3 {
		sendDM(votingArgErrorMsg)
		return nil
	}

	voteDirection := strings.ToLower(argv[2])
	if voteDirection != "keep" && voteDirection != "kick" {
		sendDM(votingArgErrorMsg)
		return nil
	}

	if account.MultisigAddress == nil || !account.MultisigAddress.Valid {
		return errors.New("User attempted to vote without a multisig address")
	}

	// Fetch their PLCR deposit
	balance, err := contracts.PLCRFetchBalance(account.MultisigAddress.String)
	if err != nil {
		return err
	}

	weight := balance
	humanWeight := contracts.GetHumanTokenAmount(weight).Int64()

	// If their PLCR deposit is 0 maybe we can help them out
	if balance.Cmp(big.NewInt(0)) == 0 {
		walletBalance, err := contracts.GetTokenBalance(account.MultisigAddress.String)
		if err != nil {
			return err
		}

		// They don't even have enough tokens in their balance, let's stop
		compare := walletBalance.Cmp(big.NewInt(defaultVoteWeight))
		if compare == -1 {
			msg := fmt.Sprintf(voteInsufficientFundsMsg, humanWeight, contracts.GetHumanTokenAmount(balance).Int64())
			sendDM(msg)
			return nil
		}

		// Cool, let's put some tokens in their account to vote with
		toDeposit := contracts.GetAtomicTokenAmount(defaultVoteWeight)
		tx, err := contracts.PLCRDeposit(account.MultisigAddress.String, toDeposit)
		if err != nil {
			sendDM(fmt.Sprintf(errorMsg, err.Error()))
			return err
		}

		log.Printf("%s did not have any tokens for voting, depositing on their behalf (tx: %s)", account.TwitterHandle, tx.Hash().Hex())
		_, err = contracts.AwaitTransactionConfirmation(tx.Hash())
		if err != nil {
			return err
		}

		// Fetch their balance again to make sure we're all good
		balance, err = contracts.PLCRFetchBalance(account.MultisigAddress.String)
		if err != nil {
			return err
		}
	}

	handle := parseHandle(argv[1])
	if len(argv) > 3 {
		intWeight, err := strconv.ParseInt(argv[3], 10, 64)
		if err != nil {
			intWeight = 50
		}

		weight = contracts.GetAtomicTokenAmount(intWeight)
	}

	// Make sure their weight is within their means
	humanWeight = contracts.GetHumanTokenAmount(weight).Int64()
	if balance.Cmp(weight) == -1 {
		msg := fmt.Sprintf(voteInsufficientFundsMsg, humanWeight, contracts.GetHumanTokenAmount(balance).Int64())
		sendDM(msg)
		return nil
	}

	// Check to make sure there is an active poll for the given listing
	listing, err := contracts.GetListingFromHandle(handle)
	if err != nil {
		return err
	} else if listing == nil {
		// Check to see if the listing is using an @
		listing, err = contracts.GetListingFromHandle("@" + handle)
		if err != nil {
			return err
		} else if listing == nil {
			sendDM(votingListingNotFoundMsg)
			return nil
		}
	}

	if listing.ChallengeID.Cmp(big.NewInt(0)) == 0 {
		sendDM(fmt.Sprintf(votingChallengeNotFoundMsg, handle))
		return nil
	}

	// Fetch the poll they want to vote on
	poll, err := contracts.GetPoll(listing.ChallengeID)
	if err != nil {
		return err
	} else if poll == nil {
		log.Printf("Poll doesn't exist for listing: %x, challenge: %d", listing.ListingHash, listing.ChallengeID)
		sendDM(votingChallengeErrorMsg)
		return err
	}

	commitEndDate := time.Unix(poll.CommitEndDate.Int64(), 0)
	revealDate := time.Unix(poll.RevealEndDate.Int64(), 0)
	fmtRevealDate := revealDate.Format(time.RFC1123)

	// Make sure we're still in the commit phase
	if commitEndDate.Before(time.Now()) {
		sendDM(fmt.Sprintf(votingEndedMsg, fmtRevealDate))
		return nil
	}

	// Make sure they aren't voting twice
	vote, err := models.FindVote(listing.ChallengeID.Int64(), account.ID)
	if err != nil {
		return err
	} else if vote != nil {
		voteValue := "keep"
		if !vote.Vote {
			voteValue = "kick"
		}
		sendDM(fmt.Sprintf(votedAlreadyMsg, voteValue))
		return nil
	}

	sendDM(votingBeginMsg)

	voteValue := voteDirection == "keep"
	salt, tx, err := contracts.PLCRCommitVote(account.MultisigAddress.String, listing.ChallengeID, weight, voteValue)
	if err != nil {
		return err
	}

	// Wait for the vote to clear
	if _, err := contracts.AwaitTransactionConfirmation(tx.Hash()); err != nil {
		return err
	}

	humanWeight = contracts.GetHumanTokenAmount(weight).Int64()
	_, err = models.CreateVote(account, listing.ChallengeID.Int64(), salt, voteValue, humanWeight)
	if err != nil {
		return err
	}

	sendDM(fmt.Sprintf(votingSuccessMsg, voteDirection, handle, humanWeight, fmtRevealDate, tx.Hash().Hex()))
	return nil
}

func handleVoteBalance(account *models.Account, argv []string, sendDM func(string)) error {
	if account.MultisigAddress == nil || !account.MultisigAddress.Valid {
		err := errors.New("User attempted to fetch PLCR balance without a multisig address")
		sendDM(fmt.Sprintf(errorMsg, err.Error()))
		return err
	}

	balance, err := contracts.PLCRFetchBalance(account.MultisigAddress.String)
	if err != nil {
		sendDM(fmt.Sprintf(errorMsg, err.Error()))
		return err
	}

	humanBalance := contracts.GetHumanTokenAmount(balance).Int64()
	sendDM(fmt.Sprintf(voteBalanceMsg, humanBalance, humanBalance))
	return nil
}

func handleVoteDeposit(account *models.Account, argv []string, sendDM func(string)) error {
	if len(argv) < 2 {
		sendDM(plcrDepositArgErrorMsg)
		return nil
	}

	if !account.MultisigAddress.Valid {
		err := errors.New("User attempted to PLCRDeposit without a multisig address")
		sendDM(fmt.Sprintf(errorMsg, err.Error()))
		return err
	}

	amount, err := strconv.ParseInt(argv[1], 10, 64)
	if err != nil {
		sendDM(invalidNumberMsg)
		return nil
	} else if amount <= 0 {
		sendDM(invalidNumberMsg)
		return nil
	}

	toDeposit := contracts.GetAtomicTokenAmount(amount)

	balance, err := contracts.GetTokenBalance(account.MultisigAddress.String)
	if err != nil {
		return err
	}

	if balance.Cmp(toDeposit) == -1 {
		msg := fmt.Sprintf(plcrDepositInsufficientFundsMsg, contracts.GetHumanTokenAmount(balance).Int64())
		sendDM(msg)
		return nil
	}

	sendDM(plcrDepositBeginMsg)
	tx, err := contracts.PLCRDeposit(account.MultisigAddress.String, toDeposit)
	if err != nil {
		sendDM(fmt.Sprintf(errorMsg, err.Error()))
		return err
	}

	// Wait for the deposit to complete
	if _, err := contracts.AwaitTransactionConfirmation(tx.Hash()); err != nil {
		return err
	}

	// Get their new balance
	plcrBalance, err := contracts.PLCRFetchBalance(account.MultisigAddress.String)
	if err != nil {
		return err
	}
	humanPLCRBalance := contracts.GetHumanTokenAmount(plcrBalance)

	walletBalance, err := contracts.GetTokenBalance(account.MultisigAddress.String)
	if err != nil {
		return err
	}
	humanWalletBalance := contracts.GetHumanTokenAmount(walletBalance)

	sendDM(fmt.Sprintf(plcrDepositSuccessMsg, humanPLCRBalance, humanWalletBalance, tx.Hash().Hex()))
	return nil
}

func handleVoteWithdraw(account *models.Account, argv []string, sendDM func(string)) error {
	if len(argv) < 2 {
		sendDM(plcrWithdrawArgErrorMsg)
		return nil
	}

	if !account.MultisigAddress.Valid {
		err := errors.New("User attempted to PLCRwithdraw without a multisig address")
		sendDM(fmt.Sprintf(errorMsg, err.Error()))
		return err
	}
	amount, err := strconv.ParseInt(argv[1], 10, 64)
	if err != nil {
		sendDM(invalidNumberMsg)
		return nil
	} else if amount <= 0 {
		sendDM(invalidNumberMsg)
		return nil
	}

	toWithdraw := contracts.GetAtomicTokenAmount(amount)

	balance, err := contracts.PLCRFetchBalance(account.MultisigAddress.String)
	if err != nil {
		return err
	}

	if balance.Cmp(toWithdraw) == -1 {
		msg := fmt.Sprintf(plcrWithdrawInsufficientFunds, contracts.GetHumanTokenAmount(balance).Int64())
		sendDM(msg)
		return nil
	}

	// Are their tokens locked up in challenges?
	lockedTokens, err := contracts.PLCRLockedTokens(account.MultisigAddress.String)
	if err != nil {
		return err
	}

	if lockedTokens.Cmp(toWithdraw) == 1 || lockedTokens.Cmp(toWithdraw) == 0 {
		challengeID, err := contracts.PLCRLockedChallengeID(account.MultisigAddress.String)
		if err != nil {
			return err
		}

		// Get the challenge
		challenge, err := models.FindRegistryChallengeByPollID(challengeID.Int64())
		if err != nil {
			return err
		}

		lockedTokenAmt := contracts.GetHumanTokenAmount(lockedTokens).Int64()
		challengeHandle := challenge.Listing.TwitterHandle
		resolveAt := humanize.Time(*challenge.RevealEndsAt)

		msg := fmt.Sprintf(plcrLockedTokensMsg, lockedTokenAmt, challengeHandle, resolveAt)
		sendDM(msg)
		return nil
	}

	sendDM(plcrWithdrawBeginMsg)
	tx, err := contracts.PLCRWithdraw(account.MultisigAddress.String, toWithdraw)
	if err != nil {
		sendDM(fmt.Sprintf(errorMsg, err.Error()))
		return err
	}

	// Wait for the deposit to complete
	if _, err := contracts.AwaitTransactionConfirmation(tx.Hash()); err != nil {
		return err
	}

	// Get their new balance
	plcrBalance, err := contracts.PLCRFetchBalance(account.MultisigAddress.String)
	if err != nil {
		return err
	}
	humanPLCRBalance := contracts.GetHumanTokenAmount(plcrBalance)

	walletBalance, err := contracts.GetTokenBalance(account.MultisigAddress.String)
	if err != nil {
		return err
	}
	humanWalletBalance := contracts.GetHumanTokenAmount(walletBalance)

	sendDM(fmt.Sprintf(plcrWithdrawSuccessMsg, humanPLCRBalance.Int64(), humanWalletBalance.Int64(), tx.Hash().Hex()))
	return nil
}
