package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateTestKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

func createTestJWT(privateKey *rsa.PrivateKey, claims Claims, exp time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub":                 claims.UserID,
		"email":               claims.Email,
		"permissions":         claims.Permissions,
		"jwt_version": claims.JwtVersion,
		"iat":                 time.Now().Unix(),
		"exp":                 exp.Unix(),
	})
	return token.SignedString(privateKey)
}

func startTestJWKSServer(publicKey *rsa.PublicKey) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		modulus := base64.RawURLEncoding.EncodeToString(toUnsignedBytes(publicKey.N))
		exponent := base64.RawURLEncoding.EncodeToString(toUnsignedBytes(big.NewInt(int64(publicKey.E))))

		jwks := map[string]interface{}{
			"keys": []map[string]string{
				{
					"kty": "RSA",
					"kid": modulus[:8],
					"use": "sig",
					"alg": "RS256",
					"n":   modulus,
					"e":   exponent,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(jwks)
	})

	return httptest.NewServer(handler)
}

func toUnsignedBytes(value *big.Int) []byte {
	bytes := value.Bytes()
	if len(bytes) > 0 && bytes[0] == 0 {
		return bytes[1:]
	}
	return bytes
}

func TestParseAndValidate_ValidToken(t *testing.T) {
	privateKey, publicKey, err := generateTestKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	server := startTestJWKSServer(publicKey)
	defer server.Close()

	jwks := NewJWKSCache(server.URL)

	claims := Claims{
		UserID:             "test-user-id",
		Email:              "test@example.com",
		Permissions:        []string{"products:read", "products:write"},
		JwtVersion: 1,
	}

	token, err := createTestJWT(privateKey, claims, time.Now().Add(15*time.Minute))
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	parsed, err := ParseAndValidate(token, jwks)
	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	if parsed.UserID != claims.UserID {
		t.Errorf("Expected Subject %s, got %s", claims.UserID, parsed.UserID)
	}

	if parsed.Email != claims.Email {
		t.Errorf("Expected Email %s, got %s", claims.Email, parsed.Email)
	}
	if parsed.JwtVersion != claims.JwtVersion {
		t.Errorf("Expected JwtVersion %d, got %d", claims.JwtVersion, parsed.JwtVersion)
	}
}

func TestParseAndValidate_ExpiredToken(t *testing.T) {
	privateKey, publicKey, err := generateTestKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	server := startTestJWKSServer(publicKey)
	defer server.Close()

	jwks := NewJWKSCache(server.URL)

	claims := Claims{
		UserID:             "test-user-id",
		Email:              "test@example.com",
		Permissions:        []string{"products:read"},
		JwtVersion: 1,
	}

	token, err := createTestJWT(privateKey, claims, time.Now().Add(-1*time.Hour))
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	_, err = ParseAndValidate(token, jwks)
	if err == nil {
		t.Fatal("Expected error for expired token, got nil")
	}
}

func TestParseAndValidate_InvalidSignature(t *testing.T) {
	privateKey, _, err := generateTestKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	_, publicKey2, err := generateTestKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	server := startTestJWKSServer(publicKey2)
	defer server.Close()

	jwks := NewJWKSCache(server.URL)

	claims := Claims{
		UserID:             "test-user-id",
		Email:              "test@example.com",
		Permissions:        []string{"products:read"},
		JwtVersion: 1,
	}

	token, err := createTestJWT(privateKey, claims, time.Now().Add(15*time.Minute))
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	_, err = ParseAndValidate(token, jwks)
	if err == nil {
		t.Fatal("Expected error for invalid signature, got nil")
	}
}

func TestExtractBearerToken(t *testing.T) {
	tests := []struct {
		name     string
		header   string
		expected string
	}{
		{"valid token", "Bearer abc123", "abc123"},
		{"missing header", "", ""},
		{"invalid format", "Basic abc123", ""},
		{"no token", "Bearer", ""},
		{"extra spaces", "Bearer  abc123  ", " abc123  "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tt.header != "" {
				req.Header.Set("Authorization", tt.header)
			}

			result := extractBearerToken(req)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestValidateClaims(t *testing.T) {
	tests := []struct {
		name    string
		claims  Claims
		wantErr bool
	}{
		{"valid claims", Claims{UserID: "id", Email: "email", JwtVersion: 1}, false},
		{"missing user id", Claims{Email: "email", JwtVersion: 1}, true},
		{"missing email", Claims{UserID: "id", JwtVersion: 1}, true},
		{"invalid version", Claims{UserID: "id", Email: "email", JwtVersion: 0}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateClaims(&tt.claims)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateClaims() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestJWKSCache_Caching(t *testing.T) {
	_, publicKey, err := generateTestKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	server := startTestJWKSServer(publicKey)
	defer server.Close()

	jwks := NewJWKSCache(server.URL)

	key1, err := jwks.GetPublicKey()
	if err != nil {
		t.Fatalf("Failed to get public key: %v", err)
	}

	key2, err := jwks.GetPublicKey()
	if err != nil {
		t.Fatalf("Failed to get public key: %v", err)
	}

	if key1.N.Cmp(key2.N) != 0 || key1.E != key2.E {
		t.Error("Expected same key from cache")
	}

}

func TestJWKSCache_Refresh(t *testing.T) {
	_, publicKey, err := generateTestKeyPair()
	if err != nil {
		t.Fatalf("Failed to generate key pair: %v", err)
	}

	server := startTestJWKSServer(publicKey)
	defer server.Close()

	jwks := NewJWKSCache(server.URL)

	_, err = jwks.GetPublicKey()
	if err != nil {
		t.Fatalf("Failed to get public key: %v", err)
	}

	err = jwks.Refresh()
	if err != nil {
		t.Fatalf("Failed to refresh JWKS: %v", err)
	}

	_, err = jwks.GetPublicKey()
	if err != nil {
		t.Fatalf("Failed to get public key after refresh: %v", err)
	}

}
