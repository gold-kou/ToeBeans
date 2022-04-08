package http

type ResponseGetComments struct {
	PostingId int64                `json:"posting_id,omitempty"`
	Comments  []ResponseGetComment `json:"comments,omitempty"`
}
