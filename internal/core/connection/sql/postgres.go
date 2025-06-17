package sql

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// openPostgres open initialize a new db connection.
func openPostgres(cf *Configuration) (*gorm.DB, error) {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s TimeZone=%s sslmode=disable",
		cf.User,
		cf.Password,
		cf.Host,
		cf.Port,
		cf.DatabaseName,
		cf.Timezone,
	)

	return gorm.Open(postgres.New(postgres.Config{
		DSN:                 dsn,
		WithoutQuotingCheck: true,
	}), defaultConfig)
}
