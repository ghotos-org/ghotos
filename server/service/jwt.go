package service

import (
	"errors"
	"fmt"
	"ghotos/model"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type Jwt struct {
}

func (j Jwt) CreateToken(user model.User) (model.Token, error) {
	var err error

	claims := jwt.MapClaims{}
	claims["id"] = user.UID
	//claims["exp"] = time.Now().Add(time.Hour * 2).Unix()
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwt := model.Token{}

	jwt.AccessToken, err = token.SignedString([]byte(os.Getenv("TOKEN_SECRET_KEY")))
	if err != nil {
		return jwt, err
	}

	return j.createRefreshToken(jwt)
}

func (Jwt) ValidateToken(accessToken string) (model.User, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("TOKEN_SECRET_KEY")), nil
	})

	user := model.User{}
	if err != nil {
		return user, err
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		user.UID = payload["id"].(string)

		return user, nil
	}

	return user, errors.New("invalid token")
}

func (j Jwt) ValidateRefreshToken(modelToken model.Token) (model.User, error) {
	token, err := jwt.Parse(modelToken.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("TOKEN_SECRET_KEY")), nil
	})

	user := model.User{}
	if err != nil {
		return user, err
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return user, errors.New("invalid token")
	}

	claims := jwt.MapClaims{}
	parser := jwt.Parser{}
	token, _, err = parser.ParseUnverified(payload["token"].(string), claims)
	if err != nil {
		return user, err
	}

	payload, ok = token.Claims.(jwt.MapClaims)
	if !ok {
		return user, errors.New("invalid token")
	}

	user.UID = payload["id"].(string)

	return user, nil
}

func (Jwt) createRefreshToken(token model.Token) (model.Token, error) {
	var err error

	claims := jwt.MapClaims{}
	claims["token"] = token.AccessToken
	claims["exp"] = time.Now().Add(time.Hour * 24 * 180).Unix()

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token.RefreshToken, err = refreshToken.SignedString([]byte(os.Getenv("TOKEN_SECRET_KEY")))
	if err != nil {
		return token, err
	}

	return token, nil
}
