package main

import (
	"context"
	"embed"
	"fmt"
	dbConn "ghotos/adapter/gorm"
	"ghotos/config"
	"ghotos/server/app"
	"ghotos/server/router"
	lr "ghotos/util/logger"
	vr "ghotos/util/validator"
	"strconv"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jessevdk/go-flags"
	"gorm.io/gorm"

	"github.com/rs/zerolog/log"

	//goose "github.com/pressly/goose"
	_ "ghotos/migrations"

	goose_v3 "github.com/pressly/goose/v3"
)

type ArgOptions struct {
	Env     string `short:"e" long:"env" description:"Environment File (default: .env)"`
	Migrate string `short:"m" long:"migrate" description:"DB mirgrate tool"  choice:"up" choice:"down" choice:"status" choice:"version" choice:"reset" choice:"up-by-one" choice:"up-to" choice:"down-to"`
}

//go:embed migrations/*
var embedMigrations embed.FS

func migrate(db *gorm.DB, logger *lr.Logger, migrateCMD string, migrateArgs []string) (doExit bool) {
	appDb, err := db.DB()
	if err != nil {
		logger.Fatal().Err(err).Msg("")
		return
	}
	goose_v3.SetBaseFS(embedMigrations)

	if err := goose_v3.SetDialect("mysql"); err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	if migrateCMD == "up" {
		if err := goose_v3.Up(appDb, "migrations"); err != nil {
			logger.Fatal().Err(err).Msg("")
		}
	}
	if migrateCMD == "down" {
		if err := goose_v3.Down(appDb, "migrations"); err != nil {
			logger.Fatal().Err(err).Msg("")
		}
	}
	if migrateCMD == "status" {
		if err := goose_v3.Status(appDb, "migrations"); err != nil {
			logger.Fatal().Err(err).Msg("")
		}
	}
	if migrateCMD == "version" {
		if err := goose_v3.Version(appDb, "migrations"); err != nil {
			logger.Fatal().Err(err).Msg("")
		}
	}
	if migrateCMD == "reset" {
		if err := goose_v3.Reset(appDb, "migrations"); err != nil {
			logger.Fatal().Err(err).Msg("")
		}
	}
	if migrateCMD == "up-by-one" {
		if err := goose_v3.Reset(appDb, "migrations"); err != nil {
			logger.Fatal().Err(err).Msg("")
		}
	}
	if migrateCMD == "up-to" || migrateCMD == "down-to" {
		if len(migrateArgs) < 2 {
			logger.Fatal().Err(err).Msg("Version fehlt")
			return true
		}
		version, err := strconv.ParseInt(migrateArgs[1], 10, 64)
		if err != nil {
			logger.Fatal().Err(err).Msg("Version Format nicht korrekt")
		}

		if migrateCMD == "up-to" {
			if err := goose_v3.UpTo(appDb, "migrations", version); err != nil {
				logger.Fatal().Err(err).Msg("")
			}
		} else {
			if err := goose_v3.DownTo(appDb, "migrations", version); err != nil {
				logger.Fatal().Err(err).Msg("")
			}
		}

	}
	if migrateCMD != "" {
		return true
	}

	if err := goose_v3.Up(appDb, "migrations"); err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	if err := goose_v3.Status(appDb, "migrations"); err != nil {
		logger.Fatal().Err(err).Msg("")
	}

	return false
}

func main() {
	var opts ArgOptions
	var err error
	var db *gorm.DB
	var args []string

	args, err = flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Error().Err(err).Msg("")
		return
	}
	appConf := config.AppConfig(opts.Env)
	logger := lr.New(appConf.Debug)

	// check DB Connection on Start 100 times
	for i := 1; i <= 100; i++ {
		db, err = dbConn.New(appConf)
		if err != nil {
			logger.Error().Stack().Err(err).Msg("")
		} else {
			break
		}
		time.Sleep(4 * time.Second)
	}
	if err != nil {
		logger.Info().Msg("Stopped Trying DB Connection")
		return
	}

	if migrate(db, logger, opts.Migrate, args) {
		logger.Info().Msg("Program exited")
		return
	}

	validator := vr.New()
	application := app.New(logger, db, validator, appConf)

	appRouter := router.New(application)
	address := fmt.Sprintf(":%d", appConf.Server.Port)
	logger.Info().Msgf("Starting server %v", address)
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

		logger.Info().Msgf("Shutting down server %v", address)

		ctx, cancel := context.WithTimeout(context.Background(), appConf.Server.TimeoutIdle)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logger.Warn().Err(err).Msg("Server shutdown failure")
		}

		sqlDB, err := db.DB()
		if err == nil {
			if err = sqlDB.Close(); err != nil {
				logger.Warn().Err(err).Msg("Db connection closing failure")
			}
		}

		close(closed)
	}()

	logger.Info().Msgf("Starting server %v", address)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal().Err(err).Msg("Server startup failure")
	}

	<-closed

}
