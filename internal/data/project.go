package data

import "github.com/google/uuid"

type Project struct {
	ProjectID   uuid.UUID `json:"projectId"`
	Name        string    `json:"name"`
	UnzipedPath string    `json:"unzip"`
	ZippedPath  string    `json:"zip"`
}
