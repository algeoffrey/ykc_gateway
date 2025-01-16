package protocols

import (
	"bytes"
	"encoding/binary"
	hex2 "encoding/hex"
	"fmt"
	"strconv"
	"ykc-proxy-server/dtos"
	"ykc-proxy-server/utils"

	log "github.com/sirupsen/logrus"
)

func CalculateChecksum(data []byte) byte {
	var checksum byte
	for _, b := range data {
		checksum += b
	}
	return checksum
}

func PackVerificationMessage(buf []byte, hex []string, header *dtos.Header) *dtos.VerificationMessage {
	//Id
	id := ""
	for _, v := range hex[6:13] {
		id += v
	}

	//type
	elcType := int(buf[13])

	//gun number
	guns := int(buf[14])

	//protocol version
	protocolVersion := int(buf[15]) / 10

	//software version
	softwareVersionBytes, _ := hex2.DecodeString(utils.MakeHexStringFromHexArray(hex[16:24]))
	softwareVersion := string(softwareVersionBytes)

	//network type
	network := int(buf[25])

	//sim
	var sim string
	for _, v := range hex[26:36] {
		sim += v
	}

	//operator
	operator := int(buf[36])

	msg := &dtos.VerificationMessage{
		Header:          header,
		Id:              id,
		ElcType:         elcType,
		Guns:            guns,
		ProtocolVersion: protocolVersion,
		SoftwareVersion: softwareVersion,
		Network:         network,
		Sim:             sim,
		Operator:        operator,
	}
	return msg
}

func PackVerificationResponseMessage(msg *dtos.VerificationResponseMessage) []byte {
	var resp bytes.Buffer
	resp.Write([]byte{StartFlag, 0x0c})
	seqStr := fmt.Sprintf("%x", msg.Header.Seq)
	seq := utils.ConvertIntSeqToReversedHexArr(seqStr)
	resp.Write(utils.HexToBytes(utils.MakeHexStringFromHexArray(seq)))
	encrypted := byte(0x00)
	if msg.Header.Encrypted {
		encrypted = byte(0x01)
	}
	resp.Write([]byte{encrypted})
	resp.Write([]byte{VerificationResponse})
	resp.Write(utils.HexToBytes(msg.Id))
	result := byte(0x01)
	if msg.Result {
		result = byte(0x00)
	}
	resp.Write([]byte{result})
	resp.Write(utils.ModbusCRC(resp.Bytes()[2:]))
	return resp.Bytes()
}

func PackBillingModelVerificationMessage(hex []string, header *dtos.Header) *dtos.BillingModelVerificationMessage {
	//Id
	id := ""
	for _, v := range hex[6:13] {
		id += v
	}

	//billing model code
	bmcode := hex[13] + hex[14]

	msg := &dtos.BillingModelVerificationMessage{
		Header:           header,
		Id:               id,
		BillingModelCode: bmcode,
	}
	return msg
}

func PackBillingModelRequestMessage(hex []string, header *dtos.Header) *dtos.BillingModelRequestMessage {
	//Id
	id := ""
	for _, v := range hex[6:13] {
		id += v
	}

	msg := &dtos.BillingModelRequestMessage{
		Header: header,
		Id:     id,
	}
	return msg
}

func PackBillingModelResponseMessage(msg *dtos.BillingModelResponseMessage) []byte {
	var resp bytes.Buffer
	resp.Write([]byte{0x68, 0x5e})
	seqStr := fmt.Sprintf("%x", utils.GenerateSeq())
	seq := utils.ConvertIntSeqToReversedHexArr(seqStr)
	resp.Write(utils.HexToBytes(utils.MakeHexStringFromHexArray(seq)))
	if msg.Header.Encrypted {
		resp.WriteByte(0x01)
	} else {
		resp.WriteByte(0x00)
	}
	resp.Write([]byte{0x0a})
	resp.Write(utils.HexToBytes(msg.Id))
	resp.Write(utils.HexToBytes(msg.BillingModelCode))
	resp.Write(utils.IntToBIN(msg.SharpUnitPrice, 4))
	resp.Write(utils.IntToBIN(msg.SharpServiceFee, 4))
	resp.Write(utils.IntToBIN(msg.PeakUnitPrice, 4))
	resp.Write(utils.IntToBIN(msg.PeakServiceFee, 4))
	resp.Write(utils.IntToBIN(msg.FlatUnitPrice, 4))
	resp.Write(utils.IntToBIN(msg.FlatServiceFee, 4))
	resp.Write(utils.IntToBIN(msg.ValleyUnitPrice, 4))
	resp.Write(utils.IntToBIN(msg.ValleyServiceFee, 4))
	resp.Write([]byte(strconv.Itoa(msg.AccrualRatio)))

	for _, v := range msg.RateList {
		resp.Write(utils.IntToBIN(v, 1))
	}

	resp.Write(utils.ModbusCRC(resp.Bytes()[2:]))
	return resp.Bytes()
}

