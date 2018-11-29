package events

import (
	"fmt"
	"log"

	"gitlab.com/alpinefresh/tcrpartybot/contracts"
	"gitlab.com/alpinefresh/tcrpartybot/models"
	"gitlab.com/alpinefresh/tcrpartybot/twitter"
)

const (
	newApplicationWithHandleTweet    = "New #TCRParty listing! @%s has nominated @%s to be on the list for %s TCRP. Challenge this application by DMing 'challenge @%s'."
	newApplicationWithoutHandleTweet = "New #TCRParty listing! @%s has been nominated to be on the list for %s TCRP. Challenge this application by DMing 'challenge @%s'."

	newChallengeTweet = "New #TCRParty challenge! @%s's listing has been put to the test. Send me a DM with 'vote %s yes/no' to determine their fate."

	initialTokenAmount = 50
)

func processMultisigWalletCreation(event *ETHEvent) error {
	instantiation, err := contracts.DecodeContractInstantiationEvent(event.Data)
	if err != nil {
		return err
	}

	account, err := models.FindAccountByMultisigFactoryIdentifier(instantiation.Identifier.Int64())
	if err != nil {
		return err
	} else if account == nil {
		log.Printf("Could not find account with identifier %d", instantiation.Identifier.Int64())
		return nil
	}

	// Link their newly created multisig address to the account
	multisigAddress := instantiation.Instantiation.Hex()
	err = account.SetMultisigAddress(multisigAddress)
	if err != nil {
		return err
	}

	log.Printf("Wallet at %s linked to %s\n", multisigAddress, account.TwitterHandle)

	// Mint them 50 tokens for voting
	atomicAmount := contracts.GetAtomicTokenAmount(initialTokenAmount)
	mintTx, err := contracts.MintTokens(multisigAddress, atomicAmount)
	if err != nil {
		return err
	}

	_, err = contracts.AwaitTransactionConfirmation(mintTx.Hash())
	if err != nil {
		return err
	}

	// And lock those tokens up into the voting contract
	_, err = contracts.PLCRDeposit(multisigAddress, atomicAmount)
	if err != nil {
		return err
	}

	return nil
}

func processNewApplication(event *ETHEvent) error {
	application, err := contracts.DecodeApplicationEvent(event.Topics, event.Data)
	if err != nil {
		return err
	}

	// See if we can find an applicant in our database
	log.Printf("New application from %s for %s (hash: %s)", application.Applicant.Hex(), application.Data, application.ListingHash)
	account, err := models.FindAccountByMultisigAddress(application.Applicant.Hex())
	if err != nil {
		return err
	}

	tweet := ""
	depositAmount := contracts.GetHumanTokenAmount(application.Deposit).String()
	if account != nil {
		tweet = fmt.Sprintf(
			newApplicationWithHandleTweet,
			account.TwitterHandle,
			application.Data,
			depositAmount,
			application.Data,
		)
	} else {
		tweet = fmt.Sprintf(
			newApplicationWithoutHandleTweet,
			application.Data,
			depositAmount,
			application.Data,
		)
	}

	return twitter.SendTweet(twitter.VIPBotHandle, tweet)
}

func processNewChallenge(event *ETHEvent) error {
	challenge, err := contracts.DecodeChallengeEvent(event.Topics, event.Data)
	if err != nil {
		return err
	}

	listing, err := contracts.GetListingFromHash(challenge.ListingHash)
	if err != nil {
		return err
	} else if listing == nil {
		return fmt.Errorf("Could not find listing for challenge %s (listing: %s)", challenge.ChallengeID, string(challenge.ListingHash[:]))
	}

	log.Printf("New challenge for %s (hash: 0x%x)", challenge.Data, challenge.ListingHash)

	tweet := fmt.Sprintf(
		newChallengeTweet,
		challenge.Data,
		challenge.Data,
	)

	return twitter.SendTweet(twitter.VIPBotHandle, tweet)
}
