package gorm

import (
	"fmt"
	"ghotos/config"
	"time"

	"ghotos/util/logger/gorm_logger"

	gosql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

	return gorm.Open(mysql.Open(cfg.FormatDSN()), &gorm.Config{
		Logger: gorm_logger.New(),
	})

}
