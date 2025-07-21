package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"github.com/CicadaHymn/guitar-shop-api/internal/db"
)

type Guitar struct {
	ID       int     `json:"id"`
	Model    string  `json:"model"`
	Price    float64 `json:"price"`
	IsCustom bool    `json:"is_custom"`
}

// GET запрос
func GetGuitars(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Pool.Query(context.Background(), "SELECT id, model, price, is_custom FROM guitars")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var guitars []Guitar		
	for rows.Next() {
		var g Guitar
		if err := rows.Scan(&g.ID, &g.Model, &g.Price, &g.IsCustom); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		guitars = append(guitars, g)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(guitars)
}

// Post запрос на апи
func AddGuitar(w http.ResponseWriter, r *http.Request) {
	var g Guitar
	if err := json.NewDecoder(r.Body).Decode(&g); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	var id int
	err := db.Pool.QueryRow(
		context.Background(),
		"INSERT INTO guitars (model, price, is_custom) VALUES ($1, $2, $3) RETURNING id",
		g.Model, g.Price, g.IsCustom,
	).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	g.ID = id
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(g)
}