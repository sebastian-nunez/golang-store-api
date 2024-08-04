package product

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sebastian-nunez/golang-store-api/types"
)

func TestProductService(t *testing.T) {
	t.Run("should successfully fetch all products", func(t *testing.T) {
		mockProductStore := &mockProductStore{}
		handler := NewHandler(mockProductStore)

		req, err := http.NewRequest(http.MethodPost, "/products", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/products", handler.handleGetProducts)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected %d and got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("should fail to fetch all products if there was a database error", func(t *testing.T) {
		mockProductStore := &mockProductStore{err: fmt.Errorf("internal DB error")}
		handler := NewHandler(mockProductStore)

		req, err := http.NewRequest(http.MethodPost, "/products", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/products", handler.handleGetProducts)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected %d and got %d", http.StatusInternalServerError, rr.Code)
		}
	})
}

type mockProductStore struct {
	err error
}

func (s *mockProductStore) GetProducts() ([]types.Product, error) {
	return nil, s.err
}
