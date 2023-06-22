package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/katerji/UserAuthKit/model"
	"github.com/katerji/UserAuthKit/service"
)

type FileDeleteRequest struct {
	Name string `json:"name"`
}

func FileDeleteHandler(c *gin.Context) {
	request := &FileDeleteRequest{}
	err := c.BindJSON(request)
	if err != nil {
		sendBadRequest(c)
		return
	}
	if request.Name == "" {
		sendBadRequestWithMessage(c, "Name is required")
		return
	}
	user := c.MustGet("user").(model.User)
	fileInput := model.FileInput{
		Name:    request.Name,
		OwnerID: user.ID,
	}
	success := service.GetGcpService().DeleteObject(fileInput)
	sendJSONResponse(c, FileResponse{success})
}
