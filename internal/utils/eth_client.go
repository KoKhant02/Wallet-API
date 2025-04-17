package util

import (
	"Wallet-API/contracts"
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthConnection struct {
	Client       *ethclient.Client
	Auth         *bind.TransactOpts
	ERC20Token   *contracts.ERC20
	ERC721Token  *contracts.ERC721
	ERC1155Token *contracts.ERC1155
}

func NewEthConnection(rpcURL, privateKeyHex, erc20Addr, erc721Addr, erc1155Addr string) (*EthConnection, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, err
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

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(11155111)) // Sepolia
	if err != nil {
		return nil, err
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice

	erc20Instance, err := contracts.NewERC20(common.HexToAddress(erc20Addr), client)
	if err != nil {
		return nil, err
	}

	erc721Instance, err := contracts.NewERC721(common.HexToAddress(erc721Addr), client)
	if err != nil {
		return nil, err
	}

	erc1155Instance, err := contracts.NewERC1155(common.HexToAddress(erc1155Addr), client)
	if err != nil {
		return nil, err
	}

	return &EthConnection{
		Client:       client,
		Auth:         auth,
		ERC20Token:   erc20Instance,
		ERC721Token:  erc721Instance,
		ERC1155Token: erc1155Instance,
	}, nil
}
