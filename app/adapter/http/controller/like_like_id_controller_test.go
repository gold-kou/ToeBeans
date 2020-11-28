package controller

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
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

var errRespDeleteLikeWithoutLikeID = `
{
  "status": 400,
  "message": "like_id: cannot be blank"
}
`
var errRespDeleteLikeNotExistingID = `
{
  "status": 400,
  "message": "not exists data error"
}
`

func TestDeleteLike(t *testing.T) {
	type args struct {
		likeID int64
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
			args:       args{likeID: dummy.Like1.ID},
			method:     http.MethodDelete,
			want:       testingHelper.RespSimpleSuccess,
			wantStatus: http.StatusOK,
		},
		{
			name:       "error empty like_id",
			args:       args{},
			method:     http.MethodDelete,
			want:       errRespDeleteLikeWithoutLikeID,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error not existing like_id",
			args:       args{likeID: 99999},
			method:     http.MethodDelete,
			want:       errRespDeleteLikeNotExistingID,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "not allowed method",
			args:       args{},
			method:     http.MethodHead,
			want:       errorNotAllowedMethod,
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
			postingRepo := repository.NewPostingRepository(db)
			err := userRepo.Create(context.Background(), &dummy.User1)
			assert.NoError(t, err)
			err = userRepo.Create(context.Background(), &dummy.User2)
			assert.NoError(t, err)
			err = postingRepo.Create(context.Background(), &dummy.Posting1)
			assert.NoError(t, err)
			err = postingRepo.Create(context.Background(), &dummy.Posting2)
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/like", strings.NewReader(successReqRegisterLike))
			assert.NoError(t, err)
			token, err := lib.GenerateToken(dummy.User1.Name)
			assert.NoError(t, err)
			req.Header.Add(helper.HeaderKeyAuthorization, "Bearer "+token)
			resp := httptest.NewRecorder()
			LikeController(resp, req)
			assert.NoError(t, err)

			// http request
			req, err = http.NewRequest(tt.method, fmt.Sprintf("/like/%v", tt.args.likeID), nil)
			assert.NoError(t, err)
			vars := map[string]string{"like_id": strconv.Itoa(int(tt.args.likeID))}
			req = mux.SetURLVars(req, vars)
			token, err = lib.GenerateToken(dummy.User1.Name)
			assert.NoError(t, err)
			req.Header.Add(helper.HeaderKeyAuthorization, "Bearer "+token)
			resp = httptest.NewRecorder()

			// test target
			LikeLikeIDController(resp, req)
			assert.NoError(t, err)

			// assert db
			if tt.wantStatus == 200 {
				likes, err := testingHelper.FindAllLikes(context.Background(), db)
				assert.NoError(t, err)
				assert.Equal(t, 0, len(likes))

				users, err := testingHelper.FindAllUsers(context.Background(), db)
				assert.NoError(t, err)
				for _, user := range users {
					if user.Name == dummy.User1.Name {
						// like
						assert.Equal(t, int64(0), users[0].LikeCount)
					}
					if user.Name == dummy.User2.Name {
						// liked
						assert.Equal(t, int64(0), users[1].LikedCount)
					}
				}

				postings, err := testingHelper.FindAllPostings(context.Background(), db)
				assert.NoError(t, err)
				assert.Equal(t, int64(0), postings[0].LikedCount)
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
