package dto

type StreamSession struct {
	ID                        string `json:"id"`
	Created                   string `json:"created"`
	RemoteAddr                string `json:"remoteAddr"`
	PeerConnectionEstablished bool   `json:"peerConnectionEstablished"`
	LocalCandidate            string `json:"localCandidate"`
	RemoteCandidate           string `json:"remoteCandidate"`
	State                     string `json:"state"`
	Path                      string `json:"path"`
	Query                     string `json:"query"`
	BytesReceived             int    `json:"bytesReceived"`
	BytesSent                 int    `json:"bytesSent"`
}

type StreamSessionList struct {
	ItemCount int             `json:"itemCount"`
	PageCount int             `json:"pageCount"`
	Items     []StreamSession `json:"items"`
}
