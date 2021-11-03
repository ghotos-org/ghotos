package app

import (
	"ghotos/config"
	"ghotos/model"
	"ghotos/util/logger"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

const (
	appErrDataAccessFailure   = "data access failure"
	appErrJsonCreationFailure = "json creation failure"

	appErrDataCreationFailure = "data creation failure"
	appErrFormDecodingFailure = "form decoding failure"

	appErrDataUpdateFailure      = "data update failure"
	appErrFormErrResponseFailure = "form error response failure"
)

type App struct {
	logger    *logger.Logger
	db        *gorm.DB
	validator *validator.Validate
	conf      *config.Conf
	user      *model.User
}

func New(
	logger *logger.Logger,
	db *gorm.DB,
	validator *validator.Validate,
	conf *config.Conf,
) *App {
	return &App{
		logger:    logger,
		db:        db,
		validator: validator,
		conf:      conf,
	}
}

func (app *App) Logger() *logger.Logger {
	return app.logger
}
func (app *App) SetUser(user model.User) {
	app.user = &user
}

func (app *App) User() *model.User {
	return app.user
}
