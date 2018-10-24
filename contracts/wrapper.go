package contracts

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var session *ethclient.Client

func getClientSession() (*ethclient.Client, error) {
	if session != nil {
		return session, nil
	}

	session, err := ethclient.Dial(os.Getenv("ETH_NODE_URI"))
	if err != nil {
		return nil, err
	}

	return session, nil
}

func GetTokenBalance(address string) (int64, error) {
	client, err := getClientSession()
	if err != nil {
		return 0, err
	}

	tokenAddress := common.HexToAddress(os.Getenv("TOKEN_ADDRESS"))
	token, err := NewTCRPartyPoints(tokenAddress, client)
	if err != nil {
		fmt.Println("something went wrong init")
		return 0, err
	}

	balance, err := token.BalanceOf(nil, common.HexToAddress(address))
	if err != nil {
		return 0, err
	}

	return balance.Int64(), nil
}

// MintTokens assigns new tokens to the given ETH address
func MintTokens(address string, amount int64) (*types.Transaction, error) {
	client, err := getClientSession()
	if err != nil {
		return nil, err
	}

	// Get private key
	privateKey, err := crypto.HexToECDSA(os.Getenv("MASTER_PRIVATE_KEY"))
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
	auth.GasLimit = uint64(500000)
	auth.GasPrice = gasPrice

	contractAddress := common.HexToAddress(os.Getenv("TOKEN_ADDRESS"))
	token, err := NewTCRPartyPoints(contractAddress, client)
	if err != nil {
		return nil, err
	}

	txAddress := common.HexToAddress(address)
	txAmount := big.NewInt(amount)

	tx, err := token.Mint(auth, txAddress, txAmount)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
