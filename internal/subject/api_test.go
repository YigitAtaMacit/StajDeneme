package subject_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"

	"github.com/YigitAtaMacit/StajDeneme/internal/db"
	"github.com/YigitAtaMacit/StajDeneme/internal/subject"
)


type MockSubjectService struct {
	GetAllSubjectsFunc func(ctx context.Context) ([]db.Subject, error)
	AddSubjectFunc     func(ctx context.Context, s db.Subject) error
	UpdateFunc         func(ctx context.Context, s db.Subject) error
	DeleteFunc         func(ctx context.Context, id string) error
	DeleteAllFunc      func(ctx context.Context) error
	GetSubjectFunc     func(ctx context.Context, id string) (db.Subject, error)
}

func (m *MockSubjectService) GetAllSubjects(ctx context.Context) ([]db.Subject, error) {
	if m.GetAllSubjectsFunc != nil {
		return m.GetAllSubjectsFunc(ctx)
	}
	return nil, nil
}
func (m *MockSubjectService) AddSubject(ctx context.Context, s db.Subject) error {
	if m.AddSubjectFunc != nil {
		return m.AddSubjectFunc(ctx, s)
	}
	return nil
}
func (m *MockSubjectService) Update(ctx context.Context, s db.Subject) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, s)
	}
	return nil
}
func (m *MockSubjectService) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}
func (m *MockSubjectService) DeleteAll(ctx context.Context) error {
	if m.DeleteAllFunc != nil {
		return m.DeleteAllFunc(ctx)
	}
	return nil
}
func (m *MockSubjectService) GetSubject(ctx context.Context, id string) (db.Subject, error) {
	if m.GetSubjectFunc != nil {
		return m.GetSubjectFunc(ctx, id)
	}
	return db.Subject{}, nil
}

func TestGetSubject(t *testing.T) {
	mockService := &MockSubjectService{
		GetAllSubjectsFunc: func(ctx context.Context) ([]db.Subject, error) {
			return []db.Subject{
				{
					ID:          "1",
					UserID:      "efe",
					DoctorName:  "Dr. Test",
					Date:        "2025-08-20",
					Time:        "10:30",
					Description: "Kontrol",
				},
			}, nil
		},
	}
	handler := subject.NewSubjectHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/subjects", nil)
	rr := httptest.NewRecorder()

	handler.GetSubject(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var subjects []db.Subject
	err := json.NewDecoder(rr.Body).Decode(&subjects)
	assert.NoError(t, err)
	assert.Len(t, subjects, 1)
	assert.Equal(t, "Dr. Test", subjects[0].DoctorName)
}

func TestPostSubject(t *testing.T) {
	mockService := &MockSubjectService{
		AddSubjectFunc: func(ctx context.Context, s db.Subject) error {

			return nil
		},
	}
	handler := subject.NewSubjectHandler(mockService)

	sub := db.Subject{
		ID:          "123",
		UserID:      "efe",
		DoctorName:  "Dr. New",
		Date:        "2025-08-21",
		Time:        "09:00",
		Description: "İlk muayene",
	}
	body, _ := json.Marshal(sub)

	req := httptest.NewRequest(http.MethodPost, "/subjects", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handler.PostSubject(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var result db.Subject
	err := json.NewDecoder(rr.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, "Dr. New", result.DoctorName)
}

func TestPutSubject(t *testing.T) {
	mockService := &MockSubjectService{
		UpdateFunc: func(ctx context.Context, s db.Subject) error {
			assert.Equal(t, "999", s.ID)
			assert.Equal(t, "Dr. Updated", s.DoctorName)
			return nil
		},
	}
	handler := subject.NewSubjectHandler(mockService)

	sub := db.Subject{
		UserID:      "efe",
		DoctorName:  "Dr. Updated",
		Date:        "2025-08-22",
		Time:        "14:45",
		Description: "Takip",
	}
	body, _ := json.Marshal(sub)

	req := httptest.NewRequest(http.MethodPut, "/subjects/999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")


	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "999")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	rr := httptest.NewRecorder()
	handler.PutSubject(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestDeleteSubject(t *testing.T) {
	mockService := &MockSubjectService{
		DeleteFunc: func(ctx context.Context, id string) error {
			assert.Equal(t, "abc", id)
			return nil
		},
	}
	handler := subject.NewSubjectHandler(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/subjects/abc", nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "abc")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	rr := httptest.NewRecorder()
	handler.DeleteSubject(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestDeleteAllSubjects(t *testing.T) {
	mockService := &MockSubjectService{
		DeleteAllFunc: func(ctx context.Context) error {
			return nil
		},
	}
	handler := subject.NewSubjectHandler(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/subjects", nil)
	rr := httptest.NewRecorder()

	handler.DeleteAllSubjects(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestGetSubjectByID(t *testing.T) {
	mockService := &MockSubjectService{
		GetSubjectFunc: func(ctx context.Context, id string) (db.Subject, error) {
			assert.Equal(t, "456", id)
			return db.Subject{
				ID:          "456",
				UserID:      "efe",
				DoctorName:  "Dr. Cardiology",
				Date:        "2025-08-23",
				Time:        "11:15",
				Description: "Tahlil",
			}, nil
		},
	}
	handler := subject.NewSubjectHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/subjects/456", nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "456")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	rr := httptest.NewRecorder()
	handler.GetByID(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var result db.Subject
	err := json.NewDecoder(rr.Body).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, "Dr. Cardiology", result.DoctorName)
}
