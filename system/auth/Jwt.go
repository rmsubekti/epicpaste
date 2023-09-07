package auth

import (
	"epicpaste/system/helper"
	"epicpaste/system/model"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func CreateLoginSignature(u *model.User, expireDay int) (string, error) {
	mySigningKey := []byte(helper.GetEnv("EPIC_JWT_SECRET_KEY", "epic_secret"))

	// Create the Claims
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add((24 * time.Duration(expireDay)) * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ID:        u.UserName,
		Issuer:    u.Name,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(mySigningKey)
}

func ParseAndVerify(tokenString string) (any, error) {
	user := model.User{}
	token := strings.SplitN(tokenString, " ", 2)

	if (len(token) != 2) && (token[0] != "Bearer") {
		return nil, errors.New("incorrect format authorization header")
	}

	key, err := jwt.ParseWithClaims(token[1], &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(helper.GetEnv("EPIC_JWT_SECRET_KEY", "epic_secret")), nil
	})

	if claim, ok := key.Claims.(*jwt.RegisteredClaims); ok && key.Valid {
		user.UserName = claim.ID
		user.Name = claim.Issuer
		// user.Name = claim.Subject

		return user, err
	}
	return user, err
}
