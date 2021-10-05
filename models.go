package main

type Video struct {
	ID        int       `json:"id"`
	Path      string    `json:"path"`
	Title     string    `json:"title"`
	Category  Category  `json:"category"`
	Thumbnail Thumbnail `json:"thumbnail"`
}

type Thumbnail struct {
	Small         []byte `json:"64x64"`
	Medium        []byte `json:"128x128"`
	Large         []byte `json:"256x256"`
	SmallEncoded  string `json:"64x64encoded"`
	MediumEncoded string `json:"64x64encoded"`
	LargeEncoded  string `json:"64x64encoded"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"category"`
}
