package model

import (
	"gorm.io/gorm"
)

type Mounts []*Mount
type Mount struct {
	gorm.Model
	Type   int
	Path   string
	Active int
}
