package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"go.uber.org/zap"

	"tokenhub-api/internal/router"
	"tokenhub-api/internal/services"
	"tokenhub-api/internal/utils"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	rpcURL := os.Getenv("SEPOLIA_RPC_URL")
	privateKey := os.Getenv("WALLET_PRIVATE_KEY")

	conn, err := utils.NewEthConnection(rpcURL, privateKey)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	tokenService := services.NewTokenService(
		conn.Client,
		conn.Auth,
	)

	nftService := services.NewNFTService(
		conn.Client,
		conn.Auth,
	)

	r := router.NewRouter(tokenService, nftService, logger)
	handler := cors.Default().Handler(r)

	logger.Info("Server running", zap.String("address", ":8080"))
	http.ListenAndServe(":8080", handler)
}
