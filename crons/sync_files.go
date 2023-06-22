package crons

import (
	"github.com/katerji/UserAuthKit/service"
)

// SyncFilesCronExpression runs every 5 minutes
const SyncFilesCronExpression = "@every 5m"

func SyncFiles() {
	userService := service.GetUserService()
	users := userService.GetUsers()
	for userID, _ := range users {
		files, ok := service.GetGcpService().ListUserObjects(userID)
		if !ok {
			continue
		}
		service.GetFileService().UpsertUserFiles(files)
	}
}
