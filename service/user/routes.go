package user

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sebastian-nunez/golang-store-api/types"
	"github.com/sebastian-nunez/golang-store-api/utils"
)

type Handler struct {
	store        types.UserStore
	hashPassword func(password string) (string, error)
}

func NewHandler(store types.UserStore, hashPassword func(password string) (string, error)) *Handler {
	return &Handler{store: store, hashPassword: hashPassword}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserRequest
	err := utils.ParseJson(r, &payload)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid request: %v", errors))
		return
	}

	user, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(
			w,
			http.StatusBadRequest,
			fmt.Errorf("user with email %s already exists", user.Email),
		)
		return
	}

	hashedPassword, err := h.hashPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusCreated, nil)
}
