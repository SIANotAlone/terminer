package handler

import (
	"terminer/internal/service"

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
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.SignUp)
		auth.POST("/sign-in", h.SignIn)
	}

	api := router.Group("/api", h.UserIdentity)
	{
		service := api.Group("/service")
		{
			service.POST("/create", h.CreateService)
			service.POST("/update", h.UpdateService)
			service.POST("/delete", h.DeleteService)
		}
		user := api.Group("/user")
		{
			user.GET("/getallusers", h.GetAllUsers)
		}
		
	}
	return router
}
