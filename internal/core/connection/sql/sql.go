package sql

import (
	"database/sql"
	"time"

	"github.com/saveblush/gofiber-v3-boilerplate/internal/core/generic"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	MysqlDriver    = "mysql"
	PostgresDriver = "postgres"
)

var (
	// RelayDatabase Database global variable database `relay`
	RelayDatabase = &gorm.DB{}
)

var (
	defaultMaxIdleConns = 10
	defaultMaxOpenConns = 15
	defaultMaxLifetime  = time.Hour
)

// gorm config
var defaultConfig = &gorm.Config{
	PrepareStmt:          true,
	DisableAutomaticPing: true,
	QueryFields:          true,
	Logger:               logger.Default.LogMode(logger.Error),
}

// Session session
type Session struct {
	Database *gorm.DB
	Conn     *sql.DB
}

// Configuration config mysql
type Configuration struct {
	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
	DriverName   string
	Charset      string
	MaxIdleConns int
	MaxOpenConns int
	MaxLifetime  time.Duration
}

// InitConnectionMysql open initialize a new db connection.
func InitConnection(cf *Configuration) (*Session, error) {
	var db *gorm.DB
	var err error

	if cf.DriverName == PostgresDriver {
		db, err = openPostgres(cf)
	} else {
		db, err = openMysql(cf)
	}
	if err != nil {
		return nil, err
	}

	// set config connection pool
	if generic.IsEmpty(cf.MaxIdleConns) {
		cf.MaxIdleConns = defaultMaxIdleConns
	}
	if generic.IsEmpty(cf.MaxOpenConns) {
		cf.MaxOpenConns = defaultMaxOpenConns
	}
	if generic.IsEmpty(cf.MaxLifetime) {
		cf.MaxLifetime = defaultMaxLifetime
	}

	// connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(cf.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cf.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cf.MaxLifetime)

	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}

	return &Session{Database: db, Conn: sqlDB}, nil
}

// CloseConnection close connection db
func CloseConnection(db *gorm.DB) error {
	c, err := db.DB()
	if err != nil {
		return err
	}

	err = c.Close()
	if err != nil {
		return err
	}

	return nil
}

// DebugRelayDatabase set debug sql
func DebugRelayDatabase() {
	RelayDatabase = RelayDatabase.Debug()
}