func PackBillingModelVerificationResponseMessage(msg *dtos.BillingModelVerificationResponseMessage) []byte {
	var resp bytes.Buffer
	resp.Write([]byte{StartFlag, 0x0e})

	seqStr := fmt.Sprintf("%x", msg.Header.Seq)
	seq := utils.ConvertIntSeqToReversedHexArr(seqStr)
	resp.Write(utils.HexToBytes(utils.MakeHexStringFromHexArray(seq)))

	encrypted := byte(0x00)
	if msg.Header.Encrypted {
		encrypted = byte(0x01)
	}
	resp.Write([]byte{encrypted})
	resp.Write([]byte{BillingModelVerificationResponse})
	resp.Write(utils.HexToBytes(msg.Id))
	resp.Write(utils.HexToBytes(msg.BillingModelCode))

	result := byte(0x01)
	if msg.Result {
		result = byte(0x00)
	}
	resp.Write([]byte{result})
	resp.Write(utils.ModbusCRC(resp.Bytes()[2:]))
	return resp.Bytes()
}

func PackHeartbeatMessage(buf []byte, header *dtos.Header) *dtos.HeartbeatMessage {
	payload := buf[6:] // Skip the header (first 6 bytes)

	// Parse fields
	signalValue := int(payload[0])
	temperature := int(payload[1])
	totalPortCount := int(payload[2])
	portStatus := make([]int, totalPortCount)
	for i := 0; i < totalPortCount; i++ {
		portStatus[i] = int(payload[3+i])
	}

	log.Debugf("Parsed Signal Value: %d", signalValue)
	log.Debugf("Parsed Temperature: %d", temperature)
	log.Debugf("Parsed Total Port Count: %d", totalPortCount)
	log.Debugf("Parsed Port Status: %v", portStatus)

	return &dtos.HeartbeatMessage{
		Header:         header,
		SignalValue:    signalValue,
		Temperature:    temperature,
		TotalPortCount: totalPortCount,
		PortStatus:     portStatus,
	}
}

func PackHeartbeatResponseMessage(msg *dtos.HeartbeatResponseMessage) []byte {
	var resp bytes.Buffer
	resp.Write(utils.HexToBytes("680d"))

	seqStr := fmt.Sprintf("%x", msg.Header.Seq)
	seq := utils.ConvertIntSeqToReversedHexArr(seqStr)
	resp.Write(utils.HexToBytes(utils.MakeHexStringFromHexArray(seq)))

	encrypted := "00"
	if msg.Header.Encrypted {
		encrypted = "01"
	}

	resp.Write(utils.HexToBytes(encrypted))
	resp.Write(utils.HexToBytes("04"))
	resp.Write(utils.HexToBytes(msg.Id))
	resp.Write(utils.HexToBytes(msg.Gun))
	resp.Write(utils.HexToBytes("00"))
	resp.Write(utils.ModbusCRC(resp.Bytes()[2:]))

	return resp.Bytes()
}

func PackRemoteBootstrapRequestMessage(msg *dtos.RemoteBootstrapRequestMessage) []byte {
	var resp bytes.Buffer
	resp.Write(utils.HexToBytes("6830"))

	seqStr := fmt.Sprintf("%x", utils.GenerateSeq())
	seq := utils.ConvertIntSeqToReversedHexArr(seqStr)
	resp.Write(utils.HexToBytes(utils.MakeHexStringFromHexArray(seq)))

	encrypted := byte(0x00)
	if msg.Header.Encrypted {
		encrypted = byte(0x01)
	}
	resp.Write([]byte{encrypted})
	resp.Write(utils.HexToBytes("34"))
	resp.Write(utils.HexToBytes(msg.TradeSeq))
	resp.Write(utils.HexToBytes(msg.Id))
	resp.Write(utils.HexToBytes(msg.GunId))
	resp.Write(utils.PadArrayWithZeros(utils.HexToBytes(msg.LogicCard), 8))
	resp.Write(utils.PadArrayWithZeros(utils.HexToBytes(msg.PhysicalCard), 8))

	balance := utils.HexToBytes(fmt.Sprintf("%x", msg.Balance))
	balance = utils.PadArrayWithZeros(balance, 4)
	resp.Write(balance)

	resp.Write(utils.ModbusCRC(resp.Bytes()[2:]))
	return resp.Bytes()
}

