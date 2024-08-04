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
const unhashablePassword = "unhashable password"

func TestUserService(t *testing.T) {
	mockUserStore := &mockUserStore{}
	handler := NewHandler(mockUserStore, mockHashPassword)

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
}

type mockUserStore struct{}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	if email == "exists@google.com" {
		return &types.User{}, nil
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
