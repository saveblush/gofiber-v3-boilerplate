package utils

import (
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils/logger"
)

type ResolutionImageInfo struct {
	Resolution uint
	Width      uint
	Height     uint
}

type ThumbnailImageInfo struct {
	File   *os.File
	Width  uint
	Height uint
	Size   int64
}

// GenFileName gen file name
// e.g. 20250608-212456-c1d06f5c-fdb2-41ac-a72a-55942f1d15eb.png
func GenFileName(origName string) string {
	return fmt.Sprintf("%s-%s%s", Now().Format("20060102-150405"), UUID(), filepath.Ext(origName))
}

// GetResolutionImage get resolution image
func GetResolutionImage(img image.Image) *ResolutionImageInfo {
	bounds := img.Bounds()
	width := uint(bounds.Dx())
	height := uint(bounds.Dy())
	resolution := uint(width * height)

	return &ResolutionImageInfo{
		Resolution: resolution,
		Width:      width,
		Height:     height,
	}
}

// GetResolutionImageByFileHeader get resolution image by multipart file header
func GetResolutionImageByFileHeader(fh *multipart.FileHeader) *ResolutionImageInfo {
	file, err := fh.Open()
	if err != nil {
		return nil
	}
	defer file.Close()

	img, err := imaging.Decode(file)
	if err != nil {
		return nil
	}

	return GetResolutionImage(img)
}

// ResizeImage resize image
func ResizeImage(img image.Image, width, height uint) image.Image {
	if width == 0 && height == 0 {
		return img
	}

	if img.Bounds().Dx() <= 0 || img.Bounds().Dy() <= 0 {
		return img
	}

	return imaging.Resize(img, int(width), int(height), imaging.Lanczos)
}

// CreateTempImage create temporary file for image
func CreateTempImage(img image.Image, fileName string) (*os.File, error) {
	// create temporary file
	// e.g. go-createimage-3880138584.png
	tmp, err := os.CreateTemp("", fmt.Sprintf("go-createimage-*%s", filepath.Ext(fileName)))
	if err != nil {
		logger.Log.Errorf("os create temp error: %s", err)
		return nil, err
	}
	defer os.Remove(tmp.Name())

	logger.Log.Debugf("successfully created temporary file %s", tmp.Name())

	// copy the uploaded file to the temporary file
	err = imaging.Save(img, tmp.Name())
	if err != nil {
		logger.Log.Errorf("imaging save error: %s", err)
		return nil, err
	}

	return tmp, nil
}

// CreateThumbnail create thumbnail
func CreateThumbnailImage(file io.Reader, fileName string, newWidth, newHeight uint) (*ThumbnailImageInfo, error) {
	logger.Log.Debug("start create thumbnail...")

	img, err := imaging.Decode(file)
	if err != nil {
		logger.Log.Errorf("imaging decode error: %s", err)
		return nil, err
	}

	var origWidth uint
	var origHeight uint
	rst := GetResolutionImage(img)
	if rst != nil {
		origWidth = rst.Width
		origHeight = rst.Height
	}

	ratio := origWidth / origHeight
	if newHeight == 0 {
		if origWidth < newWidth {
			newWidth, newHeight = origWidth, origHeight
		} else {
			newHeight = newWidth / ratio
		}
	} else {
		if origHeight < newHeight {
			newHeight = origHeight
		} else {
			newWidth = newWidth * ratio
		}
		if origWidth < newWidth {
			newWidth, newHeight = origWidth, origHeight
		} else {
			newHeight = newWidth / ratio
		}
	}

	if !(newWidth >= origWidth && newHeight >= origHeight) {
		logger.Log.Debug("start resize image...")
		img = ResizeImage(img, newWidth, newHeight)
	}

	newFile, err := CreateTempImage(img, fileName)
	if err != nil {
		logger.Log.Errorf("create temp image error: %s", err)
		return nil, err
	}

	newSize, err := GetFileSize(newFile)
	if err != nil {
		logger.Log.Errorf("get file size error: %s", err)
		return nil, err
	}

	return &ThumbnailImageInfo{
		File:   newFile,
		Width:  newWidth,
		Height: newHeight,
		Size:   newSize,
	}, nil
}
