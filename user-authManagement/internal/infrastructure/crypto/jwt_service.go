package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/lindritprekaj/user-authmanagement/internal/application/ports"
)

// normalizePEM unescapes literal "\n" sequences to real newlines.
// PEM keys delivered via env vars (docker-compose env_file, ACA secrets,
// GitHub Actions secrets) are commonly stored as a single line with
// "\n" placeholders; this lets the parser accept both forms.
func normalizePEM(s string) string {
	return strings.ReplaceAll(s, `\n`, "\n")
}

// JWTService implements ports.TokenService using RS256.
type JWTService struct {
	priv       *rsa.PrivateKey
	pub        *rsa.PublicKey
	issuer     string
	accessTTL  time.Duration
}

func NewJWTService(privPEM, pubPEM, issuer string, accessTTL time.Duration) (*JWTService, error) {
	priv, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(normalizePEM(privPEM)))
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}
	pub, err := jwt.ParseRSAPublicKeyFromPEM([]byte(normalizePEM(pubPEM)))
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}
	return &JWTService{priv: priv, pub: pub, issuer: issuer, accessTTL: accessTTL}, nil
}

type accessClaims struct {
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}

func (s *JWTService) IssueAccessToken(userID string, roles []string) (string, time.Time, error) {
	now := time.Now().UTC()
	exp := now.Add(s.accessTTL)
	claims := accessClaims{
		Roles: roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := t.SignedString(s.priv)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign jwt: %w", err)
	}
	return signed, exp, nil
}

func (s *JWTService) ParseAccessToken(token string) (*ports.AccessClaims, error) {
	parsed, err := jwt.ParseWithClaims(token, &accessClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.pub, nil
	}, jwt.WithIssuer(s.issuer), jwt.WithExpirationRequired())
	if err != nil {
		return nil, err
	}
	claims, ok := parsed.Claims.(*accessClaims)
	if !ok || !parsed.Valid {
		return nil, errors.New("invalid token")
	}
	return &ports.AccessClaims{
		UserID:    claims.Subject,
		Roles:     claims.Roles,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}

// NewRefreshToken returns 32 bytes of CSPRNG randomness, base64url-encoded.
func (s *JWTService) NewRefreshToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("read random: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// HashRefreshToken returns a hex-encoded SHA-256 of the plaintext.
// We use SHA-256 (not bcrypt) because the input is already high-entropy.
func (s *JWTService) HashRefreshToken(plaintext string) string {
	sum := sha256.Sum256([]byte(plaintext))
	return hex.EncodeToString(sum[:])
}
