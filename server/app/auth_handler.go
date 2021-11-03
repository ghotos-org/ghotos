package app

import (
	"encoding/json"
	"fmt"
	"ghotos/model"
	"ghotos/repository"
	"ghotos/server/service"
	"ghotos/util/validator"
	"log"
	"net/http"
)

func (app *App) HandleAuthLogin(w http.ResponseWriter, r *http.Request) {
	form := &model.UserLoginForm{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		app.logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrFormDecodingFailure)
		return
	}

	if err := app.validator.Struct(form); err != nil {
		app.logger.Warn().Err(err).Msg("")
		resp := validator.ToErrResponse(err)
		if resp == nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%v"}`, appErrFormErrResponseFailure)
			return
		}
		respBody, err := json.Marshal(resp)
		if err != nil {
			app.logger.Warn().Err(err).Msg("")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%v"}`, appErrJsonCreationFailure)
			return
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write(respBody)
		return
	}

	user, err := repository.LoginUser(app.db, form.Email, form.Password)
	if err != nil {
		app.logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrDataCreationFailure)
		return
	}

	jwt := service.Jwt{}
	token, err := jwt.CreateToken(*user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, "unable to create access toke")
		return
	}

	app.logger.Info().Msgf("User Login: %d", user.UID)
	app.logger.Info().Msgf("User Token: %d", token.AccessToken)

	if err := json.NewEncoder(w).Encode(token); err != nil {
		app.logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrJsonCreationFailure)
		return
	}

}

func (app *App) HandleAuthRefresh(w http.ResponseWriter, r *http.Request) {

	token := model.Token{}

	log.Printf("body: %d", r.Context())
	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		app.logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrFormDecodingFailure)
		return
	}

	jwt := service.Jwt{}
	user, err := jwt.ValidateRefreshToken(token)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error": "%v"}`, "invalid token")
		return
	}

	token, err = jwt.CreateToken(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, "unable to create access token")
		return
	}
	/*
		if err := json.NewEncoder(w).Encode(token); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%v"}`, appErrJsonCreationFailure)
			return
		}

	*/

	err = json.NewEncoder(w).Encode(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrJsonCreationFailure)
		return
	}

	/*
		w.WriteHeader(http.StatusCreated)
		app.logger.Info().Msg("TEST1")
		fmt.Fprintf(w, `{"error": "%v"}`, appErrJsonCreationFailure)
	*/
}

func (app *App) HandleAuthLogout(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusNoContent)

}
