package token

import (
	"errors"
	"log"
	"tender-backend/config"
	"time"

	"github.com/golang-jwt/jwt"
)

type Tokens struct {
	AccessToken string `json:"access_token"`
}

func GenerateJWTToken(config *config.Config, userID string) *Tokens {
	accessToken := jwt.New(jwt.SigningMethodHS256)
	claims := accessToken.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix() // Token expires in 24 hours
	access, err := accessToken.SignedString([]byte(config.SecretKey))
	if err != nil {
		log.Fatal("error while generating access token : ", err)
	}

	return &Tokens{
		AccessToken: access,
	}
}

func ExtractClaim(secretKey string, tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, errors.New("parsing token:" + err.Error())
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
