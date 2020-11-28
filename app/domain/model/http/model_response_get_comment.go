/*
 * Toe Beans
 *
 * API reference of Toe Beans
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package http

import (
	"time"
)

type ResponseGetComment struct {

	// comment id
	CommentId int64 `json:"comment_id"`

	// user_name
	UserName string `json:"user_name"`

	// commented datetime with TZ. This means created_at in postings table.
	CommentedAt time.Time `json:"commented_at"`

	// the content of comment
	Comment string `json:"comment"`
}