package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)
type Subject struct{
    ID string `json:"id"`
	Name string `json:"name"`
    Age int `json:"age"`
}
var subjects map[string]Subject

func main() {
	err := ConnectDB()
	if err != nil {
		fmt.Println("DB bağlantı hatası:", err)
		return
	}
	defer CloseDB()
	if err := CreateDB(); err != nil {
		fmt.Println("hata:", err)
	}
		err = LoadSubjectsFromDB()
	if err != nil {
		fmt.Println("Subject'ler yüklenemedi:", err)
	} else {
		fmt.Println("Subject'ler başarıyla yüklendi.")
	}
	fmt.Println("Veritabanı bağlantısı başarılı!")
	fmt.Println("http://localhost:3000")
	router :=chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/subjects",GetSubject)
	router.Post("/subjects",PostSubject)
	router.Delete("/subjects/{id}",DeleteSubject)
	router.Put("/subjects/{id}",PutSubject)
	router.Get("/subjects/{id}",GetbyID)
	router.Delete("/subjects", DeleteAllSubjects)
	http.ListenAndServe(":3000",router)
	

}


