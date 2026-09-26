package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/assethub/assethub/internal/constants"
	"github.com/assethub/assethub/internal/model"
	"github.com/assethub/assethub/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuditService writes and queries audit logs.
type AuditService struct {
	repo   *repository.AuditLogRepository
	logger *slog.Logger
}

// NewAuditService creates an AuditService.
func NewAuditService(repo *repository.AuditLogRepository, logger *slog.Logger) *AuditService {
	return &AuditService{repo: repo, logger: logger}
}

// Log writes a single audit entry.
func (s *AuditService) Log(ctx context.Context, operatorID primitive.ObjectID, action, resource, resourceID, detail, ip string) error {
	entry := &model.AuditLog{
		OperatorID: operatorID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Detail:     detail,
		IP:         ip,
		CreatedAt:  time.Now(),
	}
	if err := s.repo.Create(ctx, entry); err != nil {
		return fmt.Errorf("audit log: %w", err)
	}
	s.logger.Info(constants.LogAuditWritten, "action", action, "resource", resource, "operator", operatorID.Hex())
	return nil
}

// List returns paginated audit logs.
func (s *AuditService) List(ctx context.Context, skip, limit int64) ([]model.AuditLog, error) {
	logs, err := s.repo.List(ctx, skip, limit)
	if err != nil {
		return nil, fmt.Errorf("list audit logs: %w", err)
	}
	return logs, nil
}
