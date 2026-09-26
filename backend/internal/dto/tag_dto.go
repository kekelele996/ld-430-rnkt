package dto

// TagRequest is the create payload for tags.
type TagRequest struct {
	Name     string `json:"name" validate:"required,min=1,max=50"`
	Category string `json:"category" validate:"required,oneof=Style Color Subject Mood Technical Other"`
}
