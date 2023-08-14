package endpoints

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gocql/gocql"
	"log"
	"net/http"
)

type Transaction struct {
	ID    gocql.UUID `json:"id"`
	Name  string     `json:"name"`
	Value int        `json:"value"`
}

type api struct {
	router  *chi.Mux
	session *gocql.Session
	token   string
}

func NewAPI(session *gocql.Session, token string) *api {
	api := &api{
		router:  chi.NewRouter(),
		session: session,
		token:   token,
	}
	api.router.Use(middleware.Logger)
	api.routes()
	return api
}

func (api *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.router.ServeHTTP(w, r)
}

func (api *api) routes() {
	api.router.Route("/transaction", func(r chi.Router) {
		r.Get("/{id}", api.getTransaction)
		r.Post("/", api.createTransaction)
		//r.Put("/{id}", api.updateWorkspace)	// coming soon
		//r.Delete("/{id}", api.deleteWorkspace)
	})
}

func (api *api) getTransaction(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	transaction := Transaction{}
	if err := api.session.Query(`SELECT id, name, value FROM transaction WHERE id = ?`, id).Scan(&transaction.ID, &transaction.Name, &transaction.Value); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(404), 404)
		return
	}
	json.NewEncoder(w).Encode(transaction)
}

func (api *api) createTransaction(w http.ResponseWriter, r *http.Request) {
	transaction := Transaction{}
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(400), 400)
		return
	}
	log.Printf("Transaction.ID =%v", transaction.ID)
	if transaction.ID.String() == "00000000-0000-0000-0000-000000000000" {
		transaction.ID, _ = gocql.RandomUUID()
	}
	if err := api.session.Query(`INSERT INTO transaction (id, name, value) VALUES (?, ?, ?)`,
		transaction.ID, transaction.Name, transaction.Value).Exec(); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}
