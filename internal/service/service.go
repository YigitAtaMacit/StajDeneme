package service

import (
	"context"
	"errors"

	"github.com/YigitAtaMacit/StajDeneme/internal/db"
)

type SubjectService struct {
	Repo db.SubjectRepository
}

func NewSubjectService(repo db.SubjectRepository) *SubjectService {
	return &SubjectService{Repo: repo}
}

func (s *SubjectService) AddSubject(ctx context.Context, subject db.Subject) error {
	if subject.ID == "" {
		return errors.New("ID boş olamaz")
	}
	return s.Repo.InsertSubject(ctx, subject)
}

func (s *SubjectService) Update(ctx context.Context, subject db.Subject) error {
	return s.Repo.UpdateSubject(ctx, subject)
}

func (s *SubjectService) GetAll(ctx context.Context) ([]db.Subject, error) {
	return s.Repo.GetAllSubjects(ctx)
}

func (s *SubjectService) Delete(ctx context.Context, id string) error {
	return s.Repo.DeleteSubjectByID(ctx, id)
}

func (s *SubjectService) DeleteAll(ctx context.Context) error {
	return s.Repo.DeleteAllSubjects(ctx)
}

func (s *SubjectService) GetSubject(ctx context.Context, id string) (db.Subject, error) {
	return s.Repo.GetSubjectByID(ctx, id)
}



