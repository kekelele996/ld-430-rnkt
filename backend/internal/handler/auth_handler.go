package handler

import (
	"github.com/assethub/assethub/internal/constants"
	"github.com/assethub/assethub/internal/dto"
	"github.com/assethub/assethub/internal/errors"
	"github.com/assethub/assethub/internal/middleware"
	"github.com/assethub/assethub/internal/model"
	"github.com/assethub/assethub/internal/service"
	"github.com/assethub/assethub/internal/util"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuthHandler exposes authentication endpoints.
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates an AuthHandler.
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register handles POST /auth/register.
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := middleware.ValidateJSON(c, &req); err != nil {
		_ = c.Error(err)
		return
	}
	user, token, expiresAt, err := h.authService.Register(c.Request.Context(), req.Email, req.Password, req.Nickname)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.Created(c, dto.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt.Format("2006-01-02T15:04:05Z07:00"),
		User:      toUserDTO(user),
	})
}

// Login handles POST /auth/login.
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := middleware.ValidateJSON(c, &req); err != nil {
		_ = c.Error(err)
		return
	}
	user, token, expiresAt, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, dto.AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt.Format("2006-01-02T15:04:05Z07:00"),
		User:      toUserDTO(user),
	})
}

// Me handles GET /auth/me.
func (h *AuthHandler) Me(c *gin.Context) {
	raw := c.GetString(middleware.ContextUserIDKey)
	oid, err := primitive.ObjectIDFromHex(raw)
	if err != nil {
		_ = c.Error(errors.NewBusinessError(constants.CodeUnauthorized, constants.MsgUnauthorized))
		return
	}
	user, err := h.authService.Me(c.Request.Context(), oid)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, toUserDTO(user))
}

func toUserDTO(user *model.User) dto.UserDTO {
	return dto.UserDTO{
		ID:        user.ID.Hex(),
		Email:     user.Email,
		Nickname:  user.Nickname,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
