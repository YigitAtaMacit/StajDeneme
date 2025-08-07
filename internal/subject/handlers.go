package subject

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/YigitAtaMacit/StajDeneme/internal/db"
	"github.com/YigitAtaMacit/StajDeneme/internal/service"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Service service.SubjectService
}

func NewSubjectHandler(s service.SubjectService) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) GetSubject(w http.ResponseWriter, r *http.Request) {
	subjectList, err := h.Service.GetAllSubjects(r.Context())
	if err != nil {
		http.Error(w, "Veritabanından veri alınamadı: "+err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, subjectList)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	subject, err := h.Service.GetSubject(r.Context(), id)
	if err != nil {
		http.Error(w, "Subject bulunamadı: "+err.Error(), http.StatusNotFound)
		return
	}
	jsonResponse(w, http.StatusOK, subject)
}

func (h *Handler) PostSubject(w http.ResponseWriter, r *http.Request) {
	var s db.Subject
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Geçersiz veri", http.StatusBadRequest)
		return
	}
	if err := h.Service.AddSubject(r.Context(), s); err != nil {
		http.Error(w, "Kayıt eklenemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusCreated, s)
}

func (h *Handler) PutSubject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var s db.Subject
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "Hatalı format", http.StatusBadRequest)
		return
	}
	s.ID = id

	err := h.Service.Update(r.Context(), s)
	if err != nil {
		if strings.Contains(err.Error(), "ID bulunamadı") {
			http.Error(w, "Kayıt bulunamadı", http.StatusNotFound)
			return
		}
		http.Error(w, "Veritabanı hatası: "+err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, s)
}

func (h *Handler) DeleteSubject(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.Service.Delete(r.Context(), id); err != nil {
		http.Error(w, "Subject silinemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteAllSubjects(w http.ResponseWriter, r *http.Request) {
	if err := h.Service.DeleteAll(r.Context()); err != nil {
		http.Error(w, "Kayıtlar silinemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func jsonResponse(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
