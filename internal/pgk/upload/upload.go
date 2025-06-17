package upload

import (
	"io"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/connection/minio"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/ftpupload"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils/logger"
)

// service interface
type Service interface {
	Upload(bucketName, prefix, objectName string, object io.Reader, objectSize int64) error
}

type service struct {
	config *config.Configs
	result *config.ReturnResult
	minio  minio.Client
}

func NewService() Service {
	return &service{
		config: config.CF,
		result: config.RR,
		minio:  minio.New(),
	}
}

// Upload upload file
// มองเคส driver minio หรือ ftpupload
func (s *service) Upload(bucketName, prefix, objectName string, object io.Reader, objectSize int64) error {
	if bucketName == "" {
		bucketName = s.config.Storage.BucketName
	}

	if s.config.Storage.DriverName == config.MinioDriver {
		err := s.minio.Upload(bucketName, prefix, objectName, object, objectSize)
		if err != nil {
			logger.Log.Errorf("minio upload error: %s", err)
			return err
		}
	} else {
		err := ftpupload.Upload(bucketName, prefix, objectName, object, &ftpupload.Auth{
			Host:     s.config.Storage.Host,
			UserName: s.config.Storage.User,
			Password: s.config.Storage.Password,
		})
		if err != nil {
			logger.Log.Errorf("ftpupload upload error: %s", err)
			return err
		}
	}

	return nil
}
