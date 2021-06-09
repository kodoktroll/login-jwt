package token

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kodoktroll/login-jwt/login/redis"
	"github.com/twinj/uuid"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

func CreateToken(userID string) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error
	acc_secret := os.Getenv("ACCESS_SECRET")
	if len(acc_secret) == 0 {
		acc_secret = "cobain"
	}
	accessTokenClaims := jwt.MapClaims{}
	accessTokenClaims["authorized"] = true
	accessTokenClaims["access_uuid"] = td.AccessUuid
	accessTokenClaims["user_id"] = userID
	accessTokenClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	at, err := accessToken.SignedString([]byte(acc_secret))
	if err != nil {
		return nil, err
	}
	td.AccessToken = at

	refresh_secret := os.Getenv("REFRESH_SECRET")
	if len(refresh_secret) == 0 {
		refresh_secret = "cobain!!"
	}
	refreshTokenClaims := jwt.MapClaims{}
	refreshTokenClaims["authorized"] = true
	refreshTokenClaims["refresh_uuid"] = td.RefreshUuid
	refreshTokenClaims["user_id"] = userID
	refreshTokenClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	rt, err := refreshToken.SignedString([]byte(refresh_secret))
	if err != nil {
		return nil, err
	}
	td.RefreshToken = rt

	return td, nil
}

func CreateAuth(userID string, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAcc := redis.Client.Set(td.AccessUuid, userID, at.Sub(now)).Err()
	if errAcc != nil {
		return errAcc
	}

	errRfs := redis.Client.Set(td.RefreshUuid, userID, rt.Sub(now)).Err()
	if errRfs != nil {
		return errRfs
	}
	return nil
}

func DeleteAuth(givenUuid string) (int64, error) {
	deleted, err := redis.Client.Del(givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
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

func VerifyRefreshToken(tokenString string) (*jwt.Token, error) {
	refresh_secret := os.Getenv("REFRESH_SECRET")
	if len(refresh_secret) == 0 {
		refresh_secret = "cobain!!"
	}
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(refresh_secret), nil
	})
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}
	return parsedToken, nil
}

func TokenValid(tokenString string) error {
	token, err := VerifyToken(tokenString)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

type AccessDetails struct {
	AccessUuid string
	UserId     string
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

type RefreshDetails struct {
	RefreshUuid string
	UserId      string
}

func ExtractRefreshMetadata(tk *jwt.Token) (*RefreshDetails, error) {
	if _, ok := tk.Claims.(jwt.Claims); !ok && !tk.Valid {
		return nil, errors.New("Invalid claims")
	}
	claims, ok := tk.Claims.(jwt.MapClaims)
	if ok && tk.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string)
		if !ok {
			return nil, errors.New("Unable to parse refresh UUID")
		}
		userId, ok := claims["user_id"].(string)
		if !ok {
			return nil, errors.New("Unable to parse refresh UUID")
		}
		return &RefreshDetails{
			RefreshUuid: refreshUuid,
			UserId:      userId,
		}, nil
	}
	return nil, errors.New("Map error?")
}

func FetchAuth(authD *AccessDetails) (string, error) {
	userId, err := redis.Client.Get(authD.AccessUuid).Result()
	if err != nil {
		return "", err
	}
	return userId, nil
}
