package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/katerji/UserAuthKit/model"
	"github.com/katerji/UserAuthKit/service"
)

const FileShareLinkPath = "/file/share"

type FileShareLinkRequest struct {
	Name string `json:"name"`
}

type FileShareLinkResponse struct {
	URL string `json:"url"`
}

func FileShareLinkHandler(c *gin.Context) {
	request := &FileShareLinkRequest{}
	err := c.BindJSON(request)
	if err != nil {
		sendBadRequest(c)
		return
	}
	if request.Name == "" {
		sendBadRequest(c)
		return
	}
	user := getUserFromContext(c)
	fileInput := model.FileInput{
		Name:    request.Name,
		OwnerID: user.ID,
	}
	url, err := service.GetGcpService().SignObject(fileInput)
	if err != nil {
		sendBadRequestWithMessage(c, err.Error())
		return
	}
	response := FileShareLinkResponse{
		URL: url,
	}
	sendJSONResponse(c, response)
}
