package validate

import (
	"goserv/internal/static/constant"
	"goserv/internal/static/enum"
	"slices"
	"strings"
)

func IsValidFileType(filename string, mediaType enum.MediaType) bool {
	var validExts []string
	switch mediaType {
	case enum.MediaImage:
		validExts = constant.GetImageExts()
	case enum.MediaAudio:
		validExts = constant.GetAudioExts()
	case enum.MediaVideo:
		validExts = constant.GetVideoExts()
	default:
		return false
	}

	lastDotIndex := strings.LastIndex(filename, ".")
	if lastDotIndex == -1 {
		return false
	}

	fileExt := filename[lastDotIndex:]
	return slices.Contains(validExts, fileExt)
}
