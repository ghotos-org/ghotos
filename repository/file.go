package repository

import (
	"encoding/hex"
	"encoding/json"
	"ghotos/model"
	"ghotos/util/tools"
	"path/filepath"
	"time"

	"gorm.io/gorm"
)

var key = []byte("secret")

func GetLocalPath(mount string, useruid string, date time.Time) string {
	return mount + "/" + useruid + "/" + date.Format("2006/01")
}
func GetLocalFile(date time.Time, uid string, ext string) string {
	return date.Format("20060102_150405") + "_" + uid + ext
}
func GetLocalFilePath(file model.File) string {
	return GetLocalPath(file.Mount.Path, file.UserUID, file.Date) + "/" + GetLocalFile(file.Date, file.UID, filepath.Ext(file.OrgFilename))
}

func EncodePath(modelFile model.File) (string, error) {

	file := modelFile.ToFilePath()
	bytes, err := json.Marshal(file)
	if err != nil {
		return "", err
	}
	ct, err := tools.Encrypt(bytes, key, modelFile.UID)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(ct), nil

}
func DecodePath(encStr string) (*model.FilePath, error) {
	var file model.FilePath

	str, err := hex.DecodeString(encStr)
	if err != nil {
		return nil, err
	}

	bytes, err := tools.Decrypt(str, key)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(bytes, &file)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func LastFileUpdate(db *gorm.DB, userUID string) (*time.Time, error) {

	var lastUpdate time.Time
	db.Model(&model.File{}).Select("updated_at").Order("updated_at desc").First(&lastUpdate, "user_uid = ?", userUID)

	return &lastUpdate, nil
}

func CreateFile(db *gorm.DB, file *model.File) (*model.File, error) {
	if err := db.Create(file).Error; err != nil {
		return nil, err
	}
	return file, nil
}

func ReadFile(db *gorm.DB, uid string) (*model.File, error) {
	file := &model.File{}
	//if err := db.Where("uid = ?", uid).First(&file).Error; err != nil {
	if err := db.Preload("Mount").Find(&file, "uid = ?", uid).Error; err != nil {
		return nil, err
	}
	return file, nil
}
func ReadFileByUser(db *gorm.DB, uid string, userUID string) (*model.File, error) {
	file := &model.File{}
	//if err := db.Where("uid = ?", uid).First(&file).Error; err != nil {
	if err := db.Preload("Mount").Find(&file, "uid = ? AND  user_uid = ?", uid, userUID).Error; err != nil {
		return nil, err
	}
	return file, nil
}

func DeleteFile(db *gorm.DB, uid string, userUID string) error {
	file := &model.File{}
	if err := db.Where("uid = ? AND user_uid = ?", uid, userUID).Delete(&file).Error; err != nil {
		return err
	}
	return nil
}

func GalleryDays(db *gorm.DB, userUID string) (*model.FilesGallery, error) {
	gallery := &model.FilesGallery{}

	var sql_date_fomrat string
	if db.Config.Dialector.Name() == "mysql" {
		sql_date_fomrat = "DATE_FORMAT(date,'%Y%m%d')"
	} else {
		sql_date_fomrat = "strftime('%Y%m%d', date)"
	}

	db.Model(&model.File{}).Select(sql_date_fomrat+" as day, count(*) as count").Where("user_uid = ?", userUID).Group("day").Order(sql_date_fomrat + " desc").Scan(&gallery.Days)

	if gallery.Days != nil {
		lastUpdate, err := LastFileUpdate(db, userUID)
		if err != nil {
			return nil, err
		}
		gallery.LastUpdate = tools.MakeTimestamp(*lastUpdate)

	}

	return gallery, nil
}

func ListGalleryDayFile(db *gorm.DB, day string, userUID string) (model.Files, error) {

	files := make([]*model.File, 0)

	var sql_date_fomrat string
	if db.Config.Dialector.Name() == "mysql" {
		sql_date_fomrat = "DATE_FORMAT(date,'%Y%m%d')"
	} else {
		sql_date_fomrat = "strftime('%Y%m%d', date)"
	}

	if err := db.Find(&files, sql_date_fomrat+" = ? and user_uid = ?", day, userUID).Select("uid,id,date,created_at,updated_at,hash,org_filename,mime_type").Order("date, uid").Limit(100000).Error; err != nil {
		return nil, err
	}

	return files, nil
}
