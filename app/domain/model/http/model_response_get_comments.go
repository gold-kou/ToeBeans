/*
 * Toe Beans
 *
 * API reference of Toe Beans
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package http

// ResponseGetComments - get comments
type ResponseGetComments struct {

	// posting id
	PostingId int64 `json:"posting_id,omitempty"`

	// list of comment
	Comments []ResponseGetComment `json:"comments,omitempty"`
}