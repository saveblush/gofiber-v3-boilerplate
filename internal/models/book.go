package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Name        string `json:"name"`
	Description string `json:"description"`
	Author      string `json:"author"`
}

func (Book) TableName() string {
	return "books"
}
