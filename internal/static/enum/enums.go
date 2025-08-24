package enum

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

type TagType string

const (
	TagGeneral TagType = "General"
	TagPeople  TagType = "People"
)

func (TagType) Values() []string {
	return []string{
		string(TagGeneral),
		string(TagPeople),
	}
}
