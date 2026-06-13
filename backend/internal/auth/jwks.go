package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"
)

type jwkKey struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type jwkSet struct {
	Keys []jwkKey `json:"keys"`
}

type JWKSCache struct {
	url       string
	publicKey *rsa.PublicKey
	kid       string
	cachedAt  time.Time
	mu        sync.RWMutex
}

func NewJWKSCache(url string) *JWKSCache {
	return &JWKSCache{url: url}
}

func (c *JWKSCache) GetPublicKey() (*rsa.PublicKey, error) {
	c.mu.RLock()
	if c.publicKey != nil && time.Since(c.cachedAt) < time.Hour {
		defer c.mu.RUnlock()
		return c.publicKey, nil
	}
	c.mu.RUnlock()

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.publicKey != nil && time.Since(c.cachedAt) < time.Hour {
		return c.publicKey, nil
	}

	if err := c.fetch(); err != nil {
		return nil, err
	}
	return c.publicKey, nil
}

func (c *JWKSCache) Refresh() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.fetch()
}

func (c *JWKSCache) GetKid() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.kid
}

func (c *JWKSCache) fetch() error {
	resp, err := http.Get(c.url)
	if err != nil {
		return fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("JWKS endpoint returned status %d", resp.StatusCode)
	}

	var set jwkSet
	if err := json.NewDecoder(resp.Body).Decode(&set); err != nil {
		return fmt.Errorf("failed to decode JWKS: %w", err)
	}

	if len(set.Keys) == 0 {
		return fmt.Errorf("JWKS contains no keys")
	}

	key := set.Keys[0]
	publicKey, err := parseRSAPublicKey(key.N, key.E)
	if err != nil {
		return fmt.Errorf("failed to parse RSA public key: %w", err)
	}

	c.publicKey = publicKey
	c.kid = key.Kid
	c.cachedAt = time.Now()
	return nil
}

func parseRSAPublicKey(modulusB64, exponentB64 string) (*rsa.PublicKey, error) {
	modulusBytes, err := base64.RawURLEncoding.DecodeString(modulusB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode modulus: %w", err)
	}

	exponentBytes, err := base64.RawURLEncoding.DecodeString(exponentB64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode exponent: %w", err)
	}

	modulus := new(big.Int).SetBytes(modulusBytes)

	exponent := 0
	for _, b := range exponentBytes {
		exponent = exponent<<8 | int(b)
	}

	return &rsa.PublicKey{
		N: modulus,
		E: exponent,
	}, nil
}
