package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/YigitAtaMacit/StajDeneme/internal/auth"
	"github.com/YigitAtaMacit/StajDeneme/internal/db"
	"github.com/YigitAtaMacit/StajDeneme/internal/subject"
)

func main() {
	err := db.ConnectDB()
	if err != nil {
		fmt.Println("DB bağlantı hatası:", err)
		return
	}
	defer db.CloseDB()

	if err := db.CreateDB(); err != nil {
		fmt.Println("Subject tablosu oluşturulamadı:", err)
	}
	if err := db.CreateUserDB(); err != nil {
		fmt.Println("Kullanıcı tablosu oluşturulamadı:", err)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Post("/register", auth.RegisterHandler)
	router.Post("/login", auth.LoginHandler)

	router.Group(func(r chi.Router) {
		r.Use(auth.NewMiddleware)

		r.Get("/subjects", subject.GetSubject)
		r.Post("/subjects", subject.PostSubject)
		r.Delete("/subjects/{id}", subject.DeleteSubject)
		r.Put("/subjects/{id}", subject.PutSubject)
		r.Get("/subjects/{id}", subject.GetByID)
		r.Delete("/subjects", subject.DeleteAllSubjects)
	})

	fmt.Println("http://localhost:3000")
	http.ListenAndServe(":3000", router)
}
