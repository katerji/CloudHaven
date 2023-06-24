package main

import (
	"github.com/gin-gonic/gin"
	"github.com/katerji/UserAuthKit/crons"
	"github.com/katerji/UserAuthKit/db"
	"github.com/katerji/UserAuthKit/envs"
	"github.com/katerji/UserAuthKit/handler"
	"github.com/katerji/UserAuthKit/middleware"
	"github.com/robfig/cron"
)

func main() {
	initEnv()
	initDB()
	startCron()
	initWebServer()
}

func initDB() {
	client := db.GetDbInstance()
	err := client.Ping()
	if err != nil {
		panic(err)
	}
}

func initWebServer() {
	router := gin.Default()
	api := router.Group("/api")

	api.GET(handler.LandingPath, handler.LandingController)

	auth := api.Group("/auth")
	auth.POST(handler.RegisterPath, handler.RegisterHandler)
	auth.POST(handler.LoginPath, handler.LoginHandler)
	auth.POST(handler.RefreshTokenPath, handler.RefreshTokenHandler)

	api.Use(middleware.GetAuthMiddleware())

	api.GET(handler.UserInfoPath, handler.UserInfoHandler)

	api.GET(handler.FileListPath, handler.FileListHandler)

	fileGroup := api.Group("/file")
	fileGroup.POST(handler.FilePath, handler.FileUploadHandler)
	fileGroup.DELETE(handler.FilePath, handler.FileDeleteHandler)
	fileGroup.POST(handler.FileShareLinkPath, handler.FileShareLinkHandler)
	fileGroup.GET(handler.FileShareInfoPath, handler.FileShareInfoHandler)

	api.GET(handler.FileRouterPath, handler.FileRouterHandler)

	err := router.Run(":85")
	if err != nil {
		panic(err)
	}
}

func initEnv() {
	envs.InitEnv()
}

func startCron() {
	c := cron.New()
	c.AddFunc(crons.SyncFilesCronExpression, crons.SyncFiles)
	c.AddFunc(crons.SyncUsersCronExpression, crons.SyncUsers)
	c.AddFunc(crons.SyncOpenRatesCronExpression, crons.SyncOpenRates)
	c.Start()
}
