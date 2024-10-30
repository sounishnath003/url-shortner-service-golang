package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sounishnath003/url-shortner-service-golang/internal/core"
	"golang.org/x/crypto/bcrypt"
)

// UserLoginDto struct
//
// This struct represents the data required to login a user.
//
// - Email: The user's email address
// - Password: The user's password
//
// Example:
//
// ```json
//
//	{
//	  "email": "john.doe@example.com",
//	  "password": "password123"
//	}
type UserLoginDto struct {
	ID       int    `json:"-"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginHandler handles the user login request.
//
// This function handles the user login request and generates a JWT token for the user.
//
// The request body should contain the following fields:
//
// - email: The user's email address
// - password: The user's password
//
// If the request is successful, the response will contain the following fields:
//
// - token: The JWT token for the user
//
// If the request fails, the response will contain an error message.
//
// Example Request:
//
// ```json
//
//	{
//	  "email": "john.doe@example.com",
//	  "password": "password123"
//	}
//
// ```
//
// Example Response:
//
// ```json
//
//	{
//	  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ1cmwtc2hvcnRuZXItc2VydmljZSIsInN1YiI6ImFkbWluQGV4YW1wbGUuY29tIiwiYXVkIjpbInVzZXIiLCJ1cmwtc2hvcnRuZXItc2VydmljZSJdLCJpYXQiOjE2NzY0MjYzNjcsImV4cCI6MTY3NjUyMjM2N30.Y293464856485648564856485648564856485648564856485648564856485648564856485648564856485648564856485648564856485648564856485648564856485648564856485648564
//
// ```
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Grab the context
	co := r.Context().Value("co").(*core.Core)

	var user UserLoginDto
	json.NewDecoder(r.Body).Decode(&user)
	defer r.Body.Close()

	if user.Email == "" || user.Password == "" {
		WriteError(w, http.StatusBadRequest, errors.New("invalid credentials"))
		return
	}
	// Create JWT token for the user.
	token, err := createJwtToken(co, user.Email, user.Password)
	if err != nil {
		WriteError(w, http.StatusUnauthorized, errors.New("username or password are incorrect"))
		return
	}

	WriteJson(w, http.StatusOK, map[string]string{
		"token": token,
	})
}

// CreateNewUserDto struct
//
// This struct represents the data required to create a new user.
//
// - Name: The user's name
// - Email: The user's email address
// - Password: The user's password
//
// Example:
//
// ```json
//
//	{
//	  "name": "John Doe",
//	  "email": "john.doe@example.com",
//	  "password": "password123"
//	}
type CreateNewUserDto struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignupHandler handles the user signup request.
//
// This function handles the user signup request and creates a new user in the database.
// It also generates a JWT token for the user and returns it in the response.
//
// The request body should contain the following fields:
//
// - name: The user's name
// - email: The user's email address
// - password: The user's password
//
// If the request is successful, the response will contain the following fields:
//
// - token: The JWT token for the user
//
// If the request fails, the response will contain an error message.
//
// Example Request:
//
// ```json
//
//	{
//	  "name": "John Doe",
//	  "email": "john.doe@example.com",
//	  "password": "password123"
//	}
//
// ```
//
// Example Response:
//
// ```json
//
//	{
//	  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJ1cmwtc2hvcnRuZXItc2VydmljZSIsInN1YiI6ImFkbWluQGV4YW1wbGUuY29tIiwiYXVkIjpbInVzZXIiLCJ1cmwtc2hvcnRuZXItc2VydmljZSJdLCJpYXQiOjE2NzY0MjYzNjcsImV4cCI6MTY3NjUyMjM2N30.Y293464856485648564856485648564856485648564856485648564856485648564856485648564856485648564856485648564856485648564856485648564856485648564856485648564
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	// Grab the context
	co := r.Context().Value("co").(*core.Core)

	var user CreateNewUserDto
	json.NewDecoder(r.Body).Decode(&user)
	defer r.Body.Close()

	if user.Name == "" || user.Email == "" || user.Password == "" {
		WriteError(w, http.StatusBadRequest, errors.New("required fields are not provided"))
		return
	}
	// Generate a bcrypt hashed password
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Register the user.
	_, err = co.QueryStmts.CreateNewUser.Exec(user.Name, user.Email, hashed)
	if err != nil {
		WriteError(w, http.StatusBadRequest, errors.New("user already exists. or request is malformed"))
		return
	}

	// Create JWT token for the user.
	token, err := createJwtToken(co, user.Email, user.Password)
	if err != nil {
		WriteError(w, http.StatusBadRequest, err)
		return
	}

	WriteJson(w, http.StatusOK, map[string]string{
		"token": token,
	})
}

// createJwtToken helps to generate an access token for the user.
//
// JWT Token which validates the user credentials.
func createJwtToken(co *core.Core, email, password string) (string, error) {
	var userInDB UserLoginDto
	// Check the user present in DB
	co.QueryStmts.GetUserByEmail.QueryRow(email).Scan(&userInDB.ID, &userInDB.Email, &userInDB.Password)
	co.Lo.Info("checking user exists", "email", email)
	// If no user found.
	if userInDB.ID == 0 || userInDB.Email == "" {
		return "", errors.New("username or password are incorrect")
	}

	// Perform the salted hashed password checks using bcrypt.
	if err := bcrypt.CompareHashAndPassword([]byte(userInDB.Password), []byte(password)); err != nil {
		return "", err
	}

	// Generate the claims for the user
	claims := jwt.MapClaims{
		"iss":       "url-shortner-service",
		"sub":       userInDB.Email,
		"userEmail": userInDB.Email,
		"userID":    strconv.Itoa(userInDB.ID),
		"aud":       []string{"user", "url-shortner-service"},
		"iat":       jwt.NewNumericDate(time.Now()),
		"exp":       jwt.NewNumericDate(time.Now().Add(time.Duration(1 * time.Hour))),
	}
	// Generate a signed claimed token
	claimTok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return claimTok.SignedString([]byte(co.JwtSecret))
}
