package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
)

type config struct {
	StoragePath string `env:"STORAGE_PATH"`
	Port        string `env:"PORT"`
}

type load struct {
	Key   string                 `json:"key"`
	Value map[string]interface{} `json:"value"`
}

func getValue(w http.ResponseWriter, r *http.Request, db Storage) {
	key := chi.URLParam(r, "key")
	value, wasFound, err := db.Get([]byte(key))

	w.Header().Set("Content-Type", "application/json")
	if err == nil {
		if wasFound {
			w.Write(value)
		} else {
			http.Error(w, "Key not found.", http.StatusNotFound)
		}
	} else {
		http.Error(w, "Oops something went wrong.", http.StatusInternalServerError)
	}
}

func createValue(w http.ResponseWriter, r *http.Request, db Storage) {
	var input load
	err := json.NewDecoder(r.Body).Decode(&input)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	byteValue, err := json.Marshal(input.Value)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.Set([]byte(input.Key), byteValue)

	if err != nil {
		http.Error(w, "Oops something went wrong.", http.StatusInternalServerError)
	}
}

func startServer(s Storage, port string) {
	router := chi.NewRouter()
	router.Get("/key/{key}", func(w http.ResponseWriter, r *http.Request) {
		getValue(w, r, s)
	})
	router.Post("/key", func(w http.ResponseWriter, r *http.Request) {
		createValue(w, r, s)
	})

	log.Printf("Starting up on http://localhost:%s", port)
	http.ListenAndServe(":"+port, router)
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}

	storage := NewPebbleStorage(cfg.StoragePath)

	startServer(storage, cfg.Port)
	storage.Close()
}
