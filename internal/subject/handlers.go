package subject

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/YigitAtaMacit/StajDeneme/internal/db"
	"github.com/YigitAtaMacit/StajDeneme/internal/service"
)

var svc = service.NewSubjectService(db.NewSubjectRepo(db.DB))

func GetSubject(w http.ResponseWriter, r *http.Request) {
	subjectList, err := svc.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Veritabanından veri alınamadı: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(subjectList)
}

func GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	subject, err := svc.GetSubject(r.Context(), id)
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
	if err := svc.AddSubject(r.Context(), s); err != nil {
		http.Error(w, "Kayıt eklenemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

func PutSubject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var s db.Subject
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Hatalı format", http.StatusBadRequest)
		return
	}
	s.ID = id
	if err := svc.Update(r.Context(), s); err != nil {
		http.Error(w, "Güncelleme başarısız: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s)
}

func DeleteSubject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := svc.Delete(r.Context(), id); err != nil {
		http.Error(w, "Subject silinemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func DeleteAllSubjects(w http.ResponseWriter, r *http.Request) {
	if err := svc.DeleteAll(r.Context()); err != nil {
		http.Error(w, "Kayıtlar silinemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
