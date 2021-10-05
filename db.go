package main

import (
	"database/sql"
	b64 "encoding/base64"
	"fmt"
	"log"

	pq "github.com/lib/pq"
)

const (
	host                    = "db"
	port                    = 5432
	user                    = "postgres"
	password                = "password"
	dbname                  = "postgres"
	create_videos_table     = "CREATE TABLE videos (video_id SERIAL PRIMARY KEY, title VARCHAR NOT NULL, path VARCHAR NOT NULL, category_id INTEGER, CONSTRAINT fk_category FOREIGN KEY(category_id) REFERENCES category(category_id))"
	create_categories_table = "CREATE TABLE category (category_id SERIAL, category varchar(255), PRIMARY KEY(category_id))"
	create_thumbnail_table  = "CREATE TABLE thumbnails (thumbnail_id SERIAL PRIMARY KEY, image_small BYTEA NOT NULL, image_med BYTEA NOT NULL, image_large BYTEA NOT NULL, video_id INTEGER, CONSTRAINT fk_videos FOREIGN KEY(video_id) REFERENCES videos(video_id))"
)

func initDBConn() (db *sql.DB) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	createTable(create_categories_table, db)
	createTable(create_videos_table, db)
	createTable(create_thumbnail_table, db)
	insertInitRecordsCategory(db)

	log.Println("DB Connection Created Succesfully!")
	return db
}

func insertThumbnail(db *sql.DB, video Video) {
	insertThumbnailStatement := "INSERT INTO thumbnails (image_small, image_med, image_large, video_id) VALUES ($1, $2, $3, $4)"
	_, err := db.Exec(insertThumbnailStatement, video.Thumbnail.Small, video.Thumbnail.Medium, video.Thumbnail.Large, video.ID)
	if err != nil {
		log.Fatal(err)
	}
}

func insertInitRecordsCategory(db *sql.DB) {
	//check if init data is already there
	selectCategoryStatement := "SELECT category_id FROM category"

	row := db.QueryRow(selectCategoryStatement)
	var category_id int
	switch err := row.Scan(&category_id); err {
	case sql.ErrNoRows:
		log.Print("No categories found, init category table with default data")
		insertCategory(db, "Exercise")
		insertCategory(db, "Education")
		insertCategory(db, "Recipe")
	case nil:
		log.Println("Categories already found, no need to init table")
	default:
		log.Fatal(err)
	}
}

func insertCategory(db *sql.DB, category string) {
	initCategoriesStatement := "INSERT INTO category (category) VALUES ($1)"
	_, err := db.Exec(initCategoriesStatement, category)
	if err != nil {
		log.Fatal(err)
	}
}

func saveVideoDB(db *sql.DB, video Video) int {
	insertVideoStatement := "INSERT INTO videos (title, path, category_id) VALUES ($1, $2, $3) RETURNING video_id"
	video_id := 0
	err := db.QueryRow(insertVideoStatement, video.Title, video.Path, video.Category.ID).Scan(&video_id)
	if err != nil {
		log.Fatal(err)
	}

	return video_id
}

func createTable(tableSchema string, db *sql.DB) {
	_, err := db.Exec(tableSchema)

	if err, ok := err.(*pq.Error); ok {
		if err.Code == "42P07" { //table already exists, no issue
			log.Println("Table already exists, no need to create it")
		} else {
			log.Fatal(err)
		}
	}

}

func getAllCategoriesDB(db *sql.DB) []Category {
	getAllCategoriesStatement := "SELECT category_id, category FROM category"
	rows, err := db.Query(getAllCategoriesStatement)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var categories []Category
	for rows.Next() {
		var category_id int
		var category string

		err = rows.Scan(&category_id, &category)
		if err != nil {
			log.Fatal(err)
		}
		categories = append(categories, Category{ID: category_id, Name: category})

	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return categories

}

func getAllVideosDB(db *sql.DB) []Video {
	getAllVideosStatement := "SELECT videos.video_id, title, path, videos.category_id, category, thumbnails.image_small, thumbnails.image_med, thumbnails.image_large FROM videos, category, thumbnails WHERE videos.category_id = category.category_id AND videos.video_id = thumbnails.video_id"
	rows, err := db.Query(getAllVideosStatement)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var videos []Video
	for rows.Next() {
		var video_id int
		var title string
		var path string
		var category string
		var category_id int
		var thumbnail_small []byte
		var thumbnail_medium []byte
		var thumbnail_large []byte
		err = rows.Scan(&video_id, &title, &path, &category_id, &category, &thumbnail_small, &thumbnail_medium, &thumbnail_large)
		if err != nil {
			log.Fatal(err)
		}

		thumbnail := Thumbnail{
			Small:         thumbnail_small,
			Medium:        thumbnail_medium,
			Large:         thumbnail_large,
			SmallEncoded:  b64.StdEncoding.EncodeToString([]byte(thumbnail_small)),
			MediumEncoded: b64.StdEncoding.EncodeToString([]byte(thumbnail_medium)),
			LargeEncoded:  b64.StdEncoding.EncodeToString([]byte(thumbnail_large)),
		}
		videos = append(videos, Video{ID: video_id, Title: title, Path: path, Category: Category{ID: category_id, Name: category}, Thumbnail: thumbnail})

	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return videos
}
