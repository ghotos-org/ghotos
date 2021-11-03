package gorm

import (
	"fmt"
	"ghotos/config"
	"os"
	"time"

	gosql "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// gorm.Model definition
type ModelUID struct {
	UID       string `gorm:"primaryKey"`
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func New(conf *config.Conf) (*gorm.DB, error) {

	cfg := &gosql.Config{
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%v:%v", conf.Db.Host, conf.Db.Port),
		DBName:               conf.Db.DbName,
		User:                 conf.Db.Username,
		Passwd:               conf.Db.Password,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	var logLevel logger.LogLevel
	if conf.Debug {
		logLevel = logger.Info
	} else {
		logLevel = logger.Error
	}

	zlog := log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC822}).With().Caller().Logger()
	newLogger := logger.New(
		&zlog, // IO.writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logLevel,    // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)
	// newlogger := logger.Default.LogMode(logLevel)

	return gorm.Open(mysql.Open(cfg.FormatDSN()), &gorm.Config{
		Logger: newLogger,
	})

}
