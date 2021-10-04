package main

import (
	"github.com/gin-contrib/static"
)

func initRoutes() {
	router.GET("/", showVideosPage)
	router.GET("/upload.html", showUploadPage)
	router.POST("/upload", upload)
	router.Use(static.Serve("/videos", static.LocalFile("./videos", false)))
}
