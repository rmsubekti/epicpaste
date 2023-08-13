package middleware

import (
	"epicpaste/system/auth"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/golang-jwt/jwt/v4"

	"github.com/gin-gonic/gin"
)

// it will create context user from token/authorization header or it will set to nil if user not logged in
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		var user any

		tokenAuth := c.GetHeader("Authorization")

		if len(tokenAuth) < 1 {
			token := sessions.Default(c).Get("token")
			if token != nil {
				tokenAuth = token.(string)
			}
		}

		if len(tokenAuth) > 1 {
			if user, err = auth.ParseAndVerify(tokenAuth); err != nil {
				if strings.Contains(err.Error(), jwt.ErrTokenExpired.Error()) {
					user = nil
				} else {
					c.AbortWithStatusJSON(http.StatusExpectationFailed, gin.H{"error": err.Error()})
					return
				}
			}
			c.Set("user", user)
		} else {
			c.Set("user", nil)
		}
		c.Next()
	}
}
