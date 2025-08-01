package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"

	"github.com/YigitAtaMacit/StajDeneme/internal/db"
)

var tokenKey = []byte("anahtar")

type information struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID       string
	Username string
	Password string
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var info information
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		http.Error(w, "Geçersiz istek", http.StatusBadRequest)
		return
	}

	var storedPassword string
	err := db.DB.QueryRow(context.Background(), "SELECT password FROM users WHERE username=$1", info.Username).Scan(&storedPassword)
	if err != nil {
		http.Error(w, "Kullanıcı bulunamadı", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(info.Password))
	if err != nil {
		http.Error(w, "Şifre hatalı", http.StatusUnauthorized)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":  info.Username,
		"timelimit": time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenStr, err := token.SignedString(tokenKey)
	if err != nil {
		http.Error(w, "Token oluşturulamadı", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenStr,
	})
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Geçersiz veri", http.StatusBadRequest)
		return
	}

	user.ID = uuid.New().String()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Şifre hashlenemedi", http.StatusInternalServerError)
		return
	}

	_, err = db.DB.Exec(context.Background(),
		"INSERT INTO users (id, username, password) VALUES ($1, $2, $3)",
		user.ID, user.Username, string(hashedPassword),
	)
	if err != nil {
		fmt.Println("Kayıt hatası:", err)
		http.Error(w, "Kayıt başarısız", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Kayıt başarılı",
	})
}