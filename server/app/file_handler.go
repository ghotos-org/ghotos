package app

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"ghotos/model"
	"ghotos/repository"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/h2non/bimg"
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

func printError(app *App, w http.ResponseWriter, err error, status int, msg string) {
	app.logger.Error().Err(err).Msg("")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, `{"error": "%v"}`, appErrJsonCreationFailure)

}

func (app *App) HandleReadFile(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		app.logger.Info().Msgf("can not parse UID: %v", uid)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	modelFile, err := repository.ReadFileByUser(app.db, uid, app.User().UID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		app.logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrDataAccessFailure)
		return
	}
	src, err := repository.EncodePath(*modelFile)
	if err != nil {
		printError(app, w, err, http.StatusInternalServerError, appErrJsonCreationFailure)
		return
	}

	fileSimple := modelFile.ToDtoSimple()
	fileSimple.Src = &src
	if err := json.NewEncoder(w).Encode(fileSimple); err != nil {
		app.logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrJsonCreationFailure)
		return
	}

}

func (app *App) HandleDeleteFile(w http.ResponseWriter, r *http.Request) {
	uid := chi.URLParam(r, "uid")
	if uid == "" {
		app.logger.Info().Msgf("can not parse UID: %v", uid)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	if err := repository.DeleteFile(app.db, uid, app.User().UID); err != nil {
		app.logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, appErrDataAccessFailure)
		return
	}
	app.logger.Info().Msgf("File deleted: %d", uid)
	w.WriteHeader(http.StatusAccepted)
}

