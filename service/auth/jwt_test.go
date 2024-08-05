package auth

import "testing"

func TestCreateJWTToken(t *testing.T) {
	t.Run("should return a valid JWT token", func(t *testing.T) {
		secret := []byte("some secret")
		userId := 1234

		token, err := CreateJWTToken(secret, userId)
		if err != nil {
			t.Errorf("expected token and got error %s", err)
		}

		if token == "" {
			t.Errorf("expected a valid token and got an empty string")
		}
	})
}
