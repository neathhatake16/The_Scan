package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/neathhatake/the_Scan/internal/config"
	"github.com/neathhatake/the_Scan/internal/services"
	"github.com/neathhatake/the_Scan/pkg/apperrors"
	"github.com/neathhatake/the_Scan/pkg/respone"
)

const ctxUserID = "userID"
var ErrUnauthorized = apperrors.Unauthorized("missing or invalid Authorization header")

// Auth validates the Bearer JWT and sets userID in the Gin context.
func Auth(cfg *config.Config) gin.HandlerFunc {
	// AuthService with nil repos — only ParseAccessToken is called here,
	// which needs only cfg.JWTSecret and no DB access.
	svc := services.NewAuthService(nil, nil, cfg)

	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			respone.Error(c , ErrUnauthorized)
			c.Abort()
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")
		userID, err := svc.ParseAccessToken(token)
		if err != nil {
			respone.Error(c, ErrUnauthorized)
			c.Abort()
			return
		}
		c.Set(ctxUserID, userID)
		c.Next()
	}
}

// UserID extracts the authenticated user ID from context.
// Panics if called outside an Auth-protected route (programming error).
func UserID(c *gin.Context) uint {
	v, exists := c.Get(ctxUserID)
	if !exists {
		panic("middleware.UserID called outside of Auth middleware")
	}
	id, _ := v.(uint)
	return id
}
