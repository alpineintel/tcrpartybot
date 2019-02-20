package contracts

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

// DecodeApplicationEvent decodes data from a topic list and an ABI-encoded
// byte slice into a RegistryApplication struct
func DecodeApplicationEvent(topics []common.Hash, data []byte) (*RegistryApplication, error) {
	registryABI, err := abi.JSON(strings.NewReader(string(RegistryABI)))
	if err != nil {
		return nil, err
	}

	application := &RegistryApplication{}
	err = registryABI.Unpack(application, "_Application", data)
	if err != nil {
		return nil, err
	}

	// Load in data from topics
	copy(application.ListingHash[:], topics[1].Bytes()[0:32])
	application.Applicant = common.BytesToAddress(topics[2].Bytes())

	return application, nil
}

// DecodeWithdrawalEvent decodes data from a topic list and an ABI-encoded byte
// slice into a RegistryWithdrawal struct
func DecodeWithdrawalEvent(topics []common.Hash, data []byte) (*RegistryWithdrawal, error) {
	registryABI, err := abi.JSON(strings.NewReader(string(RegistryABI)))
	if err != nil {
		return nil, err
	}

	withdrawal := &RegistryWithdrawal{}
	err = registryABI.Unpack(withdrawal, "_Withdrawal", data)
	if err != nil {
		return nil, err
	}

	// Load in data from topics
	copy(withdrawal.ListingHash[:], topics[1].Bytes()[0:32])
	withdrawal.Owner = common.BytesToAddress(topics[2].Bytes())

	return withdrawal, nil
}

// DecodeRewardClaimedEvent decodes data from a topic list and an ABI-encoded byte
// slice into a RegistryRewardClaimed struct
func DecodeRewardClaimedEvent(topics []common.Hash, data []byte) (*RegistryRewardClaimed, error) {
	registryABI, err := abi.JSON(strings.NewReader(string(RegistryABI)))
	if err != nil {
		return nil, err
	}

	claim := &RegistryRewardClaimed{}
	err = registryABI.Unpack(claim, "_RewardClaimed", data)
	if err != nil {
		return nil, err
	}

	// Load in data from topics
	claim.ChallengeID = topics[1].Big()
	claim.Voter = common.BytesToAddress(topics[2].Bytes())

	return claim, nil
}

// DecodeChallengeEvent decodes data from a topic list and an ABI-encoded byte
// slice into a RegistryChallenge struct
func DecodeChallengeEvent(topics []common.Hash, data []byte) (*RegistryChallenge, error) {
	registryABI, err := abi.JSON(strings.NewReader(string(RegistryABI)))
	if err != nil {
		return nil, err
	}

	challenge := &RegistryChallenge{}
	err = registryABI.Unpack(challenge, "_Challenge", data)
	if err != nil {
		return nil, err
	}

	// Load in data from topics
	copy(challenge.ListingHash[:], topics[1].Bytes()[0:32])
	challenge.Challenger = common.BytesToAddress(topics[2].Bytes())

	return challenge, nil
}

// DecodeContractInstantiationEvent decodes data from an ABI-encoded byte array
// slice into a MultiSigWalletFactoryContractInstantiation struct
func DecodeContractInstantiationEvent(data []byte) (*MultiSigWalletFactoryContractInstantiation, error) {
	walletFactoryABI, err := abi.JSON(strings.NewReader(string(MultiSigWalletFactoryABI)))
	if err != nil {
		return nil, err
	}

	instantiation := &MultiSigWalletFactoryContractInstantiation{}
	err = walletFactoryABI.Unpack(instantiation, "ContractInstantiation", data)

	return instantiation, err
}

// DecodeApplicationWhitelistedEvent decodes data from an ABI-encoded byte array
// slice into a RegistryApplicationWhitelisted struct
func DecodeApplicationWhitelistedEvent(topics []common.Hash, data []byte) (*RegistryApplicationWhitelisted, error) {
	// Load in data from topics
	event := &RegistryApplicationWhitelisted{}
	copy(event.ListingHash[:], topics[1].Bytes()[0:32])

	return event, nil
}

