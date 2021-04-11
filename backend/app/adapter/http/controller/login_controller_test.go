package controller

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/gold-kou/ToeBeans/backend/app/lib"

	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"

	"github.com/gold-kou/ToeBeans/backend/testing/dummy"

	modelHttp "github.com/gold-kou/ToeBeans/backend/app/domain/model/http"
	testingHelper "github.com/gold-kou/ToeBeans/backend/testing"
	"github.com/stretchr/testify/assert"
)

var successReqLogin = `
{
  "email": "testUser1@example.com",
  "password": "Password1234"
}
`
var errReqLoginWithoutEmail = `
{
  "password": "Password1234"
}
`
var errReqLoginWithoutPassword = `
{
  "email": "testUser1@example.com"
}
`
var errReqLoginNotEmailFormat = `
{
  "email": "a",
  "password": "Password1234"
}
`
var errReqLoginWrongPassword = `
{
  "email": "testUser1@example.com",
  "password": "Password9999"
}
`
var errReqLoginNotExistingEmail = `
{
  "email": "XXXXX@example.com",
  "password": "Password9999"
}
`

var errRespLoginWithoutEmail = `
{
  "status": 400,
  "message": "email: cannot be blank."
}
`
var errRespLoginWithoutPassword = `
{
  "status": 400,
  "message": "password: cannot be blank."
}
`
var errRespLoginNotEmailFormat = `
{
  "status": 400,
  "message": "email: must be a valid email address."
}
`
var errRespLoginNotExistingEmail = `
{
  "status": 400,
  "message": "Wrong username or password"
}
`
var errRespLoginWrongPassword = `
{
  "status": 400,
  "message": "Wrong username or password"
}
`

func TestLoginController(t *testing.T) {
	type args struct {
		reqBody string
	}
	tests := []struct {
		name       string
		args       args
		method     string
		want       string
		wantStatus int
	}{
		{
			name:       "success",
			args:       args{reqBody: successReqLogin},
			method:     http.MethodPost,
			wantStatus: http.StatusOK,
		},
		{
			name:       "error empty email",
			args:       args{reqBody: errReqLoginWithoutEmail},
			method:     http.MethodPost,
			want:       errRespLoginWithoutEmail,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error empty password",
			args:       args{reqBody: errReqLoginWithoutPassword},
			method:     http.MethodPost,
			want:       errRespLoginWithoutPassword,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error not email format",
			args:       args{reqBody: errReqLoginNotEmailFormat},
			method:     http.MethodPost,
			want:       errRespLoginNotEmailFormat,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error not existing email",
			args:       args{reqBody: errReqLoginNotExistingEmail},
			method:     http.MethodPost,
			want:       errRespLoginNotExistingEmail,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error wrong password",
			args:       args{reqBody: errReqLoginWrongPassword},
			method:     http.MethodPost,
			want:       errRespLoginWrongPassword,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not allowed method",
			args:       args{},
			method:     http.MethodHead,
			want:       testingHelper.ErrNotAllowedMethod,
			wantStatus: http.StatusMethodNotAllowed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// init
			db := testingHelper.SetupDBTest()
			defer testingHelper.TeardownDBTest(db)
			testingHelper.SetTestTime()
			defer testingHelper.ResetTime()

			// insert dummy data
			userRepo := repository.NewUserRepository(db)
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dummy.User1.Password), bcrypt.DefaultCost)
			assert.NoError(t, err)
			dummy.User1.Password = string(hashedPassword)
			err = userRepo.Create(context.Background(), &dummy.User1)
			assert.NoError(t, err)

			// http request
			req, err := http.NewRequest(tt.method, "/login", strings.NewReader(tt.args.reqBody))
			assert.NoError(t, err)
			resp := httptest.NewRecorder()

			// test target
			LoginController(resp, req)

			// assert http
			respBodyByte, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)
			if tt.wantStatus == http.StatusOK {
				var respLogin modelHttp.Token
				err = json.Unmarshal(respBodyByte, &respLogin)
				assert.NoError(t, err)
				tokenUserName, err := lib.VerifyToken(respLogin.IdToken)
				assert.NoError(t, err)
				assert.Equal(t, dummy.User1.Name, tokenUserName)
			} else {
				respBody := string(respBodyByte)
				assert.Equal(t, tt.wantStatus, resp.Code)
				assert.JSONEq(t, tt.want, respBody)
			}
		})
	}
}
