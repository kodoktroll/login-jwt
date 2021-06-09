package main

func (s *server) routes() {
	s.router.POST("/profile", s.InsertProfile())
	s.router.GET("/profile", s.GetProfile())
}
