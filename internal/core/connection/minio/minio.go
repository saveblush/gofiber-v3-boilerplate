package minio

import (
	"context"
	"io"
	"mime"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils/logger"
)

var (
	connection = &minio.Client{}
	ctx        = context.Background()
)

// Configuration config minio connection
type Configuration struct {
	Host     string
	UserName string
	Password string
	Secure   bool
}

type client struct {
	client *minio.Client
}

type Client interface {
	Upload(bucketName, prefix, objectName string, object io.Reader, objectSize int64) error
}

// Init init a new minio connection
func Init(cf *Configuration) (err error) {
	connection, err = minio.New(cf.Host, &minio.Options{
		Creds:  credentials.NewStaticV4(cf.UserName, cf.Password, ""),
		Secure: cf.Secure,
	})
	if err != nil {
		return err
	}

	return nil
}

// New new client connection
func New() Client {
	return &client{
		client: connection,
	}
}

func (c *client) createBucket(bucketName string) error {
	// check exist bucket
	exists, err := c.client.BucketExists(ctx, bucketName)
	if err == nil && exists {
		logger.Log.Debugf("bucket %s already exists", bucketName)
	} else {
		// create a new bucket
		err = c.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			logger.Log.Errorf("make bucket error: %s", err)
			return err
		}
		logger.Log.Debugf("successfully created bucket %s", bucketName)
	}

	return nil
}

func (c *client) Upload(bucketName, prefix, objectName string, object io.Reader, objectSize int64) error {
	err := c.createBucket(bucketName)
	if err != nil {
		return err
	}

	if prefix != "" {
		objectName = filepath.Join(prefix, objectName)
	}

	if objectSize == 0 {
		tmp, err := utils.CreateTempFile(object, objectName)
		if err != nil {
			logger.Log.Errorf("create temp file error: %s", err)
			return err
		}
		defer tmp.Close()

		object = tmp
		objectSize, err = utils.GetFileSize(object)
		if err != nil {
			logger.Log.Errorf("get file size error: %s", err)
			return err
		}
	}

	mimeType := mime.TypeByExtension(filepath.Ext(objectName))
	info, err := c.client.PutObject(ctx, bucketName, objectName, object, objectSize, minio.PutObjectOptions{ContentType: mimeType})
	if err != nil {
		logger.Log.Errorf("put object error: %s", err)
		return err
	}
	logger.Log.Debugf("successfully uploaded %s of size %d", objectName, info.Size)

	return nil
}
