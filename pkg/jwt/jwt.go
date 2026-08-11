package jwt

import (
	"fmt"
	"time"

	gjwt "github.com/golang-jwt/jwt/v5"
)

var secret = []byte("e4f5788a-5313-4dde-856e-e5ee01ea456b")

type Claims struct {
	UserID int    `json:"user_id"` 
	Email  string `json:"email"`
	Role   string `json:"role"`
	gjwt.RegisteredClaims
}

func GenerateToken(userID int, email, role string) (string, error) { 
	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: gjwt.RegisteredClaims{
			ExpiresAt: gjwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := gjwt.NewWithClaims(gjwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := gjwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *gjwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*gjwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}