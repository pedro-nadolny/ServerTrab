package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

/* usuario:senha@/nomedobd */
const banco = "sd20171:sd20171@/sd20171"

type Baladeiro struct {
	Id      int    `json:"id"`
	Nome    string `json:"nome"`
	Email   string `json:"email"`
	Idade   int    `json:"idade"`
	Points  int    `json:"points"`
	CheckIn bool   `json:"checkIn`
	Balada  string `json:"balada"`
	Estilo  string `json:"estilo"`
}
type Id string

func handleBaladeiros(resp http.ResponseWriter, req *http.Request) {
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
	r, err := bd.Query("select * from USERBALADEIRO;")
	if err != nil {
		log.Println(err.Error())
		http.Error(resp, err.Error(),
			http.StatusInternalServerError)
		return
	}

	/* armazena o resultado da consulta */
	baladeiros := make(map[int]*Baladeiro)
	for r.Next() {
		var id int
		baladeiro := new(Baladeiro)
		err = r.Scan(&id, &baladeiro.Nome, &baladeiro.Email, &baladeiro.Idade, &baladeiro.Points, &baladeiro.CheckIn, &baladeiro.Balada, &baladeiro.Estilo)
		if err != nil {
			log.Println(err.Error())
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}
		baladeiros[id] = baladeiro
	}

	/* serializa o resultado em JSON e retorna */
	outgoingJSON, err := json.MarshalIndent(baladeiros, "", "\t")
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
	cmd, err := bd.Prepare("delete from USERBALADEIRO where id = ?;")
	if err != nil {
		log.Println(err.Error())
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	cmd.Exec(id)
	resp.WriteHeader(http.StatusNoContent)
}

func handleBaladeiro(resp http.ResponseWriter, req *http.Request) {
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
	baladeiro := new(Baladeiro)
	cmd, err := bd.Prepare(`insert into USERBALADEIRO
(id, nome, email, idade, points, checkIn, balada, estilo)
values (?, ?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		log.Println(err.Error())
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&baladeiro)
	if err != nil {
		log.Println(err.Error())
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = cmd.Exec(id,
		baladeiro.Nome,
		baladeiro.Email,
		baladeiro.Idade,
		baladeiro.Points,
		baladeiro.CheckIn,
		baladeiro.Balada,
		baladeiro.Estilo)

	if err != nil {
		resp.WriteHeader(http.StatusConflict)
		fmt.Fprintf(resp, "%s\n", string("Baladeiro já existe!"))
		return
	}
	outgoingJSON, err := json.Marshal(baladeiro)
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
	baladeiro := new(Baladeiro)
	cmd, err := bd.Prepare(`update USERBALADEIRO
	set nome = ?, email = ?, idade = ?, points = ?, checkIn = ?, balada = ?, estilo = ?,
	where id = ?;`)
	if err != nil {
		log.Println(err.Error())
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&baladeiro)
	if err != nil {
		log.Println(err.Error())
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if baladeiro.Nome == "" ||
		baladeiro.Email == "" ||
		baladeiro.Idade == 0 ||
		baladeiro.Balada == "" ||
		baladeiro.Estilo == "" {
		resp.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(resp, "%s\n", string("Não foi possível efetuar a atualização!"))
		return
	}

	package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

/* usuario:senha@/nomedobd */
const banco = "sd20171:sd20171@/sd20171"

type Baladeiro struct {
	Id      int    `json:"id"`
	Nome    string `json:"nome"`
	Email   string `json:"email"`
	Idade   int    `json:"idade"`
	Points  int    `json:"points"`
	CheckIn bool   `json:"checkIn`
	Balada  string `json:"balada"`
	Estilo  string `json:"estilo"`
}
type Id string

func handleBaladeiros(resp http.ResponseWriter, req *http.Request) {
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
	r, err := bd.Query("select * from USERBALADEIRO;")
	if err != nil {
		log.Println(err.Error())
		http.Error(resp, err.Error(),
			http.StatusInternalServerError)
		return
	}

	/* armazena o resultado da consulta */
	baladeiros := make(map[int]*Baladeiro)
	for r.Next() {
		var id int
		baladeiro := new(Baladeiro)
		err = r.Scan(&id, &baladeiro.Nome, &baladeiro.Email, &baladeiro.Idade, &baladeiro.Points, &baladeiro.CheckIn, &baladeiro.Balada, &baladeiro.Estilo)
		if err != nil {
			log.Println(err.Error())
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}
		baladeiros[id] = baladeiro
	}

	/* serializa o resultado em JSON e retorna */
	outgoingJSON, err := json.MarshalIndent(baladeiros, "", "\t")
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
	cmd, err := bd.Prepare("delete from USERBALADEIRO where id = ?;")
	if err != nil {
		log.Println(err.Error())
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	cmd.Exec(id)
	resp.WriteHeader(http.StatusNoContent)
}

func handleBaladeiro(resp http.ResponseWriter, req *http.Request) {
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
	baladeiro := new(Baladeiro)
	cmd, err := bd.Prepare(`insert into USERBALADEIRO
(id, nome, email, idade, points, checkIn, balada, estilo)
values (?, ?, ?, ?, ?, ?, ?, ?);`)
	if err != nil {
		log.Println(err.Error())
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&baladeiro)
	if err != nil {
		log.Println(err.Error())
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = cmd.Exec(id,
		baladeiro.Nome,
		baladeiro.Email,
		baladeiro.Idade,
		baladeiro.Points,
		baladeiro.CheckIn,
		baladeiro.Balada,
		baladeiro.Estilo)

	if err != nil {
		resp.WriteHeader(http.StatusConflict)
		fmt.Fprintf(resp, "%s\n", string("Baladeiro já existe!"))
		return
	}
	outgoingJSON, err := json.Marshal(baladeiro)
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
	baladeiro := new(Baladeiro)
	cmd, err := bd.Prepare(`update USERBALADEIRO
	set nome = ?, email = ?, idade = ?, points = ?, checkIn = ?, balada = ?, estilo = ?,
	where id = ?;`)
	if err != nil {
		log.Println(err.Error())
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&baladeiro)
	if err != nil {
		log.Println(err.Error())
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	if baladeiro.Nome == "" ||
		baladeiro.Email == "" ||
		baladeiro.Idade == 0 ||
		baladeiro.Balada == "" ||
		baladeiro.Estilo == "" {
		resp.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(resp, "%s\n", string("Não foi possível efetuar a atualização!"))
		return
	}

	cmd.Exec(id,
		baladeiro.Nome,
		baladeiro.Email,
		baladeiro.Idade,
		baladeiro.Points,
		baladeiro.CheckIn,
		baladeiro.Balada,
		baladeiro.Estilo)

	resp.WriteHeader(http.StatusNoContent)
}

func (self Id) Visualiza(resp http.ResponseWriter, req *http.Request, bd *sql.DB) {
	id := string(self)
	baladeiro := new(Baladeiro)
	r := bd.QueryRow(`select nome, email, idade, points, checkIn, balada, estilo from USERBALADEIRO where id = ?;`, id)
	err := r.Scan(&id, &baladeiro.Nome, &baladeiro.Email, &baladeiro.Idade, &baladeiro.Points, &baladeiro.CheckIn, &baladeiro.Balada, &baladeiro.Estilo)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(resp, "%s\n", string("Baladeiro não encontrado!"))
		return
	}
	outgoingJSON, err := json.Marshal(baladeiro)
	if err != nil {
		log.Println(err.Error())
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(resp, "%s\n", string(outgoingJSON))
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/brm", handleBaladeiros).Methods("GET")
	router.HandleFunc("/brm/{id}", handleBaladeiro).Methods("GET", "DELETE", "POST", "PUT")
	log.Fatal(http.ListenAndServe(":9999", router))
}

	cmd.Exec(randonneur.Nome, randonneur.Tempo, randonneur.Media, id)
	resp.WriteHeader(http.StatusNoContent)
}

func (self Id) Visualiza(resp http.ResponseWriter,
	req *http.Request, bd *sql.DB) {
	id := string(self)
	baladeiro := new(Baladeiro)
	r := bd.QueryRow(`select nome, email, idade, points, checkIn, balada, estilo from USERBALADEIRO where id = ?;`, id)
	err := r.Scan(&id, &baladeiro.Nome, &baladeiro.Email, &baladeiro.Idade, &baladeiro.Points, &baladeiro.CheckIn, &baladeiro.Balada, &baladeiro.Estilo)
	if err != nil {
		resp.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(resp, "%s\n", string("Baladeiro não encontrado!"))
		return
	}
	outgoingJSON, err := json.Marshal(baladeiro)
	if err != nil {
		log.Println(err.Error())
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(resp, "%s\n", string(outgoingJSON))
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/brm", handleBaladeiros).Methods("GET")
	router.HandleFunc("/brm/{id}", handleBaladeiro).Methods("GET", "DELETE", "POST", "PUT")
	log.Fatal(http.ListenAndServe(":9999", router))
}
