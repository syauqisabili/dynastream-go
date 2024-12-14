package network

type NetConn struct {
	Ip   string `json:"ip"`
	Port uint16 `json:"port"`
}

type Rtsp struct {
	Ip   string `json:"ip"`
	Port uint16 `json:"port"`
	Path string `json:"path"`
}

type MediaMtx struct {
	Http   NetConn `json:"http"`
	Rtsp   Rtsp    `json:"rtsp"`
	WebRtc NetConn `json:"webrtc"`
}

type Redis struct {
	Ip            string `json:"ip"`
	Port          uint16 `json:"port"`
	Password      string `json:"password"`
	DatabaseIndex uint8  `json:"database_index"`
}

type NetCfg struct {
	MediaMtx MediaMtx `json:"mediamtx"`
	Grpc     NetConn  `json:"grpc"`
	Redis    Redis    `json:"redis"`
}
