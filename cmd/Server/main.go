package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/YigitAtaMacit/StajDeneme/internal/service"
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


	subjectRepo := db.NewSubjectRepo(db.DB)
	subjectService := service.NewSubjectService(subjectRepo)
	subjectHandler := subject.NewSubjectHandler(subjectService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)


	r.Post("/register", auth.RegisterHandler)
	r.Post("/login", auth.LoginHandler)

	
	r.Route("/subjects", func(r chi.Router) {
		r.Use(auth.NewMiddleware)

		r.Get("/", subjectHandler.GetSubject)
		r.Get("/{id}", subjectHandler.GetByID)
		r.Post("/", subjectHandler.PostSubject)
		r.Put("/{id}", subjectHandler.PutSubject)
		r.Delete("/{id}", subjectHandler.DeleteSubject)
		r.Delete("/", subjectHandler.DeleteAllSubjects)
	})


	fmt.Println("Server başlatıldı: http://localhost:3000")
	http.ListenAndServe(":3000", r)
}