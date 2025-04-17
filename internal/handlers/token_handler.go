package handler

import (
	service "Wallet-API/internal/services"
	"encoding/json"
	"math/big"
	"net/http"
)

func GetERC20BalanceHandler(svc service.TokenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		if address == "" {
			http.Error(w, "Missing address", http.StatusBadRequest)
			return
		}
		balance, err := svc.GetERC20Balance(address)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{
			"balance": balance.String(),
		})
	}
}

func GetERC721BalanceHandler(svc service.NFTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		if address == "" {
			http.Error(w, "Missing address", http.StatusBadRequest)
			return
		}
		balance, err := svc.GetERC721Balance(address)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"balance": balance.String()})
	}
}

func GetERC1155BalanceHandler(svc service.NFTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		tokenIdStr := r.URL.Query().Get("tokenId")
		tokenId, ok := new(big.Int).SetString(tokenIdStr, 10)
		if !ok {
			http.Error(w, "Invalid tokenId", http.StatusBadRequest)
			return
		}
		balance, err := svc.GetERC1155Balance(address, tokenId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"balance": balance.String()})
	}
}

type MintRequest struct {
	To     string `json:"to"`
	Amount string `json:"amount"`
}

func MintERC20Handler(service service.TokenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req MintRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		amount, ok := new(big.Int).SetString(req.Amount, 10)
		if !ok {
			http.Error(w, "Invalid amount", http.StatusBadRequest)
			return
		}

		err := service.MintERC20(req.To, amount)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{"status": "minted"})
	}
}

type MintERC721Request struct {
	TokenURI string `json:"tokenURI"`
}

func MintERC721Handler(svc service.NFTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req MintERC721Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		if err := svc.MintERC721(req.TokenURI); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "minted"})
	}
}

type MintERC1155Request struct {
	To       string `json:"to"`
	Amount   string `json:"amount"`
	TokenURI string `json:"tokenURI"`
}

func MintERC1155Handler(svc service.NFTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req MintERC1155Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		amount, ok := new(big.Int).SetString(req.Amount, 10)
		if !ok {
			http.Error(w, "Invalid amount", http.StatusBadRequest)
			return
		}
		if err := svc.MintERC1155(req.To, amount, req.TokenURI); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "minted"})
	}
}
