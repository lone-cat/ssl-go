package storage

import (
	"encoding/pem"
)

type Pem interface {
	Load() ([]*pem.Block, error)
	Save([]*pem.Block) error
	Delete() error
}
