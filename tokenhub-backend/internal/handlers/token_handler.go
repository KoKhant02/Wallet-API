package handlers

import (
	"encoding/json"
	"math/big"
	"net/http"
	"tokenhub-api/internal/services"

	"github.com/ethereum/go-ethereum/common"
)

func HandleERC20Balance(svc services.TokenService) http.HandlerFunc {
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

		resp, err := svc.GetERC20Details(walletAddress, contractAddress)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(resp)
	}
}

type DeployERC20Request struct {
	TokenName     string `json:"tokenName"`
	TokenSymbol   string `json:"tokenSymbol"`
	InitialSupply string `json:"initialSupply"`
}

func DeployERC20Handler(svc services.TokenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DeployERC20Request
		if r.Body == nil {
			http.Error(w, "Request body is empty", http.StatusBadRequest)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		rawAmount, ok := new(big.Int).SetString(req.InitialSupply, 10)
		if !ok {
			http.Error(w, "Invalid amount", http.StatusBadRequest)
			return
		}

		addr, err := svc.DeployERC20(req.TokenName, req.TokenSymbol, rawAmount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"contractAddress": addr.Hex()})
	}
}

type MintRequest struct {
	ContractAddress string `json:"contractAddress"`
	To              string `json:"to"`
	Amount          string `json:"amount"`
}

func MintERC20Handler(svc services.TokenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req MintRequest
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
		txHash, err := svc.MintERC20(contractAddr, req.To, amount)
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

type BurnERC20Request struct {
	ContractAddress string `json:"contractAddress"`
	Amount          string `json:"amount"`
}

func BurnERC20Handler(svc services.TokenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BurnERC20Request
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
		txHash, err := svc.BurnERC20(contractAddr, amount)
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
