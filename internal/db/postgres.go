package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Subject struct {
	ID   string
	Name string
	Age  int
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
	createTable := `CREATE TABLE IF NOT EXISTS subject (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		age INT NOT NULL
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

