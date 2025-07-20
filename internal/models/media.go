package models

type MediaType string

const (
	Image MediaType = "Image"
	Video MediaType = "Video"
	Audio MediaType = "Audio"
	Book  MediaType = "Book"
)

var MediaTypes = []MediaType{Image, Video, Audio, Book}
