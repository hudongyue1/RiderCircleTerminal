package util

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("a_secret")

type Claims struct {
	UserName string
	jwt.StandardClaims
}

func ReleaseToken(username string) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserName: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "CQUTeam",
			Subject:   "user token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, err
}


func ParseToken(tokenString string) (*jwt.Token,*Claims,error){
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey,nil
	})

	return token,claims,err
}

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9
// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9