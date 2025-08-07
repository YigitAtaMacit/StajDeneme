package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/YigitAtaMacit/StajDeneme/internal/db"
	"github.com/YigitAtaMacit/StajDeneme/internal/service"
	"github.com/YigitAtaMacit/StajDeneme/internal/subject"
	"github.com/YigitAtaMacit/StajDeneme/internal/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var server *httptest.Server

func TestMain(m *testing.M) {
	_ = db.ConnectDB()
	_ = db.CreateUserDB()
	_ = db.CreateDB()

	repo := db.NewSubjectRepo(db.DB)
	svc := service.NewSubjectService(repo)
	handler := subject.NewSubjectHandler(svc)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Post("/register", auth.RegisterHandler)
	router.Post("/login", auth.LoginHandler)

	router.Route("/subjects", func(r chi.Router) {
		r.Use(auth.NewMiddleware)
		r.Post("/", handler.PostSubject)
		r.Get("/", handler.GetSubject)
	})

	server = httptest.NewServer(router)
	defer server.Close()

	code := m.Run()
	db.CloseDB()
	os.Exit(code)
}

func TestEndToEndFlow(t *testing.T) {

	registerBody := `{"username":"ata","password":"123456"}`
	resp, err := http.Post(server.URL+"/register", "application/json", bytes.NewBufferString(registerBody))
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("register failed: %v", err)
	}


	loginBody := `{"username":"ata","password":"123456"}`
	resp, err = http.Post(server.URL+"/login", "application/json", bytes.NewBufferString(loginBody))
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	defer resp.Body.Close()

	var loginResp map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&loginResp)
	token := loginResp["token"]

	subjectBody := map[string]interface{}{
		"id":   "abc123",
		"name": "EndtoendTest",
		"age":  18,
	}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(subjectBody)

	req, _ := http.NewRequest("POST", server.URL+"/subjects/", buf)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("subject post failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected subject post status: %d", resp.StatusCode)
	}


	req2, _ := http.NewRequest("GET", server.URL+"/subjects/", nil)
	req2.Header.Set("Authorization", "Bearer "+token)

	resp, err = http.DefaultClient.Do(req2)
	if err != nil {
		t.Fatalf("subject get failed: %v", err)
	}

	var subjects []map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&subjects)

	if len(subjects) == 0 {
		t.Fatalf("expected at least 1 subject")
	}
}
