/*
 * Toe Beans
 *
 * API reference of Toe Beans
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package http

// ResponseGetPostings - get postings
type ResponseGetPostings struct {

	// page
	Page int64 `json:"page,omitempty"`

	// list of posting
	Postings []ResponseGetPosting `json:"postings,omitempty"`
}
