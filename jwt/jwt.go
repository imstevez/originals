package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/micro/go-log"
)

type UserClaims struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

type EmailClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

const (
	_ = 100 * iota
	ErrTokenInvalid
	ErrTokenExpired
)

func CreateToken(claims jwt.Claims, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Logf("create token failed: %v\n", err)
		return "", err
	}
	return tokenStr, nil
}

func ParseToken(tokenStr string, secretKey string, claims jwt.Claims) uint {
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		valiErr := err.(*jwt.ValidationError)
		if valiErr.Errors == jwt.ValidationErrorExpired {
			return ErrTokenExpired
		}
		return ErrTokenInvalid
	}
	if token == nil || !token.Valid {
		return ErrTokenInvalid
	}
	return 0
}
