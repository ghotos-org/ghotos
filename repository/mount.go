package repository

import (
	"ghotos/model"

	"gorm.io/gorm"
)

func ActiveMount(db *gorm.DB) (*model.Mount, error) {
	mount := &model.Mount{}
	db.First(&mount, "active = 1")

	return mount, nil
}
