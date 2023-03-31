package controllers

import (
	"github.com/buraktabakoglu/GOLANGAPPX/api/internal/middlewares"

	_ "github.com/buraktabakoglu/GOLANGAPPX/api/pkg/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (s *Server) initializeRoutes() {

	Router := s.Router.Group("/api")
	{
		Router.Use(middlewares.CombinedAuthMiddleware())
		// use ginSwagger middleware to serve the API docs
		Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		//Read-all users
		Router.GET("/users", s.GetUsers)
		//Read user
		Router.GET("/users/:id", s.GetUser)
		//Update user
		Router.PUT("/users/:id", s.UpdateUser)
		//Delete user
		Router.DELETE("/users/:id", s.DeleteUser)

		//Create Todo
		Router.POST("/todos", s.CreateTodo)
		//Read-all todos
		Router.GET("/todos", s.GetTodos)
		//Read todo
		Router.GET("/todos/:id", s.GetTodo)
		//Update todo
		Router.PUT("/todos/:id", s.UpdateATodo)
		//Delete todo
		Router.DELETE("/todos/:id", s.DeleteATodo)

	}
	Public := s.Router.Group("/api")

	{

		//Create user
		Public.POST("/register", s.CreateUser)
		Public.GET("/activate/:token", s.ActivateUser)
		Public.POST("/password/forgot", s.ForgotPassword)
		Public.POST("/password/reset/:token", s.ResetPassword)
	}
}
