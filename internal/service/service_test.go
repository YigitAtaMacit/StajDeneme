package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/YigitAtaMacit/StajDeneme/internal/db"
	"github.com/YigitAtaMacit/StajDeneme/internal/service"
)

type MockSubjectRepo struct {
	InsertCalled      bool
	GetAllCalled      bool
	DeleteCalled      bool
	UpdateCalled      bool
	DeleteAllCalled   bool
	GetByIDCalled     bool


	ShouldFail bool
}

func (m *MockSubjectRepo) InsertSubject(ctx context.Context, s db.Subject) error {
	m.InsertCalled = true
	if m.ShouldFail {
		return errors.New("insert failed")
	}
	return nil
}

func (m *MockSubjectRepo) GetAllSubjects(ctx context.Context) ([]db.Subject, error) {
	m.GetAllCalled = true
	if m.ShouldFail {
		return nil, errors.New("get all failed")
	}
	return []db.Subject{
		{ID: "1", Name: "Test", Age: 25},
	}, nil
}

func (m *MockSubjectRepo) DeleteSubjectByID(ctx context.Context, id string) error {
	m.DeleteCalled = true
	if m.ShouldFail {
		return errors.New("delete failed")
	}
	return nil
}

func (m *MockSubjectRepo) UpdateSubject(ctx context.Context, s db.Subject) error {
	m.UpdateCalled = true
	if m.ShouldFail {
		return errors.New("update failed")
	}
	return nil
}

func (m *MockSubjectRepo) DeleteAllSubjects(ctx context.Context) error {
	m.DeleteAllCalled = true
	if m.ShouldFail {
		return errors.New("delete all failed")
	}
	return nil
}

func (m *MockSubjectRepo) GetSubjectByID(ctx context.Context, id string) (db.Subject, error) {
	m.GetByIDCalled = true
	if m.ShouldFail {
		return db.Subject{}, errors.New("get by id failed")
	}
	return db.Subject{ID: id, Name: "Example", Age: 30}, nil
}


func TestAddSubject_Success(t *testing.T) {
	mockRepo := &MockSubjectRepo{}
	svc := service.NewSubjectService(mockRepo)

	sub := db.Subject{ID: "abc", Name: "Yigit", Age: 22}
	err := svc.AddSubject(context.Background(), sub)

	assert.NoError(t, err)
	assert.True(t, mockRepo.InsertCalled)
}

func TestAddSubject_EmptyID(t *testing.T) {
	mockRepo := &MockSubjectRepo{}
	svc := service.NewSubjectService(mockRepo)

	sub := db.Subject{ID: "", Name: "Yigit", Age: 22}
	err := svc.AddSubject(context.Background(), sub)

	assert.Error(t, err)
	assert.EqualError(t, err, "ID boş olamaz")
	assert.False(t, mockRepo.InsertCalled)
}

func TestGetAllSubjects(t *testing.T) {
	mockRepo := &MockSubjectRepo{}
	svc := service.NewSubjectService(mockRepo)

	subjects, err := svc.GetAllSubjects(context.Background())

	assert.NoError(t, err)
	assert.Len(t, subjects, 1)
	assert.Equal(t, "Test", subjects[0].Name)
	assert.True(t, mockRepo.GetAllCalled)
}

func TestDeleteSubject(t *testing.T) {
	mockRepo := &MockSubjectRepo{}
	svc := service.NewSubjectService(mockRepo)

	err := svc.Delete(context.Background(), "123")

	assert.NoError(t, err)
	assert.True(t, mockRepo.DeleteCalled)
}

func TestUpdateSubject(t *testing.T) {
	mockRepo := &MockSubjectRepo{}
	svc := service.NewSubjectService(mockRepo)

	sub := db.Subject{ID: "abc", Name: "Yeni", Age: 40}
	err := svc.Update(context.Background(), sub)

	assert.NoError(t, err)
	assert.True(t, mockRepo.UpdateCalled)
}

func TestGetSubjectByID(t *testing.T) {
	mockRepo := &MockSubjectRepo{}
	svc := service.NewSubjectService(mockRepo)

	result, err := svc.GetSubject(context.Background(), "xyz")

	assert.NoError(t, err)
	assert.Equal(t, "xyz", result.ID)
	assert.Equal(t, "Example", result.Name)
	assert.True(t, mockRepo.GetByIDCalled)
}

func TestDeleteAllSubjects(t *testing.T) {
	mockRepo := &MockSubjectRepo{}
	svc := service.NewSubjectService(mockRepo)

	err := svc.DeleteAll(context.Background())

	assert.NoError(t, err)
	assert.True(t, mockRepo.DeleteAllCalled)
}
