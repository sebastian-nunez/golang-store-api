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

const errorEmail = "error@google.com"
const existingEmail = "exists@google.com"
const unhashablePassword = "unhashable password"
const correctPassword = "1234"

func TestUserService(t *testing.T) {
	mockUserStore := &mockUserStore{}
	handler := NewHandler(mockUserStore, mockHashPassword, mockComparePassword)

	t.Run("should successfully register a new user", func(t *testing.T) {
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
			t.Errorf("expected status code %d and got %d", http.StatusCreated, rr.Code)
		}
	})

	t.Run("should fail to register given an invalid user payload", func(t *testing.T) {
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
			t.Errorf("expected status code %d and got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should fail to register if unable to hash password", func(t *testing.T) {
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
			t.Errorf("expected status code %d and got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("should fail to register if unable to create user", func(t *testing.T) {
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
			t.Errorf("expected status code %d and got %d", http.StatusInternalServerError, rr.Code)
		}
	})

	t.Run("should successfully login an existing user", func(t *testing.T) {
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
			t.Errorf("expected status code %d and got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("should fail to login given an invalid payload", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/login", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/login", handler.handleLogin)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d and got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should fail to login if the email does not exists", func(t *testing.T) {

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
			t.Errorf("expected status code %d and got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should fail to login if the password is incorrect", func(t *testing.T) {
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
			t.Errorf("expected status code %d and got %d", http.StatusBadRequest, rr.Code)
		}
	})
}

type mockUserStore struct{}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	if email == existingEmail {
		return &types.User{
			FirstName: "Sebastian",
			LastName:  "Nunez",
			Email:     existingEmail,
			Password:  "hashed password",
		}, nil
	}
	return nil, fmt.Errorf("user does not exists")
}

func (m *mockUserStore) GetUserById(id int) (*types.User, error) {
	return &types.User{}, nil
}

func (m *mockUserStore) CreateUser(user types.User) error {
	if user.Email == errorEmail {
		return fmt.Errorf("unable to create user")
	}
	return nil
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
