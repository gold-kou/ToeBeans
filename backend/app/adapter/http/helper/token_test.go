package helper_test

import (
	"strings"
	"testing"

	"github.com/gold-kou/ToeBeans/backend/app/adapter/http/helper"
	testingHelper "github.com/gold-kou/ToeBeans/backend/testing"
	"github.com/gold-kou/ToeBeans/backend/testing/dummy"
	"github.com/stretchr/testify/assert"
)

// テスト順番に依存してしまうためあまり良くない
var sharedTestToken string

func TestGenerateToken(t *testing.T) {
	type args struct {
		userID   int64
		userName string
	}
	tests := []struct {
		name        string
		args        args
		environment string
		want        string
		wantErr     bool
	}{
		{
			name:        "success",
			args:        args{userID: dummy.User1.ID, userName: dummy.User1.Name},
			environment: dummy.SecretKey,
			want:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)

			// set env
			tmp := testingHelper.SetTestEnv("JWT_SECRET_KEY", tt.environment)
			defer tmp()

			// test target
			got, err := helper.GenerateToken(tt.args.userID, tt.args.userName)
			sharedTestToken = got

			if tt.wantErr {
				a.Error(err)
			} else {
				a.NoError(err)
				// just checking HEADER
				a.Equal(tt.want, strings.Split(got, ".")[0])
			}
		})
	}
}

func TestVerifyToken(t *testing.T) {
	type args struct {
		tokenString string
	}
	tests := []struct {
		name        string
		args        args
		environment string
		wantErr     bool
		watnErrMsg  string
	}{
		{
			name:        "success",
			args:        args{tokenString: sharedTestToken},
			environment: dummy.SecretKey,
			wantErr:     false,
		},
		{
			name:        "fail(expired)",
			args:        args{tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODM3MjIzNjIsImlhdCI6IjIwMjAtMDMtMDhUMTE6NTI6NDIuMjIxMjY2NCswOTowMCIsIm5hbWUiOiJ0ZXN0In0.FYMJmXo17aUhTpdaLifMovDQ0BiKSq8LnssLwxFvshI"},
			environment: dummy.SecretKey,
			wantErr:     true,
			watnErrMsg:  "token is expired",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// set env
			tmp := testingHelper.SetTestEnv("JWT_SECRET_KEY", tt.environment)
			defer tmp()

			a := assert.New(t)

			// test target
			tokenUserID, tokenUserName, err := helper.VerifyToken(tt.args.tokenString)

			// assert
			if tt.wantErr {
				a.Error(err)
				a.EqualError(err, tt.watnErrMsg)
			} else {
				a.NoError(err)
				a.Equal(dummy.User1.ID, tokenUserID)
				a.Equal(dummy.User1.Name, tokenUserName)
			}
		})
	}
}
