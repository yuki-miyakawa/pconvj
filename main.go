package main

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func convertImage(src, dest string, quality int) {
	srcImg, err := os.Open(src)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer srcImg.Close()

	img, format, err := image.Decode(srcImg)
	if err != nil {
		log.Fatalf("Error decoding image: %v file : %v", err, src)
	}

	outputImgWriter, err := os.Create(dest)
	if err != nil {
		log.Fatalf("Error creating file: %v", err)
	}
	defer outputImgWriter.Close()

	if format == "png" {
		if err := jpeg.Encode(outputImgWriter, img, &jpeg.Options{
			Quality: quality,
		}); err != nil {
			log.Fatalf("Error encoding image: %v", err)
		}
	}
}

func main() {
	original := "./original"

	err := filepath.Walk(original, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".png" {
			extensionIndex := strings.LastIndex(path, ".png")
			if extensionIndex == -1 {
				return nil
			}
			dest := fmt.Sprintf(
				"%s.jpg",
				strings.Replace(path[:extensionIndex], "original", "output", 1),
			)
			destDir := filepath.Dir(dest)
			if _, err := os.Stat(destDir); os.IsNotExist(err) {
				if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
					return fmt.Errorf("failed to create directory: %v", err)
				}
			}
			convertImage(path, dest, 10)
		}
		return nil
	})
	if err != nil {
		log.Fatalf("Error walking the path %q: %v\n", original, err)
	}
}
