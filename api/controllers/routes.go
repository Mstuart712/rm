package controllers

import "github.com/Mstuart712/rm/api/middlewares"

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

	//cHARACTER routes
	s.Router.HandleFunc("/characters", middlewares.SetMiddlewareJSON(s.CreateCharacter)).Methods("POST")
	s.Router.HandleFunc("/characters", middlewares.SetMiddlewareJSON(s.GetCharacters)).Methods("GET")
	s.Router.HandleFunc("/characters/{id}", middlewares.SetMiddlewareJSON(s.GetCharacter)).Methods("GET")
	s.Router.HandleFunc("/characters/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateCharacter))).Methods("PUT")
	s.Router.HandleFunc("/characters/{id}", middlewares.SetMiddlewareAuthentication(s.DeleteCharacter)).Methods("DELETE")
}
