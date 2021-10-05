package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var router *gin.Engine
var dbConn *sql.DB

func main() {
	dbConn = initDBConn()

	router = gin.Default()
	gin.SetMode(gin.DebugMode)
	router.LoadHTMLGlob("templates/*")

	initRoutes()

	router.Run()
	defer dbConn.Close()
}
