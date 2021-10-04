package main

import (
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {

	router = gin.Default()
	gin.SetMode(gin.ReleaseMode)
	router.LoadHTMLGlob("templates/*")

	initRoutes()

	router.Run()
}
