package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sebastian-nunez/golang-store-api/service/auth"
	"github.com/sebastian-nunez/golang-store-api/service/cart"
	"github.com/sebastian-nunez/golang-store-api/service/order"
	"github.com/sebastian-nunez/golang-store-api/service/product"
	"github.com/sebastian-nunez/golang-store-api/service/user"
)

type Server struct {
	addr string
	db   *sql.DB
}

func NewServer(addr string, db *sql.DB) *Server {
	return &Server{
		addr: addr,
		db:   db,
	}
}

func (s *Server) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	// Users
	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(
		userStore,
		auth.HashPassword,
		auth.ComparePasswords,
		auth.CreateJWTToken,
	)
	userHandler.RegisterRoutes(subrouter)

	// Products
	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore, userStore)
	productHandler.RegisterRoutes(subrouter)

	// Cart/Orders
	orderStore := order.NewStore(s.db)
	cartHandler := cart.NewHandler(productStore, orderStore, userStore)
	cartHandler.RegisterRoutes(subrouter)

	log.Println("Server: listening on port", s.addr)
	return http.ListenAndServe(s.addr, router)
}
