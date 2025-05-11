package handlers

import (
	"encoding/json"
	"math/big"
	"net/http"
	"tokenhub-api/internal/services"

	"github.com/ethereum/go-ethereum/common"
)

func HandleERC721Balance(svc services.NFTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		walletAddress := r.URL.Query().Get("walletAddress")
		if walletAddress == "" {
			http.Error(w, "walletAddress is required", http.StatusBadRequest)
			return
		}

		contractAddress := r.URL.Query().Get("contractAddress")
		if contractAddress == "" {
			http.Error(w, "contractAddress is required", http.StatusBadRequest)
			return
		}

		resp, err := svc.GetERC721Details(walletAddress, contractAddress)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(resp)
	}
}

func HandleERC1155Balance(svc services.NFTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		walletAddress := r.URL.Query().Get("walletAddress")
		if walletAddress == "" {
			http.Error(w, "walletAddress is required", http.StatusBadRequest)
			return
		}

		contractAddress := r.URL.Query().Get("contractAddress")
		if contractAddress == "" {
			http.Error(w, "contractAddress is required", http.StatusBadRequest)
			return
		}

		resp, err := svc.GetERC1155Details(walletAddress, contractAddress)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(resp)
	}
}

type DeployERC721Request struct {
	TokenName   string `json:"tokenName"`
	TokenSymbol string `json:"tokenSymbol"`
}

func DeployERC721Handler(svc services.NFTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DeployERC721Request
		if r.Body == nil {
			http.Error(w, "Request body is empty", http.StatusBadRequest)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		addr, err := svc.DeployERC721(req.TokenName, req.TokenSymbol)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"contractAddress": addr.Hex()})
	}
}

type DeployERC1155Request struct {
	TokenName   string `json:"tokenName"`
	TokenSymbol string `json:"tokenSymbol"`
}

func DeployERC1155Handler(svc services.NFTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DeployERC1155Request
		if r.Body == nil {
			http.Error(w, "Request body is empty", http.StatusBadRequest)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		addr, err := svc.DeployERC1155(req.TokenName, req.TokenSymbol)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"contractAddress": addr.Hex()})
	}
}

type MintERC721Request struct {
	ContractAddress string `json:"contractAddress"`
	TokenURI        string `json:"tokenURI"`
}

func MintERC721Handler(svc services.NFTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req MintERC721Request
		if r.Body == nil {
			http.Error(w, "Request body is empty", http.StatusBadRequest)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		contractAddr := common.HexToAddress(req.ContractAddress)
		txHash, err := svc.MintERC721(contractAddr, req.TokenURI)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]string{
			"transactionHash": txHash,
			"status":          "minted",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

type MintERC1155Request struct {
	ContractAddress string `json:"contractAddress"`
	To              string `json:"to"`
	Amount          string `json:"amount"`
	TokenURI        string `json:"tokenURI"`
}

func MintERC1155Handler(svc services.NFTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req MintERC1155Request
		if r.Body == nil {
			http.Error(w, "Request body is empty", http.StatusBadRequest)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		amount, ok := new(big.Int).SetString(req.Amount, 10)
		if !ok {
			http.Error(w, "Invalid amount", http.StatusBadRequest)
			return
		}

		contractAddr := common.HexToAddress(req.ContractAddress)
		txHash, err := svc.MintERC1155(contractAddr, req.To, amount, req.TokenURI)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]string{
			"transactionHash": txHash,
			"status":          "minted",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

type BurnERC721Request struct {
	ContractAddress string `json:"contractAddress"`
	TokenID         string `json:"tokenId"`
}

func BurnERC721Handler(svc services.NFTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BurnERC721Request
		if r.Body == nil {
			http.Error(w, "Request body is empty", http.StatusBadRequest)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		tokenId, ok := new(big.Int).SetString(req.TokenID, 10)
		if !ok {
			http.Error(w, "Invalid token ID", http.StatusBadRequest)
			return
		}

		contractAddr := common.HexToAddress(req.ContractAddress)
		txHash, err := svc.BurnERC721(contractAddr, tokenId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]string{
			"transactionHash": txHash,
			"status":          "burned",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

type BurnERC1155Request struct {
	ContractAddress string `json:"contractAddress"`
	TokenID         string `json:"tokenId"`
	Amount          string `json:"amount"`
}

func BurnERC1155Handler(svc services.NFTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BurnERC1155Request
		if r.Body == nil {
			http.Error(w, "Request body is empty", http.StatusBadRequest)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		tokenId, ok1 := new(big.Int).SetString(req.TokenID, 10)
		amount, ok2 := new(big.Int).SetString(req.Amount, 10)
		if !ok1 || !ok2 {
			http.Error(w, "Invalid token ID or amount", http.StatusBadRequest)
			return
		}

		contractAddr := common.HexToAddress(req.ContractAddress)
		txHash, err := svc.BurnERC1155(contractAddr, tokenId, amount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := map[string]string{
			"transactionHash": txHash,
			"status":          "burned",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
