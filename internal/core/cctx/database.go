package cctx

import (
	"gorm.io/gorm"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/connection/sql"
)

// GetDatabase get connection database`
func (c *Context) GetDatabase() *gorm.DB {
	return sql.Database
}
