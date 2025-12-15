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

func (h *Handler) UserIdentity(c *gin.Context) {
	// 🔥 ВАЖНО: пропускаем preflight
	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(http.StatusOK)
		return
	}

	header := c.GetHeader(authorizationHeader)
	if header == "" {
		NewErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		NewErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	userID, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, "invalid token")
		return
	}

	c.Set(userCtx, userID)
	c.Next()
}

/*************  ✨ Windsurf Command ⭐  *************/
// getUserId returns user id from context. If user id is not found or it is not uuid.UUID type, it returns error.
/*******  830e953c-5979-4def-947e-0b34eec07ef8  *******/
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
