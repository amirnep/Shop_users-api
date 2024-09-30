package app

import (
	"github.com/amirnep/shop/controllers/users"
	"github.com/amirnep/shop/middlewares"
)

func mapUrls() {
	router.POST("/Register", users.UsersController.Create)
	router.POST("/Login", users.UsersController.Login)

	protected := router.Group("/api")
	protected.Use(middlewares.JWTAuthCustomerMiddleware())

	protected.GET("/GetProfile", users.UsersController.GetProfile)
	protected.PUT("/EditProfile", users.UsersController.Update)
	protected.PATCH("/EditProfile", users.UsersController.Update)


	admin := router.Group("/api/admin")
	admin.Use(middlewares.JWTAuthMiddleware())

	admin.GET("/GetUser/:user_id", users.UsersController.Get)
	admin.DELETE("/DeleteUser/:user_id", users.UsersController.Delete)
}