package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func showVideosPage(c *gin.Context) {
	videos := getAllVideos()
	// Call the HTML method of the Context to render a template
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"index.html",
		// Pass the data that the page uses
		gin.H{
			"title":  "Home Page",
			"videos": videos,
		},
	)
}

func showUploadPage(c *gin.Context) {
	categories := getAllCategories()
	c.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"upload.html",
		// Pass the data that the page uses
		gin.H{
			"title":      "Upload Video Page",
			"categories": categories,
		},
	)
}

func upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("file err : %s", err.Error()))
		return
	}
	filename := header.Filename
	out, err := os.Create("videos/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}
	video := video{ID: 3, Path: filename, Title: c.Request.FormValue("videoTitle"), Category: c.Request.FormValue(("category"))}

	c.HTML(
		http.StatusOK,
		"uploadSuccess.html",
		gin.H{
			"title":         "Upload Video Success",
			"VideoPath":     video.Path,
			"VideoTitle":    video.Title,
			"VideoID":       video.ID,
			"VideoCategory": video.Category,
		},
	)

}

func getAllCategories() []category {
	var categoryList = []category{
		{Category: "Exercise"},
		{Category: "Education"},
		{Category: "Recipe"},
	}
	return categoryList
}

func getAllVideos() []video {
	var videoList = []video{
		{ID: 1, Path: "1.mp4", Title: "Video 1", Category: "Video 1 category"},
		{ID: 2, Path: "2.mp4", Title: "Video 2", Category: "Video 2 category"},
	}
	return videoList
}
