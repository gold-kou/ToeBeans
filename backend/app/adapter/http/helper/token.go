package helper

import (
	"os"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

var errUnexpectedSigningMethod = errors.New("unexpected signing method")
var errTokenExpired = errors.New("token is expired")
var errTokenInvalid = errors.New("token is invalid")
var errNotFoundClaims = errors.New("not found claims in token")
var errNotFoundSub = errors.New("not found sub in token")
var errNotFoundName = errors.New("not found name in token")

const GuestUserName = "guest"
const TokenExpirationHour = 24

var jwtSecretKey string

func init() {
	jwtSecretKey = os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		panic("JWT_SECRET_KEY is unset")
	}
}

func GenerateToken(userID int64, userName string) (tokenString string, err error) {
	// header
	token := jwt.New(jwt.SigningMethodHS256)

	// claims
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "ToeBeans"
	claims["sub"] = strconv.Itoa(int(userID))
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

func VerifyToken(tokenString string) (userID int64, userName string, err error) {
	// verify
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// check signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err = errUnexpectedSigningMethod
			return nil, err
		}
		return []byte(jwtSecretKey), nil
	})

	// check the result
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				err = errTokenExpired
				return
			}
			err = errTokenInvalid
			return
		}
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		err = errNotFoundClaims
		return
	}
	userIDStr, ok := claims["sub"].(string)
	if !ok {
		err = errNotFoundSub
		return
	}
	userIDInt, err := strconv.Atoi(userIDStr)
	if err != nil {
		return
	}
	userID = int64(userIDInt)
	userName, ok = claims["name"].(string)
	if !ok {
		err = errNotFoundName
		return
	}

	return
}
