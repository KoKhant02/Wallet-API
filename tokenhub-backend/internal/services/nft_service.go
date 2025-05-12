package services

import (
	"fmt"
	"log"
	"math/big"
	erc1155 "tokenhub-api/contracts/ERC1155"
	erc721 "tokenhub-api/contracts/ERC721"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type NFTService interface {
	GetERC721Details(walletAddr string, contractAddr string) (*NFTBalanceResponse, error)
	DeployERC721(name, symbol string) (common.Address, error)
	MintERC721(contractAddr common.Address, tokenURI string) (string, error)
	BurnERC721(contractAddr common.Address, tokenId *big.Int) (string, error)

	GetERC1155Details(walletAddr string, contractAddr string) (*NFTBalanceResponse, error)
	DeployERC1155(name, symbol string) (common.Address, error)
	MintERC1155(contractAddr common.Address, to string, amount *big.Int, tokenURI string) (string, error)
	BurnERC1155(contractAddr common.Address, tokenId *big.Int, amount *big.Int) (string, error)
}

type nftService struct {
	client         *ethclient.Client
	erc721         *erc721.Contracts
	erc721Address  common.Address
	erc1155        *erc1155.Contracts
	erc1155Address common.Address
	auth           *bind.TransactOpts
}

func NewNFTService(client *ethclient.Client, auth *bind.TransactOpts) NFTService {
	return &nftService{
		client: client,
		auth:   auth,
	}
}

type NFTBalanceResponse struct {
	TokenName   string    `json:"tokenName"`
	TokenSymbol string    `json:"tokenSymbol"`
	NFTItems    []NFTItem `json:"nftItems,omitempty"`
	Address     string    `json:"address"`
	TotalTokens string    `json:"totalTokens"`
}

type NFTItem struct {
	TokenURI string `json:"tokenURI,omitempty"`
	TokenID  string `json:"tokenId"`
	Amount   string `json:"amount,omitempty"`
}

func (s *nftService) GetERC721Details(walletAddr string, contractAddr string) (*NFTBalanceResponse, error) {
	addr := common.HexToAddress(walletAddr)
	contract := common.HexToAddress(contractAddr)

	erc721, err := erc721.NewContracts(contract, s.client)
	if err != nil {
		return nil, err
	}

	name, err := erc721.Name(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	symbol, err := erc721.Symbol(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	currentTokenId, err := erc721.GetCurrentTokenId(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}

	var nftItems []NFTItem
	for i := int64(1); i <= currentTokenId.Int64(); i++ {
		tokenId := big.NewInt(i)

		owner, err := erc721.OwnerOf(&bind.CallOpts{}, tokenId)
		if err != nil {
			fmt.Println("Error getting token owner:", err)
			continue
		}

		if owner != addr {
			continue
		}

		tokenURI, err := erc721.TokenURI(&bind.CallOpts{}, tokenId)
		if err != nil {
			fmt.Println("Error getting Token URI:", err)
			continue
		}

		nftItems = append(nftItems, NFTItem{
			TokenURI: tokenURI,
			TokenID:  tokenId.String(),
			Amount:   "1",
		})
	}

	balance := big.NewInt(int64(len(nftItems)))

	result := &NFTBalanceResponse{
		TokenName:   name,
		TokenSymbol: symbol,
		Address:     contract.Hex(),
		TotalTokens: balance.String(),
		NFTItems:    nftItems,
	}

	return result, nil
}

func (s *nftService) DeployERC721(name, symbol string) (common.Address, error) {
	address, tx, instance, err := erc721.DeployContracts(s.auth, s.client, name, symbol)
	if err != nil {
		return common.Address{}, err
	}
	log.Printf("ERC721 deployed at: %s (tx: %s)", address.Hex(), tx.Hash().Hex())

	s.erc721 = instance
	s.erc721Address = address
	return address, nil
}

func (s *nftService) MintERC721(contractAddr common.Address, tokenURI string) (string, error) {
	instance, err := erc721.NewContracts(contractAddr, s.client)
	if err != nil {
		return "", err
	}
	tx, err := instance.MintNFT(s.auth, tokenURI)
	if err != nil {
		return "", err
	}
	log.Println("Minted ERC721 NFT with tx:", tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

func (s *nftService) BurnERC721(contractAddr common.Address, tokenId *big.Int) (string, error) {
	instance, err := erc721.NewContracts(contractAddr, s.client)
	if err != nil {
		return "", err
	}
	tx, err := instance.BurnNFT(s.auth, tokenId)
	if err != nil {
		return "", err
	}
	log.Println("Burned ERC721 token:", tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

func (s *nftService) GetERC1155Details(walletAddr string, contractAddr string) (*NFTBalanceResponse, error) {
	addr := common.HexToAddress(walletAddr)
	contract := common.HexToAddress(contractAddr)

	erc1155, err := erc1155.NewContracts(contract, s.client)
	if err != nil {
		log.Println("Error initializing ERC1155 contract:", err)
		return nil, err
	}

	name, err := erc1155.Name(&bind.CallOpts{})
	if err != nil {
		log.Println("Error fetching name:", err)
		return nil, err
	}

	symbol, err := erc1155.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Println("Error fetching symbol:", err)
		return nil, err
	}

	var nftItems []NFTItem

	for tokenId := 0; tokenId < 100; tokenId++ {
		id := big.NewInt(int64(tokenId))

		balance, err := erc1155.BalanceOf(&bind.CallOpts{}, addr, id)
		if err != nil {
			log.Println("Error fetching balance:", err)
			continue
		}

		if balance.Cmp(big.NewInt(0)) > 0 {
			tokenURI, err := erc1155.Uri(&bind.CallOpts{}, id)
			if err != nil {
				log.Println("Error fetching token URI:", err)
				continue
			}

			nftItems = append(nftItems, NFTItem{
				TokenID:  id.String(),
				TokenURI: tokenURI,
				Amount:   balance.String(),
			})
		}
	}

	if len(nftItems) == 0 {
		return &NFTBalanceResponse{
			TokenName:   name,
			TokenSymbol: symbol,
			Address:     contract.Hex(),
			TotalTokens: "0",
			NFTItems:    []NFTItem{},
		}, nil
	}

	balance := big.NewInt(int64(len(nftItems)))

	result := &NFTBalanceResponse{
		TokenName:   name,
		TokenSymbol: symbol,
		Address:     contract.Hex(),
		TotalTokens: balance.String(),
		NFTItems:    nftItems,
	}

	return result, nil
}

func (s *nftService) DeployERC1155(name, symbol string) (common.Address, error) {
	address, tx, instance, err := erc1155.DeployContracts(s.auth, s.client, name, symbol)
	if err != nil {
		return common.Address{}, err
	}
	log.Printf("ERC1155 deployed at: %s (tx: %s)", address.Hex(), tx.Hash().Hex())

	s.erc1155 = instance
	s.erc1155Address = address
	return address, nil
}

func (s *nftService) MintERC1155(contractAddr common.Address, to string, amount *big.Int, tokenURI string) (string, error) {
	instance, err := erc1155.NewContracts(contractAddr, s.client)
	if err != nil {
		return "", err
	}
	toAddr := common.HexToAddress(to)
	tx, err := instance.Mint(s.auth, toAddr, amount, tokenURI)
	if err != nil {
		return "", err
	}
	log.Println("Minted ERC1155 with tx:", tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

func (s *nftService) BurnERC1155(contractAddr common.Address, tokenId *big.Int, amount *big.Int) (string, error) {
	instance, err := erc1155.NewContracts(contractAddr, s.client)
	if err != nil {
		return "", err
	}
	tx, err := instance.Burn(s.auth, s.auth.From, tokenId, amount)
	if err != nil {
		return "", err
	}
	log.Println("Burned ERC1155 token:", tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}
