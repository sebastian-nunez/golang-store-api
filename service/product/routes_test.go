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

	t.Run("should successfully fetch a product given valid id", func(t *testing.T) {
		mockProductStore := &mockProductStore{}
		handler := NewHandler(mockProductStore)

		req, err := http.NewRequest(http.MethodPost, "/products/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/products/{id}", handler.handleGetProductById)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected %d and got %d", http.StatusOK, rr.Code)
		}
	})

	t.Run("should fail to fetch a product given invalid", func(t *testing.T) {
		mockProductStore := &mockProductStore{}
		handler := NewHandler(mockProductStore)

		req, err := http.NewRequest(http.MethodPost, "/products/invalid", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/products/{id}", handler.handleGetProductById)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected %d and got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should return an error when fetching a product that does not exist", func(t *testing.T) {
		mockProductStore := &mockProductStore{err: fmt.Errorf("product does not exist")}
		handler := NewHandler(mockProductStore)

		req, err := http.NewRequest(http.MethodPost, "/products/1", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/products/{id}", handler.handleGetProductById)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected %d and got %d", http.StatusNotFound, rr.Code)
		}
	})
}

type mockProductStore struct {
	err error
}

func (s *mockProductStore) GetProducts() ([]types.Product, error) {
	return nil, s.err
}

func (s *mockProductStore) GetProductById(id int) (*types.Product, error) {
	return nil, s.err
}
