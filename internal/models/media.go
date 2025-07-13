package models

type MediaType string

const (
	Image MediaType = "Image"
	Video MediaType = "Video"
	Audio MediaType = "Audio"
	Book  MediaType = "Book"
)

type media struct {
	MediaTypes []MediaType
}

var MediaTypes = media{MediaTypes: []MediaType{Image, Video, Audio, Book}}