// DecodeApplicationRemovedEvent decodes data from an ABI-encoded byte array
// slice into a RegistryApplicationRemoved struct
func DecodeApplicationRemovedEvent(topics []common.Hash, data []byte) (*RegistryApplicationRemoved, error) {
	event := &RegistryApplicationRemoved{}

	// Load in data from topics
	copy(event.ListingHash[:], topics[1].Bytes()[0:32])

	return event, nil
}

// DecodeListingRemovedEvent decodes data from an ABI-encoded byte array
// slice into a RegistryListingRemoved struct
func DecodeListingRemovedEvent(topics []common.Hash, data []byte) (*RegistryListingRemoved, error) {
	event := &RegistryListingRemoved{}

	// Load in data from topics
	copy(event.ListingHash[:], topics[1].Bytes()[0:32])

	return event, nil
}

// DecodeListingWithdrawnEvent decodes data from an ABI-encoded byte array
// slice into a RegistryWithdrawn struct
func DecodeListingWithdrawnEvent(topics []common.Hash, data []byte) (*RegistryListingWithdrawn, error) {
	registryABI, err := abi.JSON(strings.NewReader(string(RegistryABI)))
	if err != nil {
		return nil, err
	}

	event := &RegistryListingWithdrawn{}
	err = registryABI.Unpack(event, "_ListingWithdrawn", data)
	if err != nil {
		return nil, err
	}

	// Load in data from topics
	copy(event.ListingHash[:], topics[1].Bytes()[0:32])
	copy(event.Owner[:], topics[2].Bytes()[0:32])

	return event, nil
}

// DecodeTouchAndRemovedEvent decodes data from an ABI-encoded byte array
// slice into a RegistryTouchAndRemoved struct
func DecodeTouchAndRemovedEvent(topics []common.Hash, data []byte) (*RegistryTouchAndRemoved, error) {
	registryABI, err := abi.JSON(strings.NewReader(string(RegistryABI)))
	if err != nil {
		return nil, err
	}

	event := &RegistryTouchAndRemoved{}
	err = registryABI.Unpack(event, "_TouchAndRemoved", data)
	if err != nil {
		return nil, err
	}

	// Load in data from topics
	copy(event.ListingHash[:], topics[1].Bytes()[0:32])

	return event, nil
}

// DecodeDepositEvent decodes data from an ABI-encoded byte array slice into a
// RegistryDeposit struct
func DecodeDepositEvent(topics []common.Hash, data []byte) (*RegistryDeposit, error) {
	registryABI, err := abi.JSON(strings.NewReader(string(RegistryABI)))
	if err != nil {
		return nil, err
	}

	event := &RegistryDeposit{}
	err = registryABI.Unpack(event, "_Deposit", data)
	if err != nil {
		return nil, err
	}

	// Load in data from topics
	copy(event.ListingHash[:], topics[1].Bytes()[0:32])
	event.Owner = common.BytesToAddress(topics[2].Bytes())

	return event, nil
}

// DecodeExitInitializedEvent decodes data from an ABI-encoded byte array slice into a
// RegistryDeposit struct
func DecodeExitInitializedEvent(topics []common.Hash, data []byte) (*RegistryExitInitialized, error) {
	registryABI, err := abi.JSON(strings.NewReader(string(RegistryABI)))
	if err != nil {
		return nil, err
	}

	event := &RegistryExitInitialized{}
	err = registryABI.Unpack(event, "_ExitInitialized", data)
	if err != nil {
		return nil, err
	}

	// Load in data from topics
	copy(event.ListingHash[:], topics[1].Bytes()[0:32])
	event.Owner = common.BytesToAddress(topics[2].Bytes())

	return event, nil
}

// DecodeChallengeSucceededEvent decodes data from an ABI-encoded byte array
// slice into a RegistryChallengeSucceeded struct
func DecodeChallengeSucceededEvent(topics []common.Hash, data []byte) (*RegistryChallengeSucceeded, error) {
	registryABI, err := abi.JSON(strings.NewReader(string(RegistryABI)))
	if err != nil {
		return nil, err
	}

	event := &RegistryChallengeSucceeded{}
	err = registryABI.Unpack(event, "_ChallengeSucceeded", data)
	if err != nil {
		return nil, err
	}

	// Load in data from topics
	copy(event.ListingHash[:], topics[1].Bytes()[0:32])
	event.ChallengeID = topics[2].Big()

	return event, nil
}

