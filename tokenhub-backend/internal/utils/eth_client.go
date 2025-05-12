package utils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthConnection struct {
	Client *ethclient.Client
	Auth   *bind.TransactOpts
}

func NewEthConnection(rpcURL, privateKeyHex string) (*EthConnection, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(11155111)) // Sepolia
	if err != nil {
		return nil, err
	}

	auth.Nonce = nil

	auth.Value = big.NewInt(0)
	auth.GasLimit = 0
	auth.GasPrice, err = client.SuggestGasPrice(auth.Context)
	if err != nil {
		return nil, err
	}

	return &EthConnection{
		Client: client,
		Auth:   auth,
	}, nil
}
