package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SubjectRepository interface{
	InsertSubject(ctx context.Context, subject Subject) error
	GetSubjectByID(ctx context.Context,id string) (Subject,error)
	GetAllSubjects(ctx context.Context) ([]Subject, error)
	UpdateSubject(ctx context.Context, subject Subject) error
	DeleteSubjectByID(ctx context.Context, id string) error
	DeleteAllSubjects(ctx context.Context) error
}

type SubjectRepo struct{
	DB *pgxpool.Pool
}

func NewSubjectRepo(db *pgxpool.Pool) SubjectRepository {
	return &SubjectRepo{DB: db}
}

func (r *SubjectRepo) InsertSubject(ctx context.Context, subject Subject) error{
	_,err := r.DB.Exec(context.Background()," INSERT INTO subject (id,name,age) VALUES ($1, $2, $3)",subject.ID,subject.Name,subject.Age)
	return err
}

func (r *SubjectRepo) UpdateSubject(ctx context.Context,subject Subject) error{
	_, err := r.DB.Exec(ctx, `UPDATE subject SET name=$1, age=$2 WHERE id=$3`, subject.Name, subject.Age, subject.ID)
	return err
}

func (r *SubjectRepo) GetAllSubjects(ctx context.Context) ([]Subject, error) {
	rows, err := r.DB.Query(ctx, `SELECT id, name, age FROM subject`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []Subject
	for rows.Next() {
		var s Subject
		if err := rows.Scan(&s.ID, &s.Name, &s.Age); err != nil {
			return nil, err
		}
		subjects = append(subjects, s)
	}
	return subjects, nil
}

func (r *SubjectRepo) DeleteSubjectByID(ctx context.Context, id string) error {
	_, err := r.DB.Exec(context.Background(), "DELETE FROM subject WHERE id=$1", id)
	return err
}

func (r *SubjectRepo) DeleteAllSubjects(ctx context.Context) error {
	_, err := r.DB.Exec(ctx, `DELETE FROM subject`)
	return err
}


func (r *SubjectRepo) GetSubjectByID(ctx context.Context, id string) (Subject, error) {
	var s Subject
	err := r.DB.QueryRow(ctx, `SELECT id, name, age FROM subject WHERE id=$1`, id).Scan(&s.ID, &s.Name, &s.Age)
	return s, err
}