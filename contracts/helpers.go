package contracts

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

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

	// Get public key
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

	txOpts, err := setupTransactionOpts(os.Getenv("MASTER_PRIVATE_KEY"), 500000)
	if err != nil {
		return nil, err
	}

	contractAddress := common.HexToAddress(multisigAddress)
	wallet, err := NewMultiSigWallet(contractAddress, client)
	if err != nil {
		return nil, err
	}

	submitTX, err := wallet.SubmitTransaction(txOpts, tx.To, tx.Value, tx.Data)

	if err != nil {
		return nil, err
	}

	return submitTX, err
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
