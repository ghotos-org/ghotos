package repository

import (
	"ghotos/model"

	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
)

func ReadUser(db *gorm.DB, uid string) (*model.User, error) {
	user := &model.User{}

	if err := db.Where("uid = ?", uid).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func ReadUserByEmail(db *gorm.DB, email string) (*model.User, error) {
	user := &model.User{}

	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func CreateUser(db *gorm.DB, user *model.User) (*model.User, error) {
	if user.UID == "" {
		user.UID = ksuid.New().String()
	}
	if err := db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func UpdateUser(db *gorm.DB, user *model.User) error {
	if err := db.Save(user).Error; err != nil {
		return err
	}
	return nil
}
