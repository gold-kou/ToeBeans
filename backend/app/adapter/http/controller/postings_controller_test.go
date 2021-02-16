package controller

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	"github.com/gold-kou/ToeBeans/backend/app/domain/repository"
	"github.com/gold-kou/ToeBeans/backend/app/lib"
	testingHelper "github.com/gold-kou/ToeBeans/backend/testing"
	"github.com/gold-kou/ToeBeans/backend/testing/dummy"
	"github.com/stretchr/testify/assert"
)

var successRespGetPostings = `
{
  "postings": [
    {
      "posting_id": 1,
      "user_name": "testUser1",
      "uploaded_at": "2020-01-01T00:00:00+09:00",
      "title": "20200101000000_testUser1_This is a sample posting.",
      "image_url": "http://minio:9000/postings/20200101000000_testUser1_This%20is%20a%20sample%20posting.",
      "liked": 0
    }
  ]
}
`
var errRespGetPostingsWithoutSinceAt = `
{
  "status": 400,
  "message": "since_at: cannot be blank."
}
`
var errRespGetPostingsWithoutLimit = `
{
  "status": 400,
  "message": "limit: cannot be blank."
}
`
var errRespGetPostingsNotExistsUserName = `
{
  "status": 400,
  "message": "not exists data error"
}
`

func TestGetPostings(t *testing.T) {
	type args struct {
		sinceAt  string
		limit    string
		userName string
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
			args:       args{sinceAt: "2020-01-01T00:00:00Z", limit: "50"},
			method:     http.MethodGet,
			want:       successRespGetPostings,
			wantStatus: http.StatusOK,
		},
		{
			name:       "success with user_name",
			args:       args{sinceAt: "2020-01-01T00:00:00Z", limit: "50", userName: dummy.User1.Name},
			method:     http.MethodGet,
			want:       successRespGetPostings,
			wantStatus: http.StatusOK,
		},
		{
			name:       "error empty since_at",
			args:       args{limit: "50"},
			method:     http.MethodGet,
			want:       errRespGetPostingsWithoutSinceAt,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error empty limit",
			args:       args{sinceAt: "2020-01-01T00:00:00Z"},
			method:     http.MethodGet,
			want:       errRespGetPostingsWithoutLimit,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "error not exists user_name",
			args:       args{sinceAt: "2020-01-01T00:00:00Z", limit: "50", userName: dummy.User2.Name},
			method:     http.MethodGet,
			want:       errRespGetPostingsNotExistsUserName,
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
			postingRepo := repository.NewPostingRepository(db)
			err := userRepo.Create(context.Background(), &dummy.User1)
			assert.NoError(t, err)
			err = postingRepo.Create(context.Background(), &dummy.Posting1)
			assert.NoError(t, err)
			err = testingHelper.UpdateNow(db, "postings")
			assert.NoError(t, err)

			// http request
			var req *http.Request
			if tt.args.userName == "" {
				req, err = http.NewRequest(tt.method, fmt.Sprintf("/postings?since_at=%s&limit=%s", tt.args.sinceAt, tt.args.limit), nil)
			} else {
				req, err = http.NewRequest(tt.method, fmt.Sprintf("/postings?since_at=%s&limit=%s&user_name=%s", tt.args.sinceAt, tt.args.limit, tt.args.userName), nil)
			}
			assert.NoError(t, err)
			token, err := lib.GenerateToken(dummy.User1.Name)
			assert.NoError(t, err)
			req.Header.Add(helper.HeaderKeyAuthorization, "Bearer "+token)
			resp := httptest.NewRecorder()

			// test target
			PostingsController(resp, req)
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
