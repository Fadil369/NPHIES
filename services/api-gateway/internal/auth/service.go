package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Service handles authentication operations
type Service struct {
	secretKey  string
	expiration time.Duration
}

// Claims represents JWT claims
type Claims struct {
	UserID string   `json:"user_id"`
	Scopes []string `json:"scopes"`
	jwt.RegisteredClaims
}

// NewService creates a new auth service
func NewService(secretKey string, expiration time.Duration) *Service {
	return &Service{
		secretKey:  secretKey,
		expiration: expiration,
	}
}

// GenerateToken generates a JWT token for a user
func (s *Service) GenerateToken(userID string, scopes []string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Scopes: scopes,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.expiration)),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "nphies-api-gateway",
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secretKey))
}

// ValidateToken validates a JWT token and returns the claims
func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken creates a new token from a valid existing token
func (s *Service) RefreshToken(tokenString string) (string, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Generate new token with same user and scopes
	return s.GenerateToken(claims.UserID, claims.Scopes)
}

// HasScope checks if the user has a specific scope
func (c *Claims) HasScope(scope string) bool {
	for _, s := range c.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}

// IsAdmin checks if the user has admin privileges
func (c *Claims) IsAdmin() bool {
	return c.HasScope("admin")
}