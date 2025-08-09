package models

type MediaType string

const (
	MediaImage MediaType = "Image"
	MediaVideo MediaType = "Video"
	MediaAudio MediaType = "Audio"
	MediaBook  MediaType = "Book"
)

func (MediaType) Values() []string {
	return []string{
		string(MediaImage),
		string(MediaVideo),
		string(MediaAudio),
		string(MediaBook),
	}
}
