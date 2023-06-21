package crons

import (
	"github.com/katerji/UserAuthKit/service"
)

// SyncFilesCronExpression runs every 5 minutes
const SyncFilesCronExpression = "0 0,5,10,15,20,25,30,35,40,45,50,55 0 ? * * *"

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
