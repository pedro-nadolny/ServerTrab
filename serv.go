package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

/* usuario:senha@/nomedobd */
const banco = "sd20171:sd20171@/sd20171"

type Randonneur struct {
	Nome  string        `json:"nome"`
	Tempo time.Duration `json:"tempo"`
	Media float32       `json:"media"`
}
type Id string

func handleRandonneurs(resp http.ResponseWriter, req *http.Request) {
	/* configura o content-type e conecta ao BD */
	/* efetua a consulta */
	/* armazena o resultado da consulta */
	/* serializa o resultado em JSON e retorna */

	/* configura o content-type e conecta ao BD */
	resp.Header().Set("Content-Type", "application/json")
	bd, err := sql.Open("mysql", banco)
	if err != nil {
		log.Println(err.Error())
		http.Error(resp, err.Error(),
			http.StatusInternalServerError)
		return
	}
	defer bd.Close()

	/* efetua a consulta */
	r, err := bd.Query("select * from WILLIAN;")
	if err != nil {
		log.Println(err.Error())
		http.Error(resp, err.Error(),
			http.StatusInternalServerError)
		return
	}

	/* armazena o resultado da consulta */
	randonneurs := make(map[int]*Randonneur)
	for r.Next() {
		var id int
		randonneur := new(Randonneur)
		err = r.Scan(&id, &randonneur.Nome, &randonneur.Tempo, &randonneur.Media)
		if err != nil {
			log.Println(err.Error())
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}
		randonneurs[id] = randonneur
	}

	/* serializa o resultado em JSON e retorna */
	outgoingJSON, err := json.MarshalIndent(randonneurs, "", "\t")
	if err != nil {
		log.Println(err.Error())
		http.Error(resp, err.Error(),
			http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(resp, "%s\n", string(outgoingJSON))

}

func (self Id) Exclui(resp http.ResponseWriter, req *http.Request, bd *sql.DB) {
	id := string(self)
	cmd, err := bd.Prepare("delete from WILLIAN where placa = ?;")
	if err != nil {
		log.Println(err.Error())
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	cmd.Exec(id)
	resp.WriteHeader(http.StatusNoContent)
}

func handleRandonneur(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	bd, err := sql.Open("mysql", banco)
	if err != nil {
		log.Println(err.Error())
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	defer bd.Close()
	vars := mux.Vars(req)
	placa := vars["id"]
	id := Id(placa)
	log.Println("Requisitou: ", id)

	switch req.Method {
	case "GET":
		id.Visualiza(resp, req, bd)
	case "DELETE":
		id.Exclui(resp, req, bd)
	case "PUT":
		id.Atualiza(resp, req, bd)
	case "POST":
		id.Insere(resp, req, bd)
	}

}

func (self Id) Insere(resp http.ResponseWriter, req *http.Request, bd *sql.DB) {
	id := string(self)
	randonneur := new(Randonneur)
	cmd, err := bd.Prepare(`insert into WILLIAN
(placa, nome, tempo, media)
values (?, ?, ?, ?);`)
	if err != nil {
		log.Println(err.Error())
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&randonneur)
	if err != nil {
		log.Println(err.Error())
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = cmd.Exec(id, randonneur.Nome, randonneur.Tempo,
		randonneur.Media)
	if err != nil {
		resp.WriteHeader(http.StatusConflict)
		fmt.Fprintf(resp, "%s\n", string("Randonneur já existe!"))
		return
	}
	outgoingJSON, err := json.Marshal(randonneur)
	if err != nil {
		log.Println(err.Error())
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	resp.WriteHeader(http.StatusCreated)
	fmt.Fprintf(resp, "%s\n", string(outgoingJSON))
}

func (self Id) Atualiza(resp http.ResponseWriter, req *http.Request, bd *sql.DB) {
	id := string(self)
	randonneur := new(Randonneur)
	cmd, err := bd.Prepare(`update WILLIAN
	set nome = ?, tempo = ?, media = ?
	where placa = ?;`)
	if err != nil {
		log.Println(err.Error())
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&randonneur)
	if err != nil {
		log.Println(err.Error())
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if randonneur.Nome == "" ||
		randonneur.Tempo == 0 ||
		randonneur.Media == 0.0 {
		resp.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(resp, "%s\n",
			string("Não foi possível efetuar a atualização!"))
		return
	}
	cmd.Exec(randonneur.Nome, randonneur.Tempo, randonneur.Media, id)
	resp.WriteHeader(http.StatusNoContent)
}

func (self Id) Visualiza(resp http.ResponseWriter,
	req *http.Request, bd *sql.DB) {
	id := string(self)
	randonneur := new(Randonneur)
	r := bd.QueryRow(`select nome, tempo, media from WILLIAN where placa = ?;`, id)
	err := r.Scan(&randonneur.Nome, &randonneur.Tempo, &randonneur.Media)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(resp, "%s\n", string("Randonneur não encontrado!"))
		return
	}
	outgoingJSON, err := json.Marshal(randonneur)
	if err != nil {
		log.Println(err.Error())
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(resp, "%s\n", string(outgoingJSON))
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/brm", handleRandonneurs).Methods("GET")
	router.HandleFunc("/brm/{id}", handleRandonneur).Methods("GET", "DELETE", "POST", "PUT")
	log.Fatal(http.ListenAndServe(":9998", router))
}
