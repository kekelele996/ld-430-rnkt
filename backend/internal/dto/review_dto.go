package dto

// ReviewRequest is the payload for asset reviews.
type ReviewRequest struct {
	Result  string `json:"result" validate:"required,oneof=Approved Rejected NeedsRevision"`
	Comment string `json:"comment"`
}
