package main

import (
	"context"
	"embed"
	"fmt"
	dbConn "ghotos/adapter/gorm"
	"ghotos/config"
	"ghotos/server/app"
	"ghotos/server/router"
	"ghotos/util/tools"
	vr "ghotos/util/validator"
	"strconv"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ghotos/util/logger/goose_logger"

	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	//goose "github.com/pressly/goose"
	_ "ghotos/migrations"

	goose_v3 "github.com/pressly/goose/v3"
)

type ArgOptions struct {
	Env     string `short:"e" long:"env" description:"Environment File (default: .env)"`
	Migrate string `short:"m" long:"migrate" description:"DB mirgrate tool"  choice:"up" choice:"down" choice:"status" choice:"version" choice:"reset" choice:"up-by-one" choice:"up-to" choice:"down-to"`
	PWHash  string `short:"p" long:"password" description:"Password Hash"`
}

//go:embed migrations/*
var embedMigrations embed.FS

func migrate(db *gorm.DB, migrateCMD string, migrateArgs []string) (doExit bool) {
	appDb, err := db.DB()
	if err != nil {
		log.Fatal(err)
		return
	}
	goose_v3.SetBaseFS(embedMigrations)

	goose_v3.SetLogger(goose_logger.New())

	if err := goose_v3.SetDialect("mysql"); err != nil {
		log.Fatal(err)
	}

	if migrateCMD == "up" {
		if err := goose_v3.Up(appDb, "migrations"); err != nil {
			log.Fatal(err)
		}
	}
	if migrateCMD == "down" {
		if err := goose_v3.Down(appDb, "migrations"); err != nil {
			log.Fatal(err)
		}
	}
	if migrateCMD == "status" {
		if err := goose_v3.Status(appDb, "migrations"); err != nil {
			log.Fatal(err)
		}
	}
	if migrateCMD == "version" {
		if err := goose_v3.Version(appDb, "migrations"); err != nil {
			log.Fatal(err)
		}
	}
	if migrateCMD == "reset" {
		if err := goose_v3.Reset(appDb, "migrations"); err != nil {
			log.Fatal(err)
		}
	}
	if migrateCMD == "up-by-one" {
		if err := goose_v3.Reset(appDb, "migrations"); err != nil {
			log.Fatal(err)
		}
	}
	if migrateCMD == "up-to" || migrateCMD == "down-to" {
		if len(migrateArgs) < 2 {
			log.Fatal("Version missing")
			return true
		}
		version, err := strconv.ParseInt(migrateArgs[1], 10, 64)
		if err != nil {
			log.WithError(err).Fatal("incorrect verison format")
		}

		if migrateCMD == "up-to" {
			if err := goose_v3.UpTo(appDb, "migrations", version); err != nil {
				log.Fatal(err)
			}
		} else {
			if err := goose_v3.DownTo(appDb, "migrations", version); err != nil {
				log.Fatal(err)
			}
		}

	}
	if migrateCMD != "" {
		return true
	}

	if err := goose_v3.Up(appDb, "migrations"); err != nil {
		log.Fatal(err)
	}

	if err := goose_v3.Status(appDb, "migrations"); err != nil {
		log.Fatal(err)
	}

	return false
}

func main() {
	var opts ArgOptions
	var err error
	var db *gorm.DB
	var args []string

	//log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
	//log.SetOutput(os.Stdout)

	args, err = flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal(err)
		return
	}

	if opts.PWHash != "" {
		hashed, err := tools.HashPassword(opts.PWHash)
		if err != nil {
			log.Fatal(err)
			return
		}
		log.Printf("Password: %s Hashed: %s", opts.PWHash, hashed)

		return
	}

	appConf := config.AppConfig(opts.Env)
	if appConf.Debug {
		log.SetLevel(log.DebugLevel)
	}

	// check DB Connection on Start 100 times
	for i := 1; i <= 100; i++ {
		db, err = dbConn.New(appConf)
		if err != nil {
			log.Error(err)
		} else {
			break
		}
		time.Sleep(4 * time.Second)
	}
	if err != nil {
		log.Info("Stopped Trying DB Connection")
		return
	}

	if opts.Migrate != "" {
		if migrate(db, opts.Migrate, args) {
			log.Info("Program exited")
			return
		}
	}

	validator := vr.New()
	application := app.New(db, validator, appConf)

	appRouter := router.New(application)
	address := fmt.Sprintf(":%d", appConf.Server.Port)
	log.Infof("Starting server %v", address)

	srv := &http.Server{
		Addr:         address,
		Handler:      appRouter,
		ReadTimeout:  appConf.Server.TimeoutRead,
		WriteTimeout: appConf.Server.TimeoutWrite,
		IdleTimeout:  appConf.Server.TimeoutIdle,
	}

	closed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Info("Shutting down server")

		ctx, cancel := context.WithTimeout(context.Background(), appConf.Server.TimeoutIdle)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.WithField("error", err).Warn("Server shutdown failure")
		}

		sqlDB, err := db.DB()
		if err == nil {
			if err = sqlDB.Close(); err != nil {
				log.WithField("error", err).Warn("Db connection closing failure")
			}
		}

		close(closed)
	}()
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.WithField("error", err).Warn("Server startup failure")
	}

	<-closed

}
