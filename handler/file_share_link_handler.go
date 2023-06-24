package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/katerji/UserAuthKit/gcp"
	"github.com/katerji/UserAuthKit/model"
	"github.com/katerji/UserAuthKit/service"
)

const FileShareLinkPath = "/share"

type FileShareLinkRequest struct {
	Name string `json:"name"`
}

type FileShareLinkResponse struct {
	FileShareID int `json:"file_share_id"`
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
	file, err := service.GetFileService().GetFile(fileInput)
	if err != nil {
		sendBadRequestWithMessage(c, err.Error())
		return
	}
	url, err := service.GetGcpService().SignObject(fileInput)
	if err != nil {
		sendBadRequestWithMessage(c, err.Error())
		return
	}
	fileShareInput := model.FileShareInput{
		FileID:    file.ID,
		URL:       url,
		ExpiresAt: gcp.GetDefaultSignExpiry(),
	}
	insertId, err := service.GetFileShareService().Insert(fileShareInput)
	if err != nil {
		sendErrorMessage(c, err.Error())
		return
	}
	response := FileShareLinkResponse{
		FileShareID: insertId,
	}
	fileShare := model.FileShare{
		ID:        insertId,
		FileID:    file.ID,
		URL:       url,
		OpenRate:  0,
		ExpiresAt: fileShareInput.ExpiresAt,
	}
	if err = service.GetFileShareService().SetCache(fileShare); err != nil {
		sendErrorMessage(c, err.Error())
		return
	}
	sendJSONResponse(c, response)
}
