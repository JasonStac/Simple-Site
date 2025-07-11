package models

type MediaType struct {
	Value Media
	Label string
}

var MediaTypes = []MediaType{
	{Image, "Image"},
	{Video, "Video"},
	{Audio, "Audio"},
	{Book, "Book"},
}
