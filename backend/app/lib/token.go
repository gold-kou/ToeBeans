package lib

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

var ErrNotFoundAuthorizationHeader = errors.New("not found Authorization header")
var ErrNotFoundBearerToken = errors.New("not found bearer token")
var ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
var ErrTokenExpired = errors.New("token is expired")
var ErrTokenInvalid = errors.New("token is invalid")
var ErrTokenInvalidNotExistingUserName = errors.New("the user name contained in token doesn't exist")
var ErrNotFoundClaims = errors.New("not found claims in token")
var ErrNotFoundName = errors.New("not found name in token")

const GuestUserName = "guest"
const TokenExpirationHour = 24

func GenerateToken(userName string) (tokenString string, err error) {
	// header
	token := jwt.New(jwt.SigningMethodHS256)

	// claims
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "ToeBeans"
	claims["name"] = userName
	claims["iat"] = time.Now()
	claims["exp"] = time.Now().Add(time.Hour * TokenExpirationHour).Unix()

	// generate token by secret key
	tokenString, err = token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyHeaderToken(r *http.Request) (userName string, err error) {
	// get jwt from header
	authHeader := r.Header.Get(helper.HeaderKeyAuthorization)
	if authHeader == "" {
		return "", ErrNotFoundAuthorizationHeader
	}

	// Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODM3MjIwNTMsImlhdCI6IjIwMjAtMDMtMDhUMTE6NDc6MzMuMTc4NjU5MyswOTowMCIsIm5hbWUiOiJ0ZXN0In0.YIyT1RJGcYbdynx1V4-6MhiosmTlHmKiyiG_GjxQeuw
	s := strings.Split(authHeader, "Bearer ")
	if len(s) != 2 {
		return "", ErrNotFoundBearerToken
	}
	bearerToken := s[1]

	// verify jwt
	userName, err = VerifyToken(bearerToken)
	if err != nil {
		return
	}
	return
}

func VerifyToken(tokenString string) (userName string, err error) {
	// verify
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err = errors.New(ErrUnexpectedSigningMethod.Error())
			return nil, err
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	// check the result
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return "", errors.New(ErrTokenExpired.Error())
			}
			return "", errors.New(ErrTokenInvalid.Error())
		}
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New(ErrNotFoundClaims.Error())
	}
	userName, ok = claims["name"].(string)
	if !ok {
		return "", errors.New(ErrNotFoundName.Error())
	}

	return
}
