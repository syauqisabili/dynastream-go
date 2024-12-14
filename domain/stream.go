package domain

type Stream struct {
	Id   string `json:"id"`
	Uuid string `json:"uuid"`
}

type StreamRepository interface {
	GetAll() ([]*Stream, error)
	FindByUuid(uuid string) *Stream
	Insert(stream *Stream) error
	Delete(uuid string) error
}
