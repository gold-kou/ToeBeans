package controller

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gold-kou/ToeBeans/app/domain/repository"
	"github.com/gold-kou/ToeBeans/app/lib"
	"github.com/gold-kou/ToeBeans/testing/dummy"

	"github.com/gold-kou/ToeBeans/app/adapter/http/helper"
	testingHelper "github.com/gold-kou/ToeBeans/testing"
	"github.com/stretchr/testify/assert"
)

var successReqPasswordReset = `
{
  "user_name": "testUser1",
  "password": "NewPassword1234",
  "password_reset_key": "fec668ed-f69e-45cd-87b9-a76f27759134"
}
`
var errReqPasswordResetWithoutUserName = `
{
  "password": "Password5678",
  "password_reset_key": "fec668ed-f69e-45cd-87b9-a76f27759134"
}
`
var errReqPasswordResetWithoutPassword = `
{
  "user_name": "testUser1",
  "password_reset_key": "fec668ed-f69e-45cd-87b9-a76f27759134"
}
`
var errReqPasswordResetWithoutPasswordResetKey = `
{
  "user_name": "testUser1",
  "password": "Password5678"
}
`
var errReqPasswordResetUserNameShort = `
{
  "user_name": "1",
  "password": "Password1234",
  "password_reset_key": "fec668ed-f69e-45cd-87b9-a76f27759134"
}
`
var errReqPasswordResetNameNotAlphanumeric = `
{
  "user_name": "test_1",
  "password": "Password1234",
  "password_reset_key": "fec668ed-f69e-45cd-87b9-a76f27759134"
}
`
var errReqPasswordResetPasswordRuleBreak = `
{
  "user_name": "testUser1",
  "password": "12345",
  "password_reset_key": "fec668ed-f69e-45cd-87b9-a76f27759134"
}
`
var errReqPasswordResetKeyNotUUID = `
{
  "user_name": "testUser1",
  "password": "Password5678",
  "password_reset_key": "a"
}
`
var errReqPasswordResetUserNotExists = `
{
  "user_name": "testUser2",
  "password": "Password5678",
  "password_reset_key": "fec668ed-f69e-45cd-87b9-a76f27759134"
}
`
var errReqPasswordResetWrongResetKey = `
{
  "user_name": "testUser1",
  "password": "NewPassword1234",
  "password_reset_key": "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
}
`

var errRespPasswordResetWithoutUserName = `
{
  "status": 400,
  "message":"user_name: cannot be blank."
}
`
var errRespPasswordResetWithoutPassword = `
{
  "status": 400,
  "message":"password: cannot be blank."
}
`
var errRespPasswordResetWithoutPasswordResetKey = `
{
  "status": 400,
  "message":"password_reset_key: cannot be blank."
}
`
var errRespPasswordResetUserNameShort = `
{
  "status": 400,
  "message": "user_name: the length must be between 2 and 255."
}
`
var errRespPasswordResetUserNameNotAlphanumeric = `
{
  "status": 400,
  "message": "user_name: must contain English letters and digits only."
}
`
var errorRespPasswordResetPasswordRuleBreak = `
{
  "status": 400,
  "message": "password: at least one upper case letter / at least one lower case letter / at least one digit / at least eight characters long."
}
`
var errRespPasswordResetKeyNotUUID = `
{
  "status": 400,
  "message": "password_reset_key: must be a valid UUID."
}
`
var errRespPasswordResetNotExists = `
{
  "status": 400,
  "message": "the user name doesn't exist or the password reset key doesn't exist or the password reset key is expired"
}
`

func TestPasswordResetController(t *testing.T) {
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
			args:       args{reqBody: successReqPasswordReset},
			method:     http.MethodPost,
			want:       testingHelper.RespSimpleSuccess,
			wantStatus: http.StatusOK,
		},
		{
			name:       "error request body without user_name",
			args:       args{reqBody: errReqPasswordResetWithoutUserName},
			method:     http.MethodPost,
			want:       errRespPasswordResetWithoutUserName,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error request body without password",
			args:       args{reqBody: errReqPasswordResetWithoutPassword},
			method:     http.MethodPost,
			want:       errRespPasswordResetWithoutPassword,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error request body without password_reset_key",
			args:       args{reqBody: errReqPasswordResetWithoutPasswordResetKey},
			method:     http.MethodPost,
			want:       errRespPasswordResetWithoutPasswordResetKey,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error request body user_name short",
			args:       args{reqBody: errReqPasswordResetUserNameShort},
			method:     http.MethodPost,
			want:       errRespPasswordResetUserNameShort,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error request body user_name is not alphanumeric",
			args:       args{reqBody: errReqPasswordResetNameNotAlphanumeric},
			method:     http.MethodPost,
			want:       errRespPasswordResetUserNameNotAlphanumeric,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error request body password short",
			args:       args{reqBody: errReqPasswordResetPasswordRuleBreak},
			method:     http.MethodPost,
			want:       errorRespPasswordResetPasswordRuleBreak,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error request body password_reset_key not uuid",
			args:       args{reqBody: errReqPasswordResetKeyNotUUID},
			method:     http.MethodPost,
			want:       errRespPasswordResetKeyNotUUID,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error request body user not exists",
			args:       args{reqBody: errReqPasswordResetUserNotExists},
			method:     http.MethodPost,
			want:       errRespPasswordResetNotExists,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error request body wrong password_reset_key",
			args:       args{reqBody: errReqPasswordResetWrongResetKey},
			method:     http.MethodPost,
			want:       errRespPasswordResetNotExists,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error password_reset_key is expired",
			args:       args{reqBody: successReqPasswordReset},
			method:     http.MethodPost,
			want:       errRespPasswordResetNotExists,
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
			err := userRepo.Create(context.Background(), &dummy.User1)
			assert.NoError(t, err)
			if tt.name == "error password_reset_key is expired" {
				err = testingHelper.UpdateResetKeyExpiresAt(db, dummy.User1.PasswordResetKey, lib.NowFunc().Add(-time.Second))
			} else {
				err = testingHelper.UpdateResetKeyExpiresAt(db, dummy.User1.PasswordResetKey, lib.NowFunc())
			}
			assert.NoError(t, err)

			// http request
			req, err := http.NewRequest(tt.method, "/password-reset", strings.NewReader(tt.args.reqBody))
			assert.NoError(t, err)
			req.Header.Add(helper.HeaderKeyContentType, "application/json")
			resp := httptest.NewRecorder()

			// test target
			PasswordResetController(resp, req)
			assert.NoError(t, err)

			// assert http
			assert.Equal(t, tt.wantStatus, resp.Code)
			respBodyByte, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)
			respBody := string(respBodyByte)
			assert.JSONEq(t, tt.want, respBody)
		})
	}
}
