package model

import (
	g "ghotos/adapter/gorm"
)

type Users []*User
type User struct {
	g.ModelUID
	Email    string
	Password string
}

type UserDtos []*UserDto
type UserDto struct {
	UID      string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"Password"`
}

type UserLoginForm struct {
	Email    string `json:"email" form:"required,max=255"`
	Password string `json:"Password"  form:"required"`
}

func (f *UserLoginForm) ToModel() (*User, error) {

	return &User{
		Email:    f.Email,
		Password: f.Password,
	}, nil
}
