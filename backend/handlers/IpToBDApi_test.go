package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/CEhresmann/Container-Monitoring/db"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/lib/pq"
)

func setupTestDB() {
	var err error
	db.DB, err = sql.Open("postgres", "host=localhost user=vk password=12345 dbname=monitoring sslmode=disable")
	if err != nil {
		panic(err)
	}

	_, err = db.DB.Exec(`CREATE TABLE IF NOT EXISTS ips (
		id SERIAL PRIMARY KEY,
		ip VARCHAR(100),
		ping_time INT,
		last_ok DATE
	);`)
	if err != nil {
		panic(err)
	}
}

func teardownTestDB() {
	db.DB.Exec("DROP TABLE IF EXISTS ips;")
	db.DB.Close()
}

func TestGetIPStatuses(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	_, err := db.DB.Exec("INSERT INTO ips (ip, ping_time, last_ok) VALUES ('192.168.1.1', 15, '2024-02-01')")
	if err != nil {
		t.Fatalf("Ошибка при подготовке теста: %v", err)
	}

	req, err := http.NewRequest("GET", "/api/ip", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetIPStatuses)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Ожидался статус-код %v, получен %v", http.StatusOK, status)
	}

	var response []IPStatus
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Ошибка десериализации ответа: %v", err)
	}

	if len(response) == 0 || response[0].IP != "192.168.1.1" {
		t.Errorf("Некорректный ответ API: %v", response)
	}
}

func TestAddIPStatus(t *testing.T) {
	setupTestDB()
	defer teardownTestDB()

	ipData := IPStatus{
		IP:       "192.168.1.100",
		PingTime: "10",
		LastOK:   "2024-02-02",
	}
	jsonData, _ := json.Marshal(ipData)

	req, err := http.NewRequest("POST", "/api/ip", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddIPStatus)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Ожидался статус-код %v, получен %v", http.StatusCreated, status)
	}

	var count int
	err = db.DB.QueryRow("SELECT COUNT(*) FROM ips WHERE ip = $1", ipData.IP).Scan(&count)
	if err != nil {
		t.Errorf("Ошибка при проверке вставки в БД: %v", err)
	}

	if count == 0 {
		t.Errorf("Данные не были добавлены в БД")
	}
}
