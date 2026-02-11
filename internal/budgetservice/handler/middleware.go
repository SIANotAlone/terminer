package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) UserIdentity(c *gin.Context) {
	// ✅ пропускаем preflight
	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(http.StatusOK)
		return
	}

	header := c.GetHeader(authorizationHeader)
	if header == "" {
		NewErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		c.Abort()
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		NewErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		c.Abort()
		return
	}

	userID, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		NewErrorResponse(c, http.StatusUnauthorized, "invalid token")
		c.Abort()
		return
	}

	c.Set(userCtx, userID)
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

func (h *Handler) BudgetAccess(c *gin.Context) {
	userID, err := getUserId(c)
	// ... проверка userID (как у тебя была) ...

	var budgetID uuid.UUID

	// 1. ПЕРВЫМ ДЕЛОМ: Проверяем ID транзакции в параметрах URL
	// В твоем роуте /api/transaction/delete/:id параметр называется "id"
	tIDParam := c.Param("id")
	if tIDParam != "" {
		tID, err := uuid.Parse(tIDParam)
		if err == nil {
			// Если нашли ID в URL, сразу выясняем реальный budget_id
			realBudgetID, err := h.services.Transaction.GetBudgetIdByTransactionID(tID)
			if err == nil {
				budgetID = realBudgetID
			}
		}
	}

	// 2. Если в URL ничего не нашли, только тогда лезем в JSON
	if budgetID == uuid.Nil {
		var body struct {
			BudgetID      uuid.UUID `json:"budget_id"`
			TransactionID uuid.UUID `json:"transaction_id"`
			ID            uuid.UUID `json:"id"`
		}
		if err := c.ShouldBindBodyWith(&body, binding.JSON); err == nil {
			tID := body.TransactionID
			if tID == uuid.Nil {
				tID = body.ID
			}

			if tID != uuid.Nil {
				realBudgetID, _ := h.services.Transaction.GetBudgetIdByTransactionID(tID)
				budgetID = realBudgetID
			} else {
				budgetID = body.BudgetID
			}
		}
	}

	// 3. Запасной вариант: если это GET запрос и budgetid в URL
	if budgetID == uuid.Nil {
		if id := c.Param("budgetid"); id != "" {
			budgetID, _ = uuid.Parse(id)
		}
	}

	// 4. Финальный вердикт
	if budgetID == uuid.Nil {
		NewErrorResponse(c, http.StatusBadRequest, "could not determine budget id")
		c.Abort()
		return
	}

	// Теперь проверяем, имеет ли юзер доступ к РЕАЛЬНОМУ бюджету транзакции
	hasAccess, err := h.services.Access.HasUserAccessToBudget(userID, budgetID)
	if err != nil || !hasAccess {
		NewErrorResponse(c, http.StatusForbidden, "access denied")
		c.Abort()
		return
	}

	c.Next()
}
