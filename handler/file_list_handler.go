package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/katerji/UserAuthKit/model"
	"github.com/katerji/UserAuthKit/service"
)

const FileListPath = "/files"

type FileListResponse struct {
	Files []model.FileOutput `json:"files"`
}

func FileListHandler(c *gin.Context) {
	user := getUserFromContext(c)
	files, err := service.GetFileService().GetUserFiles(user.ID)
	if err != nil {
		sendErrorMessage(c, err.Error())
		return
	}
	fileOutputs := []model.FileOutput{}
	for _, file := range files {
		fileOutputs = append(fileOutputs, file.ToOutput())
	}
	response := FileListResponse{Files: fileOutputs}
	sendJSONResponse(c, response)
}
