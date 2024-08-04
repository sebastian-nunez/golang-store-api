package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/sebastian-nunez/golang-store-api/types"
)

func TestParseJson(t *testing.T) {
	t.Run("should return error given no payload", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/some-endpoint", bytes.NewBuffer(nil))

		err := ParseJson(req, nil)

		if err == nil {
			t.Errorf("expected error and got %s", err)
		}
	})

	t.Run("should successfully parse a valid payload", func(t *testing.T) {
		payload := types.RegisterUserRequest{
			FirstName: "Sebastian",
			LastName:  "Nunez",
			Email:     "snunez@google.com",
			Password:  "1234",
		}
		marshalled, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPost, "/some-endpoint", bytes.NewBuffer(marshalled))
		expected := types.RegisterUserRequest{
			FirstName: "Sebastian",
			LastName:  "Nunez",
			Email:     "snunez@google.com",
			Password:  "1234",
		}

		err := ParseJson(req, &payload)
		if err != nil {
			t.Errorf("expected no error, but got %s", err)
		}

		if !reflect.DeepEqual(expected, payload) {
			t.Errorf("expected %+v and got %+v", expected, payload)
		}
	})
}

func TestWriteJson(t *testing.T) {
	type response struct {
		Message string `json:"message"`
	}

	tests := []struct {
		name     string
		status   HttpStatus
		payload  any
		expected string
	}{
		{
			name:   "success response",
			status: http.StatusOK,
			payload: response{
				Message: "Success",
			},
			expected: `{"message":"Success"}`,
		},
		{
			name:   "error response",
			status: http.StatusInternalServerError,
			payload: response{
				Message: "Internal Server Error",
			},
			expected: `{"message":"Internal Server Error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			err := WriteJson(rr, tt.status, tt.payload)
			if err != nil {
				t.Errorf("expected no error, but got %v", err)
			}

			if rr.Code != int(tt.status) {
				t.Errorf("expected status %v, but got %v", tt.status, rr.Code)
			}

			if rr.Header().Get("Content-Type") != "application/json" {
				t.Errorf("expected Content-Type application/json, but got %v", rr.Header().Get("Content-Type"))
			}

			actual := rr.Body.String()
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("expected %+v, but got %+v", tt.expected, actual)
			}
		})
	}
}
