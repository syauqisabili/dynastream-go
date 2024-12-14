package domain

type Stream struct {
	Id   string `json:"id"`
	Uuid string `json:"uuid"`
}

type StreamRepository interface {
	GetAll() ([]*Stream, error)
	Insert(stream *Stream) error
	Delete(uuid string) error
}
