package storage

type ByteMulti interface {
	Load() ([][]byte, error)
	Save([][]byte) error
	Delete() error
}
