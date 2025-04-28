package service

import (
	"Wallet-API/contracts"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/ethclient"
)

type TokenService interface {
	GetERC20Details(address string) (*ERC20BalanceResponse, error)
	MintERC20(to string, amount *big.Int) error
}

type tokenService struct {
	client       *ethclient.Client
	erc20        *contracts.ERC20
	erc20Address common.Address
	auth         *bind.TransactOpts
}

func NewTokenService(client *ethclient.Client, erc20 *contracts.ERC20, erc20Address common.Address, auth *bind.TransactOpts) TokenService {
	return &tokenService{client: client, erc20: erc20, erc20Address: erc20Address, auth: auth}
}

type ERC20BalanceResponse struct {
	TokenName   string `json:"tokenName"`
	TokenSymbol string `json:"tokenSymbol"`
	Address     string `json:"address"`
	Balance     string `json:"balance"`
}

func (s *tokenService) GetERC20Details(address string) (*ERC20BalanceResponse, error) {
	addr := common.HexToAddress(address)
	fmt.Println("Wallet Address in service:", address)

	name, err := s.erc20.Name(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	symbol, err := s.erc20.Symbol(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	decimals, err := s.erc20.Decimals(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	balance, err := s.erc20.BalanceOf(&bind.CallOpts{}, addr)
	if err != nil {
		return nil, err
	}

	convertedBalance := new(big.Float).Quo(new(big.Float).SetInt(balance), new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)))

	fmt.Println("Token Name:", name, "Token Symbol:", symbol, "Balance:", convertedBalance.String())

	return &ERC20BalanceResponse{
		TokenName:   name,
		TokenSymbol: symbol,
		Address:     s.erc20Address.Hex(),
		Balance:     convertedBalance.String(),
	}, nil
}

func (s *tokenService) MintERC20(to string, amount *big.Int) error {
	toAddr := common.HexToAddress(to)
	tx, err := s.erc20.Mint(s.auth, toAddr, amount)
	if err != nil {
		return err
	}
	log.Println("Mint tx sent:", tx.Hash().Hex())
	return nil
}

type NFTService interface {
	GetERC721Details(walletAddress string) ([]*NFTBalanceResponse, error)
	MintERC721(tokenURI string) error
	GetERC1155Details(walletAddress string) ([]*NFTBalanceResponse, error)
	MintERC1155(to string, amount *big.Int, tokenURI string) error
}

type nftService struct {
	client         *ethclient.Client
	erc721         *contracts.ERC721
	erc721Address  common.Address
	erc1155        *contracts.ERC1155
	erc1155Address common.Address
	auth           *bind.TransactOpts
}

func NewNFTService(client *ethclient.Client, erc721 *contracts.ERC721, erc721Address common.Address, erc1155 *contracts.ERC1155, erc1155Address common.Address, auth *bind.TransactOpts) NFTService {
	return &nftService{
		client:         client,
		erc721:         erc721,
		erc721Address:  erc721Address,
		erc1155:        erc1155,
		erc1155Address: erc1155Address,
		auth:           auth,
	}
}

type NFTBalanceResponse struct {
	TokenName   string    `json:"tokenName"`
	TokenSymbol string    `json:"tokenSymbol"`
	NFTItems    []NFTItem `json:"nftItems"`
	Address     string    `json:"address"`
	TotalTokens string    `json:"totalTokens"`
}

type NFTItem struct {
	TokenURI string `json:"tokenURI"`
	TokenID  string `json:"tokenId"`
}

func (s *nftService) GetERC721Details(walletAddress string) ([]*NFTBalanceResponse, error) {
	addr := common.HexToAddress(walletAddress)

	name, err := s.erc721.Name(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	symbol, err := s.erc721.Symbol(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	currentTokenId, err := s.erc721.GetCurrentTokenId(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	var nftItems []NFTItem

	for i := int64(1); i <= currentTokenId.Int64(); i++ {
		tokenId := big.NewInt(i)

		owner, err := s.erc721.OwnerOf(&bind.CallOpts{}, tokenId)
		if err != nil {
			fmt.Println("Error getting token owner:", err)
			continue
		}

		if owner != addr {
			continue
		}

		tokenURI, err := s.erc721.TokenURI(&bind.CallOpts{}, tokenId)
		if err != nil {
			fmt.Println("Error getting Token URI:", err)
			continue
		}

		nftItems = append(nftItems, NFTItem{
			TokenURI: tokenURI,
			TokenID:  tokenId.String(),
		})
	}

	balance := big.NewInt(int64(len(nftItems)))

	result := &NFTBalanceResponse{
		TokenName:   name,
		TokenSymbol: symbol,
		Address:     s.erc721Address.Hex(),
		TotalTokens: balance.String(),
		NFTItems:    nftItems,
	}

	return []*NFTBalanceResponse{result}, nil
}

func (s *nftService) MintERC721(tokenURI string) error {
	tx, err := s.erc721.MintNFT(s.auth, tokenURI)
	if err != nil {
		return err
	}
	log.Println("Minted ERC721 NFT with tx:", tx.Hash().Hex())
	return nil
}
func (s *nftService) GetERC1155Details(walletAddress string) ([]*NFTBalanceResponse, error) {
	addr := common.HexToAddress(walletAddress)
	var results []*NFTBalanceResponse

	name, err := s.erc1155.Name(&bind.CallOpts{})
	if err != nil {
		log.Println("Error fetching name:", err)
	}
	symbol, err := s.erc1155.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Println("Error fetching symbol:", err)
	}

	for tokenId := 0; tokenId < 100; tokenId++ {
		id := big.NewInt(int64(tokenId))
		balance, err := s.erc1155.BalanceOf(&bind.CallOpts{}, addr, id)
		if err != nil {
			log.Println("Error fetching balance:", err)
			continue
		}

		if balance.Cmp(big.NewInt(0)) > 0 {
			tokenURI, err := s.erc1155.Uri(&bind.CallOpts{}, id)
			if err != nil {
				log.Println("Error fetching token URI:", err)
				continue
			}

			results = append(results, &NFTBalanceResponse{
				TokenName:   name,
				TokenSymbol: symbol,
				Address:     s.erc1155Address.Hex(),
				TotalTokens: balance.String(),
				NFTItems: []NFTItem{
					{
						TokenID:  id.String(),
						TokenURI: tokenURI,
					},
				},
			})
		}
	}

	return results, nil
}

func (s *nftService) MintERC1155(to string, amount *big.Int, tokenURI string) error {
	toAddr := common.HexToAddress(to)
	tx, err := s.erc1155.Mint(s.auth, toAddr, amount, tokenURI)
	if err != nil {
		return err
	}
	log.Println("Minted ERC1155 with tx:", tx.Hash().Hex())
	return nil
}
