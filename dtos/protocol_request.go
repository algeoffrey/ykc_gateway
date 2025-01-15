package dtos

type VerificationResponseMessage struct {
	Header *Header `json:"header"`
	Id     string  `json:"id"`
	Result bool    `json:"result"`
}

type BillingModelRequestMessage struct {
	Header *Header `json:"header"`
	Id     string  `json:"Id"`
}

type RemoteBootstrapRequestMessage struct {
	Header       *Header `json:"header"`
	TradeSeq     string  `json:"tradeSeq"`
	Id           string  `json:"id"`
	GunId        string  `json:"gunId"`
	LogicCard    string  `json:"logicCard"`
	PhysicalCard string  `json:"physicalCard"`
	Balance      int     `json:"balance"`
}
type RemoteRebootRequestMessage struct {
	Header  *Header `json:"header"`
	Id      string  `json:"id"`
	Control int     `json:"control"`
}

type SetBillingModelRequestMessage struct {
	Header           *Header `json:"header"`
	Id               string  `json:"id"`
	BillingModelCode string  `json:"billingModelCode"`
	SharpUnitPrice   int     `json:"sharpUnitPrice"`
	SharpServiceFee  int     `json:"sharpServiceFee"`
	PeakUnitPrice    int     `json:"peakUnitPrice"`
	PeakServiceFee   int     `json:"peakServiceFee"`
	FlatUnitPrice    int     `json:"flatUnitPrice"`
	FlatServiceFee   int     `json:"flatServiceFee"`
	ValleyUnitPrice  int     `json:"valleyUnitPrice"`
	ValleyServiceFee int     `json:"valleyServiceFee"`
	AccrualRatio     int     `json:"accrualRatio"`
	RateList         []int   `json:"rateList"`
}

type StartChargingRequest struct {
	ClientID    string `json:"clientID"`
	Port        int    `json:"port"`
	OrderNumber string `json:"orderNumber"`
}

type StopChargingRequest struct {
	ClientID    string `json:"clientID"`
	Port        int    `json:"port"`
	OrderNumber string `json:"orderNumber"`
}
