package middleware

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/assethub/assethub/internal/constants"
	apperrors "github.com/assethub/assethub/internal/errors"
	"github.com/assethub/assethub/internal/util"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// ErrorHandler normalizes panics and handler errors into the unified response.
func ErrorHandler(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}
		err := c.Errors.Last().Err
		var bizErr *apperrors.BusinessError
		var valErr *apperrors.ValidationError
		switch {
		case errors.As(err, &bizErr):
			util.Fail(c, httpStatusForCode(bizErr.Code), bizErr.Code, bizErr.Message)
		case errors.As(err, &valErr):
			util.Fail(c, http.StatusUnprocessableEntity, constants.CodeValidationFailed, valErr.Error())
		case errors.Is(err, mongo.ErrNoDocuments):
			util.Fail(c, http.StatusNotFound, constants.CodeNotFound, constants.MsgNotFound)
		default:
			logger.Error("unhandled error", "error", err.Error(), "path", c.Request.URL.Path)
			util.Fail(c, http.StatusInternalServerError, constants.CodeInternal, constants.MsgInternalError)
		}
	}
}

func httpStatusForCode(code int) int {
	switch code {
	case constants.CodeBadRequest:
		return http.StatusBadRequest
	case constants.CodeUnauthorized:
		return http.StatusUnauthorized
	case constants.CodeForbidden:
		return http.StatusForbidden
	case constants.CodeNotFound:
		return http.StatusNotFound
	case constants.CodeConflict:
		return http.StatusConflict
	case constants.CodeValidationFailed:
		return http.StatusUnprocessableEntity
	case constants.CodeRateLimited:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}
