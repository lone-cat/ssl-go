package storage

type Byte interface {
	Load() ([]byte, error)
	Save([]byte) error
	Delete() error
}
