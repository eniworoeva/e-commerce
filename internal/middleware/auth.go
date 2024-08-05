package middleware

import (
	"e-commerce/internal/models"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AuthorizeSeller(findSellerByEmail func(string) (*models.Seller, error), tokenInBlacklist func(*string) bool) gin.HandlerFunc {
	return func(c *gin.Context) {

		var seller *models.Seller
		var errors error
		secret := os.Getenv("JWT_SECRET")
		accToken := GetTokenFromHeader(c)
		accessToken, accessClaims, err := AuthorizeToken(&accToken, &secret)
		if err != nil {
			log.Printf("authorize access token errors: %s\n", err.Error())
			RespondAndAbort(c, "", http.StatusUnauthorized, nil, []string{"unauthorized"})
			return
		}

		if tokenInBlacklist(&accessToken.Raw) || IsTokenExpired(accessClaims) {
			RespondAndAbort(c, "", http.StatusUnauthorized, nil, []string{"unauthorized"})
			return
		}

		if email, ok := accessClaims["user_email"].(string); ok {
			if seller, errors = findSellerByEmail(email); errors != nil {
				log.Printf("find Seller by email errors: %v\n", err)
				RespondAndAbort(c, "", http.StatusNotFound, nil, []string{"Seller not found"})
				return
			}
		} else {
			log.Printf("Seller email is not string\n")
			RespondAndAbort(c, "", http.StatusInternalServerError, nil, []string{"internal server errors"})
			return
		}

		// set the Seller and token as context parameters.
		c.Set("Seller", seller)
		c.Set("access_token", accessToken.Raw)

		// calling next handler
		c.Next()
	}
}

func AuthorizeUser(findUserByEmail func(string) (*models.User, error), tokenInBlacklist func(*string) bool) gin.HandlerFunc {
	return func(c *gin.Context) {

		var user *models.User
		var errors error
		secret := os.Getenv("JWT_SECRET")
		accToken := GetTokenFromHeader(c)
		accessToken, accessClaims, err := AuthorizeToken(&accToken, &secret)
		if err != nil {
			log.Printf("authorize access token errors: %s\n", err.Error())
			RespondAndAbort(c, "", http.StatusUnauthorized, nil, []string{"unauthorized"})
			return
		}

		if tokenInBlacklist(&accessToken.Raw) || IsTokenExpired(accessClaims) {
			RespondAndAbort(c, "", http.StatusUnauthorized, nil, []string{"unauthorized"})
			return
		}

		if email, ok := accessClaims["user_email"].(string); ok {
			if user, errors = findUserByEmail(email); errors != nil {
				log.Printf("find user by email errors: %v\n", err)
				RespondAndAbort(c, "", http.StatusNotFound, nil, []string{"user not found"})
				return
			}
		} else {
			log.Printf("user email is not string\n")
			RespondAndAbort(c, "", http.StatusInternalServerError, nil, []string{"internal server errors"})
			return
		}

		// set the user and token as context parameters.
		c.Set("user", user)
		c.Set("access_token", accessToken.Raw)

		// calling next handler
		c.Next()
	}
}
