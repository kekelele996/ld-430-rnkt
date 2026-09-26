package dto

// CollectionRequest is the create/update payload for collections.
type CollectionRequest struct {
	Name          string `json:"name" validate:"required,min=1,max=100"`
	Description   string `json:"description"`
	CoverImageURL string `json:"cover_image_url"`
	IsPublic      bool   `json:"is_public"`
}

// AddAssetRequest adds an asset to a collection.
type AddAssetRequest struct {
	AssetID string `json:"asset_id" validate:"required"`
}

// AddCollaboratorRequest adds a collaborator to a collection.
type AddCollaboratorRequest struct {
	UserID string `json:"user_id" validate:"required"`
}
