package middleware

import (
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gold-kou/ToeBeans/app/adapter/http/helper"
)

type Auth struct {
	UserName string
}

// GenerateToken generate and returns JWT
func GenerateToken(userName string) (tokenString string, err error) {
	// header
	token := jwt.New(jwt.SigningMethodHS256)

	// claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = userName
	claims["iat"] = time.Now()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// sign
	tokenString, err = token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerifyToken verify token and return Auth
func VerifyToken(tokenString string) (*Auth, error) {
	// verify
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err := helper.NewAuthorizationError("unexpected signing method")
			return nil, err
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	// check the result
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, helper.NewAuthorizationError("token is expired")
			}
			return nil, helper.NewAuthorizationError("token is invalid")
		}
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, helper.NewAuthorizationError("not found claims in token")
	}
	userName, ok := claims["name"].(string)
	if !ok {
		return nil, helper.NewAuthorizationError("not found name in claims")
	}

	return &Auth{
		UserName: userName,
	}, nil
}
