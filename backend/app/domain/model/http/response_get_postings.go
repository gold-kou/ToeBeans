package http

type ResponseGetPostings struct {
	Postings []ResponseGetPosting `json:"postings"`
}