func PackRemoteBootstrapResponseMessage(hex []string, header *dtos.Header) *dtos.RemoteBootstrapResponseMessage {
	//trade sequence number
	tradeSeq := ""
	for _, v := range hex[6:22] {
		tradeSeq += v
	}

	//id
	id := ""
	for _, v := range hex[22:29] {
		id += v
	}

	//gun id
	gunId := hex[29]

	//result
	result := false
	if hex[30] == "01" {
		result = true
	}

	//fail reason
	reason, _ := strconv.ParseInt(hex[31], 16, 64)

	msg := &dtos.RemoteBootstrapResponseMessage{
		Header:   header,
		TradeSeq: tradeSeq,
		Id:       id,
		GunId:    gunId,
		Result:   result,
		Reason:   int(reason),
	}
	return msg
}

func PackOfflineDataReportMessage(hex []string, raw []byte, header *dtos.Header) *dtos.OfflineDataReportMessage {
	//trade sequence number
	tradeSeq := ""
	for _, v := range hex[6:22] {
		tradeSeq += v
	}

	//id
	id := ""
	for _, v := range hex[22:29] {
		id += v
	}

	//gun id
	gunId := hex[29]

	//status
	status := utils.BINToInt([]byte{raw[30]})

	//reset
	reset := utils.BINToInt([]byte{raw[31]})

	//plugged
	plugged := utils.BINToInt([]byte{raw[32]})

	//ov
	ov := utils.BINToInt(raw[33:35])

	//oc
	oc := utils.BINToInt(raw[35:37])

	//lineTemp
	lineTemp := utils.BINToInt([]byte{raw[37]})

	//lineCode
	lineCode := utils.BINToInt(raw[38:46])

	//soc
	soc := utils.BINToInt([]byte{raw[46]})

	//bpTopTemp
	bpTopTemp := utils.BINToInt([]byte{raw[47]})

	//accumulatedChargingTime
	accumulatedChargingTime := utils.BINToInt(raw[48:50])

	//remainingTime
	remainingTime := utils.BINToInt(raw[50:52])

	//chargingDegrees
	chargingDegrees := utils.BINToInt(raw[52:56])

	//lossyChargingDegrees
	lossyChargingDegrees := utils.BINToInt(raw[56:60])

	//chargedAmount
	chargedAmount := utils.BINToInt(raw[60:64])

	//hardwareFailure
	hardwareFailure := utils.BINToInt(raw[64:66])

	msg := &dtos.OfflineDataReportMessage{
		Header:                  header,
		TradeSeq:                tradeSeq,
		Id:                      id,
		GunId:                   gunId,
		Status:                  status,
		Reset:                   reset,
		Plugged:                 plugged,
		Ov:                      ov,
		Oc:                      oc,
		LineTemp:                lineTemp,
		LineCode:                strconv.Itoa(lineCode),
		Soc:                     soc,
		BpTopTemp:               bpTopTemp,
		AccumulatedChargingTime: accumulatedChargingTime,
		RemainingTime:           remainingTime,
		ChargingDegrees:         chargingDegrees,
		LossyChargingDegrees:    lossyChargingDegrees,
		ChargedAmount:           chargedAmount,
		HardwareFailure:         hardwareFailure,
	}

	return msg
}

func PackRemoteShutdownResponseMessage(hex []string, header *dtos.Header) *dtos.RemoteShutdownResponseMessage {
	//id
	id := ""
	for _, v := range hex[6:13] {
		id += v
	}

	//gun id
	gunId := hex[13]

	//result
	result := false
	if hex[14] == "01" {
		result = true
	}

	//fail reason
	reason, _ := strconv.ParseInt(hex[15], 16, 64)

	msg := &dtos.RemoteShutdownResponseMessage{
		Header: header,
		Id:     id,
		GunId:  gunId,
		Result: result,
		Reason: int(reason),
	}
	return msg
}

func PackRemoteShutdownRequestMessage(msg *dtos.RemoteShutdownRequestMessage) []byte {
	var resp bytes.Buffer
	resp.Write(utils.HexToBytes("680c"))
	seqStr := fmt.Sprintf("%x", utils.GenerateSeq())
	seq := utils.ConvertIntSeqToReversedHexArr(seqStr)
	resp.Write(utils.HexToBytes(utils.MakeHexStringFromHexArray(seq)))
	if msg.Header.Encrypted {
		resp.WriteByte(0x01)
	} else {
		resp.WriteByte(0x00)
	}
	resp.Write(utils.HexToBytes("36"))
	resp.Write(utils.HexToBytes(msg.Id))
	resp.Write(utils.HexToBytes(msg.GunId))
	resp.Write(utils.ModbusCRC(resp.Bytes()[2:]))
	return resp.Bytes()
}

