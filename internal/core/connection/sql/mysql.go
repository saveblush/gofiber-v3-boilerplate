package sql

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// openMysql open initialize a new db connection.
func openMysql(cf *Configuration) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cf.Username,
		cf.Password,
		cf.Host,
		cf.Port,
		cf.DatabaseName,
		cf.Charset,
	)

	return gorm.Open(mysql.Open(dsn), defaultConfig)
}
