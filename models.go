package main

type video struct {
	ID       int    `json:"id"`
	Path     string `json:"path"`
	Title    string `json:"title"`
	Category string `json:"category"`
}

type category struct {
	Category string `json:"category"`
}
