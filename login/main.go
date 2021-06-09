package main

import (
	"github.com/gin-gonic/gin"
	// "github.com/kodoktroll/login-jwt/login/token"
	"github.com/kodoktroll/login-jwt/login/db"
)

var s *server

type server struct {
	route    *gin.Engine
	database *db.Database
	// token    *token.Token
}

func newServer() *server {
	s := &server{}
	// s.routes()
	return s
}

func main() {
	// s.route.POST("/login", s.Login())
	// s.route.POST("/signup", s.Signup())
	defer s.database.CloseDB()
	s.run()
}

func init() {
	s = newServer()
	s.route = gin.Default()
	db, err := db.SetupDB()
	if err != nil {
		panic(err)
	}
	s.routes()
	s.database = db
}

func (s *server) run() {
	defer func() {
		err := s.database.CloseDB()
		if err != nil {
			panic(err)
		}
	}()
	s.route.Run(":8000")
}

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var sampleUser = User{
	ID:       1,
	Username: "username",
	Password: "password",
}