func PackTransactionRecordMessage(raw []byte, hex []string, header *dtos.Header) *dtos.TransactionRecordMessage {
	//trade sequence number
	tradeSeq := ""
	for _, v := range hex[6:22] {
		tradeSeq += v
	}

	//id
	id := ""
	for _, v := range hex[22:29] {
		id += v
	}

	//gun id
	gunId := hex[29]

	//start time
	startAt := utils.Cp56time2aToUnixMilliseconds(raw[30:37])

	//end time
	endAt := utils.Cp56time2aToUnixMilliseconds(raw[37:44])

	//sharp unit price
	sharpUnitPrice := utils.BINToInt(raw[44:48])

	//sharp electric charge
	sharpElectricCharge := utils.BINToInt(raw[48:52])

	//lossy sharp electric charge
	lossySharpElectricCharge := utils.BINToInt(raw[52:56])

	//sharp price
	sharpPrice := utils.BINToInt(raw[56:60])

	//peak unit price
	peakUnitPrice := utils.BINToInt(raw[60:64])

	//peak electric charge
	peakElectricCharge := utils.BINToInt(raw[64:68])

	//lossy peak electric charge
	lossyPeakElectricCharge := utils.BINToInt(raw[68:72])

	//peak price
	peakPrice := utils.BINToInt(raw[72:76])

	//flat unit price
	flatUnitPrice := utils.BINToInt(raw[76:80])

	//flat electric charge
	flatElectricCharge := utils.BINToInt(raw[80:84])

	//lossy flat electric charge
	lossyFlatElectricCharge := utils.BINToInt(raw[84:88])

	//flat price
	flatPrice := utils.BINToInt(raw[88:92])

	//valley unit price
	valleyUnitPrice := utils.BINToInt(raw[92:96])

	//valley electric charge
	valleyElectricCharge := utils.BINToInt(raw[96:100])

	//lossy valley electric charge
	lossyValleyElectricCharge := utils.BINToInt(raw[100:104])

	//valley price
	valleyPrice := utils.BINToInt(raw[104:108])

	//initial meter reading
	initialMeterReading := utils.BINToInt(raw[108:113])

	//final meter reading
	finalMeterReading := utils.BINToInt(raw[113:118])

	//total electric charge
	totalElectricCharge := utils.BINToInt(raw[118:122])

	//lossy total electric charge
	lossyTotalElectricCharge := utils.BINToInt(raw[122:126])

	//consumption amount
	consumptionAmount := utils.BINToInt(raw[126:130])

	//vin
	vin := utils.MakeHexStringFromHexArray(hex[130:147])

	//start type
	startType := utils.BINToInt([]byte{raw[147]})

	//transaction date time
	transactionDateTime := utils.Cp56time2aToUnixMilliseconds(raw[148:155])

	//stop reason
	stopReason := utils.BINToInt([]byte{raw[155]})

	//physical card number
	physicalCardNumber := utils.MakeHexStringFromHexArray(hex[156:164])

	//fill all fields
	msg := &dtos.TransactionRecordMessage{
		Header:                    header,
		TradeSeq:                  tradeSeq,
		Id:                        id,
		GunId:                     gunId,
		StartAt:                   startAt,
		EndAt:                     endAt,
		SharpUnitPrice:            int64(sharpUnitPrice),
		SharpElectricCharge:       int64(sharpElectricCharge),
		LossySharpElectricCharge:  int64(lossySharpElectricCharge),
		SharpPrice:                int64(sharpPrice),
		PeakUnitPrice:             int64(peakUnitPrice),
		PeakElectricCharge:        int64(peakElectricCharge),
		LossyPeakElectricCharge:   int64(lossyPeakElectricCharge),
		PeakPrice:                 int64(peakPrice),
		FlatUnitPrice:             int64(flatUnitPrice),
		FlatElectricCharge:        int64(flatElectricCharge),
		LossyFlatElectricCharge:   int64(lossyFlatElectricCharge),
		FlatPrice:                 int64(flatPrice),
		ValleyUnitPrice:           int64(valleyUnitPrice),
		ValleyElectricCharge:      int64(valleyElectricCharge),
		LossyValleyElectricCharge: int64(lossyValleyElectricCharge),
		ValleyPrice:               int64(valleyPrice),
		InitialMeterReading:       int64(initialMeterReading),
		FinalMeterReading:         int64(finalMeterReading),
		TotalElectricCharge:       int64(totalElectricCharge),
		LossyTotalElectricCharge:  int64(lossyTotalElectricCharge),
		ConsumptionAmount:         int64(consumptionAmount),
		Vin:                       vin,
		StartType:                 startType,
		TransactionDateTime:       transactionDateTime,
		StopReason:                stopReason,
		PhysicalCardNumber:        physicalCardNumber,
	}
	return msg
}

