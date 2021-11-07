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
	SMTP   smtpConf
}
type serverConf struct {
	Port         int           `env:"SERVER_PORT,required"`
	PublicUrl    string        `env:"SERVER_PUBLIC_URL,required"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
	TokenKey     string        `env:"TOKEN_SECRET_KEY,required"`
}

type dbConf struct {
	Host     string `env:"DB_HOST,required"`
	Port     int    `env:"DB_PORT,required"`
	Username string `env:"DB_USER,required"`
	Password string `env:"DB_PASS,required"`
	DbName   string `env:"DB_NAME,required"`
}

type smtpConf struct {
	Host     string `env:"SMTP_HOST"`
	User     string `env:"SMTP_USER"`
	Password string `env:"SMTP_PASSWORD"`
	Port     int    `env:"SMTP_PORT"`
	Sender   string `env:"SMTP_SENDER"`
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
