package handler

import (
	service "Wallet-API/internal/services"
	"encoding/json"
	"math/big"
	"net/http"
)

type RequestHandler struct {
	WalletAddress string `json:"address"`
}

func HandleERC20Balance(svc service.TokenService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestHandler
		if r.Body == nil {
			http.Error(w, "Request body is empty", http.StatusBadRequest)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		resp, err := svc.GetERC20Details(req.WalletAddress)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(resp)
	}
}

func HandleERC721Balance(svc service.NFTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestHandler
		if r.Body == nil {
			http.Error(w, "Request body is empty", http.StatusBadRequest)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		resp, err := svc.GetERC721Details(req.WalletAddress)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(resp)
	}
}

func HandleERC1155Balance(svc service.NFTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RequestHandler
		if r.Body == nil {
			http.Error(w, "Request body is empty", http.StatusBadRequest)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		resp, err := svc.GetERC1155Details(req.WalletAddress)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(resp)
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
