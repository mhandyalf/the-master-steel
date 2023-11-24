// app/utils/jwt.go
package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenerateJWT generates a new JWT token
func GenerateJWT(claims map[string]interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	token.Claims = jwt.MapClaims{
		"username": claims["username"],
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	// Sign token
	return token.SignedString([]byte("your-secret-key"))
}

// ParseJWT parses a JWT token and returns the claims
func ParseJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("your-secret-key"), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
