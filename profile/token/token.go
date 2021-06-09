package token

import (
	"fmt"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/kodoktroll/login-jwt/profile/redis"
)

type AccessDetails struct {
	AccessUuid string
	UserId     string
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	acc_secret := os.Getenv("ACCESS_SECRET")
	if len(acc_secret) == 0 {
		acc_secret = "cobain"
	}
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(acc_secret), nil
	})
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}
	return parsedToken, nil
}

func ExtractTokenMetadata(tokenString string) (*AccessDetails, error) {
	token, err := VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userId, ok := claims["user_id"].(string)
		if !ok {
			return nil, err
		}
		return &AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}
	return nil, err
}

func FetchAuth(authD *AccessDetails) (string, error) {
	userId, err := redis.Client.Get(authD.AccessUuid).Result()
	if err != nil {
		return "", err
	}
	return userId, nil
}
