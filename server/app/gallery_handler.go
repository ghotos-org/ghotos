package app

import (
	"encoding/json"
	"fmt"
	"ghotos/model"
	"ghotos/repository"
	"net/http"

	"github.com/go-chi/chi"
)

func (app *App) HandleGallery(w http.ResponseWriter, r *http.Request) {
	gallery, err := repository.GalleryDays(app.db, app.User().UID)
	if err != nil {
		app.logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrDataAccessFailure)
		return
	}
	if gallery.Days == nil {
		fmt.Fprint(w, "[]")
		return
	}

	if err := json.NewEncoder(w).Encode(gallery); err != nil {
		app.logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrJsonCreationFailure)
		return
	}
}

func (app *App) HandleListGalleryDayFile(w http.ResponseWriter, r *http.Request) {
	day := chi.URLParam(r, "day")
	files, err := repository.ListGalleryDayFile(app.db, day, app.User().UID)
	if err != nil {
		app.logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrDataAccessFailure)
		return
	}
	if files == nil {
		fmt.Fprint(w, "[]")
		return
	}

	fileSrcs := make([]*model.FilesSrc, 0)
	for _, file := range files {
		//src, err := repository.EncodePath(*file)
		if err != nil {
			app.logger.Warn().Err(err).Msg("")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%v"}`, appErrDataAccessFailure)
			return
		}
		fileSrc := model.FilesSrc{}
		day := file.Date.Format("20060102")

		fileSrc.Day = &day
		//fileSrc.Src = &src
		fileSrc.UID = file.UID
		fileSrcs = append(fileSrcs, &fileSrc)
	}

	if err := json.NewEncoder(w).Encode(fileSrcs); err != nil {
		app.logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrJsonCreationFailure)
		return
	}

}
