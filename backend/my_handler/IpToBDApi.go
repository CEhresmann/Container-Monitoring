package my_handler

import (
	"encoding/json"
	"github.com/CEhresmann/Container-Monitoring/db"
	"go.uber.org/zap"
	"net/http"
)

type IPStatus struct {
	IP       string `json:"ip"`
	PingTime string `json:"ping_time"`
	LastOK   string `json:"last_ok"`
}

func GetIPStatuses(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT ip, ping_time, last_ok FROM ips")
	if err != nil {
		zap.L().Error("Error retrieving data", zap.Error(err))
		http.Error(w, "Error retrieving data", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var results []IPStatus
	for rows.Next() {
		var ip IPStatus
		err := rows.Scan(&ip.IP, &ip.PingTime, &ip.LastOK)
		if err != nil {
			zap.L().Error("Error scanning row", zap.Error(err))
			http.Error(w, "Error processing data", http.StatusInternalServerError)
			return
		}
		results = append(results, ip)
	}

	if err := rows.Err(); err != nil {
		zap.L().Error("Error with row iteration", zap.Error(err))
		http.Error(w, "Error processing data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		zap.L().Error("Error encoding JSON", zap.Error(err))
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func AddIPStatus(w http.ResponseWriter, r *http.Request) {
	var ip IPStatus
	err := json.NewDecoder(r.Body).Decode(&ip)
	if err != nil {
		zap.L().Error("Error decoding JSON", zap.Error(err))
		http.Error(w, "Invalid input data", http.StatusBadRequest)
		return
	}

	_, err = db.DB.Exec("INSERT INTO ips (ip, ping_time, last_ok) VALUES ($1, $2, $3)",
		ip.IP, ip.PingTime, ip.LastOK)
	if err != nil {
		zap.L().Error("Error inserting into DB", zap.Error(err))
		http.Error(w, "Error saving data", http.StatusInternalServerError)
		return
	}

	zap.L().Info("IP status added successfully", zap.String("ip", ip.IP))
	w.WriteHeader(http.StatusCreated)
}
