package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDB() error{
	conn ,err :=pgxpool.New(context.Background(), "postgres://postgres:postgres@localhost:5433/stajdb")
	if err != nil {
		return fmt.Errorf("DB bağlantı hatası: %w", err)
	}
	DB = conn
	return nil
} 

func CreateUserDB() error{

    createUserTable:= `CREATE TABLE IF NOT EXISTS users(
	     id TEXT PRIMARY KEY,
         username TEXT UNIQUE NOT NULL,
         password TEXT NOT NULL
	);`
    
	_, err:= DB.Exec(context.Background(),createUserTable)
	if err!=nil{
		return fmt.Errorf("Tablo oluşturulamadı: %w", err)
	}
	fmt.Println("kullanıcı tablosu uygun")
	return nil
}
func CreateDB() error{
	createTable:=` CREATE TABLE IF NOT EXISTS subject(
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		age INT NOT NULL
	);`
	_, err:= DB.Exec(context.Background(),createTable)
	if err!=nil{
		return fmt.Errorf("Tablo oluşturulamadı: %w", err)
	}
	fmt.Println("subject tablosu uygun")
	return nil

}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

func InsertSubject(conn *pgxpool.Pool, subject Subject) error {
	_, err := conn.Exec(context.Background(),
		"INSERT INTO subject(id,name, age) VALUES($1, $2,$3)",
		subject.ID,subject.Name, subject.Age)
		if err != nil {
		return fmt.Errorf("Subject eklenemedi: %w", err)
	}
	return nil
}

func UpdateSubject(conn *pgxpool.Pool, subject Subject) error {
	x, err := conn.Exec(
		context.Background(),
		"UPDATE subject SET name=$1, age=$2 WHERE id=$3",
		subject.Name,
		subject.Age,
		subject.ID,
	)
	if err != nil {
		return fmt.Errorf("Veritabanı güncelleme hatası: %w", err)
	}
	if x.RowsAffected() == 0 {
		return fmt.Errorf("ID bulunamadı: %s", subject.ID)
	}
	return nil
}
func GetAllSubjects(conn *pgxpool.Pool)([]Subject,error){

	rows,err:=conn.Query(context.Background(),"SELECT id,name,age FROM subject")
	if err!=nil{
		return nil,fmt.Errorf("Subject verileri alınamadı: %w", err)
	}
	defer rows.Close()
	
	var subjects []Subject

	for rows.Next(){
		var s Subject
		err:=rows.Scan(&s.ID,&s.Name,&s.Age)
		if err!=nil{
			return nil,fmt.Errorf("Veri okunamadı: %w", err)
		}
		subjects = append(subjects, s)
	}
	return subjects, nil
}

func DeleteSubjectByID(conn *pgxpool.Pool, id string) error {
	x, err := conn.Exec(
		context.Background(),
		"DELETE FROM subject WHERE id = $1",
		id,
	)
	if err != nil {
		return fmt.Errorf("subject silinemedi: %w", err)
	}
	if x.RowsAffected() == 0 {
		return fmt.Errorf("id %s bulunamadı", id)
	}
	return nil
}

func GetSubjectByID(conn *pgxpool.Pool, id string) (Subject, error) {
	var s Subject
	err := conn.QueryRow(
		context.Background(),
		"SELECT id, name, age FROM subject WHERE id = $1",
		id,
	).Scan(&s.ID, &s.Name, &s.Age)

	if err != nil {
		return Subject{}, fmt.Errorf("subject bulunamadı: %w", err)
	}
	return s, nil
}


