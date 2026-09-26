package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/assethub/assethub/internal/constants"
	"github.com/assethub/assethub/internal/errors"
	"github.com/assethub/assethub/internal/model"
	"github.com/assethub/assethub/internal/repository"
	"github.com/assethub/assethub/internal/util"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuthService handles registration, login and token issuing.
type AuthService struct {
	userRepo *repository.UserRepository
	roleRepo *repository.RoleRepository
	jwt      *util.JWTManager
	logger   *slog.Logger
}

// NewAuthService creates an AuthService.
func NewAuthService(userRepo *repository.UserRepository, roleRepo *repository.RoleRepository, jwt *util.JWTManager, logger *slog.Logger) *AuthService {
	return &AuthService{userRepo: userRepo, roleRepo: roleRepo, jwt: jwt, logger: logger}
}

// Register creates a new user (default role Uploader) and issues a token.
func (s *AuthService) Register(ctx context.Context, email, password, nickname string) (*model.User, string, time.Time, error) {
	if _, err := s.userRepo.FindByEmail(ctx, email); err == nil {
		return nil, "", time.Time{}, errors.NewBusinessError(constants.CodeConflict, constants.MsgEmailExists)
	} else if !errors.IsNotFound(err) {
		return nil, "", time.Time{}, fmt.Errorf("check email: %w", err)
	}
	hash, err := util.HashPassword(password)
	if err != nil {
		return nil, "", time.Time{}, fmt.Errorf("hash password: %w", err)
	}
	user := &model.User{
		Email:        email,
		PasswordHash: hash,
		Nickname:     nickname,
		Role:         string(constants.RoleUploader),
	}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, "", time.Time{}, err
	}
	s.logger.Info(constants.LogAuthRegister, "email", email, "user_id", user.ID.Hex())
	token, expiresAt, err := s.issueToken(user)
	if err != nil {
		return nil, "", time.Time{}, err
	}
	return user, token, expiresAt, nil
}

// Login verifies credentials and issues a token.
func (s *AuthService) Login(ctx context.Context, email, password string) (*model.User, string, time.Time, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, "", time.Time{}, errors.NewBusinessError(constants.CodeUnauthorized, constants.MsgInvalidCredentials)
		}
		return nil, "", time.Time{}, fmt.Errorf("find user: %w", err)
	}
	if !util.CheckPassword(user.PasswordHash, password) {
		return nil, "", time.Time{}, errors.NewBusinessError(constants.CodeUnauthorized, constants.MsgInvalidCredentials)
	}
	s.logger.Info(constants.LogAuthLogin, "email", email, "user_id", user.ID.Hex())
	token, expiresAt, err := s.issueToken(user)
	if err != nil {
		return nil, "", time.Time{}, err
	}
	return user, token, expiresAt, nil
}

// Me returns the user profile by id.
func (s *AuthService) Me(ctx context.Context, userID primitive.ObjectID) (*model.User, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}
	return user, nil
}

func (s *AuthService) issueToken(user *model.User) (string, time.Time, error) {
	token, expiresAt, err := s.jwt.Generate(user.ID.Hex(), user.Role, user.Email)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("issue token: %w", err)
	}
	return token, expiresAt, nil
}
