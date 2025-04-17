package service

import (
	"Wallet-API/contracts"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/ethclient"
)

type TokenService interface {
	GetERC20Balance(address string) (*big.Int, error)
	MintERC20(to string, amount *big.Int) error
}

type tokenService struct {
	client *ethclient.Client
	erc20  *contracts.ERC20
	auth   *bind.TransactOpts
}

func NewTokenService(client *ethclient.Client, erc20 *contracts.ERC20, auth *bind.TransactOpts) TokenService {
	return &tokenService{client, erc20, auth}
}

func (s *tokenService) GetERC20Balance(address string) (*big.Int, error) {
	addr := common.HexToAddress(address)
	return s.erc20.BalanceOf(&bind.CallOpts{}, addr)
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
	GetERC721Balance(address string) (*big.Int, error)
	MintERC721(tokenURI string) error
	GetERC1155Balance(address string, tokenId *big.Int) (*big.Int, error)
	MintERC1155(to string, amount *big.Int, tokenURI string) error
}

type nftService struct {
	client  *ethclient.Client
	erc721  *contracts.ERC721
	erc1155 *contracts.ERC1155
	auth    *bind.TransactOpts
}

func NewNFTService(client *ethclient.Client, erc721 *contracts.ERC721, erc1155 *contracts.ERC1155, auth *bind.TransactOpts) NFTService {
	return &nftService{client, erc721, erc1155, auth}
}

func (s *nftService) GetERC721Balance(address string) (*big.Int, error) {
	addr := common.HexToAddress(address)
	return s.erc721.BalanceOf(&bind.CallOpts{}, addr)
}

func (s *nftService) MintERC721(tokenURI string) error {
	tx, err := s.erc721.MintNFT(s.auth, tokenURI)
	if err != nil {
		return err
	}
	log.Println("Minted ERC721 NFT with tx:", tx.Hash().Hex())
	return nil
}

func (s *nftService) GetERC1155Balance(address string, tokenId *big.Int) (*big.Int, error) {
	addr := common.HexToAddress(address)
	return s.erc1155.BalanceOf(&bind.CallOpts{}, addr, tokenId)
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
