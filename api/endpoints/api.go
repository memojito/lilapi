package endpoints

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gocql/gocql"
	"log"
	"net/http"
)

type Workspace struct {
	ID       gocql.UUID   `json:"id"`
	Owner    string       `json:"owner"`
	Editor   string       `json:"editor"`
	Size     int          `json:"size"`
	Name     string       `json:"name"`
	FacetIDs []gocql.UUID `json:"facet_ids"`
}

type Facet struct {
	ID           gocql.UUID   `json:"id"`
	Name         string       `json:"name"`
	Value        string       `json:"value"`
	WorkspaceIDs []gocql.UUID `json:"workspace_ids"`
}

type api struct {
	router  *chi.Mux
	session *gocql.Session
}

func NewAPI(session *gocql.Session) *api {
	api := &api{
		router:  chi.NewRouter(),
		session: session,
	}
	api.router.Use(middleware.Logger)
	api.routes()
	return api
}

func (api *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.router.ServeHTTP(w, r)
}

func (api *api) routes() {
	api.router.Route("/workspace", func(r chi.Router) {
		r.Get("/{id}", api.getWorkspace)
		r.Post("/", api.createWorkspace)
		//r.Put("/{id}", api.updateWorkspace)	// coming soon
		//r.Delete("/{id}", api.deleteWorkspace)
	})

	api.router.Route("/facets", func(r chi.Router) {
		r.Get("/{id}", api.getFacet)
		//r.Post("/", api.createFacet)	// coming soon
		//r.Put("/{id}", api.updateFacet)
		//r.Delete("/{id}", api.deleteFacet)
	})
}

func (api *api) getWorkspace(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	workspace := Workspace{}
	if err := api.session.Query(`SELECT id, owner, editor, size, name, facet_ids FROM workspace WHERE id = ?`, id).Scan(&workspace.ID, &workspace.Owner, &workspace.Editor, &workspace.Size, &workspace.Name, &workspace.FacetIDs); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(404), 404)
		return
	}
	json.NewEncoder(w).Encode(workspace)
}

func (api *api) createWorkspace(w http.ResponseWriter, r *http.Request) {
	workspace := Workspace{}
	if err := json.NewDecoder(r.Body).Decode(&workspace); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(400), 400)
		return
	}
	if workspace.ID.String() == "" {
		workspace.ID = gocql.TimeUUID()
	}
	if err := api.session.Query(`INSERT INTO workspace (id, owner, photographer, size, name, facet_ids) VALUES (?, ?, ?, ?, ?, ?)`,
		workspace.ID, workspace.Owner, workspace.Editor, workspace.Size, workspace.Name, workspace.FacetIDs).Exec(); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(workspace)
}

func (api *api) getFacet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	facet := Facet{}
	if err := api.session.Query(`SELECT name, value, workspace_ids FROM facets WHERE id = ?`, id).Scan(&facet.Name,
		&facet.Value, &facet.WorkspaceIDs); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(404), 404)
		return
	}
	json.NewEncoder(w).Encode(facet)
}
