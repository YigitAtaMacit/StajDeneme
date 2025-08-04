package db_test
import (
	"context"
	"testing"

	
	"github.com/YigitAtaMacit/StajDeneme/internal/db"
	
	"github.com/jackc/pgx/v5/pgxpool"
)
var testDB *pgxpool.Pool

func setupDBTest(t *testing.T){
	var err error
    testDB,err =pgxpool.New(context.Background(),"postgres://postgres:postgres@localhost:5433/stajdb_test")
	createTable := `CREATE TABLE IF NOT EXISTS subject (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		age INT NOT NULL
	);`
	db.DB=testDB
    
	_, err = testDB.Exec(context.Background(), createTable)
	if err != nil {
		t.Fatalf("Tablo oluşturulamadı: %v", err)
	}
}

func TestDeleteSubjects(t *testing.T){
	setupDBTest(t)
     _,err:=testDB.Exec(context.Background(),"DELETE FROM subject")
	 if err !=nil {
		t.Fatalf("subject tablosu temizlenemedi: %v", err)
	 }
}

func TestInsertSubject(t *testing.T){
    TestDeleteSubjects(t)

	subject := db.Subject{
		ID:   "1",
		Name: "TEST Ali",
		Age:  22,
	}
    err := db.InsertSubject(testDB, subject)
	if err != nil {
		t.Fatalf("InsertSubject başarısız: %v", err)
	}


}

func TestGetAllSubjects(t *testing.T) {
	TestDeleteSubjects(t)

	subjects := []db.Subject{
		{ID: "id-1", Name: "Ahmet", Age: 20},
		{ID: "id-2", Name: "Zeynep", Age: 25},
	}
	for _, s := range subjects {
		if err := db.InsertSubject(testDB, s); err != nil {
			t.Fatalf("InsertSubject hatası: %v", err)
		}
	}

	result, err := db.GetAllSubjects(testDB)
	if err != nil {
		t.Fatalf("GetAllSubjects başarısız: %v", err)
	}

	if len(result) != len(subjects) {
		t.Errorf("Beklenen %d subject, gelen %d", len(subjects), len(result))
	}
}

func TestUpdateSubject(t *testing.T){
	TestDeleteSubjects(t)

	original := db.Subject{ID: "update1", Name: "Ali", Age: 20}
	if err := db.InsertSubject(testDB, original); err != nil {
		t.Fatalf("InsertSubject başarısız: %v", err)
	}

	updated := db.Subject{ID: "update1", Name: "UpdatedAli", Age: 30}
	if err := db.UpdateSubject(testDB, updated); err != nil {
		t.Fatalf("UpdateSubject başarısız: %v", err)
	}

	sub, err := db.GetSubjectByID(testDB, "update1")
	if err != nil {
		t.Fatalf("Güncellenen subject alınamadı: %v", err)
	}
	if sub.Name != "UpdatedAli" || sub.Age != 30 {
		t.Errorf("Güncelleme başarısız: %+v", sub)
	}

}

func TestDeleteSubjectByID(t *testing.T){
	TestDeleteSubjects(t)

	subject := db.Subject{ID: "delete1", Name: "DeleteAli", Age: 33}
	if err := db.InsertSubject(testDB, subject); err != nil {
		t.Fatalf("InsertSubject başarısız: %v", err)
	}

	if err := db.DeleteSubjectByID(testDB, subject.ID); err != nil {
		t.Fatalf("DeleteSubjectByID başarısız: %v", err)
	}

	_, err := db.GetSubjectByID(testDB, subject.ID)
	if err == nil {
		t.Errorf("Subject silinmemiş görünüyor")
	}

}

