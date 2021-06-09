package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kodoktroll/login-jwt/login/db"
	"github.com/kodoktroll/login-jwt/login/token"
)

func (s *server) Login() gin.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type response struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	type error struct {
		Error string `json:"error"`
	}
	return func(c *gin.Context) {
		var u request
		if err := c.ShouldBindJSON(&u); err != nil {
			c.JSON(http.StatusUnprocessableEntity, error{Error: "Invalid json provided"})
			return
		}
		user, err := s.database.GetUser(u.Username)
		log.Print(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		if comparePassword(user.Password, []byte(u.Password)) == false {
			// false == password ga sama
			c.JSON(http.StatusUnauthorized, error{Error: "Invalid Credentials"})
			return
		}
		tokens, err := token.CreateToken(user.ID.Hex())
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, error{Error: err.Error()})
			return
		}
		saveErr := token.CreateAuth(user.ID.Hex(), tokens)
		if saveErr != nil {
			c.JSON(http.StatusUnprocessableEntity, error{saveErr.Error()})
		}
		res := &response{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		}
		c.JSON(http.StatusOK, res)
	}
}

func (s *server) Signup() gin.HandlerFunc {
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type response struct {
		Message string `json:"message"`
	}
	type error struct {
		Error string `json:"error"`
	}
	return func(c *gin.Context) {
		var u request
		if err := c.ShouldBindJSON(&u); err != nil {
			c.JSON(http.StatusUnprocessableEntity, error{Error: "Invalid json provided"})
			return
		}
		user := db.User{
			Username: u.Username,
			Password: hashAndSalt(u.Password),
		}
		err := s.database.InsertUser(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, error{Error: err.Error()})
			return
		}
		res := response{
			Message: "Successfully inserted",
		}
		c.JSON(http.StatusOK, res)
	}
}

func (s *server) Logout() gin.HandlerFunc {
	type response struct {
		Message string `json:"message"`
	}
	type error struct {
		Error string `json:"error"`
	}
	return func(c *gin.Context) {
		tokenString := extractToken(c.Request)
		auth, err := token.ExtractTokenMetadata(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, &error{Error: "Unauthorized: " + err.Error()})
			return
		}
		deleted, delErr := token.DeleteAuth(auth.AccessUuid)
		if delErr != nil || deleted == 0 {
			c.JSON(http.StatusUnauthorized, &error{Error: "Unauthorized: You are not logged in"})
			return
		}
		c.JSON(http.StatusOK, &response{Message: "Successfully logged out"})

	}
}

func (s *server) Refresh() gin.HandlerFunc {
	type request struct {
		RefreshToken string `json:"refresh_token"`
	}
	type error struct {
		Error string `json:"error"`
	}
	type response struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}
	return func(c *gin.Context) {
		var rt request
		if err := c.ShouldBindJSON(&rt); err != nil {
			c.JSON(http.StatusUnprocessableEntity, error{Error: "Invalid json provided"})
			return
		}
		tk, err := token.VerifyRefreshToken(rt.RefreshToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, error{Error: "Unable to verify token: Refresh token might be expired or invalid"})
			return
		}
		refreshDetails, err := token.ExtractRefreshMetadata(tk)
		if err != nil {
			c.JSON(http.StatusUnauthorized, error{Error: err.Error()})
			return
		}
		deleted, delErr := token.DeleteAuth(refreshDetails.RefreshUuid)
		if delErr != nil || deleted == 0 {
			c.JSON(http.StatusUnauthorized, error{Error: "Unauthorized"})
			return
		}
		newTk, createErr := token.CreateToken(refreshDetails.UserId)
		if createErr != nil {
			c.JSON(http.StatusForbidden, error{Error: createErr.Error()})
			return
		}
		saveErr := token.CreateAuth(refreshDetails.RefreshUuid, newTk)
		if saveErr != nil {
			c.JSON(http.StatusForbidden, error{Error: saveErr.Error()})
			return
		}
		c.JSON(http.StatusOK, response{AccessToken: newTk.AccessToken, RefreshToken: newTk.RefreshToken})
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
