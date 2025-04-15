package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/usama1031/book-management/pkg/config"
	"github.com/usama1031/book-management/pkg/helpers"
	"github.com/usama1031/book-management/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func HashPassword(password string) string {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(hashedBytes)
}

func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))

	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("Email or passowrd is incorrect!")
		check = false
	}

	return check, msg
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Invalid JSON:"+err.Error(), http.StatusBadRequest)
		return
	}

	if err := validate.Struct(user); err != nil {
		http.Error(w, "Validation Error:"+err.Error(), http.StatusBadRequest)
		return
	}

	var existingUser models.User

	result := config.GetDB().Where("email = ? OR phone = ?", user.Email, user.Phone).First(&existingUser)

	if result.RowsAffected > 0 {
		http.Error(w, "Email or phone already exists", http.StatusConflict)
		return
	}

	hashedPassword := HashPassword(*user.Password)
	user.Password = &hashedPassword

	now := time.Now()

	user.Created_at = now
	user.Updated_at = now

	user.User_id = uuid.New().String()

	token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, user.User_id)

	user.Token = &token

	user.Refresh_Token = &refreshToken

	if err := config.GetDB().Create(&user).Error; err != nil {
		http.Error(w, "Failed to create user:"+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "User created successfully",
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	var foundUser models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input:"+err.Error(), http.StatusBadRequest)
		return
	}

	result := config.GetDB().Where("email = ?", user.Email).First(&foundUser)

	if result.RowsAffected == 0 {
		http.Error(w, "Email or password is incorrect!", http.StatusUnauthorized)
		return
	}

	isValid, msg := VerifyPassword(*user.Password, *foundUser.Password)

	if !isValid {
		http.Error(w, msg, http.StatusUnauthorized)
		return
	}

	token, refreshToken, _ := helpers.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, foundUser.User_id)

	foundUser.Token = &token
	foundUser.Refresh_Token = &refreshToken

	foundUser.Updated_at = time.Now()

	config.GetDB().Save(&foundUser)

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   3600,
	})

	json.NewEncoder(w).Encode(foundUser)

}
