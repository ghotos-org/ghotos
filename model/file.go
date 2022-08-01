package model

import (
	"time"

	g "ghotos/adapter/gorm"
	"ghotos/util/tools"
)

type Files []*File
type File struct {
	g.ModelUID
	MountID     int
	Mount       Mount
	Status      int
	Date        time.Time
	MimeType    string
	Size        int
	OrgFilename string
	Filename    string
	Path        string
	Extension   string
	UserUID     string
	User        User
	Art         int
	Type        *string
	Hash        string
	Width       *int
	Height      *int
	Meta        *string
}

type FilesDayCount struct {
	Day   string `json:"day"`
	Count int    `json:"count"`
}

type FilesGallery struct {
	Days       []FilesDayCount `json:"days"`
	LastUpdate int64           `json:"updated"`
}
type FilePaths []*FilePath
type FilePath struct {
	UID         string `json:"uid"`
	Date        int64  `json:"date"`
	CreatedAt   int64  `json:"created"`
	UpdatedAt   int64  `json:"updated"`
	Hash        string `json:"Hash"`
	OrgFilename string `json:"filename"`
	MimeType    string `json:"mime_type"`
}

type FilesSrcs []*FilesSrc
type FilesSrc struct {
	UID string  `json:"id"`
	Src *string `json:"src,omitempty"`
	Day *string `json:"day,omitempty"`
}

type FileSimples []*FileSimple
type FileSimple struct {
	UID    string  `json:"id"`
	Src    *string `json:"src,omitempty"`
	Day    *string `json:"day,omitempty"`
	Height *int    `json:"height,omitempty"`
	Width  *int    `json:"width,omitempty"`
}

func (f File) ToFilePath() *FilePath {
	return &FilePath{

		UID:         f.UID,
		Date:        tools.MakeTimestamp(f.Date),
		CreatedAt:   tools.MakeTimestamp(f.CreatedAt),
		UpdatedAt:   tools.MakeTimestamp(f.UpdatedAt),
		Hash:        f.Hash,
		OrgFilename: f.OrgFilename,
		MimeType:    f.MimeType,
	}
}
func (fs Files) ToFilePath() FilePaths {
	dtos := make([]*FilePath, len(fs))
	for i, f := range fs {
		dtos[i] = f.ToFilePath()
	}
	return dtos
}

func (f File) ToDtoSimple() *FileSimple {
	day := f.Date.Format("20060102")
	return &FileSimple{
		UID:    f.UID,
		Day:    &day,
		Height: f.Height,
		Width:  f.Width,
	}
}
