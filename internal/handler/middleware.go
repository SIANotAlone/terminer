package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

// AuthMiddleware — middleware авторизації
func (h *Handler) UserIdentity(c *gin.Context) {
	// ✅ 1. Пропускаємо CORS preflight
	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(http.StatusOK)
		return
	}

	// ✅ 2. Читаємо Authorization header
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		NewErrorResponse(c, http.StatusUnauthorized, "authorization header is empty")
		c.Abort()
		return
	}

	// ✅ 3. Перевіряємо формат "Bearer <token>"
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		NewErrorResponse(c, http.StatusUnauthorized, "invalid authorization header format")
		c.Abort()
		return
	}

	// ✅ 4. Парсимо токен
	userID, err := h.services.Authorization.ParseToken(parts[1])
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, "invalid token")
		c.Abort()
		return
	}

	// ✅ 5. Кладемо userID в context
	c.Set(userCtx, userID)

	// ✅ 6. Продовжуємо ланцюжок middleware
	c.Next()
}

func getUserId(c *gin.Context) (uuid.UUID, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return uuid.UUID{}, errors.New("user id not found")
	}
	id_uuid, ok := id.(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("user id is not uuid.UUID type")
	}
	return id_uuid, nil
}
