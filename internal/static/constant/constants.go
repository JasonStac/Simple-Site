package constant

const ThumbnailExt = ".jpg"

func GetImageExts() []string {
	return []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".avif"}
}

func GetVideoExts() []string {
	return []string{".mp4", ".webm"}
}

func GetAudioExts() []string {
	return []string{".m4a", ".mp3", ".wav", ".ogg", ".opus", ".flac"}
}
