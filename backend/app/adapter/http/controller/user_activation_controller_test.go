package controller

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"

	testingHelper "github.com/gold-kou/ToeBeans/backend/testing"
	"github.com/gold-kou/ToeBeans/backend/testing/dummy"
	"github.com/stretchr/testify/assert"
)

var errRespUserActivationWithoutUserName = `
{
  "status": 400,
  "message": "user_name cannot be blank"
}
`
var errRespUserActivationWithoutActivationKey = `
{
  "status": 400,
  "message": "activation_key cannot be blank"
}
`
var errRespUserActivationNameShort = `
{
  "status": 400,
  "message": "the length must be between 2 and 255"
}
`
var errRespUserActivationNameNotAlphanumeric = `
{
  "status": 400,
  "message": "must contain English letters and digits only"
}
`
var errRespUserActivationActivationKeyNotUUID = `
{
  "status": 400,
  "message": "must be a valid UUID"
}
`
var errRespUserActivationNotFound = `
{
  "status": 400,
  "message": "no such user_name and activation_key"
}
`

func TestUserActivationController(t *testing.T) {
	type args struct {
		userName      string
		activationKey string
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
			args:       args{userName: dummy.User1.Name, activationKey: dummy.User1.ActivationKey},
			method:     http.MethodGet,
			want:       testingHelper.RespSimpleSuccess,
			wantStatus: http.StatusOK,
		},
		{
			name:       "error user_name is empty",
			args:       args{activationKey: dummy.User1.ActivationKey},
			method:     http.MethodGet,
			want:       errRespUserActivationWithoutUserName,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error user_name is short",
			args:       args{userName: "a", activationKey: dummy.User1.ActivationKey},
			method:     http.MethodGet,
			want:       errRespUserActivationNameShort,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error user_name is not alphanumeric",
			args:       args{userName: "test_user1", activationKey: dummy.User1.ActivationKey},
			method:     http.MethodGet,
			want:       errRespUserActivationNameNotAlphanumeric,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error activation_key is empty",
			args:       args{userName: dummy.User1.Name},
			method:     http.MethodGet,
			want:       errRespUserActivationWithoutActivationKey,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error activation_key not UUID",
			args:       args{userName: dummy.User1.Name, activationKey: "abc"},
			method:     http.MethodGet,
			want:       errRespUserActivationActivationKeyNotUUID,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error not existing user_name",
			args:       args{userName: "testUser0", activationKey: dummy.User1.ActivationKey},
			method:     http.MethodGet,
			want:       errRespUserActivationNotFound,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error not existing activation_key",
			args:       args{userName: dummy.User1.Name, activationKey: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"},
			method:     http.MethodGet,
			want:       errRespUserActivationNotFound,
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

			// http request
			req, err := http.NewRequest(tt.method, fmt.Sprintf("/user/%s/%s", tt.args.userName, tt.args.activationKey), nil)
			assert.NoError(t, err)
			vars := map[string]string{"user_name": tt.args.userName, "activation_key": tt.args.activationKey}
			req = mux.SetURLVars(req, vars)
			resp := httptest.NewRecorder()

			// test target
			UserActivationController(resp, req)
			assert.NoError(t, err)

			// assert db
			if tt.wantStatus == 200 {
				users, err := testingHelper.FindAllUsers(context.Background(), db)
				assert.NoError(t, err)
				assert.Equal(t, true, users[0].EmailVerified)
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
