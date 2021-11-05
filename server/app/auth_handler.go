package app

import (
	"encoding/json"
	"fmt"
	"ghotos/model"
	"ghotos/repository"
	"ghotos/server/service"
	"ghotos/util/mail"
	"ghotos/util/validator"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"net/http"
)

func (app *App) HandleAuthLogin(w http.ResponseWriter, r *http.Request) {
	form := &model.UserLoginForm{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		log.Warn(err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error.message": "%v"}`, appErrFormDecodingFailure)
		return
	}

	if err := app.validator.Struct(form); err != nil {
		log.Warn(err)
		resp := validator.ToErrResponse(err, nil)
		if resp == nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error.message": "%v"}`, appErrFormErrResponseFailure)
			return
		}
		respBody, err := json.Marshal(resp)
		if err != nil {
			log.Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error.message": "%v"}`, appErrJsonCreationFailure)
			return
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(respBody)
		return
	}

	user, err := repository.LoginUser(app.db, form.Email, form.Password)
	if err != nil {
		//log.Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error.message": "%v"}`, appErrDataCreationFailure)
		return
	}

	jwt := service.Jwt{}
	token, err := jwt.CreateToken(*user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error.message": "%v"}`, "unable to create access toke")
		return
	}

	if err := json.NewEncoder(w).Encode(token); err != nil {
		log.Warn(err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error.message": "%v"}`, appErrJsonCreationFailure)
		return
	}

}

func (app *App) HandleAuthSignup(w http.ResponseWriter, r *http.Request) {

	form := &model.UserRegisterForm{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		log.Warn(err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error.message": "%v"}`, appErrFormDecodingFailure)
		return
	}

	if err := app.validator.Struct(form); err != nil {
		log.Warn(err)
		resp := validator.ToErrResponse(err, nil)
		if resp == nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error.message": "%v"}`, appErrFormErrResponseFailure)
			return
		}
		respBody, err := json.Marshal(resp)
		if err != nil {
			log.Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error.message": "%v"}`, appErrJsonCreationFailure)
			return
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(respBody)
		return
	}

	userModel, err := form.ToModel()
	if err != nil {
		log.Warn(err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error.message": "%v"}`, appErrFormDecodingFailure)
		return
	}

	user, err := repository.ReadUserByEmail(app.db, userModel.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Warn(err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error.message": "%v"}`, appErrDataAccessFailure)
		return
	}

	if user == nil {
		user, err = repository.CreateUser(app.db, userModel)
		if err != nil {
			log.Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error.message": "%v"}`, appErrDataCreationFailure)
			return
		}
	} else {
		if user.Status != 0 {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, `{"error.message": "%v"}`, "User exists")
			return
		}

		err = repository.UpdateUser(app.db, user)
		if err != nil {
			log.Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error.message": "%v"}`, appErrDataUpdateFailure)
			return
		}
	}

	log.Infof("New User created/updated: %d", user.UID)
	link := app.conf.Server.PublicUrl + "/auth/activate/test"
	email := mail.SendMail{}
	email.From = app.conf.SMTP.Sender
	email.To = []string{user.Email}
	email.Html = "Hi User<br>Activate <a href=" + link + ">Link</a>"
	email.Subject = "ghotos | New Accout Registration"
	err = mail.Send(email, mail.GetConf(app.conf.SMTP.Host, app.conf.SMTP.Port, app.conf.SMTP.User, app.conf.SMTP.Password))
	if err != nil {
		log.Warn(err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error.message": "%v"}`, "mailversandt")
		return
	}
	w.WriteHeader(http.StatusCreated)

	/*
		user, err := repository.LoginUser(app.db, form.Email, form.Password)
		if err != nil {
			//log.Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error.message": "%v"}`, appErrDataCreationFailure)
			return
		}

		jwt := service.Jwt{}
		token, err := jwt.CreateToken(*user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error.message": "%v"}`, "unable to create access toke")
			return
		}

		if err := json.NewEncoder(w).Encode(token); err != nil {
			log.Warn(err)
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error.message": "%v"}`, appErrJsonCreationFailure)
			return
		}
	*/

}

func (app *App) HandleAuthRefresh(w http.ResponseWriter, r *http.Request) {

	token := model.Token{}

	log.Printf("body: %d", r.Context())
	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		log.Warn(err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error.message": "%v"}`, appErrFormDecodingFailure)
		return
	}

	jwt := service.Jwt{}
	user, err := jwt.ValidateRefreshToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error.message": "%v"}`, "invalid token")
		return
	}

	token, err = jwt.CreateToken(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error.message": "%v"}`, "unable to create access token")
		return
	}
	/*
		if err := json.NewEncoder(w).Encode(token); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error.message": "%v"}`, appErrJsonCreationFailure)
			return
		}

	*/

	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error.message": "%v"}`, appErrJsonCreationFailure)
		return
	}

	/*
		w.WriteHeader(http.StatusCreated)
		app.logger.Info().Msg("TEST1")
		fmt.Fprintf(w, `{"error.message": "%v"}`, appErrJsonCreationFailure)
	*/
}

func (app *App) HandleAuthLogout(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusNoContent)

}
