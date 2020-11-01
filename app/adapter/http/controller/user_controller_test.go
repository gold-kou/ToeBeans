package controller

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gold-kou/ToeBeans/testing/dummy"

	testingHelper "github.com/gold-kou/ToeBeans/testing"
	"github.com/stretchr/testify/assert"
)

var successRegisterUserReq = `
{
  "user_name": "test1",
  "email": "test1@example.com",
  "password": "Password1234"
}
`
var successRegisterUserResp = `
{
  "status": 200,
  "message": "success"
}
`

var errorRegisterUserReqWithoutUserName = `
{
  "email": "test1@example.com",
  "password": "Password1234"
}
`
var errorRegisterUserRespWithoutUserName = `
{
  "status": 400,
  "message": "user_name: cannot be blank."
}
`

var errorRegisterUserReqWithoutEmail = `
{
  "user_name": "test1",
  "password": "Password1234"
}
`
var errorRegisterUserRespWithoutEmail = `
{
  "status": 400,
  "message": "email: cannot be blank."
}
`

var errorRegisterUserReqWithoutPassword = `
{
  "user_name": "test1",
  "email": "test1@example.com"
}
`
var errorRegisterUserRespWithoutPassword = `
{
  "status": 400,
  "message":"password: cannot be blank."
}
`

var errorRegisterUserReqNotEmailFormat = `
{
  "user_name": "test1",
  "email": "hoge",
  "password": "Password1234"
}
`
var errorRegisterUserRespNotEmailFormat = `
{
  "status": 400,
  "message": "email: must be a valid email address."
}
`

var errorRegisterUserUserNameShort = `
{
  "user_name": "1",
  "email": "test1@example.com",
  "password": "Password1234"
}
`
var errorRegisterUserRespUserNameShort = `
{
  "status": 400,
  "message": "user_name: the length must be between 2 and 255."
}
`

var errorRegisterUserPasswordShort = `
{
  "user_name": "test1",
  "email": "test1@example.com",
  "password": "12345"
}
`
var errorRegisterUserRespPassword = `
{
  "status": 400,
  "message": "password: at least one upper case letter / at least one lower case letter / at least one digit / at least eight characters long."
}
`

func TestRegisterUser(t *testing.T) {
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
			args:       args{reqBody: successRegisterUserReq},
			method:     http.MethodPost,
			want:       successRegisterUserResp,
			wantStatus: 200,
		},
		{
			name:       "error request body without user_name",
			args:       args{reqBody: errorRegisterUserReqWithoutUserName},
			method:     http.MethodPost,
			want:       errorRegisterUserRespWithoutUserName,
			wantStatus: 400,
		},
		{
			name:       "error request body without email",
			args:       args{reqBody: errorRegisterUserReqWithoutEmail},
			method:     http.MethodPost,
			want:       errorRegisterUserRespWithoutEmail,
			wantStatus: 400,
		},
		{
			name:       "error request body without password",
			args:       args{reqBody: errorRegisterUserReqWithoutPassword},
			method:     http.MethodPost,
			want:       errorRegisterUserRespWithoutPassword,
			wantStatus: 400,
		},
		{
			name:       "error request body not email format",
			args:       args{reqBody: errorRegisterUserReqNotEmailFormat},
			method:     http.MethodPost,
			want:       errorRegisterUserRespNotEmailFormat,
			wantStatus: 400,
		},
		{
			name:       "error request body user_name short",
			args:       args{reqBody: errorRegisterUserUserNameShort},
			method:     http.MethodPost,
			want:       errorRegisterUserRespUserNameShort,
			wantStatus: 400,
		},
		{
			name:       "error request body password short",
			args:       args{reqBody: errorRegisterUserPasswordShort},
			method:     http.MethodPost,
			want:       errorRegisterUserRespPassword,
			wantStatus: 400,
		},
		{
			name:       "not allowed method",
			args:       args{reqBody: errorRegisterUserPasswordShort},
			method:     http.MethodHead,
			want:       errorNotAllowedMethod,
			wantStatus: 405,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// init
			db := testingHelper.SetupDBTest()
			defer testingHelper.TeardownDBTest(db)
			testingHelper.SetTestTime()
			defer testingHelper.ResetTime()

			// http request
			req, err := http.NewRequest(tt.method, "", strings.NewReader(tt.args.reqBody))
			assert.NoError(t, err)
			resp := httptest.NewRecorder()

			// test target
			UserController(resp, req)
			assert.NoError(t, err)

			// db check
			if tt.wantStatus == 200 {
				users, err := testingHelper.FindAllUsers(context.Background(), db)
				assert.NoError(t, err)
				if len(users) != 1 || !assert.Equal(t, users[0].Name, dummy.User1.Name) || !assert.Equal(t, users[0].Email, dummy.User1.Email) {
					t.Errorf("want %+v, but got %+v", dummy.User1, users)
				}
			}

			// assert http
			assert.Equal(t, tt.wantStatus, resp.Code)
			respBodyByte, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)
			respBody := string(respBodyByte)
			assert.JSONEq(t, tt.want, respBody)
		})
	}
}
