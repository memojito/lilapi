package model

import "github.com/gocql/gocql"

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
