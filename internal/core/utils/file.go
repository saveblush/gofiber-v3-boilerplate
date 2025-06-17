package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils/logger"
)

// CreateTempImage create temporary file
func CreateTempFile(file io.Reader, fileName string) (*os.File, error) {
	// create temporary file
	// e.g. go-createtmpfile-3880138584.png
	tmp, err := os.CreateTemp("", fmt.Sprintf("go-createtmpfile-*%s", filepath.Ext(fileName)))
	if err != nil {
		logger.Log.Errorf("os create temp error: %s", err)
		return nil, err
	}
	defer os.Remove(tmp.Name())

	logger.Log.Debugf("successfully created temporary file %s", tmp.Name())

	buf := &bytes.Buffer{}
	buf.ReadFrom(file)
	b := buf.Bytes()
	os.WriteFile(tmp.Name(), b, os.ModeAppend)

	return tmp, nil
}

// GetFileSize get file size
func GetFileSize(r io.Reader) (int64, error) {
	if f, ok := r.(*os.File); ok {
		fi, err := f.Stat()
		if err != nil {
			return 0, err
		}
		return fi.Size(), nil
	}
	return 0, errors.New("not an *os.File")
}
