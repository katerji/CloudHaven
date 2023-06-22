package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/katerji/UserAuthKit/model"
)

const UserInfoPath = "/user"

type UserInfoResponse struct {
	User model.UserOutput `json:"user"`
}

func UserInfoHandler(c *gin.Context) {
	user := getUserFromContext(c)
	response := UserInfoResponse{
		User: user.ToOutput(),
	}
	sendJSONResponse(c, response)
	return
}
