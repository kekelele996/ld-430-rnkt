package dto

// CategoryRequest is the create/update payload for categories.
type CategoryRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	ParentID    string `json:"parent_id"`
	Icon        string `json:"icon"`
	SortOrder   int    `json:"sort_order"`
	Description string `json:"description"`
}
