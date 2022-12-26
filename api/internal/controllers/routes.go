package controllers

import (
	"github.com/buraktabakoglu/GOLANGAPPX/api/internal/middlewares"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_"github.com/buraktabakoglu/GOLANGAPPX/docs"
)

func (s *Server) initializeRoutes() {

	rou := s.Router.Group("/api")
	 {


		


		// Reset password:
		rou.POST("/password/forgot",s.ForgotPassword)
		rou.POST("/password/reset",s.ResetPassword)




		// use ginSwagger middleware to serve the API docs
		rou.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		//Login
		rou.POST("/login",s.Login)
		//Logout
		rou.DELETE("/logout",s.Logout)

		


		//Create user
		rou.POST("/users",s.CreateUser)
		//Read-all users
		rou.GET("/users",s.GetUsers)
		//Read user
		rou.GET("/users/:id",s.GetUser)
		//Update user 
		rou.PUT("/users/:id",middlewares.TokenAuthMiddleware(),s.UpdateUser)
		//Delete user
		rou.DELETE("/users/:id",middlewares.TokenAuthMiddleware(),s.DeleteUser)

		rou.POST("/todos",middlewares.TokenAuthMiddleware(),s.CreateTodo)


		//Read-all todos
		rou.GET("/todos",s.GetTodos)
		//Read todo
		rou.GET("/todos/:id",s.GetTodo)
		//Update todo
		rou.PUT("/todos/:id",middlewares.TokenAuthMiddleware(),s.UpdateATodo)
		//Delete todo
		rou.DELETE("/todos/:id",middlewares.TokenAuthMiddleware(),s.DeleteATodo)
		




			}
}