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
    if err := CreateUserDB(); err != nil {
	fmt.Println("Kullanıcı tablosu oluşturulamadı:", err)
}
	fmt.Println("Veritabanı bağlantısı başarılı!")
	fmt.Println("http://localhost:3000")

	router :=chi.NewRouter()
	router.Use(middleware.Logger)
	router.Post("/register", RegisterHandler)
	router.Post("/login", loginHandler)
	

	

    router.Group(func(r chi.Router) {
	r.Use(NewMiddleware) 

	r.Get("/subjects",GetSubject)
	r.Post("/subjects",PostSubject)
	r.Delete("/subjects/{id}",DeleteSubject)
	r.Put("/subjects/{id}",PutSubject)
	r.Get("/subjects/{id}",GetbyID)
	r.Delete("/subjects", DeleteAllSubjects)
})

	
	
	http.ListenAndServe(":3000",router)
	

}


