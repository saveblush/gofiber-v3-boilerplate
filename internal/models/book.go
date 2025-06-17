package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Author      string     `json:"author"`
	Display     *BookFiles `json:"display" gorm:"foreignKey:BookID"`
}

func (Book) TableName() string {
	return "books"
}

type BookFiles struct {
	gorm.Model
	BookID          uint   `json:"-" gorm:"index"`
	AttachType      string `json:"-" gorm:"type:varchar(15)"`
	FileOrigName    string `json:"file_orig_name"`
	FileName        string `json:"-" gorm:"type:varchar(250)"`
	FilePath        string `json:"-" gorm:"type:varchar(100)"`
	FileMimeType    string `json:"-" gorm:"type:varchar(100)"`
	FileWidth       uint   `json:"file_width"`
	FileHeight      uint   `json:"file_height"`
	FileSize        int64  `json:"file_size"`
	FileDescription string `json:"file_description,omitempty"`
	ThumbnailName   string `json:"-" gorm:"type:varchar(250)"`
	ThumbnailPath   string `json:"-" gorm:"type:varchar(100)"`
	ThumbnailWidth  uint   `json:"thumbnail_width"`
	ThumbnailHeight uint   `json:"thumbnail_height"`
	Url             string `json:"url" gorm:"-"`
	ThumbnailUrl    string `json:"thumbnail_url" gorm:"-"`
}

func (BookFiles) TableName() string {
	return "books_files"
}
