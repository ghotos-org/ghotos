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
	"time"

	"github.com/go-chi/chi"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"net/http"
)

func getUserfromLink(app *App, w http.ResponseWriter, r *http.Request) (*model.User, error) {
	// get & check email from form
	userformEnc := chi.URLParam(r, "userform")
	userBytesEnc, err := hex.DecodeString(userformEnc)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	userBytes, err := tools.Decrypt(userBytesEnc, []byte(app.conf.Server.TokenKey))
	userForm := &model.UserRegisterEmailForm{}
	json.Unmarshal(userBytes, &userForm)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	userModel, err := userForm.ToModel()
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if userModel.Email == "" {
		printError(app, w, http.StatusInternalServerError, "link is not valid, no mail informations", nil)
		return nil, errors.New("user is null")
	}

	return userModel, nil
}

func (app *App) HandleAuthLogin(w http.ResponseWriter, r *http.Request) {
	form := &model.UserLoginForm{}
	if app.checkForm(form, w, r) {
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

func (app *App) HandleSignUpLink2Mail(w http.ResponseWriter, r *http.Request) {
	form := &model.UserRegisterEmailForm{}
	if app.checkForm(form, w, r) {
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

func (app *App) HandleSignUpCheckLink(w http.ResponseWriter, r *http.Request) {

	// get & check email from form

	userModel, err := getUserfromLink(app, w, r)
	if err != nil {
		printError(app, w, http.StatusInternalServerError, "link is invalid", err)
		return
	}
	//database check
	_, err = repository.ReadUserByEmail(app.db, userModel.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		printError(app, w, http.StatusInternalServerError, "link is invalid", err)
		return
	}

	if err == nil { // User Allready inseretd
		printError(app, w, http.StatusGone, "link is expired", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *App) HandleSignUpCreateUser(w http.ResponseWriter, r *http.Request) {

	userModel, err := getUserfromLink(app, w, r)
	if err != nil {
		printError(app, w, http.StatusInternalServerError, "link is invalid", err)
		return
	}
	//database check
	_, err = repository.ReadUserByEmail(app.db, userModel.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		printError(app, w, http.StatusInternalServerError, "link is invalid", err)
		return
	}

	if err == nil { // User Allready inseretd
		printError(app, w, http.StatusGone, "link is expired", err)
		return
	}

	// get & check password from form
	passwordForm := &model.UserRegisterPasswordForm{}
	if app.checkForm(passwordForm, w, r) {
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

func (app *App) HandleNewPasswordLink2Mail(w http.ResponseWriter, r *http.Request) {

	form := &model.UserRegisterEmailForm{}
	if app.checkForm(form, w, r) {
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

	if err != nil && err == gorm.ErrRecordNotFound {
		w.WriteHeader(http.StatusCreated)
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
	link := app.conf.Server.PublicUrl + "/password/" + hex.EncodeToString(textByte)
	email := mail.SendMail{}
	email.From = app.conf.SMTP.Sender
	email.To = []string{form.Email}
	email.Html = "Hi User<br>new password creation link:  <a href=" + link + ">Link</a>"
	email.Subject = "ghotos | New Accout Registration"
	err = mail.Send(email, mail.GetConf(app.conf.SMTP.Host, app.conf.SMTP.Port, app.conf.SMTP.User, app.conf.SMTP.Password))
	if err != nil {
		printError(app, w, http.StatusInternalServerError, "mailversandt", err)
		return
	}
	w.WriteHeader(http.StatusCreated)

}

func (app *App) HandleNewPasswordCheckLink(w http.ResponseWriter, r *http.Request) {

	userModel, err := getUserfromLink(app, w, r)
	if err != nil {
		printError(app, w, http.StatusInternalServerError, "link is invalid", err)
		return
	}
	//database check
	_, err = repository.ReadUserByEmail(app.db, userModel.Email)
	if err != nil {
		printError(app, w, http.StatusInternalServerError, "link exipred", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *App) HandleNewPasswordCreate(w http.ResponseWriter, r *http.Request) {

	userModel, err := getUserfromLink(app, w, r)
	if err != nil {
		printError(app, w, http.StatusInternalServerError, "link is invalid", err)
		return
	}
	//database check
	user, err := repository.ReadUserByEmail(app.db, userModel.Email)
	if err != nil {
		printError(app, w, http.StatusInternalServerError, "link exipred", err)
		return
	}

	// get & check password from form
	passwordForm := &model.UserRegisterPasswordForm{}
	if app.checkForm(passwordForm, w, r) {
		return
	}

	user.Password = passwordForm.Password
	err = repository.UpdateUser(app.db, user)
	if err != nil {
		printError(app, w, http.StatusInternalServerError, appErrCreationFailure, err)
		return
	}

	log.Infof("User updated %s", user.UID)
	w.WriteHeader(http.StatusCreated)
}
