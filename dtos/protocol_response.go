package dtos

type VerificationMessage struct {
	Header          *Header `json:"header"`
	Id              string  `json:"Id"`
	ElcType         int     `json:"elcType"`
	Guns            int     `json:"guns"`
	ProtocolVersion int     `json:"protocolVersion"`
	SoftwareVersion string  `json:"softwareVersion"`
	Network         int     `json:"network"`
	Sim             string  `json:"sim"`
	Operator        int     `json:"operator"`
}

type BillingModelVerificationMessage struct {
	Header           *Header `json:"header"`
	Id               string  `json:"Id"`
	BillingModelCode string  `json:"billingModelCode"`
}

type BillingModelResponseMessage struct {
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

type BillingModelVerificationResponseMessage struct {
	Header           *Header `json:"header"`
	Id               string  `json:"id"`
	BillingModelCode string  `json:"billingModelCode"`
	Result           bool    `json:"result"`
}

type HeartbeatMessage struct {
	Header         *Header `json:"header"`
	SignalValue    int     `json:"signalValue"`
	Temperature    int     `json:"temperature"`
	TotalPortCount int     `json:"totalPortCount"`
	PortStatus     []int   `json:"portStatus"`
}

type HeartbeatResponseMessage struct {
	Header   *Header `json:"header"`
	Id       string  `json:"Id"`
	Gun      string  `json:"gun"`
	Response int     `json:"response"`
}

type RemoteBootstrapResponseMessage struct {
	Header   *Header `json:"header"`
	TradeSeq string  `json:"tradeSeq"`
	Id       string  `json:"id"`
	GunId    string  `json:"gunId"`
	Result   bool    `json:"result"`
	Reason   int     `json:"reason"`
}

type OfflineDataReportMessage struct {
	Header                  *Header `json:"header"`
	TradeSeq                string  `json:"tradeSeq"`
	Id                      string  `json:"id"`
	GunId                   string  `json:"gunId"`
	Status                  int     `json:"status"`
	Reset                   int     `json:"reset"`
	Plugged                 int     `json:"plugged"`
	Ov                      int     `json:"ov"`
	Oc                      int     `json:"oc"`
	LineTemp                int     `json:"lineTemp"`
	LineCode                string  `json:"lineCode"`
	Soc                     int     `json:"soc"`
	BpTopTemp               int     `json:"bpTopTemp"`
	AccumulatedChargingTime int     `json:"accumulatedChargingTime"`
	RemainingTime           int     `json:"remainingTime"`
	ChargingDegrees         int     `json:"chargingDegrees"`
	LossyChargingDegrees    int     `json:"lossyChargingDegrees"`
	ChargedAmount           int     `json:"chargedAmount"`
	HardwareFailure         int     `json:"hardwareFailure"`
}

type RemoteShutdownResponseMessage struct {
	Header *Header `json:"header"`
	Id     string  `json:"id"`
	GunId  string  `json:"gunId"`
	Result bool    `json:"result"`
	Reason int     `json:"reason"`
}

type RemoteShutdownRequestMessage struct {
	Header *Header `json:"header"`
	Id     string  `json:"id"`
	GunId  string  `json:"gunId"`
}

type TransactionRecordMessage struct {
	Header                    *Header `json:"header"`
	TradeSeq                  string  `json:"tradeSeq"`
	Id                        string  `json:"id"`
	GunId                     string  `json:"gunId"`
	StartAt                   int64   `json:"startAt"`
	EndAt                     int64   `json:"endAt"`
	SharpUnitPrice            int64   `json:"sharpUnitPrice"`
	SharpElectricCharge       int64   `json:"sharpElectricCharge"`
	LossySharpElectricCharge  int64   `json:"lossySharpElectricCharge"`
	SharpPrice                int64   `json:"sharpPrice"`
	PeakUnitPrice             int64   `json:"peakUnitPrice"`
	PeakElectricCharge        int64   `json:"peakElectricCharge"`
	LossyPeakElectricCharge   int64   `json:"lossyPeakElectricCharge"`
	PeakPrice                 int64   `json:"peakPrice"`
	FlatUnitPrice             int64   `json:"flatUnitPrice"`
	FlatElectricCharge        int64   `json:"flatElectricCharge"`
	LossyFlatElectricCharge   int64   `json:"lossyFlatElectricCharge"`
	FlatPrice                 int64   `json:"flatPrice"`
	ValleyUnitPrice           int64   `json:"valleyUnitPrice"`
	ValleyElectricCharge      int64   `json:"valleyElectricCharge"`
	LossyValleyElectricCharge int64   `json:"lossyValleyElectricCharge"`
	ValleyPrice               int64   `json:"valleyPrice"`
	InitialMeterReading       int64   `json:"initialMeterReading"`
	FinalMeterReading         int64   `json:"finalMeterReading"`
	TotalElectricCharge       int64   `json:"totalElectricCharge"`
	LossyTotalElectricCharge  int64   `json:"lossyTotalElectricCharge"`
	ConsumptionAmount         int64   `json:"consumptionAmount"`
	Vin                       string  `json:"vin"`
	StartType                 int     `json:"startType"`
	TransactionDateTime       int64   `json:"transactionDateTime"`
	StopReason                int     `json:"stopReason"`
	PhysicalCardNumber        string  `json:"physicalCardNumber"`
}

type TransactionRecordConfirmedMessage struct {
	Header   *Header `json:"header"`
	Id       string  `json:"id"`
	TradeSeq string  `json:"tradeSeq"`
	Result   int     `json:"result"`
}

type RemoteRebootResponseMessage struct {
	Header *Header `json:"header"`
	Id     string  `json:"id"`
	Result int     `json:"result"`
}

type SetBillingModelResponseMessage struct {
	Header *Header `json:"header"`
	Id     string  `json:"id"`
	Result int     `json:"result"`
}

type ChargingFinishedMessage struct {
	Header                           *Header `json:"header"`
	TradeSeq                         string  `json:"tradeSeq"`
	Id                               string  `json:"id"`
	GunId                            string  `json:"gunId"`
	BmsSoc                           int     `json:"bmsSoc"`
	BmsBatteryPackLowestVoltage      int     `json:"bmsBatteryPackLowestVoltage"`
	BmsBatteryPackHighestVoltage     int     `json:"bmsBatteryPackHighestVoltage"`
	BmsBatteryPackLowestTemperature  int     `json:"bmsBatteryPackLowestTemperature"`
	BmsBatteryPackHighestTemperature int     `json:"bmsBatteryPackHighestTemperature"`
	CumulativeChargingDuration       int     `json:"cumulativeChargingDuration"`
	OutputPower                      int     `json:"outputPower"`
	ChargingUnitId                   int     `json:"chargingUnitId"`
}

type DeviceLoginMessage struct {
	Header          *Header `json:"Header"`
	IMEI            string  `json:"imei"`
	DevicePortCount int     `json:"devicePortCount"`
	HardwareVersion string  `json:"hardwareVersion"`
	SoftwareVersion string  `json:"softwareVersion"`
	CCID            string  `json:"ccid"`
	SignalValue     int     `json:"signalValue"`
	LoginReason     int     `json:"loginReason"`
}

type RemoteStartMessage struct {
	Header      *Header `json:"header"`
	Port        int     `json:"port"`
	OrderNumber uint32  `json:"orderNumber"`
	StartMode   int     `json:"startMode"`
	StartResult int     `json:"startResult"`
}

type RemoteStopMessage struct {
	Header      *Header `json:"header"`
	Port        int     `json:"port"`
	OrderNumber uint32  `json:"orderNumber"`
}

type RemoteStopResponseMessage struct {
	Header      *Header `json:"header"`
	Port        int     `json:"port"`
	OrderNumber uint32  `json:"orderNumber"`
	Result      byte    `json:"result"`
}

type SubmitFinalStatusMessage struct {
	Header           *Header  `json:"header"`
	Port             byte     `json:"port"`
	OrderNumber      uint32   `json:"orderNumber"`
	ChargingTime     uint32   `json:"chargingTime"`
	ElectricityUsage uint32   `json:"electricityUsage"`
	UsageCost        uint32   `json:"usageCost"`
	StopReason       byte     `json:"stopReason"`
	StopPower        uint16   `json:"stopPower"`
	CardID           uint32   `json:"cardId"`
	SegmentCount     byte     `json:"segmentCount"`
	SegmentDurations []uint16 `json:"segmentDurations"`
	SegmentPrices    []uint16 `json:"segmentPrices"`
	Reserved         []byte   `json:"reserved"`
}

type SubmitFinalStatusResponse struct {
	Header *Header `json:"header"`
	Result byte    `json:"result"`
}

type ChargingPortDataMessage struct {
	Header          *Header `json:"header"`
	Reserved        byte    `json:"reserved"`
	PortCount       byte    `json:"portCount"`
	Voltage         uint16  `json:"voltage"`
	Temperature     byte    `json:"temperature"`
	ActivePort      byte    `json:"activePort"`
	CurrentTier     byte    `json:"currentTier"`
	CurrentRate     uint16  `json:"currentRate"`
	CurrentPower    uint16  `json:"currrentPower"`
	UsageTime       uint32  `json:"usageTime"`
	UsedAmount      uint16  `json:"usedAmount"`
	EnergyUsed      uint32  `json:"energyUsed"`
	PortTemperature byte    `json:"portTemperature"`
}
