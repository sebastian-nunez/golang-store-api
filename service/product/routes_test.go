package product

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

func TestProductService(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		method     string
		endpoint   string
		payload    interface{}
		mockErr    error
		wantStatus int
	}{
		{
			name:       "should successfully fetch all products",
			method:     http.MethodGet,
			endpoint:   "/products",
			wantStatus: http.StatusOK,
		},
		{
			name:       "should fail to fetch all products if there was a database error",
			method:     http.MethodGet,
			endpoint:   "/products",
			mockErr:    fmt.Errorf("internal DB error"),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "should successfully fetch a product given valid id",
			method:     http.MethodGet,
			endpoint:   "/products/1",
			wantStatus: http.StatusOK,
		},
		{
			name:       "should fail to fetch a product given invalid id",
			method:     http.MethodGet,
			endpoint:   "/products/invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "should return an error when fetching a product that does not exist",
			method:     http.MethodGet,
			endpoint:   "/products/1",
			mockErr:    fmt.Errorf("product does not exist"),
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "should successfully create a new product",
			method:     http.MethodPost,
			endpoint:   "/products",
			payload:    types.CreateProductRequest{Name: "Jordans", Price: 125, Quantity: 5},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "should fail to create a new product given invalid request payload",
			method:     http.MethodPost,
			endpoint:   "/products",
			payload:    types.CreateProductRequest{Price: 125, Quantity: 5},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "should return error if unable to create valid product within the database",
			method:     http.MethodPost,
			endpoint:   "/products",
			payload:    types.CreateProductRequest{Name: "Jordans", Price: 125, Quantity: 5},
			mockErr:    fmt.Errorf("internal DB error"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockProductStore := &mockProductStore{err: tc.mockErr}
			handler := NewHandler(mockProductStore)

			var bodyBytes []byte
			if tc.payload != nil {
				var err error
				bodyBytes, err = json.Marshal(tc.payload)
				if err != nil {
					t.Fatal(err)
				}
			}

			req, err := http.NewRequest(tc.method, tc.endpoint, bytes.NewBuffer(bodyBytes))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			router := mux.NewRouter()
			router.HandleFunc("/products", handler.handleGetProducts).Methods(http.MethodGet)
			router.HandleFunc("/products/{id}", handler.handleGetProductByID).Methods(http.MethodGet)
			router.HandleFunc("/products", handler.handleCreateProduct).Methods(http.MethodPost)
			router.ServeHTTP(rr, req)

			if rr.Code != tc.wantStatus {
				t.Errorf("expected status code %d and got %d", tc.wantStatus, rr.Code)
			}
		})
	}
}

type mockProductStore struct {
	err error
}

func (s *mockProductStore) GetProducts() ([]types.Product, error) {
	return nil, s.err
}

func (s *mockProductStore) GetProductByID(id int) (*types.Product, error) {
	return nil, s.err
}

func (s *mockProductStore) GetProductsByID(productIDs []int) ([]types.Product, error) {
	return nil, s.err
}

func (s *mockProductStore) CreateProduct(product types.CreateProductRequest) (int, error) {
	return 1, s.err
}

func (s *mockProductStore) UpdateProduct(product types.Product) error {
	return s.err
}
