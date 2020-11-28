package controller

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"

	"github.com/gold-kou/ToeBeans/app/domain/repository"

	"github.com/gold-kou/ToeBeans/app/adapter/http/helper"
	"github.com/gold-kou/ToeBeans/app/lib"
	"github.com/gold-kou/ToeBeans/testing/dummy"

	testingHelper "github.com/gold-kou/ToeBeans/testing"
	"github.com/stretchr/testify/assert"
)

var errRespDeletePostingWithoutPostingID = `
{
  "status": 400,
  "message": "posting_id: cannot be blank"
}
`
var errRespDeletePostingNotExistingID = `
{
  "status": 400,
  "message": "not exists data error"
}
`

func TestDeletePosting(t *testing.T) {
	type args struct {
		postingID int64
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
			args:       args{postingID: dummy.Posting1.ID},
			method:     http.MethodDelete,
			want:       testingHelper.RespSimpleSuccess,
			wantStatus: http.StatusOK,
		},
		{
			name:       "error empty posting_id",
			args:       args{},
			method:     http.MethodDelete,
			want:       errRespDeletePostingWithoutPostingID,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error not existing posting_id",
			args:       args{postingID: 99999},
			method:     http.MethodDelete,
			want:       errRespDeletePostingNotExistingID,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error forbidden guest user",
			args:       args{postingID: dummy.Posting1.ID},
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
			postingRepo := repository.NewPostingRepository(db)
			err := userRepo.Create(context.Background(), &dummy.User1)
			assert.NoError(t, err)
			err = postingRepo.Create(context.Background(), &dummy.Posting1)
			assert.NoError(t, err)

			// http request
			req, err := http.NewRequest(tt.method, fmt.Sprintf("/posting/%v", tt.args.postingID), nil)
			assert.NoError(t, err)
			vars := map[string]string{"posting_id": strconv.Itoa(int(tt.args.postingID))}
			req = mux.SetURLVars(req, vars)
			var token string
			if tt.name == "error forbidden guest user" {
				token, err = lib.GenerateToken(lib.GuestUserName)
			} else {
				token, err = lib.GenerateToken(dummy.User1.Name)
			}
			assert.NoError(t, err)
			req.Header.Add(helper.HeaderKeyAuthorization, "Bearer "+token)
			resp := httptest.NewRecorder()

			// test target
			PostingPostingIDController(resp, req)
			assert.NoError(t, err)

			// assert db
			if tt.wantStatus == 200 {
				postings, err := testingHelper.FindAllPostings(context.Background(), db)
				assert.NoError(t, err)
				assert.Equal(t, 0, len(postings))
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
