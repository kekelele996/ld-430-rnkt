package constants

// Role enumerates RBAC roles.
type Role string

const (
	RoleAdmin     Role = "Admin"
	RoleModerator Role = "Moderator"
	RoleUploader  Role = "Uploader"
	RoleViewer    Role = "Viewer"
)

// AllRoles lists every defined role.
var AllRoles = []Role{RoleAdmin, RoleModerator, RoleUploader, RoleViewer}
