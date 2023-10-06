package store

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Tenant struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	ApiKeys []string `json:"apiKeys"`
}

func NewTenant(name string) *Tenant {
	return &Tenant{
		Id:   uuid.NewString(),
		Name: name,
	}
}

func (t *Tenant) ToJSON() ([]byte, error) {
	return json.Marshal(t)
}
