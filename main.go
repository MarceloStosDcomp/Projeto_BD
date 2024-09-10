package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	users = make(map[int]User)
	mu    sync.Mutex
)

type User struct {
	CPF            int       `json:"cpf"`
	Nome           string    `json:"nome"`
	DataNascimento time.Time `json:"dataNascimento"`
}

func main() {
	http.HandleFunc("/user", userHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handlePost(w, r)
	case http.MethodGet:
		handleGet(w, r)
	default:
		http.Error(w, "Metodo de Requisão não suportado.", http.StatusMethodNotAllowed)
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Requisição Inválida", http.StatusBadRequest)
		return
	}

	mu.Lock()
	users[user.CPF] = user
	mu.Unlock()

	w.WriteHeader(http.StatusCreated)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	cpfStr := r.URL.Query().Get("cpf")
	if cpfStr == "" {
		http.Error(w, "CPF é obrigatório", http.StatusBadRequest)
		return
	}

	cpf, err := strconv.Atoi(cpfStr)
	if err != nil {
		http.Error(w, "CPF Inválido", http.StatusBadRequest)
		return
	}

	mu.Lock()
	user, existe := users[cpf]
	mu.Unlock()

	if !existe {
		http.Error(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
