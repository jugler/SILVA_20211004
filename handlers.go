package main

import (
	b64 "encoding/base64"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
)

const MAX_UPLOAD_SIZE = 200 * 1024 * 1024 //200MB -> byte

func showVideosPage(c *gin.Context) {
	videos := getAllVideosDB(dbConn)
	c.HTML(
		http.StatusOK,
		"index.html",
		gin.H{
			"title":  "Home Page",
			"videos": videos,
		},
	)
}

func showUploadPage(c *gin.Context) {
	categories := getAllCategoriesDB(dbConn)
	c.HTML(
		http.StatusOK,
		"upload.html",
		gin.H{
			"title":      "Upload Video Page",
			"categories": categories,
		},
	)
}

func upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if c.Request.ContentLength > MAX_UPLOAD_SIZE {
		//received file is larger than our Max upload size, return an error message
		c.HTML(
			http.StatusOK,
			"uploadSuccess.html",
			gin.H{
				"title": "Upload Video Failure",
				"error": "File exceeded max upload size of 200mb",
			},
		)
		return
	}

	newFileName := "videos/" + file.Filename
	//TODO: Add check to see if filename already exist, do not replace file
	//Use DB generated sequence as filename to avoid clash of filenames
	if err := c.SaveUploadedFile(file, newFileName); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	intVar, err := strconv.Atoi(c.Request.FormValue("category"))
	if err != nil {
		log.Fatal(err)
	}
	category := Category{ID: intVar}
	video := Video{Path: newFileName, Title: c.Request.FormValue("videoTitle"), Category: category}

	video.ID = saveVideoDB(dbConn, video)

	video.Thumbnail.Small = generateThumbnails(newFileName, 64)
	video.Thumbnail.SmallEncoded = b64.StdEncoding.EncodeToString([]byte(video.Thumbnail.Small))
	video.Thumbnail.Medium = generateThumbnails(newFileName, 128)
	video.Thumbnail.MediumEncoded = b64.StdEncoding.EncodeToString([]byte(video.Thumbnail.Medium))
	video.Thumbnail.Large = generateThumbnails(newFileName, 256)
	video.Thumbnail.LargeEncoded = b64.StdEncoding.EncodeToString([]byte(video.Thumbnail.Large))

	insertThumbnail(dbConn, video)

	c.HTML(
		http.StatusOK,
		"uploadSuccess.html",
		gin.H{
			"title": "Upload Video Success",
			"video": video,
		},
	)

}
