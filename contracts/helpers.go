package contracts

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func GetListingHash(twitterHandle string) [32]byte {
	if twitterHandle == "obstropolos" {
		twitterHandle = "Obstropolos"
	}

	listingHash := sha256.Sum256([]byte(twitterHandle))

	// Convert that hash into the type it needs to be
	var txListingHash [32]byte
	copy(txListingHash[:], listingHash[0:32])
	return txListingHash
}

func getPublicAddress(privateKeyString string) (common.Address, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		return common.Address{}, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Address{}, errors.New("Could not convert public key to ECDSA")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return address, nil
}

func setupTransactionOpts(privateKeyHex string, gasLimit int64) (*bind.TransactOpts, error) {
	client, err := GetClientSession()
	if err != nil {
		return nil, err
	}

	// Get private key
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("Could not convert public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return nil, err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	// Add a bit more to the gas price to increase the chances the tx will get
	// picked up.
	gasModifierStr := os.Getenv("GAS_MODIFIER")
	if gasModifierStr != "" {
		gasModifier, err := strconv.ParseInt(gasModifierStr, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("could not parse GAS_MODIFIER; check your environment")
		}

		gasPrice.Add(gasPrice, big.NewInt(gasModifier))
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(gasLimit)
	auth.GasPrice = gasPrice

	return auth, nil
}

func submitTransaction(multisigAddress string, tx *proxiedTransaction) (*types.Transaction, error) {
	client, err := GetClientSession()
	if err != nil {
		return nil, err
	}

	contractAddress := common.HexToAddress(multisigAddress)
	wallet, err := NewMultiSigWallet(contractAddress, client)
	if err != nil {
		return nil, err
	}

	// Try the transaction until it goes through
	return ensureTransactionSubmission(func() (*types.Transaction, error) {
		txOpts, err := setupTransactionOpts(os.Getenv("MASTER_PRIVATE_KEY"), 500000)
		if err != nil {
			return nil, err
		}

		return wallet.SubmitTransaction(txOpts, tx.To, tx.Value, tx.Data)
	})
}

type proxiedTransaction struct {
	To    common.Address
	Value *big.Int
	Data  []byte
}

func newProxiedTransaction(to common.Address, abiString string, method string, args ...interface{}) (*proxiedTransaction, error) {
	parsed, err := abi.JSON(strings.NewReader(abiString))
	data, err := parsed.Pack(method, args...)

	tx := &proxiedTransaction{
		To:    to,
		Value: big.NewInt(0),
		Data:  data,
	}

	return tx, err
}

type txSubmitter func() (*types.Transaction, error)

func ensureTransactionSubmission(submit txSubmitter) (*types.Transaction, error) {
	var tx *types.Transaction
	var err error
	for {
		tx, err = submit()
		if err != nil && err.Error() == core.ErrReplaceUnderpriced.Error() {
			// Underpriced transaction, let's try again in a bit
			log.Println("Underpriced tx, trying again in 5s")
			time.Sleep(5 * time.Second)
			continue
		} else if err != nil && err.Error() == core.ErrNonceTooLow.Error() {
			// Underpriced transaction, let's try again in a bit
			log.Println("Nonce too low, trying again in 5s")
			time.Sleep(5 * time.Second)
			continue
		} else if err != nil {
			// Some other error
			return nil, err
		} else {
			// Success!
			break
		}
	}

	return tx, err
}
