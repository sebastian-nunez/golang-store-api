package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sebastian-nunez/golang-store-api/types"
)

const (
	errorEmail         = "error@google.com"
	existingEmail      = "exists@google.com"
	unhashablePassword = "unhashable password"
	correctPassword    = "1234"
	badJwtEmail        = "badjwt@google.com"
	badUserId          = 999
)

func TestUserService(t *testing.T) {
	t.Parallel()

	t.Run("should successfully register a new user", func(t *testing.T) {
		mockUserStore := &mockUserStore{}
		handler := NewHandler(
			mockUserStore,
			mockHashPassword,
			mockComparePassword,
			mockCreateJWTToken,
		)

		payload := types.RegisterUserRequest{
			FirstName: "Sebastian",
			LastName:  "Nunez",
			Email:     "snunez@gmail.com", // new email
			Password:  "1234",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("want status code %d and got %d", http.StatusCreated, rr.Code)
		}
	})

	t.Run("should fail to register given an invalid user payload", func(t *testing.T) {
		mockUserStore := &mockUserStore{}
		handler := NewHandler(
			mockUserStore,
			mockHashPassword,
			mockComparePassword,
			mockCreateJWTToken,
		)

		invalidEmail := "invalid"
		payload := types.RegisterUserRequest{
			FirstName: "Sebastian",
			LastName:  "Nunez",
			Email:     invalidEmail,
			Password:  "1234",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("want status code %d and got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should fail to register if unable to hash password", func(t *testing.T) {
		mockUserStore := &mockUserStore{}
		handler := NewHandler(
			mockUserStore,
			mockHashPassword,
			mockComparePassword,
			mockCreateJWTToken,
		)

		payload := types.RegisterUserRequest{
			FirstName: "Sebastian",
			LastName:  "Nunez",
			Email:     "snunez@google.com",
			Password:  unhashablePassword,
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("want status code %d and got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("should fail to register if unable to create user", func(t *testing.T) {
		mockUserStore := &mockUserStore{}
		handler := NewHandler(
			mockUserStore,
			mockHashPassword,
			mockComparePassword,
			mockCreateJWTToken,
		)

		payload := types.RegisterUserRequest{
			FirstName: "Sebastian",
			LastName:  "Nunez",
			Email:     errorEmail,
			Password:  "1234",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("want status code %d and got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("should successfully login an existing user", func(t *testing.T) {
		mockUserStore := &mockUserStore{}
		handler := NewHandler(
			mockUserStore,
			mockHashPassword,
			mockComparePassword,
			mockCreateJWTToken,
		)

		payload := types.LoginUserRequest{
			Email:    existingEmail,
			Password: correctPassword,
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/login", handler.handleLogin)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("want status code %d and got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("should fail to login given an invalid payload", func(t *testing.T) {
		mockUserStore := &mockUserStore{}
		handler := NewHandler(
			mockUserStore,
			mockHashPassword,
			mockComparePassword,
			mockCreateJWTToken,
		)

		req, err := http.NewRequest(http.MethodPost, "/login", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/login", handler.handleLogin)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("want status code %d and got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should fail to login if the email does not exists", func(t *testing.T) {
		mockUserStore := &mockUserStore{}
		handler := NewHandler(
			mockUserStore,
			mockHashPassword,
			mockComparePassword,
			mockCreateJWTToken,
		)

		payload := types.LoginUserRequest{
			Email:    "does not exist",
			Password: "some password",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/login", handler.handleLogin)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("want status code %d and got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should fail to login if the password is incorrect", func(t *testing.T) {
		mockUserStore := &mockUserStore{}
		handler := NewHandler(
			mockUserStore,
			mockHashPassword,
			mockComparePassword,
			mockCreateJWTToken,
		)

		payload := types.LoginUserRequest{
			Email:    existingEmail,
			Password: "incorrect password",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/login", handler.handleLogin)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("want status code %d and got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should fail to login if unable to create JWT token", func(t *testing.T) {
		mockUserStore := &mockUserStore{}
		handler := NewHandler(
			mockUserStore,
			mockHashPassword,
			mockComparePassword,
			mockCreateJWTToken,
		)

		payload := types.LoginUserRequest{
			Email:    badJwtEmail,
			Password: correctPassword,
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/login", handler.handleLogin)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("want status code %d and got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("should successfully fetch all users", func(t *testing.T) {
		mockUserStore := &mockUserStore{}
		handler := NewHandler(
			mockUserStore,
			mockHashPassword,
			mockComparePassword,
			mockCreateJWTToken,
		)

		req, err := http.NewRequest(http.MethodGet, "/users", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/users", handler.handleGetUsers)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("want status code %d and got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("should fail to fetch all users given internal DB error", func(t *testing.T) {
		mockUserStore := &mockUserStore{err: fmt.Errorf("internal DB error")}
		handler := NewHandler(
			mockUserStore,
			mockHashPassword,
			mockComparePassword,
			mockCreateJWTToken,
		)

		req, err := http.NewRequest(http.MethodGet, "/users", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/users", handler.handleGetUsers)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("want status code %d and got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("should successfully fetch a user given a valid id", func(t *testing.T) {
		mockUserStore := &mockUserStore{}
		handler := NewHandler(
			mockUserStore,
			mockHashPassword,
			mockComparePassword,
			mockCreateJWTToken,
		)

		req, err := http.NewRequest(http.MethodGet, "/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/users/{id}", handler.handleGetUserById)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("want status code %d and got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("should fail to fetch a user given a invalid id", func(t *testing.T) {
		mockUserStore := &mockUserStore{}
		handler := NewHandler(
			mockUserStore,
			mockHashPassword,
			mockComparePassword,
			mockCreateJWTToken,
		)

		req, err := http.NewRequest(http.MethodGet, "/users/invalid", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/users/{id}", handler.handleGetUserById)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("want status code %d and got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("should fail to fetch a user given an internal DB error", func(t *testing.T) {
		mockUserStore := &mockUserStore{err: fmt.Errorf("internal DB error")}
		handler := NewHandler(
			mockUserStore,
			mockHashPassword,
			mockComparePassword,
			mockCreateJWTToken,
		)

		req, err := http.NewRequest(http.MethodGet, "/users/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/users/{id}", handler.handleGetUserById)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("want status code %d and got %d", http.StatusInternalServerError, rr.Code)
		}
	})
}

type mockUserStore struct {
	err error
}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	if email == existingEmail {
		return &types.User{
			ID:        1,
			FirstName: "Sebastian",
			LastName:  "Nunez",
			Email:     existingEmail,
			Password:  "hashed password",
		}, nil
	}
	if email == badJwtEmail {
		return &types.User{
			ID:        badUserId,
			FirstName: "Sebastian",
			LastName:  "Nunez",
			Email:     existingEmail,
			Password:  "hashed password",
		}, nil
	}
	return nil, fmt.Errorf("user does not exists")
}

func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	return &types.User{}, m.err
}

func (m *mockUserStore) CreateUser(user types.User) (int, error) {
	if user.Email == errorEmail {
		return 0, fmt.Errorf("unable to create user")
	}
	return 1, nil
}

func (m *mockUserStore) GetUsers() ([]types.User, error) {
	return nil, m.err
}

func mockHashPassword(password string) (string, error) {
	if password == unhashablePassword {
		return "", fmt.Errorf("unable to hash password")
	}
	return "hashed", nil
}

func mockComparePassword(password string, plain string) bool {
	return plain == correctPassword
}

func mockCreateJWTToken(secret []byte, userId int) (string, error) {
	if userId == badUserId {
		return "", fmt.Errorf("bad user id, unable to create token")
	}
	return "some token", nil
}
