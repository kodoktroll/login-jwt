package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kodoktroll/login-jwt/profile/db"
)

var s *server

type server struct {
	router *gin.Engine
	db     *db.Database
}

func newServer() *server {
	s := &server{}
	return s
}

func init() {
	s = newServer()
	s.router = gin.Default()
	s.routes()
	db, err := db.SetupDB()
	if err != nil {
		panic(err)
	}
	s.db = db
}

func main() {
	defer func() {
		err := s.db.CloseDB()
		if err != nil {
			panic(err)
		}
	}()
	s.run()
}

func (s *server) run() {
	s.router.Run(":8001")
}
