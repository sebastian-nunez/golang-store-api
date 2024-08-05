package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/sebastian-nunez/golang-store-api/config"
	"github.com/sebastian-nunez/golang-store-api/types"
	"github.com/sebastian-nunez/golang-store-api/utils"
)

type Handler struct {
	store           types.UserStore
	hashPassword    func(password string) (string, error)
	comparePassword func(hashed string, plain string) bool
	createJwtToken  func(secret []byte, userId int) (string, error)
}

func NewHandler(
	store types.UserStore,
	hashPassword func(password string) (string, error),
	comparePassword func(hashed string, plain string) bool,
	createJwtToken func(secret []byte, userId int) (string, error),
) *Handler {
	return &Handler{
		store:           store,
		hashPassword:    hashPassword,
		comparePassword: comparePassword,
		createJwtToken:  createJwtToken,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods(http.MethodPost)
	router.HandleFunc("/register", h.handleRegister).Methods(http.MethodPost)

	// Admin only routes.
	// TODO(sebastian-nunez): add JWT auth
	router.HandleFunc("/users", h.handleGetUsers).Methods(http.MethodGet)
	router.HandleFunc("/users/{id}", h.handleGetUserById).Methods(http.MethodGet)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserRequest
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
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid email or password"))
		return
	}

	if !h.comparePassword(user.Password, payload.Password) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid email or password"))
		return
	}

	secret := []byte(config.Envs.JwtSecret)
	jwtToken, err := h.createJwtToken(secret, user.Id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, map[string]string{"token": jwtToken})
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

	id, err := h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusCreated, map[string]int{"id": id})
}

func (h *Handler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.store.GetUsers()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, users)
}

func (h *Handler) handleGetUserById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strId, ok := vars["id"]
	if !ok {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("missing user id"))
		return
	}

	id, err := strconv.Atoi(strId)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	user, err := h.store.GetUserById(id)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJson(w, http.StatusOK, user)
}
