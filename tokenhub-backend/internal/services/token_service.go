package services

import (
	"context"
	"fmt"
	"log"
	"math/big"
	erc20 "tokenhub-api/contracts/ERC20"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/ethclient"
)

type TokenService interface {
	GetERC20Details(walletAddr string, contractAddr string) (*ERC20BalanceResponse, error)
	DeployERC20(name, symbol string, initialSupply *big.Int) (*ERC20DeployResponse, error)
	MintERC20(contractAddr common.Address, to string, amount *big.Int) (string, error)
	BurnERC20(contractAddr common.Address, amount *big.Int) (string, error)
}

type tokenService struct {
	client       *ethclient.Client
	erc20        *erc20.Contracts
	erc20Address common.Address
	auth         *bind.TransactOpts
}

func NewTokenService(client *ethclient.Client, auth *bind.TransactOpts) TokenService {
	return &tokenService{client: client, auth: auth}
}

type ERC20BalanceResponse struct {
	TokenName   string `json:"tokenName"`
	TokenSymbol string `json:"tokenSymbol"`
	Address     string `json:"address"`
	Balance     string `json:"balance"`
}

func (s *tokenService) GetERC20Details(walletAddr string, contractAddr string) (*ERC20BalanceResponse, error) {
	wallet := common.HexToAddress(walletAddr)
	contract := common.HexToAddress(contractAddr)

	fmt.Println("Wallet Address:", walletAddr)
	fmt.Println("Contract Address:", contractAddr)

	erc20, err := erc20.NewContracts(contract, s.client)
	if err != nil {
		return nil, err
	}

	name, err := erc20.Name(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	symbol, err := erc20.Symbol(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	decimals, err := erc20.Decimals(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	balance, err := erc20.BalanceOf(&bind.CallOpts{}, wallet)
	if err != nil {
		return nil, err
	}

	convertedBalance := new(big.Float).Quo(
		new(big.Float).SetInt(balance),
		new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)),
	)

	fmt.Println("Token Name:", name, "Token Symbol:", symbol, "Balance:", convertedBalance.String())

	return &ERC20BalanceResponse{
		TokenName:   name,
		TokenSymbol: symbol,
		Address:     contract.Hex(),
		Balance:     convertedBalance.String(),
	}, nil
}

type ERC20DeployResponse struct {
	TokenName   string `json:"tokenName"`
	TokenSymbol string `json:"tokenSymbol"`
	Address     string `json:"address"`
	TotalSupply string `json:"totalSupply"`
}

func (s *tokenService) DeployERC20(name, symbol string, initialSupply *big.Int) (*ERC20DeployResponse, error) {
	decimals := 18
	scaleFactor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	scaledSupply := new(big.Int).Mul(initialSupply, scaleFactor)

	address, tx, instance, err := erc20.DeployContracts(s.auth, s.client, name, symbol, scaledSupply)
	if err != nil {
		return nil, err
	}

	log.Printf("ERC20 deployed at: %s (tx: %s)", address.Hex(), tx.Hash().Hex())

	_, err = bind.WaitDeployed(context.Background(), s.client, tx)
	if err != nil {
		return nil, fmt.Errorf("deployment tx failed: %v", err)
	}

	// Fetch on-chain data
	tokenName, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("error fetching name: %v", err)
	}

	tokenSymbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("error fetching symbol: %v", err)
	}

	totalSupply, err := instance.TotalSupply(&bind.CallOpts{})
	if err != nil {
		return nil, fmt.Errorf("error fetching total supply: %v", err)
	}

	humanReadableSupply := new(big.Float).Quo(
		new(big.Float).SetInt(totalSupply),
		new(big.Float).SetInt(scaleFactor),
	)

	humanReadableSupplyStr := humanReadableSupply.Text('f', 0)

	s.erc20 = instance
	s.erc20Address = address

	return &ERC20DeployResponse{
		TokenName:   tokenName,
		TokenSymbol: tokenSymbol,
		Address:     address.Hex(),
		TotalSupply: humanReadableSupplyStr,
	}, nil
}

func (s *tokenService) MintERC20(contractAddr common.Address, to string, amount *big.Int) (string, error) {
	instance, err := erc20.NewContracts(contractAddr, s.client)
	if err != nil {
		return "", err
	}

	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		return "", err
	}

	multiplier := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
	scaledAmount := new(big.Int).Mul(amount, multiplier)

	toAddr := common.HexToAddress(to)
	tx, err := instance.Mint(s.auth, toAddr, scaledAmount)
	if err != nil {
		return "", err
	}

	log.Println("Minted ERC20 token:", tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

func (s *tokenService) BurnERC20(contractAddr common.Address, amount *big.Int) (string, error) {
	instance, err := erc20.NewContracts(contractAddr, s.client)
	if err != nil {
		return "", err
	}

	scaledAmount := new(big.Int).Mul(amount, new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))

	tx, err := instance.Burn(s.auth, scaledAmount)
	if err != nil {
		return "", err
	}
	log.Println("Burned ERC20:", tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}
