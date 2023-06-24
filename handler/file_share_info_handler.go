package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/katerji/UserAuthKit/model"
	"github.com/katerji/UserAuthKit/service"
	"strconv"
)

const FileShareInfoPath = "/share/info/:file_id"

type FileShareInfoResponse struct {
	FileShares    []model.FileShareOutput `json:"file_shares"`
	TotalOpenRate int                     `json:"total_open_rate"`
}

func FileShareInfoHandler(c *gin.Context) {
	fileIDStr := c.Param("file_id")
	fileID, err := strconv.Atoi(fileIDStr)
	if err != nil || fileID == 0 {
		sendBadRequestWithMessage(c, "Invalid or missing file ID")
		return
	}
	user := getUserFromContext(c)
	fileShares, err := service.GetFileShareService().GetFileShares(fileID, user.ID)
	if err != nil {
		sendErrorMessage(c, err.Error())
		return
	}

	totalOpenRate := 0
	fileShareOutputs := []model.FileShareOutput{}
	for _, fileShare := range fileShares {
		totalOpenRate += fileShare.OpenRate
		fileShareOutputs = append(fileShareOutputs, fileShare.ToOutput())
	}
	response := FileShareInfoResponse{
		FileShares:    fileShareOutputs,
		TotalOpenRate: totalOpenRate,
	}

	sendJSONResponse(c, response)
}
