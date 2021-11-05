package app

import (
	"encoding/json"
	"fmt"
	"ghotos/model"
	"ghotos/repository"
	"net/http"
	"strconv"

	"ghotos/util/validator"

	log "github.com/sirupsen/logrus"

	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

func (app *App) HandleListBooks(w http.ResponseWriter, r *http.Request) {
	books, err := repository.ListBooks(app.db)
	if err != nil {
		log.Warn(err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error.message": "%v"}`, appErrDataAccessFailure)
		return
	}
	if books == nil {
		fmt.Fprint(w, "[]")
		return
	}
	dtos := books.ToDto()
	if err := json.NewEncoder(w).Encode(dtos); err != nil {
		log.Warn(err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error.message": "%v"}`, appErrJsonCreationFailure)
		return
	}
}
func (app *App) HandleReadBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)
	if err != nil || id == 0 {
		log.Infof("can not parse ID: %v", id)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	book, err := repository.ReadBook(app.db, uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Warn(err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error.message": "%v"}`, appErrDataAccessFailure)
		return
	}
	dto := book.ToDto()
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		log.Warn(err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error.message": "%v"}`, appErrJsonCreationFailure)
		return
	}
}

func (app *App) HandleDeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)
	if err != nil || id == 0 {
		log.Infof("can not parse ID: %v", id)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	if err := repository.DeleteBook(app.db, uint(id)); err != nil {
		log.Warn(err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error.message": "%v"}`, appErrDataAccessFailure)
		return
	}
	log.Infof("Book deleted: %d", id)
	w.WriteHeader(http.StatusAccepted)
}

func (app *App) HandleCreateBook(w http.ResponseWriter, r *http.Request) {
	form := &model.BookForm{}
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

	bookModel, err := form.ToModel()
	if err != nil {
		log.Warn(err)

		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error.message": "%v"}`, appErrFormDecodingFailure)
		return
	}
	book, err := repository.CreateBook(app.db, bookModel)
	if err != nil {
		log.Warn(err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error.message": "%v"}`, appErrDataCreationFailure)
		return
	}
	log.Infof("New book created: %d", book.ID)
	w.WriteHeader(http.StatusCreated)
}

func (app *App) HandleUpdateBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)
	if err != nil || id == 0 {
		log.Infof("can not parse ID: %v", id)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	form := &model.BookForm{}
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

	bookModel, err := form.ToModel()
	if err != nil {
		log.Warn(err)

		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error.message": "%v"}`, appErrFormDecodingFailure)
		return
	}
	bookModel.ID = uint(id)
	if err := repository.UpdateBook(app.db, bookModel); err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		log.Warn(err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error.message": "%v"}`, appErrDataUpdateFailure)
		return
	}
	log.Infof("Book updated: %d", id)
	w.WriteHeader(http.StatusAccepted)
}
