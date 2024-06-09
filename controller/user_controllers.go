package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/temuka-api-service/config"
	"github.com/temuka-api-service/models"
	"gorm.io/gorm"
)

var db *gorm.DB = config.GetDBInstance()

func SearchUsers(w http.ResponseWriter, r *http.Request) {
	var users []models.User

	if err := db.Find(&users).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(users) == 0 {
		http.Error(w, "No user found", http.StatusNotFound)
		return
	}

	response := struct {
		Message string        `json:"message"`
		Data    []models.User `json:"data"`
	}{
		Message: "User has been created",
		Data:    users,
	}
	json.NewEncoder(w).Encode(response)
}

func GetUserDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDstr := vars["id"]

	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	var user models.User
	res := db.First(&user, userID)
	if res.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	db.Find(&user)
	response := struct {
		Message string      `json:"message"`
		Data    models.User `json:"data"`
	}{
		Message: "User detail has been retrieved",
		Data:    user,
	}

	json.NewEncoder(w).Encode(response)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	newUser := models.User{
		Username: requestBody.Username,
		Email:    requestBody.Email,
		Password: requestBody.Password,
	}

	db.Create(&newUser)
	response := struct {
		Message string      `json:"message"`
		Data    models.User `json:"data"`
	}{
		Message: "User has been created",
		Data:    newUser,
	}
	json.NewEncoder(w).Encode(response)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDstr := vars["user_id"]

	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	var user models.User
	res := db.First(&user, userID)
	if res.Error != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	db.Save(&user)
	response := struct {
		Message string      `json:"message"`
		Data    models.User `json:"data"`
	}{
		Message: "User has been updated",
		Data:    user,
	}

	json.NewEncoder(w).Encode(response)
}
