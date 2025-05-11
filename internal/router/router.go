package router

import (
	"net/http"
	"tokenhub-api/internal/handlers"
	"tokenhub-api/internal/middleware"
	"tokenhub-api/internal/services"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewRouter(tokenSvc services.TokenService, nftSvc services.NFTService, logger *zap.Logger) http.Handler {
	r := mux.NewRouter()
	r.Use(middleware.NewZapLoggerMiddleware(logger))

	api := r.PathPrefix("/api").Subrouter()

	balance := api.PathPrefix("/balance").Subrouter()
	balance.HandleFunc("/erc20", handlers.HandleERC20Balance(tokenSvc)).Methods("GET")
	balance.HandleFunc("/erc721", handlers.HandleERC721Balance(nftSvc)).Methods("GET")
	balance.HandleFunc("/erc1155", handlers.HandleERC1155Balance(nftSvc)).Methods("GET")

	deploy := api.PathPrefix("/deploy").Subrouter()
	deploy.HandleFunc("/erc20", handlers.DeployERC20Handler(tokenSvc)).Methods("POST")
	deploy.HandleFunc("/erc721", handlers.DeployERC721Handler(nftSvc)).Methods("POST")
	deploy.HandleFunc("/erc1155", handlers.DeployERC1155Handler(nftSvc)).Methods("POST")

	mint := api.PathPrefix("/mint").Subrouter()
	mint.HandleFunc("/erc20", handlers.MintERC20Handler(tokenSvc)).Methods("POST")
	mint.HandleFunc("/erc721", handlers.MintERC721Handler(nftSvc)).Methods("POST")
	mint.HandleFunc("/erc1155", handlers.MintERC1155Handler(nftSvc)).Methods("POST")

	burn := api.PathPrefix("/burn").Subrouter()
	burn.HandleFunc("/erc20", handlers.BurnERC20Handler(tokenSvc)).Methods("POST")
	burn.HandleFunc("/erc721", handlers.BurnERC721Handler(nftSvc)).Methods("POST")
	burn.HandleFunc("/erc1155", handlers.BurnERC1155Handler(nftSvc)).Methods("POST")

	return r
}
