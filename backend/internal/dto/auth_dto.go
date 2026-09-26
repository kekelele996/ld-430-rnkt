package dto

// RegisterRequest is the payload for user registration.
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Nickname string `json:"nickname" validate:"required,min=1,max=50"`
}

// LoginRequest is the payload for login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// AuthResponse is returned after successful login/register.
type AuthResponse struct {
	Token     string  `json:"token"`
	ExpiresAt string  `json:"expires_at"`
	User      UserDTO `json:"user"`
}

// UserDTO is the public user representation.
type UserDTO struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Nickname  string `json:"nickname"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}
