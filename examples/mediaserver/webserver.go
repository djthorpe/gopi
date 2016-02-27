package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

func indexHandler(db *Database, w http.ResponseWriter, r *http.Request) {
	js, err := json.Marshal(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	w.Write([]byte{ '\n' })
}

func ListenAndServeInBackground(port uint, eq chan error, db *Database) {

	// index handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		indexHandler(db, w, r)
	})

	// serve in background
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
		if err != nil {
			eq <- err
		}
	}()
}