func (app *App) HandleShowFile(w http.ResponseWriter, r *http.Request) {
	srcParam := chi.URLParam(r, "src")

	fileOptions := strings.Join(strings.Split(srcParam, "=")[1:], "=")
	if fileOptions == "" {
		printError(app, w, nil, http.StatusInternalServerError, "fileOptions ist LEER!")
		return

	}

	re := regexp.MustCompile(`w[0-9]+`)
	widthStr := strings.Replace(re.FindString(fileOptions), "w", "", -1)

	re = regexp.MustCompile(`h[0-9]+`)
	heightStr := strings.Replace(re.FindString(fileOptions), "h", "", -1)

	width, err := strconv.Atoi(widthStr)
	if err != nil {
		printError(app, w, err, http.StatusInternalServerError, "Width nicht vorhanden")
		return
	}

	height, err := strconv.Atoi(heightStr)
	if err != nil {
		printError(app, w, err, http.StatusInternalServerError, "Heiht nicht vorhanden")
		return

	}
	src := strings.Split(srcParam, "=")[0]
	filePath, err := repository.DecodePath(src)
	if err != nil {
		printError(app, w, err, http.StatusInternalServerError, appErrJsonCreationFailure)
	}

	file, err := repository.ReadFile(app.db, filePath.UID)
	if err != nil {
		printError(app, w, err, http.StatusInternalServerError, appErrJsonCreationFailure)
	}

	eTag := file.UpdatedAt.Format("20060102150405000") + "_" + strconv.Itoa(width) + "_" + strconv.Itoa(height) + "_" + file.UID

	bufferOrg, err := bimg.Read(file.Path + "/" + file.Filename)
	imageOrg := bimg.NewImage(bufferOrg)
	if !bimg.IsTypeNameSupported(imageOrg.Type()) {
		printError(app, w, err, http.StatusInternalServerError, "Dateityp wird nicht unterstützt!")
	}

	quality := 100

	if width < 100 {
		quality = 50
	} else if width < 200 {
		quality = 80
	} else if width < 300 {
		quality = 90
	} else if width < 600 {
		quality = 95
	} else if width < 900 {
		quality = 100
	}

	options := bimg.Options{
		Height:        height,
		Width:         width,
		StripMetadata: true,
		Enlarge:       true,
		Crop:          false,
		Quality:       quality,
		Interlace:     true,
	}

	image, err := imageOrg.Process(options)

	if !bimg.IsTypeNameSupported(imageOrg.Type()) {
		printError(app, w, err, http.StatusInternalServerError, "Dateityp wird nicht unterstützt!")
	}

	e := `"` + eTag + `"`

	w.Header().Set("Cache-Control", "private, max-age=2592000, no-transform") // 30 days
	w.Header().Set("Content-Length", strconv.Itoa(len(image)))
	w.Header().Set("Expire", "Fri, 01 Jan 1990 00:00:00 GMT")
	w.Header().Set("Expire", "Fri, 01 Jan 1990 00:00:00 GMT")
	w.Header().Set("x-xss-protection", "0")
	w.Header().Set("access-control-expose-headers", "Content-Length")

	if match := r.Header.Get("If-None-Match"); match != "" {
		if strings.Contains(match, e) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	w.Header().Set("Content-Type", file.MimeType)
	w.Header().Set("content-disposition", "inline;filename=\""+file.Filename+"\"")

	w.Write(image)

}
func (app *App) HandleCreateFile(w http.ResponseWriter, r *http.Request) {
	uploadedFile, multipartFileHeader, err := r.FormFile("file")
	if err != nil {
		printError(app, w, err, http.StatusInternalServerError, appErrJsonCreationFailure)
		return
	}
	// Create a buffer to store the header of the file in
	fileHeader := make([]byte, 512)
	// Copy the headers into the FileHeader buffer
	if _, err := uploadedFile.Read(fileHeader); err != nil {
		printError(app, w, err, http.StatusInternalServerError, appErrJsonCreationFailure)
		return
	}
	// set position back to start.
	if _, err := uploadedFile.Seek(0, 0); err != nil {
		printError(app, w, err, http.StatusInternalServerError, appErrJsonCreationFailure)
		return
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, uploadedFile); err != nil {
		printError(app, w, err, http.StatusInternalServerError, appErrJsonCreationFailure)
		return
	}

	orgImage := bimg.NewImage(buf.Bytes())
	if !bimg.IsTypeNameSupported(orgImage.Type()) {
		printError(app, w, err, http.StatusInternalServerError, "Dateityp wird nicht unterstützt!")
		return
	}

	meta, _ := bimg.Metadata(orgImage.Image())
	metaJSON, _ := json.Marshal(meta)
	metaStr := string(metaJSON)

	date := time.Now()
	hash := sha256.Sum256(orgImage.Image())

	if orgImage.Type() == bimg.ImageTypes[bimg.JPEG] && meta.EXIF.Datetime != "" {
		exif_date_layout := "2006:01:02 15:04:05"
		date, err = time.Parse(exif_date_layout, meta.EXIF.Datetime)
		if err != nil {
			date = time.Now()
		}

	} else if r.PostFormValue("time") != "" {
		dateTimeJS, _ := strconv.ParseInt(r.PostFormValue("time"), 10, 64)
		date = time.Unix(dateTimeJS/1000, dateTimeJS-(dateTimeJS/1000))
	}

	mount, err := repository.ActiveMount(app.db)
	if err != nil {
		printError(app, w, err, http.StatusInternalServerError, appErrJsonCreationFailure)
		return
	}
	if mount.ID == 0 {
		printError(app, w, err, http.StatusInternalServerError, "No Mount defined!")
		return
	}

	fileModel := &model.File{}
	fileModel.UserUID = app.User().UID
	fileModel.Mount = *mount
	fileModel.Art = 1
	fileModel.UID = ksuid.New().String()
	fileModel.Date = date
	fileModel.Hash = hex.EncodeToString(hash[:])
	fileModel.Meta = &metaStr
	fileType := orgImage.Type()
	fileModel.Type = &fileType
	fileModel.MimeType = http.DetectContentType(fileHeader)
	fileModel.OrgFilename = multipartFileHeader.Filename
	fileModel.Path = repository.GetLocalPath(fileModel.Mount.Path, fileModel.UserUID, fileModel.Date)
	fileModel.Filename = repository.GetLocalFile(fileModel.Date, fileModel.UID, filepath.Ext(fileModel.OrgFilename))
	// REsize
	options := bimg.Options{
		Width:         app.conf.File.MaxImageWidth,
		Height:        app.conf.File.MaxImageWidth * 2,
		StripMetadata: false,
		Crop:          false,
		Enlarge:       true,
		Quality:       95,
		Interlace:     true,
	}
	newImage, err := orgImage.Process(options)
	if err != nil {
		printError(app, w, err, http.StatusInternalServerError, appErrJsonCreationFailure)
		return
	}

	if _, err = os.Stat(fileModel.Path); os.IsNotExist(err) {
		//os.MkdirAll(pfile.Path, os.ModeSticky|os.ModePerm)
		os.MkdirAll(fileModel.Path, os.ModePerm)
	}

	filepath := fileModel.Path + "/" + fileModel.Filename
	if err = bimg.Write(filepath, newImage); err != nil {
		printError(app, w, err, http.StatusInternalServerError, appErrJsonCreationFailure)
		return
	}

	metaC, _ := bimg.Metadata(newImage)
	fileModel.Width = &metaC.Size.Width
	fileModel.Height = &metaC.Size.Height
	fileModel.Size = len(newImage)
	//pfileModel.Vpath = uuid.NewV4().String() + uuid.NewV4().String() + uuid.NewV4().String() + uuid.NewV4().String() + uuid.NewV4().String() + uuid.NewV4().String() + uuid.NewV4().String() + ksuid.New().String() + uuid.NewV4().String() + hashCode + uuid.NewV4().String() + date.Format("20060102_150405") + ksuid.New().String()

	file, err := repository.CreateFile(app.db, fileModel)

	if err != nil {
		printError(app, w, err, http.StatusInternalServerError, appErrDataCreationFailure)
		return
	}
	app.logger.Info().Msgf("New File created: %d", file.UID)

	/*
		str, err1 := repository.EncodePath(*file)
		if err1 != nil {
			app.logger.Warn().Err(err1).Msg("")
		}
		app.logger.Info().Msgf("path: %s", str)
	*/
	w.WriteHeader(http.StatusCreated)
}
