package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Subject struct {
	ID          string `json:"id"`
	UserID      string `json:"userId"`
	DoctorName  string `json:"doctorName"`
	Date        string `json:"date"`
	Time        string `json:"time"`
	Description string `json:"description"`
}

var DB *pgxpool.Pool

func ConnectDB() error {
	conn, err := pgxpool.New(context.Background(), "postgres://postgres:postgres@localhost:5433/stajdb")
	if err != nil {
		return fmt.Errorf("DB bağlantı hatası: %w", err)
	}
	DB = conn
	return nil
}

func CreateUserDB() error {
	createTable := `CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	)`
	_, err := DB.Exec(context.Background(), createTable)
	if err != nil {
		return fmt.Errorf("kullanıcı tablosu oluşturulamadı: %w", err)
	}
	return nil
}

func CreateDB() error {
	createTable := `CREATE TABLE IF NOT EXISTS appointments (
		id TEXT PRIMARY KEY,
        user_id TEXT NOT NULL,
        doctor_name TEXT NOT NULL,
        date TEXT NOT NULL,
        time TEXT NOT NULL,
        description TEXT
	)`
	_, err := DB.Exec(context.Background(), createTable)
	if err != nil {
		return fmt.Errorf("Subject tablosu oluşturulamadı: %w", err)
	}
	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}
