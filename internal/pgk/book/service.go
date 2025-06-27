package book

import (
	"fmt"
	"io"
	"mime"
	"net/url"
	"path"

	"github.com/jinzhu/copier"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/cctx"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/config"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/connection/cache"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/generic"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/utils/logger"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/models"
	"github.com/saveblush/gofiber-v3-boilerplate/internal/pgk/upload"
)

var (
	keyCache        = "book"
	pathFileDisplay = "/book"
)

// service interface
type Service interface {
	Find(c *cctx.Context, req *Request) (interface{}, error)
	FindAll(c *cctx.Context, req *Request) (interface{}, error)
	FindAllPage(c *cctx.Context, req *RequestPage) (interface{}, error)
	FindByID(c *cctx.Context, req *RequestID) (interface{}, error)
	Create(c *cctx.Context, req *RequestCreate) (interface{}, error)
	Update(c *cctx.Context, req *RequestUpdate) (interface{}, error)
	Delete(c *cctx.Context, req *RequestID) error
	Script(c *cctx.Context) error
}

type service struct {
	config     *config.Configs
	repository Repository
	cache      cache.Client
	upload     upload.Service
}

func NewService() Service {
	return &service{
		config:     config.CF,
		repository: NewRepository(),
		cache:      cache.New(),
		upload:     upload.NewService(),
	}
}

// Find find
func (s *service) Find(c *cctx.Context, req *Request) (interface{}, error) {
	res, err := s.repository.Find(c.GetDatabase(), req)
	if err != nil {
		return nil, err
	}

	// url display
	if !generic.IsEmpty(res.Display) {
		url, err := url.Parse(s.config.Storage.URL)
		if err != nil {
			return nil, err
		}
		if res.Display.FileName != "" {
			res.Display.Url = url.JoinPath(res.Display.FilePath, res.Display.FileName).String()
		}
		if res.Display.ThumbnailName != "" {
			res.Display.ThumbnailUrl = url.JoinPath(res.Display.ThumbnailPath, res.Display.ThumbnailName).String()
		}
	}

	return res, nil
}

// FindAll find all
func (s *service) FindAll(c *cctx.Context, req *Request) (interface{}, error) {
	res, err := s.repository.FindAll(c.GetDatabase(), req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// FindAll find all page
func (s *service) FindAllPage(c *cctx.Context, req *RequestPage) (interface{}, error) {
	res, err := s.repository.FindAllPage(c.GetDatabase(), req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// FindByID find by id
func (s *service) FindByID(c *cctx.Context, req *RequestID) (interface{}, error) {
	key := fmt.Sprintf("%s-%d", keyCache, req.ID)
	res := &models.Book{}
	err := s.cache.Get(key, res)

	// ถ้าไม่เจอ cache
	if err != nil {
		fetch, err := s.Find(c, &Request{ID: req.ID})
		if err != nil {
			return nil, err
		}

		if !generic.IsEmpty(fetch) {
			copier.Copy(res, fetch)

			// เก็บใน cache
			_ = s.cache.Set(key, res, s.config.Cache.ExprieTime.Default)
		}
	}

	return res, nil
}

// Create create
func (s *service) Create(c *cctx.Context, req *RequestCreate) (interface{}, error) {
	data := &models.Book{
		Name: req.Name,
	}

	err := s.repository.Create(c.GetDatabase(), data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Update update
func (s *service) Update(c *cctx.Context, req *RequestUpdate) (interface{}, error) {
	data := &models.Book{}
	copier.Copy(data, &req)

	err := s.repository.Update(c.GetDatabase(), data, data)
	if err != nil {
		return nil, err
	}

	// แนบรูปโปรไฟล์
	_ = s.AttachDisplay(c, &RequestAttach{req.RequestID})

	// เคลีย cache
	key := fmt.Sprintf("%s-%d", keyCache, req.ID)
	_ = s.cache.Delete(key)

	return data, nil
}

// AttachDisplay attach display
func (s *service) AttachDisplay(c *cctx.Context, req *RequestAttach) error {
	files, err := c.FormFile("display_attachment")
	if err != nil {
		return err
	}

	mainPath := pathFileDisplay
	fileOrigName := files.Filename
	newFileName := utils.GenFileName(fileOrigName)
	ext := path.Ext(fileOrigName)
	mimeType := mime.TypeByExtension(ext)

	// get width, height for file image
	var origWidth uint
	var origHeight uint
	rst := utils.GetResolutionImageByFileHeader(files)
	if rst != nil {
		origWidth = rst.Width
		origHeight = rst.Height
	}

	file, err := files.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	newPath := mainPath
	//err = s.upload.Upload(s.config.Storage.BucketName, newPath, newFileName, file, files.Size)	// ตัวอย่าง กรณีมีค่า file size
	err = s.upload.Upload(s.config.Storage.BucketName, newPath, newFileName, file, 0) // ตัวอย่าง กรณีไม่มีค่า file size
	if err != nil {
		logger.Log.Errorf("upload error: %s", err)
		return err
	}

	// รีเซ็ตตำแหน่งของไฟล์กลับไปที่เริ่มต้น เพื่อเอาไปใช้กับ thumbnail ต่อ
	_, _ = file.Seek(0, io.SeekStart)

	// thumbnail
	var thumbnailName string
	var thumbnailPath string
	var thumbnailWidth uint
	var thumbnailHeight uint
	tmb, err := utils.CreateThumbnailImage(file, fileOrigName, s.config.Image.ThumbnailWidth, s.config.Image.ThumbnailHeight)
	if err == nil {
		defer tmb.File.Close()

		newPath := path.Join(mainPath, "/thumbnail")
		err := s.upload.Upload(s.config.Storage.BucketName, newPath, newFileName, tmb.File, tmb.Size)
		if err == nil {
			thumbnailName = newFileName
			thumbnailPath = newPath
			thumbnailWidth = tmb.Width
			thumbnailHeight = tmb.Height
		}
	}

	// db
	// เคลียรูปเดิมออก
	err = s.repository.DeleteFile(c.GetDatabase(), &RequestAttach{req.RequestID})
	if err != nil {
		return err
	}

	// insert
	data := &models.BookFiles{
		BookID:          req.ID,
		AttachType:      "1", // type display
		FileOrigName:    fileOrigName,
		FileName:        newFileName,
		FilePath:        mainPath,
		FileMimeType:    mimeType,
		FileWidth:       origWidth,
		FileHeight:      origHeight,
		FileSize:        files.Size,
		ThumbnailName:   thumbnailName,
		ThumbnailPath:   thumbnailPath,
		ThumbnailWidth:  thumbnailWidth,
		ThumbnailHeight: thumbnailHeight,
	}
	err = s.repository.Create(c.GetDatabase(), data)
	if err != nil {
		return err
	}

	return nil
}

// Delete delete
func (s *service) Delete(c *cctx.Context, req *RequestID) error {
	data := &models.Book{}
	copier.Copy(data, &req)

	err := s.repository.Delete(c.GetDatabase(), data)
	if err != nil {
		return err
	}

	// เคลีย cache
	key := fmt.Sprintf("%s-%d", keyCache, req.ID)
	_ = s.cache.Delete(key)

	return nil
}

// test cronjob
func (s *service) Script(c *cctx.Context) error {
	req := &RequestUpdate{
		RequestID: RequestID{ID: 2},
		Name:      "อัพเดทจาก cronjob",
	}
	data := &models.Book{}
	copier.Copy(data, req)

	err := s.repository.Update(c.GetDatabase(), data, data)
	if err != nil {
		return err
	}

	return nil
}
