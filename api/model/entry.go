package model

import "github.com/gocql/gocql"

type Workspace struct {
	ID       gocql.UUID
	Owner    string
	Editor   string
	Size     int
	Name     string
	FacetIDs []gocql.UUID
}

type Facet struct {
	ID           gocql.UUID
	Name         string
	Value        string
	WorkspaceIDs []gocql.UUID
}