func TestGetSubjectByID(t *testing.T){
	TestDeleteSubjects(t)

	subject := db.Subject{ID: "get1", Name: "Ali", Age: 55}
	if err := db.InsertSubject(testDB, subject); err != nil {
		t.Fatalf("InsertSubject başarısız: %v", err)
	}

	result, err := db.GetSubjectByID(testDB, subject.ID)
	if err != nil {
		t.Fatalf("GetSubjectByID başarısız: %v", err)
	}

	if result.ID != subject.ID || result.Name != subject.Name || result.Age != subject.Age {
		t.Errorf("Beklenen %+v, gelen %+v", subject, result)
	}


}
/* package main

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func setupDBTest(t *testing.T) *pgxpool.Pool {
	conn, err := pgxpool.New(context.Background(), "postgres://postgres:postgres@localhost:5433/stajdb")
	if err != nil {
		t.Fatalf("Veritabanı bağlantısı kurulamadı: %v", err)
	}
	NewTable := `CREATE TABLE IF NOT EXISTS subject(
	    id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		age INT NOT NULL
	
	)`
	_, err = conn.Exec(context.Background(), NewTable)
	if err != nil {
		t.Fatalf("Tablo oluşturulamadı: %v", err)
	}
	return conn
}



func TestInsertSubject(t *testing.T) {
	
	conn := setupDBTest(t)
	defer conn.Close()
    
	testSubject := Subject{
		ID:   "test123",
		Name: "TestUser",
		Age:  123,
	}
    conn.Exec(context.Background(), "DELETE FROM subject WHERE id=$1", testSubject.ID)
	err := InsertSubject(conn, testSubject)
	if err != nil {
		t.Fatalf("Subject eklenemedi: %v", err)
	}

	var name string
	err = conn.QueryRow(context.Background(), "SELECT name FROM subject WHERE id=$1", testSubject.ID).Scan(&name)
	if err != nil {
		t.Fatalf("Subject veritabanında bulunamadı: %v", err)
	}
	if name != testSubject.Name {
		t.Errorf("Beklenen isim %s, gelen %s", testSubject.Name, name)
	}
	
}

func TestUpdateSubject(t *testing.T) {
	conn := setupDBTest(t)
	defer conn.Close()

	org := Subject{ID: "update1", Name: "Initial", Age: 20}
	_ = InsertSubject(conn, org)

	updated := Subject{ID: "update1", Name: "Updated", Age: 30}
	err := UpdateSubject(conn, updated)
	if err != nil {
		t.Fatalf("Update işlemi başarısız: %v", err)
	}

	var name string
	var age int
	err = conn.QueryRow(context.Background(), "SELECT name, age FROM subject WHERE id=$1", updated.ID).Scan(&name, &age)
	if err != nil {
		t.Fatalf("Kayıt bulunamadı: %v", err)
	}
	if name != "Updated" || age != 30 {
		t.Errorf("Güncelleme hatalı: name=%s, age=%d", name, age)
	}
}

func TestDeleteSubjectByID(t *testing.T) {
	conn := setupDBTest(t)
	defer conn.Close()

	org := Subject{ID: "delete1", Name: "ToDelete", Age: 99}
	_ = InsertSubject(conn, org)

	err := DeleteSubjectByID(conn, org.ID)
	if err != nil {
		t.Fatalf("Silinemedi: %v", err)
	}

	var exists bool
	err = conn.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM subject WHERE id=$1)", org.ID).Scan(&exists)
	if err != nil {
		t.Fatalf("Exist kontrol hatası: %v", err)
	}
	if exists {
		t.Errorf("Kayıt hâlâ mevcut, silinememiş.")
	}
}

func TestGetSubjectByID(t *testing.T) {
	conn := setupDBTest(t)
	defer conn.Close()

	sub := Subject{ID: "get1", Name: "Fetch", Age: 55}
	_ = InsertSubject(conn, sub)

	got, err := GetSubjectByID(conn, sub.ID)
	if err != nil {
		t.Fatalf("GetSubjectByID başarısız: %v", err)
	}

	if got.ID != sub.ID || got.Name != sub.Name || got.Age != sub.Age {
		t.Errorf("Beklenmeyen değerler: %+v", got)
	}
} */