package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)



type SubjectRepository interface {
	InsertSubject(ctx context.Context, subject Subject) error
	GetSubjectByID(ctx context.Context, id string) (Subject, error)
	GetAllSubjects(ctx context.Context) ([]Subject, error)
	UpdateSubject(ctx context.Context, subject Subject) error
	DeleteSubjectByID(ctx context.Context, id string) error
	DeleteAllSubjects(ctx context.Context) error
}

type SubjectRepo struct {
	DB *pgxpool.Pool
}

func NewSubjectRepo(db *pgxpool.Pool) SubjectRepository {
	return &SubjectRepo{DB: db}
}

func (r *SubjectRepo) InsertSubject(ctx context.Context, subject Subject) error {
	_, err := r.DB.Exec(ctx,
		`INSERT INTO appointments (id, user_id, doctor_name, date, time, description)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		subject.ID, subject.UserID, subject.DoctorName, subject.Date, subject.Time, subject.Description)
	return err
}

func (r *SubjectRepo) UpdateSubject(ctx context.Context, subject Subject) error {
	_, err := r.DB.Exec(ctx,
		`UPDATE appointments 
		 SET doctor_name=$1, date=$2, time=$3, description=$4 
		 WHERE id=$5`,
		subject.DoctorName, subject.Date, subject.Time, subject.Description, subject.ID)
	return err
}

func (r *SubjectRepo) GetAllSubjects(ctx context.Context) ([]Subject, error) {
	rows, err := r.DB.Query(ctx,
		`SELECT id, user_id, doctor_name, date, time, description FROM appointments`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subjects []Subject
	for rows.Next() {
		var s Subject
		err := rows.Scan(&s.ID, &s.UserID, &s.DoctorName, &s.Date, &s.Time, &s.Description)
		if err != nil {
			return nil, err
		}
		subjects = append(subjects, s)
	}
	return subjects, nil
}

func (r *SubjectRepo) DeleteSubjectByID(ctx context.Context, id string) error {
	_, err := r.DB.Exec(ctx, `DELETE FROM appointments WHERE id=$1`, id)
	return err
}

func (r *SubjectRepo) DeleteAllSubjects(ctx context.Context) error {
	_, err := r.DB.Exec(ctx, `DELETE FROM appointments`)
	return err
}

func (r *SubjectRepo) GetSubjectByID(ctx context.Context, id string) (Subject, error) {
	var s Subject
	err := r.DB.QueryRow(ctx,
		`SELECT id, user_id, doctor_name, date, time, description FROM appointments WHERE id=$1`,
		id).Scan(&s.ID, &s.UserID, &s.DoctorName, &s.Date, &s.Time, &s.Description)
	return s, err
}
