package repository

import (
	"ghotos/model"

	"gorm.io/gorm"
)

func LoginUser(db *gorm.DB, email string, password string) (*model.User, error) {
	user := &model.User{}

	if err := db.Where("email = ? AND password = ?", email, password).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
