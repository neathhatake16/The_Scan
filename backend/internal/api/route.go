package api

import (
	"github.com/gin-gonic/gin"
	"github.com/neathhatake/the_Scan/internal/config"
	"github.com/neathhatake/the_Scan/internal/handlers"
	"github.com/neathhatake/the_Scan/internal/middleware"
	"github.com/neathhatake/the_Scan/internal/repositories/repoimpl"
	"github.com/neathhatake/the_Scan/internal/services"
	"github.com/neathhatake/the_Scan/pkg/logger"
	"gorm.io/gorm"
)


func Run(cfg *config.Config, db *gorm.DB) {
	r := gin.Default()
	r.SetTrustedProxies([]string{"192.168.1.3"})

	r.Use(gin.Recovery())
	r.Use(middleware.RequestLogger())

	authSvc := services.NewAuthService(repoimpl.NewUserRepository(db), repoimpl.NewTokenRepository(db), cfg)
	authHandler := handlers.NewAuthHandler(authSvc)

	userSvc := services.NewUserService(repoimpl.NewUserRepository(db), authSvc)
	userHandler := handlers.NewUserHandler(userSvc)

	docSvc := services.NewDocumentService(repoimpl.NewDocumentRepository(db) , repoimpl.NewUserRepository(db), cfg)
	docHandler := handlers.NewDocumentHandler(docSvc)


	// Auth endpoints
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.Refresh)
		auth.POST("/logout", authHandler.Logout)
	}


	//user endpoints
	user := r.Group("/users")
	user.Use(middleware.Auth(cfg))
	{
		user.GET("/me", userHandler.GetProfile)
		user.PATCH("/me", userHandler.UpdateProfile)
		user.POST("/me/change-password", userHandler.ChangePassword)
		user.GET("/me/storage", userHandler.GetStorage)
	}


	// Document endpoints
	doc := r.Group("/documents")
	doc.Use(middleware.Auth(cfg))
	{
		doc.POST("/scan", docHandler.Scan)
		doc.GET("/list", docHandler.ListDocuments)
		doc.GET("/download/:id", docHandler.DownloadDocuments)
		doc.PUT("/rename/:id", docHandler.RenameDocument)
		doc.DELETE("/delete/:id", docHandler.DeleteDocuments)
	}




	// Run server
	logger.Log.Infow("server starting", "port", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		logger.Log.Errorw("failed to start server", "error", err)
	}

}