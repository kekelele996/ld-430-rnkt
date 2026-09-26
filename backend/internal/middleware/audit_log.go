package middleware

import (
	"time"

	"github.com/assethub/assethub/internal/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AuditLogger returns a middleware that records the request as an audit log entry.
// Only requests routed through groups using this middleware are captured.
func AuditLogger(db *mongo.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		action := c.GetString("audit_action")
		if action == "" {
			return
		}
		operatorID := primitive.NilObjectID
		if raw, ok := c.Get(ContextUserIDKey); ok {
			if id, err := primitive.ObjectIDFromHex(raw.(string)); err == nil {
				operatorID = id
			}
		}
		entry := &model.AuditLog{
			OperatorID: operatorID,
			Action:     action,
			Resource:   c.Request.Method + " " + c.Request.URL.Path,
			ResourceID: c.Param("id"),
			Detail:     c.GetString("audit_detail"),
			IP:         c.ClientIP(),
			CreatedAt:  time.Now(),
		}
		_, _ = db.Collection("audit_logs").InsertOne(c.Request.Context(), entry)
	}
}

// SetAuditAction attaches audit metadata to the request context.
func SetAuditAction(c *gin.Context, action, detail string) {
	c.Set("audit_action", action)
	c.Set("audit_detail", detail)
}
