package models

import (
	"time"

	"github.com/usama1031/book-management/pkg/config"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	First_name    *string   `json:"first_name" validate:"required,min=2,max=100"`
	Last_name     *string   `json:"last_name" validate:"required,min=2,max=100"`
	Password      *string   `json:"Password" validate:"required,min=6"`
	Email         *string   `json:"email" validate:"required"`
	Avatar        *string   `json:"avatar"`
	Phone         *string   `json:"phone" validate:"required"`
	Token         *string   `json:"token"`
	Refresh_Token *string   `json:"refresh_token"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"updated_at"`
	User_id       string    `json:"user_id"`
	User_type     *string   `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
}

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&User{})
}
