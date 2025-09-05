package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"backend/db"
	"backend/middleware"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Username == "" || req.Password == "" {
		middleware.WriteJSONError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		middleware.WriteJSONError(w, http.StatusInternalServerError, "Hash error")
		return
	}

	_, err = db.DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", req.Username, string(hashed))
	if err != nil {
		middleware.WriteJSONError(w, http.StatusBadRequest, "Username taken")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered"})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteJSONError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	var id int
	var hashed string
	err := db.DB.QueryRow("SELECT id, password FROM users WHERE username=$1", req.Username).Scan(&id, &hashed)
	if err != nil {
		middleware.WriteJSONError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(hashed), []byte(req.Password)) != nil {
		middleware.WriteJSONError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, _ := createJWT(id)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func createJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(middleware.JwtSecret)
}
