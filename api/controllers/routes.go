package controllers

import "../middlewares"

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	// Login Route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Cats routes
	s.Router.HandleFunc("/breeds/{name}", middlewares.SetMiddlewareJSON(s.GetCat)).Methods("GET")
}
