package middleware

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

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TestAuthMiddlewareValidatesThunderIDTokenContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	const keyID = "auth-test-key"
	t.Setenv("THUNDERID_ISSUER", "https://identity.example")
	jwkServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"keys": []map[string]any{{
			"kty": "RSA", "kid": keyID, "use": "sig", "alg": "RS256",
			"n": base64.RawURLEncoding.EncodeToString(privateKey.PublicKey.N.Bytes()),
			"e": base64.RawURLEncoding.EncodeToString(big.NewInt(int64(privateKey.PublicKey.E)).Bytes()),
		}}})
	}))
	defer jwkServer.Close()

	if err := InitJWKS(jwkServer.URL); err != nil {
		t.Fatalf("InitJWKS() error = %v", err)
	}

	now := time.Now()
	validClaims := Claims{
		Sub: "11111111-1111-1111-1111-111111111111", Email: "admin@example.com",
		Roles: StringOrSlice{"admin"},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "https://identity.example",
			Audience:  jwt.ClaimStrings{"openschool-api"},
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	tests := []struct {
		name   string
		claims Claims
		method jwt.SigningMethod
		key    any
		want   int
	}{
		{name: "valid", claims: validClaims, method: jwt.SigningMethodRS256, key: privateKey, want: http.StatusNoContent},
		{name: "wrong issuer", claims: withIssuer(validClaims, "https://attacker.example"), method: jwt.SigningMethodRS256, key: privateKey, want: http.StatusUnauthorized},
		{name: "expired", claims: withExpiry(validClaims, now.Add(-time.Minute)), method: jwt.SigningMethodRS256, key: privateKey, want: http.StatusUnauthorized},
		{name: "wrong algorithm", claims: validClaims, method: jwt.SigningMethodHS256, key: []byte("not-an-rsa-key"), want: http.StatusUnauthorized},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			token := jwt.NewWithClaims(test.method, test.claims)
			token.Header["kid"] = keyID
			signed, err := token.SignedString(test.key)
			if err != nil {
				t.Fatal(err)
			}

			router := gin.New()
			router.Use(AuthMiddleware())
			router.GET("/protected", func(c *gin.Context) {
				if c.GetString("userID") != validClaims.Sub {
					c.AbortWithStatus(http.StatusInternalServerError)
					return
				}
				c.Status(http.StatusNoContent)
			})
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/protected", nil)
			request.Header.Set("Authorization", "Bearer "+signed)
			router.ServeHTTP(recorder, request)
			if recorder.Code != test.want {
				t.Fatalf("status = %d, want %d; body=%s", recorder.Code, test.want, recorder.Body.String())
			}
		})
	}
}

func TestInitJWKSRejectsMissingSecurityConfiguration(t *testing.T) {
	t.Setenv("THUNDERID_ISSUER", "")
	if err := InitJWKS(""); err == nil {
		t.Fatal("expected missing JWKS URL to fail")
	}
	if err := InitJWKS("https://identity.example/.well-known/jwks.json"); err == nil {
		t.Fatal("expected missing issuer to fail")
	}
}

func TestAuthMiddlewareRejectsMissingAndMalformedHeaders(t *testing.T) {
	tests := []struct {
		name   string
		header string
	}{
		{name: "missing"},
		{name: "malformed", header: "Token abc"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := gin.New()
			router.Use(AuthMiddleware())
			router.GET("/protected", func(c *gin.Context) { c.Status(http.StatusNoContent) })
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if test.header != "" {
				request.Header.Set("Authorization", test.header)
			}
			router.ServeHTTP(recorder, request)
			if recorder.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
			}
		})
	}
}

func withIssuer(claims Claims, issuer string) Claims {
	claims.Issuer = issuer
	return claims
}

func withExpiry(claims Claims, expiry time.Time) Claims {
	claims.ExpiresAt = jwt.NewNumericDate(expiry)
	return claims
}
