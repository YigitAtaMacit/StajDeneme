package main

import (
	"encoding/json"
	"net/http"
    "fmt"
	"github.com/go-chi/chi/v5"
	"strings"
	"context"
)
func GetSubject(w http.ResponseWriter, r *http.Request){
	subjects, err := GetAllSubjects(DB)
	if err != nil {
		http.Error(w, "Veritabanından veri alınamadı: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subjects)
}
func GetbyID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	subject, err := GetSubjectByID(DB, id)
	if err != nil {
		http.Error(w, "Subject bulunamadı: "+err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subject)
}

func DeleteSubject(w http.ResponseWriter,r *http.Request){
	id:=chi.URLParam(r,"id")
    err := DeleteSubjectByID(DB, id)
	if err != nil {
		http.Error(w, "Veritabanından silinemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}
	delete(subjects,id)
	w.WriteHeader(http.StatusNoContent)
}
func PostSubject(w http.ResponseWriter, r *http.Request){
	var s Subject;
    if err:= json.NewDecoder(r.Body).Decode(&s); err!=nil{
		return
	}
	if s.ID==""{
		http.Error(w,"Boş id",http.StatusBadRequest)
		return 
	}
	if _, exists := subjects[s.ID]; exists {
	    http.Error(w, "Aynı id zaten var", http.StatusBadRequest)
	    return
    }
	err := InsertSubject(DB, s)
	subjects[s.ID]=s
	
	if err != nil {
		fmt.Println("Kayıt eklenemedi:", err)
	} else {
		fmt.Println("Kayıt başarıyla eklendi.")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)

}
func PutSubject(w http.ResponseWriter,r *http.Request){
	id :=chi.URLParam(r,"id")
	if _,exists:=subjects[id]; !exists{
		http.Error(w,"Böyle bir id yok",http.StatusNotFound)
		return
	}
    var subject Subject
	if err :=json.NewDecoder(r.Body).Decode(&subject); err!=nil{
		http.Error(w,"Hatalı format",http.StatusBadRequest)
		return
	}

	subject.ID = id 

	err := UpdateSubject(DB, subject)
	if err != nil {
		if strings.Contains(err.Error(), "ID bulunamadı") {
			http.Error(w, "Kayıt bulunamadı", http.StatusNotFound)
			return
		}
		http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
		fmt.Println("Güncelleme hatası:", err)
		return
	}


	subjects[id]=subject
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subject)

}

func LoadSubjectsFromDB() error {
	dbSubjects, err := GetAllSubjects(DB)
	if err != nil {
		return fmt.Errorf("Veritabanından subjectler yüklenemedi: %w", err)
	}

	subjects = make(map[string]Subject)
	for _, s := range dbSubjects {
		subjects[s.ID] = s
	}
	return nil
}

func DeleteAllSubjects(w http.ResponseWriter, r *http.Request) {
	_, err := DB.Exec(context.Background(), "DELETE FROM subject")
	if err != nil {
		http.Error(w, "Kayıtlar silinemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}

	subjects = make(map[string]Subject) // in-memory map'i de temizle

	w.WriteHeader(http.StatusNoContent)
}