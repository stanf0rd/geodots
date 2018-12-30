package main

import (
	"encoding/json"
	"fmt"
	"geomap/models"
	"log"
	"net/http"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "geomap"
	password = "geomap"
	dbname   = "geomap"
)

func main() {
	dbInfo := fmt.Sprintf(
		"host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	models.InitDB(dbInfo)

	http.HandleFunc("/dots", addDots)
	http.HandleFunc("/dots/all", getAllDots)
	http.HandleFunc("/dots/area", getAreaDots)

	log.Fatal(http.ListenAndServe(":9700", nil))
}

func addDots(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/dots" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method.", 405)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var dots []models.Dot
	err := decoder.Decode(&dots)
	if err != nil {
		panic(err)
	}

	models.AddDots(dots)
}

func getAllDots(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/dots/all" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method.", 405)
		return
	}

	dots, err := models.AllDots()
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(dots)
}

func getAreaDots(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/dots/area" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method.", 405)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var area models.Area
	err := decoder.Decode(&area)
	if err != nil {
		panic(err)
	}

	areaDots := models.AreaDots(area)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(areaDots)
}
