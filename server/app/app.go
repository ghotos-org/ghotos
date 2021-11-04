package app

import (
	"fmt"
	"ghotos/config"
	"ghotos/model"
	"net/http"
	"runtime"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
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
	db        *gorm.DB
	validator *validator.Validate
	conf      *config.Conf
	user      *model.User
}

func New(
	db *gorm.DB,
	validator *validator.Validate,
	conf *config.Conf,
) *App {
	return &App{
		db:        db,
		validator: validator,
		conf:      conf,
	}
}

func (app *App) SetUser(user model.User) {
	app.user = &user
}

func (app *App) User() *model.User {
	return app.user
}

func printError(app *App, w http.ResponseWriter, err error, status int, msg string) {
	_, fn, line, _ := runtime.Caller(1)
	log.WithFields(log.Fields{
		"func": fmt.Sprintf("%s", fn),
		"line": fmt.Sprintf("%d", line),
	}).Error(err)

	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, `{"error": "%v"}`, appErrJsonCreationFailure)

}
