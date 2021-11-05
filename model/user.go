package model

import (
	g "ghotos/adapter/gorm"
)

type Users []*User
type User struct {
	g.ModelUID
	Status   int
	Email    string
	Password string
}

type UserDtos []*UserDto
type UserDto struct {
	UID      string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginForm struct {
	Email    string `json:"email" form:"required,max=255,email"`
	Password string `json:"password"  form:"required"`
}
type UserRegisterForm struct {
	Email    string `json:"email" form:"required,max=255,email"`
	Password string `json:"password"  form:"required,min=8,max=20"`
}

func (f *UserLoginForm) ToModel() (*User, error) {

	return &User{
		Email:    f.Email,
		Password: f.Password,
	}, nil
}
func (f *UserRegisterForm) ToModel() (*User, error) {

	return &User{
		Email:    f.Email,
		Password: f.Password,
	}, nil
}
