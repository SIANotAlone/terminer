package handler

import (
	"terminer/internal/budgetservice/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes(router *gin.Engine) {

	api := router.Group("/api", h.UserIdentity)
	{
		budget := api.Group("/budget")
		{
			budget.POST("/create", h.CreateBudget)
			budget.PUT("/update", h.UpdateBudget)
			budget.DELETE("/delete", h.DeleteBudget)
			budget.GET("/getall", h.GetAvailableBudgets)
			budget.PUT("/archive", h.ArchiveBudget)
			budget.PUT("/unarchive", h.UnArchiveBudget)
			budget.GET("/types", h.GetBudgetTypes)
			budget.GET("/currencies", h.GetCurrencies)
		}

		goal := api.Group("/goal")
		{
			goal.POST("/create", h.CreateGoal)
			goal.PUT("/update", h.UpdateGoal)
			goal.DELETE("/delete", h.DeleteGoal)
			goal.GET("/getavailablegoals", h.GetAvailableGoals)
			goal.GET("/getgoalstransactions/:goalid", h.GetGoalsTransactions)

		}
		transaction := api.Group("/transaction", h.BudgetAccess)
		{
			transaction.POST("/create", h.CreateTransaction)
			transaction.PUT("/update", h.UpdateTransaction)
			transaction.DELETE("/delete/:id", h.DeleteTransaction)
			transaction.GET("/getbybudget/:budgetid", h.GetTransactionsByBudget)
		}

		category := api.Group("/category")
		{
			category.POST("/create", h.CreateCategory)
			category.PUT("/update", h.UpdateCategory)
			category.DELETE("/delete", h.DeleteCategory)
			category.GET("/getavailablecategories", h.GetAvaliableCategories)
		}

		access := api.Group("/access")
		{
			access.POST("/sharebudget", h.ShareBudget)
			access.DELETE("/revokeaccess", h.RevokeAccess)
			access.GET("/getbudgetaccesslist", h.GetBudgetAccessList) //не используем
			access.GET("/getallusers", h.GetAllUsers)
			access.GET("/getbudgetusers/:budgetid", h.GetBudgetUsers)
		}

	}

}
