package app

import (
	"github.com/amirnep/shop/controllers/users"
	"github.com/amirnep/shop/middlewares"
)

func mapUrls() {
	public := router.Group("/api")

	public.POST("/Register", users.UsersController.Create)
	public.PUT("/EditProfile/:user_id", users.UsersController.Update)
	public.PATCH("/EditProfile/:user_id", users.UsersController.Update)
	public.POST("/Login", users.UsersController.Login)

	protected := router.Group("/api/admin")
	protected.Use(middlewares.JWTAuthMiddleware())

	protected.GET("/GetUser/:user_id", users.UsersController.Get)
	protected.DELETE("/DeleteUser/:user_id", users.UsersController.Delete)
}