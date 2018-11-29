package contracts

import (
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