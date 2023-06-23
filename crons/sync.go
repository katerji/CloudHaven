package crons

import (
	"github.com/katerji/UserAuthKit/service"
)

const (
	SyncFilesCronExpression = "@every 5m"
	SyncUsersCronExpression = "@every 5s"
)

func SyncFiles() {
	userService := service.GetUserService()
	users := userService.GetUsers()
	for userID, _ := range users {
		go service.GetFileService().SyncUserFiles(userID)
	}
}

func SyncUsers() {
	service.GetUserService().SyncUsers()
}
