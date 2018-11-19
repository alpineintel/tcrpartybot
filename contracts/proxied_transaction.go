package contracts

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

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
