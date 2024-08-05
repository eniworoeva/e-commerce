package api

import (
	"e-commerce/internal/models"
	"e-commerce/internal/util"

	"github.com/gin-gonic/gin"
)

func (u *HTTPHandler) Logout(c *gin.Context) {
	blacklistTokens := &models.BlacklistTokens{}

	// Extract the token from the request header
	token, err := u.GetTokenFromContext(c)
	if err != nil {
		util.Response(c, "Error getting token from header", 500, err.Error(), nil)
		return
	}
	blacklistTokens.Token = token

	// Blacklist the token
	err = u.Repository.BlacklistToken(blacklistTokens)
	if err != nil {
		util.Response(c, "Could not blacklist token", 500, err.Error(), nil)
		return
	}
	util.Response(c, "Logged out successfully", 200, nil, nil)

}
