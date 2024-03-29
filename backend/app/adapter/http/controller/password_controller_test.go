package controller

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	httpContext "github.com/gold-kou/ToeBeans/backend/app/adapter/http/context"

	"github.com/gold-kou/ToeBeans/backend/app/lib"

	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	testingHelper "github.com/gold-kou/ToeBeans/backend/testing"
	"github.com/gold-kou/ToeBeans/backend/testing/dummy"
	"github.com/stretchr/testify/assert"
)

var successReqPasswordChange = `
{
  "old_password": "Password1234",
  "new_password": "Password5678"
}
`
var errReqPasswordChangeNoOldPassword = `
{
  "new_password": "Password5678"
}
`
var errReqPasswordChangeNoNewPassword = `
{
  "old_password": "Password1234"
}
`
var errReqPasswordChangeInvalidPassword = `
{
  "old_password": "password",
  "new_password": "password"
}
`
var errReqPasswordChangeWrongOldPassword = `
{
  "old_password": "Hoge1234",
  "new_password": "Password5678"
}
`

var errRespPasswordChangeNoOldPassword = `
{
  "status": 400,
  "message": "old_password: cannot be blank."
}
`
var errRespPasswordChangeNoNewPassword = `
{
  "status": 400,
  "message": "new_password: cannot be blank."
}
`
var errRespPasswordChangeInvalidPassword = `
{
  "status": 400,
  "message": "new_password: Your password must be at least 8 characters long, contain at least one number and have a mixture of uppercase and lowercase letters; old_password: Your password must be at least 8 characters long, contain at least one number and have a mixture of uppercase and lowercase letters."
}
`
var errRespPasswordChangeWrongOldPassword = `
{
  "status": 400,
  "message": "not correct password"
}
`

func TestPasswordChange(t *testing.T) {
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
		//{
		//	name:       "success",
		//	args:       args{reqBody: successReqPasswordChange},
		//	method:     http.MethodPut,
		//	want:       testingHelper.RespSimpleSuccess,
		//	wantStatus: http.StatusOK,
		//},
		{
			name:       "error no old password",
			args:       args{reqBody: errReqPasswordChangeNoOldPassword},
			method:     http.MethodPut,
			want:       errRespPasswordChangeNoOldPassword,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error no new password",
			args:       args{reqBody: errReqPasswordChangeNoNewPassword},
			method:     http.MethodPut,
			want:       errRespPasswordChangeNoNewPassword,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error invalid password",
			args:       args{reqBody: errReqPasswordChangeInvalidPassword},
			method:     http.MethodPut,
			want:       errRespPasswordChangeInvalidPassword,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error old password is wrong",
			args:       args{reqBody: errReqPasswordChangeWrongOldPassword},
			method:     http.MethodPut,
			want:       errRespPasswordChangeWrongOldPassword,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error forbidden guest user",
			args:       args{reqBody: successReqPasswordChange},
			method:     http.MethodPut,
			want:       testingHelper.ErrForbidden,
			wantStatus: http.StatusForbidden,
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

			// http request
			req, err := http.NewRequest(tt.method, "/password", strings.NewReader(tt.args.reqBody))
			assert.NoError(t, err)
			if tt.name == "error forbidden guest user" {
				req = req.WithContext(httpContext.SetTokenUserName(req.Context(), helper.GuestUserName))
			} else {
				req = req.WithContext(httpContext.SetTokenUserName(req.Context(), dummy.User1.Name))
			}
			resp := httptest.NewRecorder()

			// test target
			PasswordController(resp, req)
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

var successReqPasswordResetEmail = `
{
  "email": "testUser1@example.com"
}
`
var errReqPasswordResetEmailWithoutEmail = `
{
}
`
var errReqPasswordResetEmailNotEmailFormat = `
{
  "email": "abc"
}
`
var errReqPasswordResetEmailNotExistingEmail = `
{
  "email": "testUser0@example.com"
}
`

var errRespPasswordResetEmailWithoutEmail = `
{
  "status": 400,
  "message":"email: cannot be blank."
}
`
var errRespPasswordResetEmailNotEmailFormat = `
{
  "status": 400,
  "message":"email: must be a valid email address."
}
`
var errRespPasswordResetEmailNotExistingEmail = `
{
  "status": 400,
  "message":"not exists data error"
}
`
var errRespPasswordResetEmailOverLimitCount = `
{
  "status": 400,
  "message": "you can't reset password as it exceeds limit counts"
}
`

func TestPasswordResetEmail(t *testing.T) {
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
			args:       args{reqBody: successReqPasswordResetEmail},
			method:     http.MethodPost,
			want:       testingHelper.RespSimpleSuccess,
			wantStatus: http.StatusOK,
		},
		{
			name:       "error request body without email",
			args:       args{reqBody: errReqPasswordResetEmailWithoutEmail},
			method:     http.MethodPost,
			want:       errRespPasswordResetEmailWithoutEmail,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error request body not email format",
			args:       args{reqBody: errReqPasswordResetEmailNotEmailFormat},
			method:     http.MethodPost,
			want:       errRespPasswordResetEmailNotEmailFormat,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error request body not existing email",
			args:       args{reqBody: errReqPasswordResetEmailNotExistingEmail},
			method:     http.MethodPost,
			want:       errRespPasswordResetEmailNotExistingEmail,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error over password reset limit count",
			args:       args{reqBody: successReqPasswordResetEmail},
			method:     http.MethodPost,
			want:       errRespPasswordResetEmailOverLimitCount,
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
			err = testingHelper.CreatePasswordReset(db, lib.NowFunc())
			assert.NoError(t, err)
			if tt.name == "error over password reset limit count" {
				err = testingHelper.UpdatePasswordResetEmailCount(db)
				assert.NoError(t, err)
			}

			// http request
			req, err := http.NewRequest(tt.method, "/password-reset-email", strings.NewReader(tt.args.reqBody))
			assert.NoError(t, err)
			req.Header.Add(helper.HeaderKeyContentType, "application/json")
			resp := httptest.NewRecorder()

			// test target
			PasswordController(resp, req)
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
  "message": "password: Your password must be at least 8 characters long, contain at least one number and have a mixture of uppercase and lowercase letters."
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
  "message": "not exists"
}
`
var errRespPasswordResetWrongKey = `
{
  "status": 400,
  "message": "password reset key is wrong"
}
`
var errRespPasswordExpiredKey = `
{
  "status": 400,
  "message": "password reset key is expired"
}
`

func TestPasswordReset(t *testing.T) {
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
			want:       errRespPasswordResetWrongKey,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error password_reset_key is expired",
			args:       args{reqBody: successReqPasswordReset},
			method:     http.MethodPost,
			want:       errRespPasswordExpiredKey,
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
			err = testingHelper.CreatePasswordReset(db, lib.NowFunc())
			assert.NoError(t, err)
			if tt.name == "error password_reset_key is expired" {
				err = testingHelper.UpdatePasswordResetExpiresAt(db, lib.NowFunc().Add(-time.Second))
			}
			assert.NoError(t, err)

			// http request
			req, err := http.NewRequest(tt.method, "/password-reset", strings.NewReader(tt.args.reqBody))
			assert.NoError(t, err)
			req.Header.Add(helper.HeaderKeyContentType, "application/json")
			resp := httptest.NewRecorder()

			// test target
			PasswordController(resp, req)
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
