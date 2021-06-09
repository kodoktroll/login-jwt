package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kodoktroll/login-jwt/profile/db"
	"github.com/kodoktroll/login-jwt/profile/token"
)

func (s *server) InsertProfile() gin.HandlerFunc {
	type request struct {
		Username    string `json:"username"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		PictureLink string `json:"picture_link"`
	}
	type error struct {
		Error string `json:"error"`
	}
	type response struct {
		Message string `json:"message"`
	}
	return func(c *gin.Context) {
		var p request
		if err := c.ShouldBindJSON(&p); err != nil {
			c.JSON(http.StatusUnprocessableEntity, error{Error: "Invalid json provided"})
			return
		}
		tkString := extractToken(c.Request)
		auth, err := token.ExtractTokenMetadata(tkString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, error{Error: err.Error()})
			return
		}
		profile := db.Profile{
			UserID:      auth.UserId,
			Username:    p.Username,
			FirstName:   p.FirstName,
			LastName:    p.LastName,
			PictureLink: p.PictureLink,
		}
		err = s.db.InsertProfile(profile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, error{Error: err.Error()})
			return
		}
		c.JSON(http.StatusOK, response{Message: "Successfully inserted profile"})
	}
}

func extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func (s *server) GetProfile() gin.HandlerFunc {
	type error struct {
		Error string `json:"error"`
	}
	type response struct {
		UserId      string `json:"user_id"`
		Username    string `json:"username"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		PictureLink string `json:"picture_link"`
	}
	return func(c *gin.Context) {
		tkString := extractToken(c.Request)
		auth, err := token.ExtractTokenMetadata(tkString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, error{Error: err.Error()})
			return
		}
		// prof, err := s.db.FindOne(auth.UserId)
		// log.Print(prof)
		// log.Print(&prof)
		prof, err := s.db.GetProfile(auth.UserId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, error{Error: err.Error()})
			return
		}
		log.Print(prof)
		res := &response{
			UserId:      prof.UserID,
			Username:    prof.Username,
			FirstName:   prof.FirstName,
			LastName:    prof.LastName,
			PictureLink: prof.PictureLink,
		}
		c.JSON(http.StatusOK, res)
	}
}
