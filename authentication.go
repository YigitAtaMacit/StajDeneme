package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
    "fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"context"
	"github.com/google/uuid"
)

var tokenkey = []byte("anahtar")

type information struct{
    Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID       string
	Username string
	Password string
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var info information
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		http.Error(w, "Geçersiz istek", http.StatusBadRequest)
		return
	}


	var storedPassword string
	err := DB.QueryRow(context.Background(), "SELECT password FROM users WHERE username=$1", info.Username).Scan(&storedPassword)
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

	tokenStr, err := token.SignedString(tokenkey)
	if err != nil {
		http.Error(w, "Token oluşturulamadı", http.StatusInternalServerError)
		return
	}


	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenStr,
	})
}



func NewMiddleware(next http.Handler)http.Handler{
    return http.HandlerFunc(func(w http.ResponseWriter,r*http.Request){

		authHeader:= r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Yetki hatası",http.StatusUnauthorized)
			return
		}

		splitToken:= strings.Split(authHeader, " ")

		if len(splitToken)!=2 || splitToken[0] !="Bearer"{
			http.Error(w,"Geçersiz yetki idsi",http.StatusUnauthorized)
	        return
		}

		tokenstring :=splitToken[1]
		token, err := jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("beklenmeyen imzalama yöntemi: %v", token.Header["alg"])
			}
			return tokenkey, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Geçersiz token", http.StatusUnauthorized)
			return
		}

        next.ServeHTTP(w, r)




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

    _, err = DB.Exec(context.Background(),
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