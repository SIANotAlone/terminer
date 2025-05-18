package handler

import (
	"terminer/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true, // Разрешить запросы с любых источников
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", // Разрешить все методы
		},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"}, // Разрешить указанные заголовки
	}))
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
			service.POST("/create_promo", h.CreatePromoService)
			service.POST("/update", h.UpdateService)
			service.POST("/delete", h.DeleteService)
			service.GET("/gettypes", h.GetTypes)
			service.POST("/createservicetype", h.CreateServiceType)
			service.GET("/getmyservices", h.GetMyServices)
			service.GET("/available", h.GetAvailableServices)
			service.POST("/availabletime", h.GetAvailableTime)
			service.POST("/validate_promo", h.ValidatePromoCode)
			service.POST("/activate_promo", h.ActivatePromoCode)
			service.GET("/getmyactualservices", h.GetMyActualServices)
			service.POST("/getmyhistory", h.GetMyHistory)

		}
		record := api.Group("/record")
		{
			record.POST("/create", h.CreateRecord)
			record.POST("/done", h.DoneRecord)
			record.POST("/confirm", h.ConfirmRecord)
			record.POST("/termins", h.GetTerminsFromService)
		}

		user := api.Group("/user")
		{
			user.GET("/getallusers", h.GetAllUsers)
		}
		comment := api.Group("/comment")
		{
			comment.POST("/create", h.CreateComment)
			comment.POST("/update", h.UpdateComment)
			comment.POST("/delete", h.DeleteComment)
			comment.POST("/getcomments", h.GetCommentsOnRecord)

		}
		Termin := api.Group("/termin")
		{
			Termin.GET("/getallperformertermins", h.GetAllPerformerTermins)
			Termin.GET("/getallusertermins", h.GetAllUserTermins)
		}

	}

	return router
}
