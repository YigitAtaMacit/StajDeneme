package subject_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
    "context"
	"fmt"
	"github.com/YigitAtaMacit/StajDeneme/internal/auth"
	"github.com/YigitAtaMacit/StajDeneme/internal/db"
	"github.com/YigitAtaMacit/StajDeneme/internal/subject"
	"github.com/go-chi/chi/v5"
)

func setupRouter() http.Handler  {

	_ = db.ConnectDB()
	_ = db.CreateDB()

	r := chi.NewRouter()

	r.Post("/register", auth.RegisterHandler)
	r.Post("/login", auth.LoginHandler)

	r.Group(func(r chi.Router) {
		r.Use(auth.NewMiddleware)
		r.Get("/subjects", subject.GetSubject)
		r.Post("/subjects", subject.PostSubject)
		r.Get("/subjects/{id}", subject.GetByID)
		r.Put("/subjects/{id}", subject.PutSubject)
		r.Delete("/subjects/{id}", subject.DeleteSubject)
		r.Delete("/subjects", subject.DeleteAllSubjects)
	})
	return r
}

func TestUnauthorizedAccess(t *testing.T) {
	router := setupRouter()

	req := httptest.NewRequest("GET", "/subjects", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("JWT olmadan 401 bekleniyordu, gelen: %d", w.Code)
	}
}

func TestPostSubjectAPI(t *testing.T) {
	router := setupRouter()

	_, _ = db.DB.Exec(context.Background(), "DELETE FROM users")
	_, _ = db.DB.Exec(context.Background(), "DELETE FROM subject")


	token := registerAndLogin(t, router, "apiuser", "apipass")


	subject := db.Subject{
    ID:   "api-1",
    Name: "Test",
    Age:  33,
    }
	payload, _ := json.Marshal(subject)

	req := httptest.NewRequest("POST", "/subjects", bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Beklenen 201 Created, gelen: %d", w.Code)
	}

	var created db.Subject
	if err := json.NewDecoder(w.Body).Decode(&created); err != nil {
		t.Fatalf("Yanıt çözülemedi: %v", err)
	}
	if created.ID != "api-1" || created.Name != "Test" || created.Age != 33 {
		t.Errorf("Yanıt verisi beklenenden farklı: %+v", created)
	}
}

/* func TestGetSubjectByIDAPI(t *testing.T) {
	router := setupRouter()
	

	token := registerAndLogin(t, router, "getiduser", "getidpass")

	_ = db.InsertSubject(db.DB, db.Subject{ID: "id-20", Name: "Zeynep", Age: 25})

	
	req := httptest.NewRequest("GET", "/subjects/id-20", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Beklenen 200 OK, gelen: %d", w.Code)
	}
}
func TestPutSubjectAPI(t *testing.T) {
	router := setupRouter()

	token := registerAndLogin(t, router, "putuser", "putpass")

	_ = db.InsertSubject(db.DB, db.Subject{ID: "put-id", Name: "Ayşe", Age: 20})

	updated := map[string]interface{}{
		"name": "Yeni Ayşe ",
		"age":  28,
	}
	payload, _ := json.Marshal(updated)

	req := httptest.NewRequest("PUT", "/subjects/put-id", bytes.NewBuffer(payload))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Güncelleme başarısız. Kod: %d", w.Code)
	}
}
func TestDeleteSubjectAPI(t *testing.T) {
	router := setupRouter()

	token := registerAndLogin(t, router, "deleteuser", "deletepass")

	_ = db.InsertSubject(db.DB, db.Subject{ID: "del-id", Name: "Silinecek", Age: 99})

	req := httptest.NewRequest("DELETE", "/subjects/del-id", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("Silme başarısız. Kod: %d", w.Code)
	}
}


func TestGetSubjectsAPI(t *testing.T) {
	router := setupRouter()
	token := registerAndLogin(t, router, "getuser", "getpass")
    

    sub := db.Subject{ID: "get-id", Name: "Ali", Age: 30}
	_ = db.InsertSubject(db.DB, sub)

	req := httptest.NewRequest("GET", "/subjects", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)


	if w.Code != http.StatusOK {
		t.Errorf("GET /subjects başarısız. Beklenen 200, gelen: %d", w.Code)
	}
}

func TestDeleteAllSubjectsAPI(t *testing.T) {
	router := setupRouter()

	token := registerAndLogin(t, router, "deletetestuser", "deletetestpass")

	_ = db.InsertSubject(db.DB, db.Subject{ID: "s1", Name: "A", Age: 20})
	_ = db.InsertSubject(db.DB, db.Subject{ID: "s2", Name: "B", Age: 30})

	req := httptest.NewRequest("DELETE", "/subjects", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("Tüm kayıtları silme başarısız. Kod: %d", w.Code)
	}

	
	remaining, err := db.GetAllSubjects(db.DB)
	if err != nil {
		t.Fatalf("Kontrol sırasında hata: %v", err)
	}
	if len(remaining) != 0 {
		t.Errorf("Tüm kayıtlar silinmedi, kalan: %d", len(remaining))
	}
} */


func registerAndLogin(t *testing.T, router http.Handler, username, password string) string {
	_, _ = db.DB.Exec(context.Background(), "DELETE FROM users WHERE username=$1", username)

	registerBody := fmt.Sprintf(`{"username":"%s", "password":"%s"}`, username, password)
	req1 := httptest.NewRequest("POST", "/register", bytes.NewBuffer([]byte(registerBody)))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	loginBody := fmt.Sprintf(`{"username":"%s", "password":"%s"}`, username, password)
	req2 := httptest.NewRequest("POST", "/login", bytes.NewBuffer([]byte(loginBody)))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusOK {
		t.Fatalf("Login başarısız. Kod: %d", w2.Code)
	}
	var response struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(w2.Body).Decode(&response); err != nil {
		t.Fatalf("Token parse hatası: %v", err)
	}
	return response.Token
}


