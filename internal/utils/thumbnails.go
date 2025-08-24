package utils

import (
	"fmt"
	"goserv/internal/static/constant"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/disintegration/imaging"
)

const width = 400
const thumbTimestamp = 5

func CreateImageThumbnail(dir string, filename string, fileExt string) error {
	image, err := imaging.Open(filepath.Join(dir, filename+fileExt))
	if err != nil {
		log.Println("open error")
		return err
	}

	thumbnailDir := filepath.Join("thumbnail", filename[0:2], filename[2:4])
	if err := os.MkdirAll(thumbnailDir, 0755); err != nil {
		log.Printf("failed to create thumbnail dir: %v\n", err)
		return err
	}

	thumbnail := imaging.Resize(image, width, 0, imaging.Linear)
	return imaging.Save(thumbnail, filepath.Join(thumbnailDir, filename+constant.ThumbnailExt))
}

func ExctractVideoThumbnail(filename string, fileExt string) error {
	videoPath := filepath.Join("content", filename[0:2], filename[2:4], filename+fileExt)
	tmpPath := filepath.Join("tmp", filename+constant.ThumbnailExt)
	cmd := exec.Command("ffmpeg",
		"-y",
		"-i", videoPath,
		"-ss", fmt.Sprintf("%d", thumbTimestamp),
		"-vframes", "1",
		tmpPath,
	)

	err := cmd.Run()
	if err != nil {
		return err
	}

	err = CreateImageThumbnail("tmp", filename, constant.ThumbnailExt)
	if err != nil {
		log.Printf("Failed to create thumbnail for: %s\n", videoPath)
	}

	err = os.Remove(tmpPath)
	if err != nil {
		log.Printf("Failed to remove temp file: %s\n", tmpPath)
	}
	return nil
}
