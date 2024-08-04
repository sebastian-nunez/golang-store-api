package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sebastian-nunez/golang-store-api/service/auth"
	"github.com/sebastian-nunez/golang-store-api/service/product"
	"github.com/sebastian-nunez/golang-store-api/service/user"
)

type ApiServer struct {
	addr string
	db   *sql.DB
}

func NewApiServer(addr string, db *sql.DB) *ApiServer {
	return &ApiServer{
		addr: addr,
		db:   db,
	}
}

func (s *ApiServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	// Users
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(
		userStore,
		auth.HashPassword,
		auth.ComparePasswords,
		auth.CreateJwtToken,
	)
	userHandler.RegisterRoutes(subrouter)

	// Products
	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore)
	productHandler.RegisterRoutes(subrouter)

	log.Println("Server: listening on port", s.addr)
	return http.ListenAndServe(s.addr, router)
}
