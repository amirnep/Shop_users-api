package app

import (
	"github.com/amirnep/shop/src/controllers"
	"github.com/amirnep/shop/src/middlewares"
)

func mapUrls() {
	router.POST("/Register", controllers.UsersController.Create)
	router.POST("/Login", controllers.UsersController.Login)

	protected := router.Group("/api")
	protected.Use(middlewares.JWTAuthCustomerMiddleware())

	protected.GET("/GetProfile", controllers.UsersController.GetProfile)
	protected.PUT("/EditProfile", controllers.UsersController.Update)
	protected.PATCH("/EditProfile", controllers.UsersController.Update)


	admin := router.Group("/api/admin")
	admin.Use(middlewares.JWTAuthMiddleware())

	admin.GET("/GetUsers", controllers.UsersController.GetUsers)
	admin.GET("/GetUser/:user_id", controllers.UsersController.Get)
	admin.DELETE("/DeleteUser/:user_id", controllers.UsersController.Delete)
}