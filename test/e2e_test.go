package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
    "context"
	"github.com/YigitAtaMacit/StajDeneme/internal/auth"
	"github.com/YigitAtaMacit/StajDeneme/internal/db"
	"github.com/YigitAtaMacit/StajDeneme/internal/service"
	"github.com/YigitAtaMacit/StajDeneme/internal/subject"
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
		r.Delete("/{id}", handler.DeleteSubject)
	})

	server = httptest.NewServer(router)
	defer server.Close()

	code := m.Run()
	db.CloseDB()
	os.Exit(code)
}

func TestEndToEndFlow(t *testing.T) {

	registerBody := `{"username":"test1234","password":"123456"}`
	resp, err := http.Post(server.URL+"/register", "application/json", bytes.NewBufferString(registerBody))
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("register failed: %v", err)
	}
	resp.Body.Close()


	loginBody := `{"username":"test1234","password":"123456"}`
	resp, err = http.Post(server.URL+"/login", "application/json", bytes.NewBufferString(loginBody))
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	var loginResp map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&loginResp)
	resp.Body.Close()
	token := loginResp["token"]

	subjectBody := map[string]interface{}{
		    "id":          "ete123",
			"userid":      "test1234",
            "doctorName":  "Dr. E2E",
            "date":        "2025-08-25",
            "time":        "10:15",
            "description": "E2E randevu",
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
	resp.Body.Close()

	req2, _ := http.NewRequest("GET", server.URL+"/subjects/", nil)
	req2.Header.Set("Authorization", "Bearer "+token)
	resp, err = http.DefaultClient.Do(req2)
	if err != nil {
		t.Fatalf("subject get failed: %v", err)
	}
	var subjects []map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&subjects)
	resp.Body.Close()

	found := false
	for _, s := range subjects {
		if s["id"] == "ete123" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("added subject not found in GET")
	}


	req3, _ := http.NewRequest("DELETE", server.URL+"/subjects/ete123", nil)
	req3.Header.Set("Authorization", "Bearer "+token)
	resp, err = http.DefaultClient.Do(req3)
	if err != nil {
		t.Fatalf("subject delete failed: %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("unexpected delete status: %d", resp.StatusCode)
	}
	_, err = db.DB.Exec(context.Background(), `DELETE FROM users WHERE username = $1`, "test1234")
    if err != nil {
    t.Fatalf("cleanup failed: %v", err)
    }
	resp.Body.Close()
}
