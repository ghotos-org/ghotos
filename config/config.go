package config

import (
	"log"
	"os"
	"time"

	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
)

type Conf struct {
	Debug  bool `env:"DEBUG,required"`
	Server serverConf
	Db     dbConf
	File   fileConf
}
type serverConf struct {
	Port         int           `env:"SERVER_PORT,required"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
}

type dbConf struct {
	Host     string `env:"DB_HOST,required"`
	Port     int    `env:"DB_PORT,required"`
	Username string `env:"DB_USER,required"`
	Password string `env:"DB_PASS,required"`
	DbName   string `env:"DB_NAME,required"`
}

type fileConf struct {
	DefaultMount  string `env:"MOUNT_DEFAULT,required"`
	MaxImageWidth int    `env:"MAX_IMAGE_WIDTH,required"`
}

func AppConfig(envFile string) *Conf {

	if envFile != "" {
		godotenv.Load(envFile)
	}

	if os.Getenv("DEBUG") == "" {
		godotenv.Load(".env")
	}

	var c Conf

	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}

	return &c
}
