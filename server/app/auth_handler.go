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

func getUserfromLink(app *App, w http.ResponseWriter, r *http.Request) (*model.User, *model.UserRegisterEmailForm, error) {
	// get & check email from form
	userformEnc := chi.URLParam(r, "userform")
	userBytesEnc, err := hex.DecodeString(userformEnc)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}
	userBytes, err := tools.Decrypt(userBytesEnc, []byte(app.conf.Server.TokenKey))
	userForm := &model.UserRegisterEmailForm{}
	json.Unmarshal(userBytes, &userForm)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	userModel, err := userForm.ToModel()
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}

	if userModel.Email == "" {
		printError(app, w, http.StatusInternalServerError, "link is not valid, no mail informations", nil)
		return nil, nil, errors.New("user is null")
	}

	return userModel, userForm, nil
}

func (app *App) HandleAuthLogin(w http.ResponseWriter, r *http.Request) {
	form := &model.UserLoginForm{}
	if app.checkForm(form, w, r) {
		return
	}

	user, err := repository.ReadUserByEmail(app.db, form.Email)
	if err != nil {
		printError(app, w, http.StatusInternalServerError, "user & password not matched", err)
		return
	}
	if !tools.CheckPasswordHash(form.Password, user.Password) {
		printError(app, w, http.StatusInternalServerError, "user & password not matched", nil)
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

	userModel, _, err := getUserfromLink(app, w, r)
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

	userModel, _, err := getUserfromLink(app, w, r)
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

	userModel.Password, err = tools.HashPassword(passwordForm.Password)
	if err != nil {
		printError(app, w, http.StatusInternalServerError, "Password not allowd", err)
		return
	}

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

	user, err := repository.ReadUserByEmail(app.db, userModel.Email)

	if err != nil && err != gorm.ErrRecordNotFound {
		printError(app, w, http.StatusInternalServerError, appErrDataAccessFailure, err)
		return
	}

	if err != nil && err == gorm.ErrRecordNotFound {
		log.Warnf("Fake Response, user not exists: %v", form.Email)
		w.WriteHeader(http.StatusCreated)
		return
	}

	if user.NewPasswordRequest != nil {
		duration := float64(5)
		log.Info(time.Since(*user.NewPasswordRequest).Minutes())
		if time.Since(*user.NewPasswordRequest).Minutes() < duration {
			printError(app, w, http.StatusInternalServerError, "new password link already send, pleas check your email or try later (5min)", nil)
			return
		}

	}

	now := time.Now()
	user.NewPasswordRequest = &now

	err = repository.UpdateUser(app.db, user)
	if err != nil {
		printError(app, w, http.StatusInternalServerError, appErrDataAccessFailure, err)
		return
	}
	// todo:
	// check if email existes, if, then send mail with password forogot

	form.Date = now
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

	userModel, userForm, err := getUserfromLink(app, w, r)
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

	if user.NewPasswordRequest == nil {
		log.Error("no user.NewPasswordRequest was set")
		printError(app, w, http.StatusInternalServerError, "link is invalid", nil)
		return
	}
	userDate := *user.NewPasswordRequest

	if userForm.Date.Format(tools.DATE_FORMAT_COMPARE) != userDate.Format(tools.DATE_FORMAT_COMPARE) {
		log.Error("userForm.Date != *user.NewPasswordRequest")
		printError(app, w, http.StatusInternalServerError, "link exipred", nil)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *App) HandleNewPasswordCreate(w http.ResponseWriter, r *http.Request) {

	userModel, userForm, err := getUserfromLink(app, w, r)
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

	if user.NewPasswordRequest == nil {
		log.Error("no user.NewPasswordRequest was set")
		printError(app, w, http.StatusInternalServerError, "link is invalid", nil)
		return
	}
	userDate := *user.NewPasswordRequest

	if userForm.Date.Format(tools.DATE_FORMAT_COMPARE) != userDate.Format(tools.DATE_FORMAT_COMPARE) {
		log.Error("userForm.Date != *user.NewPasswordRequest")
		printError(app, w, http.StatusInternalServerError, "link exipred", nil)
		return
	}

	// get & check password from form
	passwordForm := &model.UserRegisterPasswordForm{}
	if app.checkForm(passwordForm, w, r) {
		return
	}

	user.NewPasswordRequest = nil
	user.Password, err = tools.HashPassword(passwordForm.Password)
	if err != nil {
		printError(app, w, http.StatusInternalServerError, "Password not allowd", err)
		return
	}
	err = repository.UpdateUser(app.db, user)
	if err != nil {
		printError(app, w, http.StatusInternalServerError, appErrCreationFailure, err)
		return
	}

	log.Infof("User updated %s", user.UID)
	w.WriteHeader(http.StatusCreated)
}