func PackTransactionRecordConfirmedMessage(msg *dtos.TransactionRecordConfirmedMessage) []byte {
	var resp bytes.Buffer
	resp.Write([]byte{0x68, 0x15})
	seqStr := fmt.Sprintf("%x", utils.GenerateSeq())
	seq := utils.ConvertIntSeqToReversedHexArr(seqStr)
	resp.Write(utils.HexToBytes(utils.MakeHexStringFromHexArray(seq)))
	if msg.Header.Encrypted {
		resp.WriteByte(0x01)
	} else {
		resp.WriteByte(0x00)
	}
	resp.Write([]byte{0x40})
	resp.Write(utils.HexToBytes(msg.TradeSeq))
	resp.Write([]byte(strconv.Itoa(msg.Result)))
	resp.Write(utils.ModbusCRC(resp.Bytes()[2:]))
	return resp.Bytes()
}

func PackRemoteRebootResponseMessage(hex []string, header *dtos.Header) *dtos.RemoteRebootResponseMessage {
	//id
	id := ""
	for _, v := range hex[6:13] {
		id += v
	}

	//result 0-fail 1-success
	result := 1
	if hex[13] == "00" {
		result = 0
	}

	msg := &dtos.RemoteRebootResponseMessage{
		Header: header,
		Id:     id,
		Result: result,
	}
	return msg
}

func PackRemoteRebootRequestMessage(msg *dtos.RemoteRebootRequestMessage) []byte {
	var resp bytes.Buffer
	resp.Write([]byte{0x68, 0x0c})
	seqStr := fmt.Sprintf("%x", utils.GenerateSeq())
	seq := utils.ConvertIntSeqToReversedHexArr(seqStr)
	resp.Write(utils.HexToBytes(utils.MakeHexStringFromHexArray(seq)))
	if msg.Header.Encrypted {
		resp.WriteByte(0x01)
	} else {
		resp.WriteByte(0x00)
	}
	resp.Write([]byte{0x92})
	resp.Write(utils.HexToBytes(msg.Id))
	resp.Write([]byte(strconv.Itoa(msg.Control)))
	resp.Write(utils.ModbusCRC(resp.Bytes()[2:]))
	return resp.Bytes()
}

func PackSetBillingModelRequestMessage(msg *dtos.SetBillingModelRequestMessage) []byte {
	var resp bytes.Buffer
	resp.Write([]byte{0x68, 0x5e})
	seqStr := fmt.Sprintf("%x", utils.GenerateSeq())
	seq := utils.ConvertIntSeqToReversedHexArr(seqStr)
	resp.Write(utils.HexToBytes(utils.MakeHexStringFromHexArray(seq)))
	if msg.Header.Encrypted {
		resp.WriteByte(0x01)
	} else {
		resp.WriteByte(0x00)
	}
	resp.Write([]byte{0x0a})
	resp.Write(utils.HexToBytes(msg.Id))
	resp.Write(utils.HexToBytes(msg.BillingModelCode))
	resp.Write(utils.IntToBIN(msg.SharpUnitPrice, 4))
	resp.Write(utils.IntToBIN(msg.SharpServiceFee, 4))
	resp.Write(utils.IntToBIN(msg.PeakUnitPrice, 4))
	resp.Write(utils.IntToBIN(msg.PeakServiceFee, 4))
	resp.Write(utils.IntToBIN(msg.FlatUnitPrice, 4))
	resp.Write(utils.IntToBIN(msg.FlatServiceFee, 4))
	resp.Write(utils.IntToBIN(msg.ValleyUnitPrice, 4))
	resp.Write(utils.IntToBIN(msg.ValleyServiceFee, 4))
	resp.Write([]byte(strconv.Itoa(msg.AccrualRatio)))

	for _, v := range msg.RateList {
		resp.Write(utils.IntToBIN(v, 1))
	}

	resp.Write(utils.ModbusCRC(resp.Bytes()[2:]))
	return resp.Bytes()
}

func PackSetBillingModelResponseMessage(hex []string, header *dtos.Header) *dtos.SetBillingModelResponseMessage {
	//id
	id := ""
	for _, v := range hex[6:13] {
		id += v
	}

	//result 0-fail 1-success
	result := 1
	if hex[13] == "00" {
		result = 0
	}

	msg := &dtos.SetBillingModelResponseMessage{
		Header: header,
		Id:     id,
		Result: result,
	}
	return msg
}

