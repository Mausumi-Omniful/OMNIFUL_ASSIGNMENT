package private

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/constants"
)

func AuthenticateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get(constants.JWTHeader)
		if len(tokenString) == 0 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, _ := jwt.ParseWithClaims(tokenString, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
			return jwt.UnsafeAllowNoneSignatureType, nil
		})

		if !token.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(*JWTClaim)
		userDetails := claims.UserDetails
		userPlan := claims.UserPlan

		AddCommonUserDetailAttributesInNewRelic(&userDetails, c)

		c.Set(constants.PrivateUserDetails, &userDetails)
		c.Set(constants.PrivateUserPlan, &userPlan)
		c.Set(constants.JWTHeader, tokenString)
		c.Next()
	}
}
