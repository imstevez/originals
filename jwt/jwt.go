package jwt

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/micro/go-log"
)

type UserClaims struct {
	UserId   int64  `json:"user_id"`
	Email    string `json:"email"`
	Nickname string `json:"user_name"`
	Mobile   string `json:"mobile"`
	ImageUrl string `json:"image_url"`
	jwt.StandardClaims
}

type EmailClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

var (
	ErrTokenInvalid = errors.New("invalid token")
	ErrTokenExpired = errors.New("expired token")
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

func ParseToken(tokenStr string, secretKey string, claims jwt.Claims) error {
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
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
	return nil
}