func PackChargingFinishedMessage(hex []string, header *dtos.Header) *dtos.ChargingFinishedMessage {
	//trade sequence number
	tradeSeq := ""
	for _, v := range hex[6:22] {
		tradeSeq += v
	}

	//id
	id := ""
	for _, v := range hex[22:29] {
		id += v
	}

	//gun id
	gunId := hex[29]

	//soc
	soc, _ := strconv.ParseInt(hex[30], 16, 64)

	bmsBatteryPackLowestVoltage, _ := strconv.ParseInt(utils.MakeHexStringFromHexArray(hex[31:33]), 16, 64)
	bmsBatteryPackHighestVoltage, _ := strconv.ParseInt(utils.MakeHexStringFromHexArray(hex[33:35]), 16, 64)
	bmsBatteryPackLowestTemperature, _ := strconv.ParseInt(hex[35], 16, 64)
	bmsBatteryPackHighestTemperature, _ := strconv.ParseInt(hex[36], 16, 64)
	cumulativeChargingDuration, _ := strconv.ParseInt(utils.MakeHexStringFromHexArray(hex[37:39]), 16, 64)
	outputPower, _ := strconv.ParseInt(utils.MakeHexStringFromHexArray(hex[39:41]), 16, 64)
	chargingUnitId, _ := strconv.ParseInt(utils.MakeHexStringFromHexArray(hex[41:45]), 16, 64)

	msg := &dtos.ChargingFinishedMessage{
		Header:                           header,
		TradeSeq:                         tradeSeq,
		Id:                               id,
		GunId:                            gunId,
		BmsSoc:                           int(soc),
		BmsBatteryPackLowestVoltage:      int(bmsBatteryPackLowestVoltage),
		BmsBatteryPackHighestVoltage:     int(bmsBatteryPackHighestVoltage),
		BmsBatteryPackLowestTemperature:  int(bmsBatteryPackLowestTemperature),
		BmsBatteryPackHighestTemperature: int(bmsBatteryPackHighestTemperature),
		CumulativeChargingDuration:       int(cumulativeChargingDuration),
		OutputPower:                      int(outputPower),
		ChargingUnitId:                   int(chargingUnitId),
	}
	return msg
}

func PackDeviceLoginMessage(buf []byte, header *dtos.Header) *dtos.DeviceLoginMessage {
	if len(buf) < 6 {
		log.Error("Message too short to process")
		return nil
	}

	// Start reading from the payload (after the 5-byte header)
	payload := buf[6:] // Skip header bytes

	// Parse fields
	imeiHex := payload[:15]
	imei := utils.HexToASCII(utils.MakeHexStringFromHexArray(utils.BytesToHex(imeiHex)))
	log.Debugf("Parsed IMEI: %s", imei)

	devicePortCount := int(payload[15])
	log.Debugf("Parsed Device Port Count: %d", devicePortCount)

	hardwareHex := payload[16:32]
	hardwareVersion := utils.HexToASCII(utils.MakeHexStringFromHexArray(utils.BytesToHex(hardwareHex)))
	log.Debugf("Parsed Hardware Version: %s", hardwareVersion)

	softwareHex := payload[32:48]
	softwareVersion := utils.HexToASCII(utils.MakeHexStringFromHexArray(utils.BytesToHex(softwareHex)))
	log.Debugf("Parsed Software Version: %s", softwareVersion)

	ccidHex := payload[48:68]
	ccid := utils.HexToASCII(utils.MakeHexStringFromHexArray(utils.BytesToHex(ccidHex)))
	log.Debugf("Parsed CCID: %s", ccid)

	signalValue := int(payload[68])
	log.Debugf("Parsed Signal Value: %d", signalValue)

	loginReason := int(payload[69])
	log.Debugf("Parsed Login Reason: %d", loginReason)

	return &dtos.DeviceLoginMessage{
		Header:          header,
		IMEI:            imei,
		DevicePortCount: devicePortCount,
		HardwareVersion: hardwareVersion,
		SoftwareVersion: softwareVersion,
		CCID:            ccid,
		SignalValue:     signalValue,
		LoginReason:     loginReason,
	}
}