// DecodeChallengeFailedEvent decodes data from an ABI-encoded byte array
// slice into a RegistryChallengeFailed struct
func DecodeChallengeFailedEvent(topics []common.Hash, data []byte) (*RegistryChallengeFailed, error) {
	registryABI, err := abi.JSON(strings.NewReader(string(RegistryABI)))
	if err != nil {
		return nil, err
	}

	event := &RegistryChallengeFailed{}
	err = registryABI.Unpack(event, "_ChallengeFailed", data)
	if err != nil {
		return nil, err
	}

	// Load in data from topics
	copy(event.ListingHash[:], topics[1].Bytes()[0:32])
	event.ChallengeID = topics[2].Big()

	return event, nil
}

// DecodeTransferEvent decodes data from an ABI-encoded byte array slice into a
// TCRPartyPointsTransfer struct
func DecodeTransferEvent(topics []common.Hash, data []byte) (*TCRPartyPointsTransfer, error) {
	tokenABI, err := abi.JSON(strings.NewReader(string(TCRPartyPointsABI)))
	if err != nil {
		return nil, err
	}

	event := &TCRPartyPointsTransfer{}
	err = tokenABI.Unpack(event, "Transfer", data)
	if err != nil {
		return nil, err
	}

	// Load in data from topics
	event.From = common.BytesToAddress(topics[1].Bytes())
	event.To = common.BytesToAddress(topics[2].Bytes())

	return event, nil
}

// DecodeApprovalEvent decodes data from an ABI-encoded byte array slice into a
// TCRPartyPointsApproval struct
func DecodeApprovalEvent(topics []common.Hash, data []byte) (*TCRPartyPointsApproval, error) {
	tokenABI, err := abi.JSON(strings.NewReader(string(TCRPartyPointsABI)))
	if err != nil {
		return nil, err
	}

	event := &TCRPartyPointsApproval{}
	err = tokenABI.Unpack(event, "Approval", data)
	if err != nil {
		return nil, err
	}

	// Load in data from topics
	event.Owner = common.BytesToAddress(topics[1].Bytes())
	event.Spender = common.BytesToAddress(topics[2].Bytes())

	return event, nil
}

// DecodeMintEvent decodes data from an ABI-encoded byte array slice into a
// TCRPartyPointsMint struct
func DecodeMintEvent(topics []common.Hash, data []byte) (*TCRPartyPointsMint, error) {
	tokenABI, err := abi.JSON(strings.NewReader(string(TCRPartyPointsABI)))
	if err != nil {
		return nil, err
	}

	event := &TCRPartyPointsMint{}
	err = tokenABI.Unpack(event, "Mint", data)
	if err != nil {
		return nil, err
	}

	// Load in data from topics
	event.To = common.BytesToAddress(topics[1].Bytes())

	return event, nil
}

// DecodeMintFinishedEvent decodes data from an ABI-encoded byte array slice into a
// TCRPartyPointsMintFinished struct
func DecodeMintFinishedEvent(topics []common.Hash, data []byte) (*TCRPartyPointsMintFinished, error) {
	tokenABI, err := abi.JSON(strings.NewReader(string(TCRPartyPointsABI)))
	if err != nil {
		return nil, err
	}

	event := &TCRPartyPointsMintFinished{}
	err = tokenABI.Unpack(event, "MintFinished", data)
	if err != nil {
		return nil, err
	}

	return event, nil
}

// DecodeVoteCommittedEvent decodes data from an ABI-encoded byte array slice into a
// PCLRVotingVoteCommitted struct
func DecodeVoteCommittedEvent(topics []common.Hash, data []byte) (*PLCRVotingVoteCommitted, error) {
	plcrABI, err := abi.JSON(strings.NewReader(string(PLCRVotingABI)))
	if err != nil {
		return nil, err
	}

	event := &PLCRVotingVoteCommitted{}
	err = plcrABI.Unpack(event, "_VoteCommitted", data)
	if err != nil {
		return nil, err
	}

	// Load in data from topics
	event.PollID = new(big.Int)
	event.PollID.SetBytes(topics[1].Bytes()[0:32])
	event.Voter = common.BytesToAddress(topics[2].Bytes())

	return event, nil
}

