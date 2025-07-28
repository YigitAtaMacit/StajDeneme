package main

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"
)

func TestMain(m *testing.M) {

	err := ConnectDB()
	if err != nil {
		panic(fmt.Sprintf("Veritabanına bağlanılamadı: %v", err))
	}
	defer CloseDB()

	if err := CreateDB(); err != nil {
		panic(fmt.Sprintf("Tablo oluşturulamadı: %v", err))
	}


	subjects = make(map[string]Subject)

	if err := LoadSubjectsFromDB(); err != nil {
		fmt.Println("Uyarı: In-memory map'e veriler yüklenemedi:", err)
	}


	m.Run()
}

func TestGet(t *testing.T) {

	rout := httptest.NewRequest("GET", "/subjects", nil)
	rr := httptest.NewRecorder()
	GetSubject(rr, rout)
	if rr.Code != http.StatusOK {
		t.Errorf("Beklenen 200 , gelen %d", rr.Code)
	}
	var got []Subject
	err := json.Unmarshal(rr.Body.Bytes(), &got)
	if err != nil {
		t.Errorf("JSON çevirilemedi %v", err)
	}

	if len(got) == 0 {
		t.Errorf("Subject verisi gelmedi, beklenenden az kayıt var.")
	}
}
func TestPut(t *testing.T) {
	InsertSubject(DB, Subject{ID: "1", Name: "Ali", Age: 20})
	subjects = map[string]Subject{
		"1": {ID: "1", Name: "Ali", Age: 20},
	}

	UpdatedSubject := Subject{
		ID:   "1",
		Name: "Veli",
		Age:  25,
	}

	jsonSubject, err := json.Marshal(UpdatedSubject)
	if err != nil {
		t.Fatalf("JSON'a çevrilemedi")
	}

	rout := httptest.NewRequest("PUT", "/subjects/1", bytes.NewReader(jsonSubject))

	rout.Header.Set("Content-Type", "application/json")

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "1")
	rout = rout.WithContext(context.WithValue(rout.Context(), chi.RouteCtxKey, ctx))

	rr := httptest.NewRecorder()

	PutSubject(rr, rout)

	if rr.Code != http.StatusOK {
		t.Errorf("Beklenen 200, gelen %d", rr.Code)

	}

	result := subjects["1"]
	if result.Name != "Veli" || result.Age != 25 {
		t.Errorf("Subject güncellenmedi. Gelen: %+v", result)
	}

}
func TestPost(t *testing.T) {

	testSubject := Subject{
		ID:   "3",
		Name: "Mahmut",
		Age:  32,
	}

	jsonSubject, err := json.Marshal(testSubject)
	if err != nil {
		t.Fatalf("JSON'a çevrilemedi: %v", err)
	}

	rout := httptest.NewRequest("POST", "/subjects", bytes.NewReader(jsonSubject))
	rout.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	PostSubject(rr, rout)

	if rr.Code != http.StatusCreated {
		t.Errorf("Beklenen 201 (Created), gelen %v", rr.Code)
	}

	added, exists := subjects["3"]
	if !exists {
		t.Fatalf("Subject eklenemedi, map içinde yok.")
	}

	if added.Name != testSubject.Name || added.Age != testSubject.Age {
		t.Fatalf("Yanlış değerler: %+v", added)
	}
}

func TestDelete(t *testing.T){
	InsertSubject(DB, Subject{ID: "5", Name: "Test", Age: 50})
    subjects = map[string]Subject{
		"5": {ID: "5", Name: "Test", Age: 50},
	}
	rout :=httptest.NewRequest("DELETE","/subjects/5",nil)

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "5")
	rout = rout.WithContext(context.WithValue(rout.Context(), chi.RouteCtxKey, ctx))

	rr := httptest.NewRecorder()

	DeleteSubject(rr,rout)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Beklenen 204, gelen %d", rr.Code)
	}
    
		if _, exists := subjects["5"]; exists {
		t.Errorf("Subject silinemedi. Hala var.")
	}


}

func TestGetById(t *testing.T){
	sub := Subject{ID: "10", Name: "Ayşe", Age: 30}


	_, _ = DB.Exec(context.Background(), "DELETE FROM subject WHERE id=$1", sub.ID)

	err := InsertSubject(DB, sub)
	if err != nil {
		t.Fatalf("Veritabanına subject eklenemedi: %v", err)
	}


	req := httptest.NewRequest("GET", "/subjects/10", nil)

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "10")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))


	rr := httptest.NewRecorder()
	GetbyID(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Beklenen 200, gelen %d", rr.Code)
	}


	var result Subject
	err = json.NewDecoder(rr.Body).Decode(&result)
	if err != nil {
		t.Errorf("JSON çözülemedi: %v", err)
	}

	
	if result.ID != sub.ID || result.Name != sub.Name || result.Age != sub.Age {
		t.Errorf("Subject yanlış: %+v", result)
	}
}

