package app

import (
	"encoding/json"
	"fmt"
	"ghotos/config"
	"ghotos/model"
	"net/http"
	"runtime"

	val "ghotos/util/validator"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	appErr                    = "app error"
	appErrCreationFailure     = "error createn failure"
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

func printError(app *App, w http.ResponseWriter, status int, msg string, err error) {
	if err != nil && app.conf.Debug {
		_, fn, line, _ := runtime.Caller(1)
		log.WithFields(log.Fields{
			"func": fn,
			"line": fmt.Sprintf("%d", line),
		}).Error(err)
	} else {
		log.Warn(err)
	}

	w.WriteHeader(status)
	message := ""

	if msg == "" {
		message = appErr
	} else {
		message = msg
	}
	errorObj := val.ErrorMsg(message)
	errorJson, err := json.Marshal(errorObj)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error.message": "%v"}`, appErrCreationFailure)
		return
	}
	w.Write(errorJson)

}
