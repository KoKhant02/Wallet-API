package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	handler "Wallet-API/internal/handlers"
	service "Wallet-API/internal/services"
	utils "Wallet-API/internal/utils"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	rpcURL := os.Getenv("SEPOLIA_RPC_URL")
	privateKey := os.Getenv("WALLET_PRIVATE_KEY")
	erc20Address := os.Getenv("ERC20_CONTRACT_ADDRESS")
	erc721Address := os.Getenv("ERC721_CONTRACT_ADDRESS")
	erc1155Address := os.Getenv("ERC1155_CONTRACT_ADDRESS")

	conn, err := utils.NewEthConnection(rpcURL, privateKey, erc20Address, erc721Address, erc1155Address)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	tokenService := service.NewTokenService(
		conn.Client,
		conn.ERC20Token,
		common.HexToAddress(erc20Address),
		conn.Auth,
	)

	nftService := service.NewNFTService(
		conn.Client,
		conn.ERC721Token,
		common.HexToAddress(erc721Address),
		conn.ERC1155Token,
		common.HexToAddress(erc1155Address),
		conn.Auth,
	)
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/balance/erc20", handler.HandleERC20Balance(tokenService)).Methods("GET")
	api.HandleFunc("/balance/erc721", handler.HandleERC721Balance(nftService)).Methods("GET")
	api.HandleFunc("/balance/erc1155", handler.HandleERC1155Balance(nftService)).Methods("GET")

	api.HandleFunc("/mint/erc20", handler.MintERC20Handler(tokenService)).Methods("POST")
	api.HandleFunc("/mint/erc721", handler.MintERC721Handler(nftService)).Methods("POST")
	api.HandleFunc("/mint/erc1155", handler.MintERC1155Handler(nftService)).Methods("POST")

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
