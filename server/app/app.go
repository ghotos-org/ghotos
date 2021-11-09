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
		if err != nil {
			log.Warn(err)
		}
	}
	if msg != "" {
		log.Error(msg)
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

func (app *App) checkForm(form interface{}, w http.ResponseWriter, r *http.Request) (stop bool) {
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		printError(app, w, http.StatusUnprocessableEntity, appErrFormDecodingFailure, err)
		return true
	}

	if err := app.validator.Struct(form); err != nil {
		log.Warn(err)
		resp := val.ToErrResponse(err, nil)
		if resp == nil {
			printError(app, w, http.StatusInternalServerError, appErrFormErrResponseFailure, err)
			return true
		}
		respBody, err := json.Marshal(resp)
		if err != nil {
			printError(app, w, http.StatusInternalServerError, appErrJsonCreationFailure, err)
			return true
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(respBody)
		return true
	}

	return false
}
