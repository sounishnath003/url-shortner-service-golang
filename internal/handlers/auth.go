package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sounishnath003/url-shortner-service-golang/internal/core"
	"golang.org/x/crypto/bcrypt"
)

type UserLoginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

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

type CreateNewUserDto struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

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
	co.QueryStmts.GetUserByEmail.QueryRow(email).Scan(&userInDB.Email, &userInDB.Password)
	co.Lo.Info("checking user exists", "email", email)
	// If no user found.
	if userInDB.Email == "" {
		return "", errors.New("username or password are incorrect")
	}

	// Perform the salted hashed password checks using bcrypt.
	if err := bcrypt.CompareHashAndPassword([]byte(userInDB.Password), []byte(password)); err != nil {
		return "", err
	}

	// Generate the claims for the user
	claims := jwt.RegisteredClaims{
		Issuer:    "url-shortner-service",
		Subject:   userInDB.Email,
		Audience:  []string{"user", "url-shortner-service"},
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(1 * time.Hour))),
	}
	// Generate a signed claimed token
	claimTok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return claimTok.SignedString([]byte(co.JwtSecret))
}
