package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// 🔍 Verifica se o banco existe (evita erro silencioso)
	if _, err := os.Stat("../data/db1"); err != nil {
		log.Fatal("DB not found:", err)
	}

	// 🔌 Conecta ao SQLite
	db, err := sql.Open("sqlite3", "../data/db1")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// ================================
	// 📊 Partidas por dia
	// ================================
	http.HandleFunc("/analytics/matches-per-day", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`
			SELECT substr(created_at, 1, 10) AS date, COUNT(*) 
			FROM matches
			GROUP BY date
			ORDER BY date
		`)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		type Row struct {
			Date  string `json:"date"`
			Total int    `json:"total"`
		}

		result := []Row{}
		for rows.Next() {
			var row Row
			rows.Scan(&row.Date, &row.Total)
			result = append(result, row)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	// ================================
	// 🏆 Vitórias por jogador
	// ================================
	http.HandleFunc("/analytics/wins-per-player", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`
			SELECT winner, COUNT(*) AS total
			FROM matches
			WHERE winner IS NOT NULL
			GROUP BY winner
			ORDER BY total DESC
		`)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		defer rows.Close()

		type Row struct {
			Player string `json:"player"`
			Total  int    `json:"total"`
		}

		result := []Row{}
		for rows.Next() {
			var row Row
			rows.Scan(&row.Player, &row.Total)
			result = append(result, row)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	// ================================
	// ⏱️ Tempo médio de partida
	// ================================
	http.HandleFunc("/analytics/average-duration", func(w http.ResponseWriter, r *http.Request) {
		row := db.QueryRow(`
			SELECT AVG(duration)
			FROM matches
		`)

		var avg sql.NullFloat64
		err := row.Scan(&avg)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		result := map[string]float64{
			"average_duration": 0,
		}

		if avg.Valid {
			result["average_duration"] = avg.Float64
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	// ================================
	// 🚀 Start server
	// ================================
	log.Println("Analytics running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
