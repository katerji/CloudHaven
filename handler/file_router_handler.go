package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/katerji/UserAuthKit/model"
	"github.com/katerji/UserAuthKit/service"
	"net/http"
	"strconv"
)

const FileRouterPath = "/s/:id"

func FileRouterHandler(c *gin.Context) {
	fileShareIDStr := c.Param("id")
	if fileShareIDStr == "" {
		sendBadRequestWithMessage(c, "Invalid file share id")
		return
	}
	fileShareID, err := strconv.Atoi(fileShareIDStr)
	if err != nil {
		sendBadRequestWithMessage(c, "Invalid file share id")
		return
	}
	fileShareInput := model.FileShareInput{
		ID: fileShareID,
	}
	url, err := service.GetFileService().GetFileShareURL(fileShareInput)
	if err != nil {
		sendBadRequestWithMessage(c, err.Error())
		return
	}
	go service.GetFileService().IncrementOpenRate(fileShareInput)
	c.Redirect(http.StatusMovedPermanently, url)
}
