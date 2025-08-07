package service

import (
	"context"
	"errors"

	"github.com/YigitAtaMacit/StajDeneme/internal/db"
)

type SubjectService interface {
	AddSubject(ctx context.Context, subject db.Subject) error
	Update(ctx context.Context, subject db.Subject) error
	GetAllSubjects(ctx context.Context) ([]db.Subject, error)
	Delete(ctx context.Context, id string) error
	DeleteAll(ctx context.Context) error
	GetSubject(ctx context.Context, id string) (db.Subject, error)
}

type subjectServiceImpl struct {
	Repo db.SubjectRepository
}

func NewSubjectService(repo db.SubjectRepository) SubjectService {
	return &subjectServiceImpl{Repo: repo}
}

func (s *subjectServiceImpl) AddSubject(ctx context.Context, subject db.Subject) error {
	if subject.ID == "" {
		return errors.New("ID boş olamaz")
	}
	return s.Repo.InsertSubject(ctx, subject)
}

func (s *subjectServiceImpl) Update(ctx context.Context, subject db.Subject) error {
	return s.Repo.UpdateSubject(ctx, subject)
}

func (s *subjectServiceImpl) GetAllSubjects(ctx context.Context) ([]db.Subject, error) {
	return s.Repo.GetAllSubjects(ctx)
}

func (s *subjectServiceImpl) Delete(ctx context.Context, id string) error {
	return s.Repo.DeleteSubjectByID(ctx, id)
}

func (s *subjectServiceImpl) DeleteAll(ctx context.Context) error {
	return s.Repo.DeleteAllSubjects(ctx)
}

func (s *subjectServiceImpl) GetSubject(ctx context.Context, id string) (db.Subject, error) {
	return s.Repo.GetSubjectByID(ctx, id)
}
