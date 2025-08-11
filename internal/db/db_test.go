package db_test

import (
	"context"
	"testing"

	"github.com/YigitAtaMacit/StajDeneme/internal/db"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var testRepo db.SubjectRepository

func init() {
	
	if err := db.ConnectDB(); err != nil {
		panic("Veritabanına bağlanılamadı: " + err.Error())
	}

	createAppointments := `
	CREATE TABLE IF NOT EXISTS appointments (
		id TEXT PRIMARY KEY,
		user_id TEXT NOT NULL,
		doctor_name TEXT NOT NULL,
		date TEXT NOT NULL,
		time TEXT NOT NULL,
		description TEXT
	);`
	if _, err := db.DB.Exec(context.Background(), createAppointments); err != nil {
		panic("appointments tablosu oluşturulamadı: " + err.Error())
	}

	testRepo = db.NewSubjectRepo(db.DB)
}

func TestInsertAndGetSubjectByID(t *testing.T) {
	ctx := context.Background()

	sub := db.Subject{
		ID:          uuid.New().String(),
		UserID:      "efe",          
		DoctorName:  "Dr. House",
		Date:        "2025-08-12",
		Time:        "10:30",
		Description: "Kontrol randevusu",
	}

	err := testRepo.InsertSubject(ctx, sub)
	assert.NoError(t, err)

	got, err := testRepo.GetSubjectByID(ctx, sub.ID)
	assert.NoError(t, err)

	assert.Equal(t, sub.UserID, got.UserID)
	assert.Equal(t, sub.DoctorName, got.DoctorName)
	assert.Equal(t, sub.Date, got.Date)
	assert.Equal(t, sub.Time, got.Time)
	assert.Equal(t, sub.Description, got.Description)
}

func TestUpdateSubject(t *testing.T) {
	ctx := context.Background()

	sub := db.Subject{
		ID:          uuid.New().String(),
		UserID:      "efe",
		DoctorName:  "Dr. Strange",
		Date:        "2025-08-15",
		Time:        "14:00",
		Description: "İlk muayene",
	}
	assert.NoError(t, testRepo.InsertSubject(ctx, sub))


	sub.DoctorName  = "Dr. Strange (Güncel)"
	sub.Date        = "2025-08-16"
	sub.Time        = "15:45"
	sub.Description = "Takip muayenesi"

	err := testRepo.UpdateSubject(ctx, sub)
	assert.NoError(t, err)

	updated, err := testRepo.GetSubjectByID(ctx, sub.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Dr. Strange (Güncel)", updated.DoctorName)
	assert.Equal(t, "2025-08-16", updated.Date)
	assert.Equal(t, "15:45", updated.Time)
	assert.Equal(t, "Takip muayenesi", updated.Description)
}

func TestDeleteAllSubjectsAndGetAll(t *testing.T) {
	ctx := context.Background()


	_ = testRepo.DeleteAllSubjects(ctx)

	sub1 := db.Subject{
		ID:          uuid.New().String(),
		UserID:      "efe",
		DoctorName:  "Dr. Who",
		Date:        "2025-08-20",
		Time:        "09:00",
		Description: "Genel kontrol",
	}
	sub2 := db.Subject{
		ID:          uuid.New().String(),
		UserID:      "efe",
		DoctorName:  "Dr. Watson",
		Date:        "2025-08-21",
		Time:        "11:30",
		Description: "Göz muayenesi",
	}

	assert.NoError(t, testRepo.InsertSubject(ctx, sub1))
	assert.NoError(t, testRepo.InsertSubject(ctx, sub2))

	all, err := testRepo.GetAllSubjects(ctx)
	assert.NoError(t, err)
	assert.Len(t, all, 2)

	err = testRepo.DeleteAllSubjects(ctx)
	assert.NoError(t, err)

	all, err = testRepo.GetAllSubjects(ctx)
	assert.NoError(t, err)
	assert.Len(t, all, 0)
}

func TestDeleteSubjectByID(t *testing.T) {
	ctx := context.Background()

	sub := db.Subject{
		ID:          uuid.New().String(),
		UserID:      "efe",
		DoctorName:  "Dr. Brown",
		Date:        "2025-08-22",
		Time:        "13:15",
		Description: "Aşı",
	}
	assert.NoError(t, testRepo.InsertSubject(ctx, sub))

	err := testRepo.DeleteSubjectByID(ctx, sub.ID)
	assert.NoError(t, err)

	_, err = testRepo.GetSubjectByID(ctx, sub.ID)
	assert.Error(t, err) 
}
