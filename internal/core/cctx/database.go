package cctx

import (
	"gorm.io/gorm"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/connection/sql"
)

// GetRelayDatabase get connection database `ralay`
func (c *Context) GetRelayDatabase() *gorm.DB {
	return sql.RelayDatabase
}
