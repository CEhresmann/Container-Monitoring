package handlers

import (
	"encoding/json"
	"github.com/CEhresmann/Container-Monitoring/db"
	"log"
	"net/http"
)

type IPStatus struct {
	IP       string `sql:"ip"`
	PingTime string `sql:"ping_time"`
	LastOK   string `sql:"last_ok"`
}

func GetIPStatuses(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT ip, ping_time, last_ok FROM ips")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	var results []IPStatus
	for rows.Next() {
		var ip IPStatus
		err := rows.Scan(&ip.IP, &ip.PingTime, &ip.LastOK)
		if err != nil {
			log.Println(err)
		}
		results = append(results, ip)
	}
	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		log.Println(err)
	}
}

func AddIPStatus(w http.ResponseWriter, r *http.Request) {
	var ip IPStatus
	err := json.NewDecoder(r.Body).Decode(&ip)
	if err != nil {
		log.Println(err)
	}
	_, err = db.DB.Exec("INSERT INTO ips (ip, ping_time, last_ok) VALUES ($1, $2, $3)",
		ip.IP, ip.PingTime, ip.LastOK)
	if err != nil {
		log.Println(err)
	}
	w.WriteHeader(http.StatusCreated)
}
