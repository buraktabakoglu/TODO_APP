package controllers

import "github.com/buraktabakoglu/GOLANGAPPX/api/middlewares"

func (s *Server) initializeRoutes() {

	rou := s.Router.Group("api/")
	 {

		rou.POST("/login",s.Login)


		rou.POST("/users",s.CreateUser)
		rou.GET("/users",s.GetUsers)
		rou.GET("/users/:id",s.GetUser)
		rou.PUT("/users/:id",middlewares.TokenAuthMiddleware(),s.UpdateUser)
		rou.DELETE("/users/id",middlewares.TokenAuthMiddleware(),s.DeleteUser)

		rou.POST("/todos",middlewares.TokenAuthMiddleware(),s.CreateTodo)
		rou.GET("/todos",s.GetTodos)
		rou.GET("/todos/:id",s.GetTodo)
		rou.PUT("/todos/:id",middlewares.TokenAuthMiddleware(),s.UpdateATodo)
		rou.DELETE("/todos/:id",middlewares.TokenAuthMiddleware(),s.DeleteATodo)
			}
}