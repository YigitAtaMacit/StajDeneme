package subject

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/YigitAtaMacit/StajDeneme/internal/db"
)

var subjects map[string]db.Subject

func GetSubject(w http.ResponseWriter, r *http.Request) {
	subjectList, err := db.GetAllSubjects(db.DB)
	if err != nil {
		http.Error(w, "Veritabanından veri alınamadı: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subjectList)
}

func GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	subject, err := db.GetSubjectByID(db.DB, id)
	if err != nil {
		http.Error(w, "Subject bulunamadı: "+err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subject)
}

func PostSubject(w http.ResponseWriter, r *http.Request) {
	var s db.Subject
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Geçersiz veri", http.StatusBadRequest)
		return
	}
	if s.ID == "" {
		http.Error(w, "Boş ID", http.StatusBadRequest)
		return
	}
	if _, exists := subjects[s.ID]; exists {
		http.Error(w, "Aynı ID zaten var", http.StatusBadRequest)
		return
	}
	if err := db.InsertSubject(db.DB, s); err != nil {
		http.Error(w, "Kayıt eklenemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}
	subjects[s.ID] = s
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

func PutSubject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, exists := subjects[id]; !exists {
		http.Error(w, "Böyle bir ID yok", http.StatusNotFound)
		return
	}
	var subject db.Subject
	if err := json.NewDecoder(r.Body).Decode(&subject); err != nil {
		http.Error(w, "Hatalı format", http.StatusBadRequest)
		return
	}
	subject.ID = id
	if err := db.UpdateSubject(db.DB, subject); err != nil {
		if strings.Contains(err.Error(), "ID bulunamadı") {
			http.Error(w, "Kayıt bulunamadı", http.StatusNotFound)
			return
		}
		http.Error(w, "Veritabanı hatası", http.StatusInternalServerError)
		return
	}
	subjects[id] = subject
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subject)
}

func DeleteSubject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := db.DeleteSubjectByID(db.DB, id); err != nil {
		http.Error(w, "Subject silinemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}
	delete(subjects, id)
	w.WriteHeader(http.StatusNoContent)
}

func DeleteAllSubjects(w http.ResponseWriter, r *http.Request) {
	_, err := db.DB.Exec(context.Background(), "DELETE FROM subject")
	if err != nil {
		http.Error(w, "Kayıtlar silinemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}
	subjects = make(map[string]db.Subject)
	w.WriteHeader(http.StatusNoContent)
}

func LoadSubjectsFromDB() error {
	dbSubjects, err := db.GetAllSubjects(db.DB)
	if err != nil {
		return fmt.Errorf("veritabanından subjectler yüklenemedi: %w", err)
	}
	subjects = make(map[string]db.Subject)
	for _, s := range dbSubjects {
		subjects[s.ID] = s
	}
	return nil
}