package main

func (s *server) routes() {
	s.route.POST("/login", s.Login())
	s.route.POST("/signup", s.Signup())
	s.route.POST("/logout", s.Logout())
	s.route.POST("/token/refresh", s.Refresh())
}
