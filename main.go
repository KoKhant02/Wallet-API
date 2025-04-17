package main

import (
	"log"
	"net/http"
	"os"

	handler "Wallet-API/internal/handlers"
	service "Wallet-API/internal/services"
	utils "Wallet-API/internal/utils"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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

	tokenService := service.NewTokenService(conn.Client, conn.ERC20Token, conn.Auth)
	nftService := service.NewNFTService(conn.Client, conn.ERC721Token, conn.ERC1155Token, conn.Auth)
	r := mux.NewRouter()
	r.HandleFunc("/balance/erc20", handler.GetERC20BalanceHandler(tokenService)).Methods("GET")
	r.HandleFunc("/mint/erc20", handler.MintERC20Handler(tokenService)).Methods("POST")
	r.HandleFunc("/balance/erc721", handler.GetERC721BalanceHandler(nftService)).Methods("GET")
	r.HandleFunc("/balance/erc1155", handler.GetERC1155BalanceHandler(nftService)).Methods("GET")
	r.HandleFunc("/mint/erc721", handler.MintERC721Handler(nftService)).Methods("POST")
	r.HandleFunc("/mint/erc1155", handler.MintERC1155Handler(nftService)).Methods("POST")
	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
