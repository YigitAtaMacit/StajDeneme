package main
/* package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

type TestSubject struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestEndToEndFlow(t *testing.T) {
	baseURL := "http://localhost:3000/subjects"
	sub := Subject{ID: "e2e123", Name: "E2ETest", Age: 30}

	// POST
	body, _ := json.Marshal(sub)
	resp, err := http.Post(baseURL, "application/json", bytes.NewBuffer(body))
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("POST başarısız: %v, status: %d", err, resp.StatusCode)
	}

	// GET
	getResp, err := http.Get(baseURL + "/e2e123")
	if err != nil || getResp.StatusCode != http.StatusOK {
		t.Fatalf("GET başarısız: %v", err)
	}

	var got Subject
	json.NewDecoder(getResp.Body).Decode(&got)
	if got.Name != sub.Name {
		t.Errorf("Beklenen isim: %s, gelen: %s", sub.Name, got.Name)
	}
}
 */