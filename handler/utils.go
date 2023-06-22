package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/katerji/UserAuthKit/model"
)

func getUserFromContext(c *gin.Context) model.User {
	user := c.MustGet("user")
	return user.(model.User)
}
