package protocols

const (
	//start flag
	StartFlag = byte(0x68)

	//device -> platform
	Verification                = byte(0x01)
	Heartbeat                   = byte(0x82)
	BillingModelVerification    = byte(0x05)
	BillingModelRequest         = byte(0x09)
	OfflineDataReport           = byte(0x13)
	ChargingHandshake           = byte(0x15)
	Configuration               = byte(0x17)
	ChargingFinished            = byte(0x19)
	ErrorReport                 = byte(0x1b)
	BmsInterrupted              = byte(0x1d)
	ChargingPileInterrupted     = byte(0x21)
	ChargingMetrics             = byte(0x23)
	BmsInformation              = byte(0x25)
	ActiveChargingRequest       = byte(0x31)
	RemoteBootstrapResponse     = byte(0x33)
	RemoteShutdownResponse      = byte(0x35)
	TransactionRecord           = byte(0x3b)
	BalanceUpdateResponse       = byte(0x41)
	CardSynchronizationResponse = byte(0x43)
	CardClearingResponse        = byte(0x45)
	CardQueryingResponse        = byte(0x47)
	SetWorkingParamsResponse    = byte(0x51)
	NtpResponse                 = byte(0x55)
	SetBillingModelResponse     = byte(0x57)
	FloorLockDataUpload         = byte(0x61)
	Response                    = byte(0x63)
	RemoteRebootResponse        = byte(0x91)
	OtaResponse                 = byte(0x93)

	// platform -> device
	VerificationResponse             = byte(0x02)
	HeartbeatResponse                = byte(0x04)
	BillingModelVerificationResponse = byte(0x06)
	BillingModelResponse             = byte(0x0a)
	RealTimeDataRequest              = byte(0x12)
	ChargingRequestConfirmed         = byte(0x32)
	RemoteBootstrapRequest           = byte(0x34)
	RemoteShutdownRequest            = byte(0x36)
	TransactionRecordConfirmed       = byte(0x40)
	AccountBalanceRemoteUpdate       = byte(0x42)
	CardSynchronizationRequest       = byte(0x44)
	CardClearingRequest              = byte(0x46)
	CardQueryingRequest              = byte(0x48)
	SetWorkingParamsRequest          = byte(0x52)
	NtpRequest                       = byte(0x56)
	SetBillingModelRequest           = byte(0x58)
	UpDownFloorLock                  = byte(0x62)
	RemoteRebootRequest              = byte(0x92)
	OtaRequest                       = byte(0x94)

	//Handle protocol from Huaping Power
	DeviceLogin       = byte(0x81)
	RemoteStart       = byte(0x83)
	RemoteStop        = byte(0x84)
	SubmitFinalStatus = byte(0x85)
)