func PackDeviceLoginResponseMessage(msg *dtos.DeviceLoginResponseMessage) []byte {
	var resp bytes.Buffer

	// Frame Header (5AA5)
	resp.Write(utils.HexToBytes("5AA5"))

	// Data Length (12 bytes)
	resp.Write(utils.HexToBytes("0C00"))

	// Command (81)
	resp.Write([]byte{0x81, 0x00})

	// Time (7 bytes, BCD format)
	resp.Write([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	// Heartbeat Interval (1 byte)
	resp.Write([]byte{byte(msg.HeartbeatPeriod)})

	// Login Result (1 byte)
	resp.Write([]byte{byte(msg.Result)})

	// Checksum Result (1 byte)
	resp.Write([]byte{0x7D})

	return resp.Bytes()
}

func PackRemoteStartMessage(buf []byte, header *dtos.Header) *dtos.RemoteStartMessage {
	payload := buf[6:] // Skip header and imei bytes

	return &dtos.RemoteStartMessage{
		Header:      header,
		Port:        int(payload[0]),
		OrderNumber: hex2.EncodeToString(payload[1:5]),
		StartMode:   int(payload[5]),
		StartResult: int(payload[6]),
	}
}

func PackRemoteStopMessage(buf []byte, header *dtos.Header) *dtos.RemoteStopMessage {
	payload := buf[6:] // Skip the header and imei

	return &dtos.RemoteStopMessage{
		Header:      header,
		Port:        int(payload[0]),
		OrderNumber: hex2.EncodeToString(payload[1:5]),
	}
}

func PackSubmitFinalStatusMessage(buf []byte, header *dtos.Header) *dtos.SubmitFinalStatusMessage {

	if len(buf) < 24 { // Make sure the buffer is large enough to hold the expected data
		log.Error("Message too short to process CMD0x85")
		return nil
	}

	payload := buf[6:]

	// Parsing fields and logging the values
	log.Debugf("Parsed Port: %d", payload[0])
	port := payload[0]

	log.Debugf("Parsed Order Number: %d", hex2.EncodeToString(payload[1:5]))
	orderNumber := hex2.EncodeToString(payload[1:5])

	log.Debugf("Parsed Charging Time: %d", binary.BigEndian.Uint32(payload[5:9]))
	chargingTime := binary.LittleEndian.Uint32(payload[5:9])

	log.Debugf("Parsed Electricity Usage: %d", binary.BigEndian.Uint32(payload[9:13]))
	electricityUsage := binary.LittleEndian.Uint32(payload[9:13])

	log.Debugf("Parsed Usage Cost: %d", binary.BigEndian.Uint32(payload[13:17]))
	usageCost := binary.LittleEndian.Uint32(payload[13:17])

	stopReason := payload[17]
	log.Debugf("Parsed Stop Reason: %d", stopReason)

	log.Debugf("Parsed Stop Power: %d", binary.BigEndian.Uint16(payload[18:20]))
	stopPower := binary.LittleEndian.Uint16(payload[18:20])

	log.Debugf("Parsed Card ID: %d", binary.BigEndian.Uint32(payload[20:24]))
	cardID := binary.LittleEndian.Uint32(payload[20:24])

	segmentCount := payload[24]
	log.Debugf("Parsed Segment Count: %d", segmentCount)

	// Parsing segment durations and prices using the segment count
	segmentDurations := parseSegments(payload[25:], int(segmentCount))
	log.Debugf("Parsed Segment Durations: %v", segmentDurations)

	segmentPrices := parseSegments(payload[25+int(segmentCount)*2:], int(segmentCount))
	log.Debugf("Parsed Segment Prices: %v", segmentPrices)

	// Return the parsed message
	return &dtos.SubmitFinalStatusMessage{
		Header:           header,
		Port:             port,
		OrderNumber:      orderNumber,
		ChargingTime:     chargingTime,
		ElectricityUsage: electricityUsage,
		UsageCost:        usageCost,
		StopReason:       stopReason,
		StopPower:        stopPower,
		CardID:           cardID,
		SegmentCount:     segmentCount,
		SegmentDurations: segmentDurations,
		SegmentPrices:    segmentPrices,
	}

}

func parseSegments(data []byte, count int) []uint16 {
	segments := make([]uint16, count)
	for i := 0; i < count; i++ {
		segments[i] = binary.BigEndian.Uint16(data[i*2 : i*2+2])
	}
	return segments
}

func PackSubmitFinalStatusResponse(IMEI []byte, hexPort []byte, hexOrderNumber []byte) []byte {

	resp := &bytes.Buffer{}
	resp.Write([]byte{0x5A, 0xA5})

	resp.Write([]byte{0x07, 0x00})
	resp.Write([]byte{0x85, 0x00})
	// resp.Write(IMEI)
	resp.Write(hexPort)
	// resp.Write([]byte{
	// 	0x00, 0x12, 0x34, 0x56, // order number (00123456)
	// })
	resp.Write(hexOrderNumber)
	resp.Write([]byte{0xEA})
	// Print accumulated buffer in hex format
	log.Debugf("Submit Final Status Response buffer (hex): %X", resp.Bytes())

	return resp.Bytes()
}

func ParseStartChargingRequest(IMEI string, hexPort []byte, orderNumberHex []byte) []byte {
	var resp bytes.Buffer

	// Frame Header (5AA5)
	resp.Write([]byte{0x5A, 0xA5})

	// Data Length (12 bytes)
	resp.Write([]byte{0x16, 0x00})

	// Command (83)
	resp.Write([]byte{RemoteStart, 0x00})

	//IMEI
	// resp.Write(imei)

	//PORT
	resp.Write(hexPort)
	resp.Write(orderNumberHex)

	resp.Write([]byte{
		0x01,                   // payment
		0x01, 0x01, 0x01, 0x01, // card number
		0x01,                   // chargin mode
		0x3C, 0x00, 0x00, 0x00, // time based (60 seconds)
		0x01, 0x01, 0x01, 0x01, // available amount
		0xED, // checksum
	})

	return resp.Bytes()
}

func ParseStopChargingRequest(IMEI string, orderNumberHex []byte, hexPort []byte) []byte {
	var resp bytes.Buffer

	// imei := utils.ASCIIToHex(IMEI)
	// Frame Header (5AA5)
	resp.Write([]byte{0x5A, 0xA5})

	// Data Length (12 bytes)
	resp.Write([]byte{0x12, 0x00})

	// Command (84)
	resp.Write([]byte{RemoteStop, 0x00})

	//IMEI
	// resp.Write(imei)

	//PORT
	resp.Write(hexPort)

	resp.Write(orderNumberHex)

	resp.Write([]byte{
		0x8F, // checksum
	})

	return resp.Bytes()
}

func PackChargingPortDataMessage(buf []byte, header *dtos.Header) *dtos.ChargingPortDataMessage {
    if len(buf) < 6 {
        log.Error("Message too short to process CMD088")
        return nil
    }

    payload := buf[6:]
    log.Debugf("Payload: %x", payload)

    portCount := payload[0]
    log.Debugf("Parsed Port Count: %d", portCount)

    voltage := binary.LittleEndian.Uint16(payload[1:3]) // Voltage in 0.1V
    log.Debugf("Parsed Voltage: %.1fV", float64(voltage)*0.1)

    temperature := payload[3]
    log.Debugf("Parsed Temperature: %d°C", temperature)

    ports := []dtos.PortData{}
    offset := 4
    for i := 0; i < int(portCount); i++ {
        if len(payload[offset:]) < 17 {
            log.Errorf("Insufficient data for port %d", i+1)
            break
        }

        log.Debugf("Parsing data for port %d starting at offset %d", i+1, offset)

        portID := payload[offset]
        currentTier := payload[offset+1]
        currentRate := binary.LittleEndian.Uint16(payload[offset+2 : offset+4]) // Current rate in 0.01 Yuan
        currentPower := binary.LittleEndian.Uint16(payload[offset+4 : offset+6]) // Power in Watts
        usageTime := binary.LittleEndian.Uint32(payload[offset+6 : offset+10]) // Time in seconds
        usedAmount := binary.LittleEndian.Uint16(payload[offset+10 : offset+12]) // Used amount in 0.01 Yuan
        energyUsed := binary.LittleEndian.Uint32(payload[offset+12 : offset+16]) // Energy used in 0.01 kWh
        portTemperature := payload[offset+16]

        log.Debugf("Parsed Port Data - ID: %d, Tier: %d, Rate: %.2f Yuan, Power: %dW, Time: %d sec, Amount: %.2f Yuan, Energy: %.2f kWh, Temp: %d°C",
            portID, currentTier, float64(currentRate)*0.01, currentPower, usageTime, float64(usedAmount)*0.01, float64(energyUsed)*0.01, portTemperature)

        ports = append(ports, dtos.PortData{
            PortID:          portID,
            CurrentTier:     currentTier,
            CurrentRate:     currentRate,
            CurrentPower:    currentPower,
            UsageTime:       usageTime,
            UsedAmount:      usedAmount,
            EnergyUsed:      energyUsed,
            PortTemperature: portTemperature,
        })

        offset += 17 // Adjust offset to 17 bytes per port
    }

    log.Debugf("All ports parsed successfully. Total ports: %d", len(ports))

    return &dtos.ChargingPortDataMessage{
        Header:      header,
        PortCount:   portCount,
        Voltage:     voltage,
        Temperature: temperature,
        Ports:       ports,
    }
}




