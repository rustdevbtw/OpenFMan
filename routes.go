package main

import (
	"github.com/gin-gonic/gin"
	"rajdeepm.xyz/ofm/handlers"
)

func setupRoutes(r *gin.Engine) {
	// Public routes
	r.POST("/login", handlers.Login)

	// Routes requiring authentication
	authorized := r.Group("/", handlers.AuthMiddleware())
	authorized.GET("/get/*filename", handlers.GetFile)
	authorized.GET("/list/*dir", handlers.Listing)
	authorized.POST("/upload/*filename", handlers.UploadFile)
	authorized.DELETE("/delete/*filename", handlers.DeleteFile)
	authorized.POST("/rename/*oldname", handlers.RenameFile)
	authorized.POST("/logout", handlers.Logout)
}
