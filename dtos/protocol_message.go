package dtos

type DeviceLoginResponseMessage struct {
	Header          *Header `json:"header"`
	HeartbeatPeriod int     `json:"heartbeatPeriod"` // Heartbeat interval in seconds
	Result          byte    `json:"result"`          // Login Result (0x00 = success, 0x01 = illegal module, 0xF0 = protocol upgrade)
}

type RemoteStartResponseMessage struct {
	Header      *Header `json:"header"`
	Port        int     `json:"port"`
	OrderNumber uint32  `json:"orderNumber"`
	StartMethod int     `json:"startMethod"`
	Result      int     `json:"result"`
}
