package lib

import (
	"os"
	"time"

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

var jwtSecretKey string

func init() {
	jwtSecretKey = os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		panic("JWT_SECRET_KEY is unset")
	}
}

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
	tokenString, err = token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (userName string, err error) {
	// verify
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err = errors.New(ErrUnexpectedSigningMethod.Error())
			return nil, err
		}
		return []byte(jwtSecretKey), nil
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
