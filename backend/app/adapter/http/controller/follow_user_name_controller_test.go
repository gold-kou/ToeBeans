package controller

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"

	"github.com/gold-kou/ToeBeans/app/domain/repository"

	"github.com/gold-kou/ToeBeans/app/adapter/http/helper"
	"github.com/gold-kou/ToeBeans/app/lib"
	"github.com/gold-kou/ToeBeans/testing/dummy"

	testingHelper "github.com/gold-kou/ToeBeans/testing"
	"github.com/stretchr/testify/assert"
)

var errRespDeleteFollowWithoutUserName = `
{
  "status": 400,
  "message": "followed_user_name: cannot be blank."
}
`

func TestDeleteFollow(t *testing.T) {
	type args struct {
		followedUserName string
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
			args:       args{followedUserName: dummy.User2.Name},
			method:     http.MethodDelete,
			want:       testingHelper.RespSimpleSuccess,
			wantStatus: http.StatusOK,
		},
		{
			name:       "error empty user_name",
			args:       args{},
			method:     http.MethodDelete,
			want:       errRespDeleteFollowWithoutUserName,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error forbidden guest user",
			args:       args{followedUserName: dummy.User2.Name},
			method:     http.MethodDelete,
			want:       testingHelper.ErrForbidden,
			wantStatus: http.StatusForbidden,
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
			err = userRepo.Create(context.Background(), &dummy.User2)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/follow", strings.NewReader(successReqRegisterFollow))
			assert.NoError(t, err)
			token, err := lib.GenerateToken(dummy.User1.Name)
			assert.NoError(t, err)
			req.Header.Add(helper.HeaderKeyAuthorization, "Bearer "+token)
			resp := httptest.NewRecorder()
			FollowController(resp, req)
			assert.NoError(t, err)

			// http request
			req, err = http.NewRequest(tt.method, fmt.Sprintf("/follow/%s", tt.args.followedUserName), nil)
			assert.NoError(t, err)
			vars := map[string]string{"followed_user_name": tt.args.followedUserName}
			req = mux.SetURLVars(req, vars)
			if tt.name == "error forbidden guest user" {
				token, err = lib.GenerateToken(lib.GuestUserName)
			} else {
				token, err = lib.GenerateToken(dummy.User1.Name)
			}
			assert.NoError(t, err)
			req.Header.Add(helper.HeaderKeyAuthorization, "Bearer "+token)
			resp = httptest.NewRecorder()

			// test target
			FollowUserNameController(resp, req)
			assert.NoError(t, err)

			// assert db
			if tt.wantStatus == 200 {
				follows, err := testingHelper.FindAllFollows(context.Background(), db)
				assert.NoError(t, err)
				assert.Equal(t, 0, len(follows))

				users, err := testingHelper.FindAllUsers(context.Background(), db)
				assert.NoError(t, err)
				for _, user := range users {
					if user.Name == dummy.User1.Name {
						// following
						assert.Equal(t, int64(0), users[0].FollowCount)
					}
					if user.Name == dummy.User2.Name {
						// followed
						assert.Equal(t, int64(0), users[1].FollowedCount)
					}
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
