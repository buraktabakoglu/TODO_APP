package controllers

import "github.com/buraktabakoglu/GOLANGAPPX/api/middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//todos routes
	s.Router.HandleFunc("/todos", middlewares.SetMiddlewareJSON(s.CreateTodo)).Methods("POST")
	s.Router.HandleFunc("/todos", middlewares.SetMiddlewareJSON(s.GetTodos)).Methods("GET")
	s.Router.HandleFunc("/todos/{id}", middlewares.SetMiddlewareJSON(s.GetTodo)).Methods("GET")
	s.Router.HandleFunc("/todos/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateATodo))).Methods("PUT")
	s.Router.HandleFunc("/todos/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteATodo)).Methods("DELETE")
}
