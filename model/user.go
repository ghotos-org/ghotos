package model

import (
	g "ghotos/adapter/gorm"
	"time"
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

type UserRegisterEmailFormDto struct {
	Email string `json:"email"`
}
type UserRegisterEmailForm struct {
	Email string    `json:"email" form:"required,max=255,email"`
	Date  time.Time `json:"date"`
}

func (f *UserLoginForm) ToModel() (*User, error) {
	return &User{
		Email:    f.Email,
		Password: f.Password,
	}, nil
}
func (f *UserRegisterEmailForm) ToModel() (*User, error) {
	return &User{
		Email: f.Email,
	}, nil
}

func (u UserRegisterEmailForm) ToDto() *UserRegisterEmailFormDto {
	return &UserRegisterEmailFormDto{
		Email: u.Email,
	}
}
