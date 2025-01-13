package dtos

type Header struct {
	Length    int    `json:"length"`
	Seq       int    `json:"seq"`
	Encrypted bool   `json:"encrypted"`
	FrameId   string `json:"frameId"`
}
type ClientInfo struct {
	IPAddress string
	IMEI      string
}
