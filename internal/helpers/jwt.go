package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type signedDetails struct {
	Uid      string `json:"uid,omitempty"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateTokens(uid string, username string) (signedToken string, signedRefreshToken string, err error) {
	claims := signedDetails{
		Uid:      uid,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(1)).Unix(),
		},
	}

	refreshClaim := signedDetails{
		Uid:      uid,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	SECRET_KEY := os.Getenv("secret_key")
	access_token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}

	refresh_token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}
	return access_token, refresh_token, nil
}

func ValidateToken(signedToken string) (claims *signedDetails, err error) {
	SECRET_KEY := os.Getenv("secret_key")
	token, err := jwt.ParseWithClaims(
		signedToken,
		&signedDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*signedDetails)
	if !ok {
		return nil, fmt.Errorf("please login to get on or via refresh_token")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, fmt.Errorf("token expired")
	}
	return claims, nil
}