// DecodeVoteRevealedEvent decodes data from an ABI-encoded byte array slice into a
// PCLRVotingVoteCommitted struct
func DecodeVoteRevealedEvent(topics []common.Hash, data []byte) (*PLCRVotingVoteRevealed, error) {
	plcrABI, err := abi.JSON(strings.NewReader(string(PLCRVotingABI)))
	if err != nil {
		return nil, err
	}

	event := &PLCRVotingVoteRevealed{}
	err = plcrABI.Unpack(event, "_VoteRevealed", data)
	if err != nil {
		return nil, err
	}

	// Load in data from topics
	event.PollID = new(big.Int)
	event.PollID.SetBytes(topics[1].Bytes()[0:32])
	event.Choice = new(big.Int)
	event.Choice.SetBytes(topics[2].Bytes()[0:32])
	event.Voter = common.BytesToAddress(topics[3].Bytes())

	return event, nil
}

// DecodePollCreatedEvent decodes data from an ABI-encoded byte array slice into a
// PLCRVotingPollCreated struct
func DecodePollCreatedEvent(topics []common.Hash, data []byte) (*PLCRVotingPollCreated, error) {
	plcrABI, err := abi.JSON(strings.NewReader(string(PLCRVotingABI)))
	if err != nil {
		return nil, err
	}

	event := &PLCRVotingPollCreated{}
	err = plcrABI.Unpack(event, "_PollCreated", data)
	if err != nil {
		return nil, err
	}

	// Load in data from topics
	event.PollID = new(big.Int)
	event.PollID.SetBytes(topics[1].Bytes()[0:32])
	event.Creator = common.BytesToAddress(topics[2].Bytes())

	return event, nil
}

// DecodeVotingRightsGrantedEvent decodes data from an ABI-encoded byte array slice into a
// PLCRVotingVotingRightsGranted struct
func DecodeVotingRightsGrantedEvent(topics []common.Hash, data []byte) (*PLCRVotingVotingRightsGranted, error) {
	plcrABI, err := abi.JSON(strings.NewReader(string(PLCRVotingABI)))
	if err != nil {
		return nil, err
	}

	event := &PLCRVotingVotingRightsGranted{}
	err = plcrABI.Unpack(event, "_VotingRightsGranted", data)
	if err != nil {
		return nil, err
	}

	// Load in data from topics
	event.Voter = common.BytesToAddress(topics[1].Bytes())

	return event, nil
}

// DecodeVotingRightsWithdrawnEvent decodes data from an ABI-encoded byte array slice into a
// PLCRVotingVotingRightsWithdrawn struct
func DecodeVotingRightsWithdrawnEvent(topics []common.Hash, data []byte) (*PLCRVotingVotingRightsWithdrawn, error) {
	plcrABI, err := abi.JSON(strings.NewReader(string(PLCRVotingABI)))
	if err != nil {
		return nil, err
	}

	event := &PLCRVotingVotingRightsWithdrawn{}
	err = plcrABI.Unpack(event, "_VotingRightsWithdrawn", data)
	if err != nil {
		return nil, err
	}

	// Load in data from topics
	event.Voter = common.BytesToAddress(topics[1].Bytes())

	return event, nil
}

// DecodeTokensRescuedEvent decodes data from an ABI-encoded byte array slice into a
// PLCRVotingVotingRightsWithdrawn struct
func DecodeTokensRescuedEvent(topics []common.Hash, data []byte) (*PLCRVotingTokensRescued, error) {
	plcrABI, err := abi.JSON(strings.NewReader(string(PLCRVotingABI)))
	if err != nil {
		return nil, err
	}

	event := &PLCRVotingTokensRescued{}
	err = plcrABI.Unpack(event, "_TokensRescued", data)
	if err != nil {
		return nil, err
	}

	// Load in data from topics
	event.PollID = new(big.Int)
	event.PollID.SetBytes(topics[1].Bytes()[0:32])
	event.Voter = common.BytesToAddress(topics[2].Bytes())

	return event, nil
}
