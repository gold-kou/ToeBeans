/*
 * Toe Beans
 *
 * API reference of Toe Beans
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package http

// ResponseGetNotifications - get notifications
type ResponseGetNotifications struct {

	// acted user name
	VisitedName string `json:"visited_name,omitempty"`

	// actions
	Actions []ResponseGetNotification `json:"actions,omitempty"`
}
