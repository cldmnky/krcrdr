package store

import (
	"encoding/json"

	"github.com/google/uuid"
)

const (
	KVPut KVWatchOp = iota
	KVDelete
	KVPurge
)

type (
	Tenant struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	}

	KVWatchOp int64
)

func NewTenant(name string) *Tenant {
	return &Tenant{
		Id:   uuid.NewString(),
		Name: name,
	}
}

func (t *Tenant) ToJSON() ([]byte, error) {
	return json.Marshal(t)
}

func (k KVWatchOp) String() string {
	switch k {
	case KVPut:
		return "put"
	case KVDelete:
		return "delete"
	case KVPurge:
		return "purge"
	default:
		return "unknown"
	}
}
