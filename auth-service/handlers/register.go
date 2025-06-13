package handlers

import (
	"auth-service/helpers"
	"auth-service/middleware"
	data "auth-service/models"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/asaskevich/govalidator"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.ErrorJSON(w, http.ErrNotSupported, http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := helpers.ReadJSON(w, r, &req); err != nil {
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// 1. Validating the input
	if req.Email == "" || req.Password == "" {
		helpers.ErrorJSON(w, errors.New("email and password are required"), http.StatusBadRequest)
		return
	}

	if !govalidator.IsEmail(req.Email) {
		helpers.ErrorJSON(w, errors.New("invalid email address"), http.StatusBadRequest)
		return
	}

	// 2. Hashing the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// 3. Storing in database
	user := data.User{
		Email:           req.Email,
		Password:        string(hashedPassword),
		EmailVerifiedAt: time.Time{},
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Get database connection from context
	db, ok := r.Context().Value(middleware.DBContextKey).(*sql.DB)
	if !ok {
		helpers.ErrorJSON(w, errors.New("database connection not found"), http.StatusInternalServerError)
		return
	}

	userRepo := data.NewUserRepository(db)
	_, err = userRepo.CreateUser(r.Context(), user)
	if err != nil {
		helpers.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// 4. giving a response
	response := RegisterResponse{
		Success: true,
		Message: "User registered successfully",
	}

	helpers.WriteJSON(w, http.StatusOK, response, nil)
}
