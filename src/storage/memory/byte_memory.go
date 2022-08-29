package memory

type byteMemory struct {
	data []byte
}

func NewByteMemory() *byteMemory {
	return &byteMemory{}
}

func (s *byteMemory) Load() (data []byte, err error) {
	if s.data == nil {
		return
	}

	data = make([]byte, len(s.data))
	copy(data, s.data)

	return
}

func (s *byteMemory) Save(data []byte) (err error) {
	if data == nil {
		err = s.Delete()
		return
	}

	s.data = make([]byte, len(data))
	copy(s.data, data)

	return
}

func (s *byteMemory) Delete() (err error) {
	s.data = nil
	return
}
