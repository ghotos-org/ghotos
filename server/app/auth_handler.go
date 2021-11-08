package app

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"ghotos/model"
	"ghotos/repository"
	"ghotos/server/service"
	"ghotos/util/mail"
	"ghotos/util/tools"
	"ghotos/util/validator"
	"time"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"net/http"
)

func (app *App) HandleAuthLogin(w http.ResponseWriter, r *http.Request) {
	form := &model.UserLoginForm{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		printError(app, w, http.StatusUnprocessableEntity, appErrFormDecodingFailure, err)
		return
	}

	if err := app.validator.Struct(form); err != nil {
		log.Warn(err)
		resp := validator.ToErrResponse(err, nil)
		if resp == nil {
			printError(app, w, http.StatusInternalServerError, appErrFormErrResponseFailure, err)
			return
		}
		respBody, err := json.Marshal(resp)
		if err != nil {
			printError(app, w, http.StatusInternalServerError, appErrJsonCreationFailure, err)
			return
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(respBody)
		return
	}

	user, err := repository.LoginUser(app.db, form.Email, form.Password)
	if err != nil {
		printError(app, w, http.StatusInternalServerError, "user & password not matched", err)
		return
	}

	jwt := service.Jwt{}
	token, err := jwt.CreateToken(*user)
	if err != nil {
		printError(app, w, http.StatusInternalServerError, "unable to create access token", err)
		return
	}

	if err := json.NewEncoder(w).Encode(token); err != nil {
		printError(app, w, http.StatusInternalServerError, appErrJsonCreationFailure, err)
		return
	}

}

func (app *App) HandleSignUpCheckMail(w http.ResponseWriter, r *http.Request) {
	form := &model.UserRegisterEmailForm{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		log.Warn(err)
		printError(app, w, http.StatusUnprocessableEntity, appErrFormDecodingFailure, err)
		return
	}
	if err := app.validator.Struct(form); err != nil {
		log.Warn(err)
		resp := validator.ToErrResponse(err, nil)
		if resp == nil {
			printError(app, w, http.StatusInternalServerError, appErrFormErrResponseFailure, err)
			return
		}
		respBody, err := json.Marshal(resp)
		if err != nil {
			printError(app, w, http.StatusInternalServerError, appErrJsonCreationFailure, err)
			return
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(respBody)
		return
	}

	userModel, err := form.ToModel()
	if err != nil {
		printError(app, w, http.StatusUnprocessableEntity, appErrFormDecodingFailure, err)
		return
	}

	_, err = repository.ReadUserByEmail(app.db, userModel.Email)

	if err != nil && err != gorm.ErrRecordNotFound {
		printError(app, w, http.StatusInternalServerError, appErrDataAccessFailure, err)
		return
	}

	if err == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error" : { "fields" : { "email": "%v"}}}`, "email already taken")
		return
	}

	// todo:
	// check if email existes, if, then send mail with password forogot

	form.Date = time.Now()
	formBytes, err := json.Marshal(form)
	if err != nil {
		printError(app, w, http.StatusUnprocessableEntity, "", err)
		return
	}
	textByte, err := tools.Encrypt(formBytes, []byte(app.conf.Server.TokenKey), "")
	if err != nil {
		printError(app, w, http.StatusUnprocessableEntity, "error on creating link", err)
		return
	}
	link := app.conf.Server.PublicUrl + "/register/" + hex.EncodeToString(textByte)
	email := mail.SendMail{}
	email.From = app.conf.SMTP.Sender
	email.To = []string{form.Email}
	email.Html = "Hi User<br>Activate <a href=" + link + ">Link</a>"
	email.Subject = "ghotos | New Accout Registration"
	err = mail.Send(email, mail.GetConf(app.conf.SMTP.Host, app.conf.SMTP.Port, app.conf.SMTP.User, app.conf.SMTP.Password))
	if err != nil {
		printError(app, w, http.StatusInternalServerError, "mailversandt", err)
		return
	}
	w.WriteHeader(http.StatusCreated)

}

func checkLink(app *App, w http.ResponseWriter, r *http.Request) (*model.User, error) {
	// get & check email from form
	userformEnc := chi.URLParam(r, "userform")
	userBytesEnc, err := hex.DecodeString(userformEnc)
	if err != nil {
		printError(app, w, http.StatusInternalServerError, "check link cannot read", err)
		return nil, err
	}
	userBytes, err := tools.Decrypt(userBytesEnc, []byte(app.conf.Server.TokenKey))
	userForm := &model.UserRegisterEmailForm{}
	json.Unmarshal(userBytes, &userForm)
	if err != nil {
		printError(app, w, http.StatusUnprocessableEntity, "link is not valid ", err)
		return nil, err
	}

	userModel, err := userForm.ToModel()
	if err != nil {
		printError(app, w, http.StatusInternalServerError, appErrFormDecodingFailure, err)
		return nil, err
	}

	if userModel == nil {
		printError(app, w, http.StatusInternalServerError, "link is not valid, no mail informations", nil)
		return nil, errors.New("user is null")
	}

	//database check
	_, err = repository.ReadUserByEmail(app.db, userModel.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		printError(app, w, http.StatusInternalServerError, appErrDataAccessFailure, err)
		return nil, err
	}

	if err == nil { // User Allready inseretd
		printError(app, w, http.StatusGone, "link is expired", err)
		return nil, errors.New("user is null")
	}

	return userModel, nil
}

func (app *App) HandleSignUpCheckLink(w http.ResponseWriter, r *http.Request) {

	// get & check email from form
	_, err := checkLink(app, w, r)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (app *App) HandleSignUpCreateUser(w http.ResponseWriter, r *http.Request) {

	userModel, err := checkLink(app, w, r)
	if err != nil {
		return
	}

	// get & check password from form
	passwordForm := &model.UserRegisterPasswordForm{}
	if err := json.NewDecoder(r.Body).Decode(passwordForm); err != nil {
		log.Warn(err)
		printError(app, w, http.StatusUnprocessableEntity, appErrFormDecodingFailure, err)
		return
	}
	if err := app.validator.Struct(passwordForm); err != nil {
		log.Warn(err)
		resp := validator.ToErrResponse(err, nil)
		if resp == nil {
			printError(app, w, http.StatusInternalServerError, appErrFormErrResponseFailure, err)
			return
		}
		respBody, err := json.Marshal(resp)
		if err != nil {
			printError(app, w, http.StatusInternalServerError, appErrJsonCreationFailure, err)
			return
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(respBody)
		return
	}

	userModel.Password = passwordForm.Password
	user, err := repository.CreateUser(app.db, userModel)
	if err != nil {
		printError(app, w, http.StatusInternalServerError, appErrCreationFailure, err)
		return
	}

	log.Infof("User created %s", user.UID)
	w.WriteHeader(http.StatusCreated)
}

func (app *App) HandleAuthRefresh(w http.ResponseWriter, r *http.Request) {

	token := model.Token{}

	log.Printf("body: %d", r.Context())
	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		printError(app, w, http.StatusUnprocessableEntity, appErrFormDecodingFailure, err)
		return
	}

	jwt := service.Jwt{}
	user, err := jwt.ValidateRefreshToken(token)
	if err != nil {
		printError(app, w, http.StatusUnauthorized, "invalid token", err)
		return
	}

	token, err = jwt.CreateToken(user)
	if err != nil {
		printError(app, w, http.StatusInternalServerError, "unable to create access token", err)
		return
	}

	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error.message": "%v"}`, appErrJsonCreationFailure)
		return
	}
}

func (app *App) HandleAuthLogout(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusNoContent)

}
