package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// helper to setup in-memory DB and router
func setupTestRouter(t *testing.T) *gin.Engine {
	// initialize in-memory sqlite
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	// run migrations
	if err := db.AutoMigrate(&Product{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}

	// assign to package-level database var used by DAL
	database = db

	// reduce test noise
	gin.SetMode(gin.TestMode)

	return newRouter()
}

// decode envelope into map for assertions
func decodeEnvelope(t *testing.T, resp *httptest.ResponseRecorder) map[string]interface{} {
	var env map[string]interface{}
	if err := json.Unmarshal(resp.Body.Bytes(), &env); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}
	return env
}

func TestGETProductsPagination(t *testing.T) {
	r := setupTestRouter(t)

	// seed 25 products
	for i := 1; i <= 25; i++ {
		p := Product{Code: "code" + strconv.Itoa(i), Price: uint(i)}
		_ = database.Create(&p)
	}

	req := httptest.NewRequest(http.MethodGet, "/products?page=2&per_page=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	env := decodeEnvelope(t, w)
	if success, ok := env["success"].(bool); !ok || !success {
		t.Fatalf("expected success=true envelope, got %v", env["success"])
	}

	// data should be array length 10
	data, ok := env["data"].([]interface{})
	if !ok {
		// Try decoding 'data' that may be encoded as []map[string]interface{}
		var raw struct {
			Data []map[string]interface{} `json:"data"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &raw); err != nil {
			t.Fatalf("unexpected data shape: %v", err)
		}
		if len(raw.Data) != 10 {
			t.Fatalf("expected 10 products, got %d", len(raw.Data))
		}
	} else {
		if len(data) != 10 {
			t.Fatalf("expected 10 products, got %d", len(data))
		}
	}

	meta, ok := env["meta"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected meta object in response")
	}
	if int(meta["page"].(float64)) != 2 {
		t.Fatalf("expected page=2 in meta, got %v", meta["page"])
	}
	if int(meta["per_page"].(float64)) != 10 {
		t.Fatalf("expected per_page=10 in meta, got %v", meta["per_page"])
	}
	if int(meta["total"].(float64)) != 25 {
		t.Fatalf("expected total=25 in meta, got %v", meta["total"])
	}
	if int(meta["total_pages"].(float64)) != 3 {
		t.Fatalf("expected total_pages=3 in meta, got %v", meta["total_pages"])
	}
}

func TestPOSTCreatesProduct(t *testing.T) {
	r := setupTestRouter(t)

	payload := map[string]interface{}{"code": "new-code", "price": 42}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/product", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	env := decodeEnvelope(t, w)
	if success, ok := env["success"].(bool); !ok || !success {
		t.Fatalf("expected success=true envelope, got %v", env["success"])
	}

	// Ensure returned data has Code and Price
	var raw struct {
		Data Product `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &raw); err != nil {
		t.Fatalf("failed to decode product data: %v", err)
	}
	if raw.Data.Code != "new-code" || raw.Data.Price != 42 {
		t.Fatalf("created product mismatch: %+v", raw.Data)
	}
}

func TestPUTUpdatesProduct(t *testing.T) {
	r := setupTestRouter(t)

	p := Product{Code: "orig", Price: 5}
	_ = database.Create(&p)

	payload := map[string]interface{}{"code": "updated", "price": 99}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPut, "/product/"+strconv.FormatUint(uint64(p.ID), 10), bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var raw struct {
		Data Product `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &raw); err != nil {
		t.Fatalf("failed to decode updated product: %v", err)
	}
	if raw.Data.Code != "updated" || raw.Data.Price != 99 {
		t.Fatalf("updated product mismatch: %+v", raw.Data)
	}
}

func TestGETProductByID(t *testing.T) {
	r := setupTestRouter(t)

	p := Product{Code: "byid", Price: 7}
	_ = database.Create(&p)

	req := httptest.NewRequest(http.MethodGet, "/product/"+strconv.FormatUint(uint64(p.ID), 10), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var raw struct {
		Data Product `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &raw); err != nil {
		t.Fatalf("failed to decode product by id: %v", err)
	}
	if raw.Data.ID != p.ID {
		t.Fatalf("expected id %d, got %d", p.ID, raw.Data.ID)
	}
}

func TestGETLatestProduct(t *testing.T) {
	r := setupTestRouter(t)

	// create two products; getLatestProduct currently returns first by primary key
	p1 := Product{Code: "first", Price: 1}
	p2 := Product{Code: "second", Price: 2}
	_ = database.Create(&p1)
	_ = database.Create(&p2)

	req := httptest.NewRequest(http.MethodGet, "/product/latest", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var raw struct {
		Data Product `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &raw); err != nil {
		t.Fatalf("failed to decode latest product: %v", err)
	}
	if raw.Data.ID != p1.ID {
		t.Fatalf("expected latest product id %d, got %d", p1.ID, raw.Data.ID)
	}
}

func TestPOSTLocationHeaderAndBadRequest(t *testing.T) {
	r := setupTestRouter(t)

	// Bad request: missing price
	payload := map[string]interface{}{"code": "x"}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/product", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for bad request, got %d", w.Code)
	}

	// Good request: check Location header
	payload2 := map[string]interface{}{"code": "loc", "price": 10}
	b2, _ := json.Marshal(payload2)
	req2 := httptest.NewRequest(http.MethodPost, "/product", bytes.NewReader(b2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusCreated {
		t.Fatalf("expected 201 for create, got %d", w2.Code)
	}
	loc := w2.Header().Get("Location")
	if loc == "" {
		t.Fatalf("expected Location header to be set")
	}
}

func TestPerPageLimit(t *testing.T) {
	r := setupTestRouter(t)

	// Request with excessive per_page
	req := httptest.NewRequest(http.MethodGet, "/products?per_page=1000", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 when per_page too large, got %d", w.Code)
	}
}

func TestPerPageEnvOverride(t *testing.T) {
	// set env override
	_ = os.Setenv("MAX_PER_PAGE", "5")
	defer os.Unsetenv("MAX_PER_PAGE")

	r := setupTestRouter(t)

	// Request with per_page greater than env max
	req := httptest.NewRequest(http.MethodGet, "/products?per_page=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 when per_page too large via env override, got %d", w.Code)
	}

	env := decodeEnvelope(t, w)
	errObj, ok := env["error"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected error object in envelope")
	}
	details, ok := errObj["details"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected details map in error")
	}
	if int(details["max_per_page"].(float64)) != 5 {
		t.Fatalf("expected max_per_page=5 in details, got %v", details["max_per_page"])
	}
}

func TestPUTNotFoundAndDELETENotFound(t *testing.T) {
	r := setupTestRouter(t)

	// Attempt to update non-existent product id 999
	payload := map[string]interface{}{"code": "x", "price": 1}
	b, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPut, "/product/999", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for updating non-existent product, got %d", w.Code)
	}
	env := decodeEnvelope(t, w)
	errObj, ok := env["error"].(map[string]interface{})
	if !ok || errObj["code"].(string) != CodeProductNotFound {
		t.Fatalf("expected %s error code, got %v", CodeProductNotFound, env["error"])
	}

	// Attempt to delete non-existent product id 999
	req2 := httptest.NewRequest(http.MethodDelete, "/product/999", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for deleting non-existent product, got %d", w2.Code)
	}
	env2 := decodeEnvelope(t, w2)
	errObj2, ok := env2["error"].(map[string]interface{})
	if !ok || errObj2["code"].(string) != CodeProductNotFound {
		t.Fatalf("expected %s error code for delete, got %v", CodeProductNotFound, env2["error"])
	}
}
