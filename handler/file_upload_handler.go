package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/katerji/UserAuthKit/model"
	"github.com/katerji/UserAuthKit/service"
	"io"
)

const FilePath = ""

type FileResponse struct {
	Success bool `json:"success"`
}

func FileUploadHandler(c *gin.Context) {
	err := c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		fmt.Println(err)
		sendBadRequestWithMessage(c, "File should be less than 32MB")
		return
	}

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Println(err)
		sendError(c)
		return
	}
	bytes := make([]byte, fileHeader.Size)
	_, err = io.ReadFull(file, bytes)
	if err != nil {
		sendError(c)
		return
	}
	user := getUserFromContext(c)
	fileInput := model.FileInput{
		Name:        fileHeader.Filename,
		OwnerID:     user.ID,
		Content:     bytes,
		ContentType: fileHeader.Header.Get("Content-Type"),
	}
	success := service.GetGcpService().CreateObject(fileInput)
	go service.GetFileService().SyncUserFiles(user.ID)
	sendJSONResponse(c, FileResponse{success})
}
