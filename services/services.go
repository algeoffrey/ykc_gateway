package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"ykc-proxy-server/dtos"
	"ykc-proxy-server/protocols"
	"ykc-proxy-server/utils"

	log "github.com/sirupsen/logrus"
)

func Verification(opt *dtos.Options, buf []byte, hex []string, header *dtos.Header, conn net.Conn) *dtos.VerificationMessage {
	msg := protocols.PackVerificationMessage(buf, hex, header)

	log.WithFields(log.Fields{
		"id":               msg.Id,
		"elc_type":         msg.ElcType,
		"guns":             msg.Guns,
		"protocol_version": msg.ProtocolVersion,
		"software_version": msg.SoftwareVersion,
		"network":          msg.Network,
		"sim":              msg.Sim,
		"operator":         msg.Operator,
	}).Debug("[01] Verification message")
	utils.StoreClient(dtos.ClientInfo{IPAddress: msg.Id}, conn)

	return msg

}

func SendHeartbeatResponse(conn net.Conn, header *dtos.Header) error {
	resp := &bytes.Buffer{}

	// Frame Header
	resp.Write(utils.HexToBytes("5AA5"))

	// Data Length
	resp.Write(utils.HexToBytes("0400"))

	// Command
	resp.Write([]byte{0x82})

	// Reserved Field
	resp.Write([]byte{0x00})

	// Checksum
	checksum := protocols.CalculateChecksum(resp.Bytes()[2:])
	resp.Write([]byte{checksum})

	_, err := conn.Write(resp.Bytes())
	if err != nil {
		log.Errorf("Failed to send Heartbeat Response: %v", err)
		return err
	}
	log.Debug("Sent Heartbeat Response successfully")
	return nil
}

func Hearthbeat(buf []byte, header *dtos.Header, conn net.Conn) *dtos.HeartbeatMessage {
	IPAddress := conn.RemoteAddr().String()
	conn, imei, err := utils.GetClientByIPAddress(IPAddress)
	if err != nil {
		return nil
	}
	fmt.Println(imei)
	fmt.Println(IPAddress)
	msg := protocols.PackHeartbeatMessage(buf, header)
	if msg == nil {
		log.Error("Failed to parse Heartbeat message")
		return nil
	}

	log.WithFields(log.Fields{
		"header":         msg.Header,
		"signalValue":    msg.SignalValue,
		"temperature":    msg.Temperature,
		"totalPortCount": msg.TotalPortCount,
		"portStatus":     msg.PortStatus,
	}).Debug("[82] Heartbeat message")
	return msg
}

func BillingModelVerification(opt *dtos.Options, hex []string, header *dtos.Header, conn net.Conn) *dtos.BillingModelVerificationMessage {
	msg := protocols.PackBillingModelVerificationMessage(hex, header)
	log.WithFields(log.Fields{
		"id":                 msg.Id,
		"billing_model_code": msg.BillingModelCode,
	}).Debug("[05] BillingModelVerification message")

	return msg

}

func BillingModelRequestMessage(opt *dtos.Options, hex []string, header *dtos.Header, conn net.Conn) *dtos.BillingModelRequestMessage {
	msg := protocols.PackBillingModelRequestMessage(hex, header)
	log.WithFields(log.Fields{
		"id": msg.Id,
	}).Debug("[09] BillingModelRequest message")

	return msg

}
func SendBillingModelResponseMessage(req *dtos.BillingModelResponseMessage) error {
	c, _, err := utils.GetClientByIPAddress(req.Id)
	if err != nil {
		return err
	}
	resp := protocols.PackBillingModelResponseMessage(req)
	_, _ = c.Write(resp)
	log.WithFields(log.Fields{
		"id":      req.Id,
		"request": utils.BytesToHex(resp),
	}).Debug("[0a] BillingModelResponse message sent")
	return nil
}

func ResponseToBillingModelVerification(req *dtos.BillingModelVerificationResponseMessage) error {
	c, _, err := utils.GetClientByIPAddress(req.Id)
	if err != nil {
		return err
	}
	resp := protocols.PackBillingModelVerificationResponseMessage(req)
	_, err = c.Write(resp)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"id":       req.Id,
		"response": utils.BytesToHex(resp),
	}).Debug("[06] BillingModelVerificationResponse message sent")
	return nil
}

func ResponseToVerification(req *dtos.VerificationResponseMessage) error {
	c, _, err := utils.GetClientByIPAddress(req.Id)
	if err != nil {
		return err
	}
	resp := protocols.PackVerificationResponseMessage(req)
	_, _ = c.Write(resp)
	log.WithFields(log.Fields{
		"id":       req.Id,
		"response": utils.BytesToHex(resp),
	}).Debug("[02] VerificationResponse message sent")
	return nil
}

func ResponseToHeartbeat(req *dtos.HeartbeatResponseMessage) error {
	c, _, err := utils.GetClientByIPAddress(req.Id)
	if err != nil {
		return err
	}
	resp := protocols.PackHeartbeatResponseMessage(req)
	_, _ = c.Write(resp)
	log.WithFields(log.Fields{
		"id":       req.Id,
		"response": utils.BytesToHex(resp),
	}).Debug("[04] HeartbeatResponse message sent")
	return nil
}

func SendRemoteBootstrapRequest(req *dtos.RemoteBootstrapRequestMessage) error {
	c, _, err := utils.GetClientByIPAddress(req.Id)
	if err != nil {
		return err
	}
	resp := protocols.PackRemoteBootstrapRequestMessage(req)
	_, _ = c.Write(resp)
	log.WithFields(log.Fields{
		"id":      req.Id,
		"request": utils.BytesToHex(resp),
	}).Debug("[34] RemoteBootstrapRequest message sent")
	return nil
}

func RemoteBootstrapResponseRouter(opt *dtos.Options, hex []string, header *dtos.Header) {
	msg := protocols.PackRemoteBootstrapResponseMessage(hex, header)
	log.WithFields(log.Fields{
		"id":                    msg.Id,
		"trade_sequence_number": msg.TradeSeq,
		"gun_id":                msg.GunId,
		"result":                msg.Result,
		"reason":                msg.Reason,
	}).Debug("[33] RemoteBootstrapResponse message")

	//forward
	if opt.MessageForwarder != nil {
		//convert msg to json string bytes
		b, _ := json.Marshal(msg)
		_ = opt.MessageForwarder.Publish("33", b)
	}
}

func OfflineDataReportMessageRouter(opt *dtos.Options, raw []byte, hex []string, header *dtos.Header) {
	msg := protocols.PackOfflineDataReportMessage(hex, raw, header)
	log.WithFields(log.Fields{
		"id":                               msg.Id,
		"trade_sequence_number":            msg.TradeSeq,
		"gun_id":                           msg.GunId,
		"status":                           msg.Status,
		"reset":                            msg.Reset,
		"plugged":                          msg.Plugged,
		"output_voltage":                   msg.Ov,
		"output_current":                   msg.Oc,
		"gun_line_temperature":             msg.LineTemp,
		"gun_line_encoding":                msg.LineCode,
		"battery_pack_highest_temperature": msg.BpTopTemp,
		"accumulated_charging_time":        msg.AccumulatedChargingTime,
		"remaining_time":                   msg.RemainingTime,
		"charging_degrees":                 msg.ChargingDegrees,
		"lossy_charging_degrees":           msg.LossyChargingDegrees,
		"charged_amount":                   msg.ChargedAmount,
		"hardware_failure":                 msg.HardwareFailure,
	}).Debug("[13] OfflineDataReport message")

	//forward
	if opt.MessageForwarder != nil {
		//convert msg to json string bytes
		b, _ := json.Marshal(msg)
		_ = opt.MessageForwarder.Publish("13", b)
	}
}

func RemoteShutdownResponseRouter(opt *dtos.Options, hex []string, header *dtos.Header) {
	msg := protocols.PackRemoteShutdownResponseMessage(hex, header)
	log.WithFields(log.Fields{
		"id":     msg.Id,
		"gun_id": msg.GunId,
		"result": msg.Result,
		"reason": msg.Reason,
	}).Debug("[35] RemoteShutdownResponse message")

	//forward
	if opt.MessageForwarder != nil {
		//convert msg to json string bytes
		b, _ := json.Marshal(msg)
		_ = opt.MessageForwarder.Publish("35", b)
	}
}

func TransactionRecordMessageRouter(opt *dtos.Options, raw []byte, hex []string, header *dtos.Header) {
	msg := protocols.PackTransactionRecordMessage(raw, hex, header)
	msgJson, _ := json.Marshal(msg)
	log.WithFields(log.Fields{
		"msg": string(msgJson),
	}).Debug("[3b] TransactionRecord message")

	if opt.AutoTransactionRecordConfirm {
		m := &dtos.TransactionRecordConfirmedMessage{
			Header: &dtos.Header{
				Seq:       0,
				Encrypted: false,
			},
			Id:       msg.Id,
			TradeSeq: msg.TradeSeq,
			Result:   0,
		}
		_ = SendTransactionRecordConfirmed(m)
		return
	}

	//forward
	if opt.MessageForwarder != nil {
		//convert msg to json string bytes
		b, _ := json.Marshal(msg)
		_ = opt.MessageForwarder.Publish("3b", b)
	}
}

func RemoteRebootResponseMessageRouter(opt *dtos.Options, hex []string, header *dtos.Header) {
	msg := protocols.PackRemoteRebootResponseMessage(hex, header)
	log.WithFields(log.Fields{
		"id":     msg.Id,
		"result": msg.Result,
	}).Debug("[91] RemoteRebootResponse message")

	//forward
	if opt.MessageForwarder != nil {
		//convert msg to json string bytes
		b, _ := json.Marshal(msg)
		_ = opt.MessageForwarder.Publish("91", b)
	}
}

func SetBillingModelResponseMessageRouter(opt *dtos.Options, hex []string, header *dtos.Header) {
	msg := protocols.PackSetBillingModelResponseMessage(hex, header)
	log.WithFields(log.Fields{
		"id":     msg.Id,
		"result": msg.Result,
	}).Debug("[57] SetBillingModelResponse message")

	//forward
	if opt.MessageForwarder != nil {
		//convert msg to json string bytes
		b, _ := json.Marshal(msg)
		_ = opt.MessageForwarder.Publish("57", b)
	}
}

func ChargingFinishedMessageRouter(opt *dtos.Options, hex []string, header *dtos.Header) {
	msg := protocols.PackChargingFinishedMessage(hex, header)
	log.WithFields(log.Fields{
		"id": msg.Id,
	}).Debug("[19] ChargingFinished message")

	//forward
	if opt.MessageForwarder != nil {
		//convert msg to json string bytes
		b, _ := json.Marshal(msg)
		_ = opt.MessageForwarder.Publish("19", b)
	}
}

func DeviceLogin(opt *dtos.Options, buf []byte, header *dtos.Header, conn net.Conn) (*dtos.DeviceLoginMessage, []byte) {
	// Unpack Device Login Message
	msg := protocols.PackDeviceLoginMessage(buf, header)
	utils.StoreClient(dtos.ClientInfo{IPAddress: conn.RemoteAddr().String(), IMEI: msg.IMEI}, conn)
	if msg == nil {
		log.Error("Failed to parse Device Login message due to checksum mismatch or invalid buffer")
		return nil, nil
	}
	// Log the extracted details
	log.WithFields(log.Fields{
		"imei":            msg.IMEI,
		"devicePortCount": msg.DevicePortCount,
		"hardwareVersion": msg.HardwareVersion,
		"softwareVersion": msg.SoftwareVersion,
		"ccid":            msg.CCID,
		"signalValue":     msg.SignalValue,
		"loginReason":     msg.LoginReason,
	}).Debug("[81] Device Login message")

	// Auto response preparation
	heartbeatPeriod := 30 // Default heartbeat interval (10 seconds)
	if heartbeatPeriod < 10 || heartbeatPeriod > 250 {
		heartbeatPeriod = 30 // Enforce valid range (10-250 seconds)
	}

	resp := &dtos.DeviceLoginResponseMessage{
		Header: &dtos.Header{
			Seq:       header.Seq,
			Encrypted: false,
		},
		HeartbeatPeriod: heartbeatPeriod, // Valid interval
		Result:          0x00,            // Login successful
	}

	// // Pack the response message
	data := protocols.PackDeviceLoginResponseMessage(resp)

	log.Debug("Sent Device Login response successfully")
	return msg, data

}

func RemoteStart(buf []byte, header *dtos.Header, conn net.Conn) {
	msg := protocols.PackRemoteStartMessage(buf, header)
	if msg == nil {
		log.Error("Failed to parse Remote Start message")
	}

	log.WithFields(log.Fields{
		"port":        msg.Port,
		"orderNumber": msg.OrderNumber,
		"startMode":   msg.StartMode,
		"startResult": msg.StartResult,
	}).Debug("[83] Remote Start message")

}

func RemoteStop(buf []byte, header *dtos.Header, conn net.Conn) {
	msg := protocols.PackRemoteStopMessage(buf, header)
	if msg == nil {
		log.Error("Failed to parse Remote Stop message")
	}
	log.WithFields(log.Fields{
		"port":        msg.Port,
		"orderNumber": msg.OrderNumber,
	}).Debug("[84] Remote Stop message")

}

func SubmitFinalStatus(opt *dtos.Options, buf []byte, header *dtos.Header, conn net.Conn) []byte {
	msg := protocols.PackSubmitFinalStatusMessage(buf, header)
	if msg == nil {
		log.Error("Failed to parse Submit Final Status message")
		return nil
	}
	log.WithFields(log.Fields{
		"port":             msg.Port,
		"orderNumber":      msg.OrderNumber,
		"chargingTime":     msg.ChargingTime,
		"electricityUsage": msg.ElectricityUsage,
		"usageCost":        msg.UsageCost,
		"stopReason":       msg.StopReason,
		"stopPower":        msg.StopPower,
		"segmentCount":     msg.SegmentCount,
		"segmentDurations": msg.SegmentDurations,
		"segmentPrices":    msg.SegmentPrices,
	}).Debug("[85] Submit Final Status message")

	data := protocols.PackSubmitFinalStatusResponse()
	return data

}

func ChargingPortData(opt *dtos.Options, buf []byte, header *dtos.Header, conn net.Conn) *dtos.ChargingPortDataMessage {

	msg := protocols.PackChargingPortDataMessage(buf, header)
	if msg == nil {
		log.Error("Failed to parse Charging Port Data message")
		return nil
	}

	log.WithFields(log.Fields{
		"header":          msg.Header,
		"reserved":        msg.Reserved,
		"portCount":       msg.PortCount,
		"voltage":         msg.Voltage,
		"temperature":     msg.Temperature,
		"activePort":      msg.ActivePort,
		"currentTier":     msg.CurrentTier,
		"currentRate":     msg.CurrentRate,
		"currentPower":    msg.CurrentPower,
		"usageTime":       msg.UsageTime,
		"usedAmount":      msg.UsedAmount,
		"energyUsed":      msg.EnergyUsed,
		"portTemperature": msg.PortTemperature,
	}).Debug("[88] Charging Port Data message parsed successfully")

	return msg
}

func SendRemoteShutdownRequest(req *dtos.RemoteShutdownRequestMessage) error {
	c, _, err := utils.GetClientByIPAddress(req.Id)
	if err != nil {
		return err
	}
	resp := protocols.PackRemoteShutdownRequestMessage(req)
	_, _ = c.Write(resp)
	log.WithFields(log.Fields{
		"id":      req.Id,
		"request": utils.BytesToHex(resp),
	}).Debug("[36] RemoteShutdownRequest message sent")
	return nil
}

func SendTransactionRecordConfirmed(req *dtos.TransactionRecordConfirmedMessage) error {
	c, _, err := utils.GetClientByIPAddress(req.Id)
	if err != nil {
		return err
	}
	resp := protocols.PackTransactionRecordConfirmedMessage(req)
	_, _ = c.Write(resp)
	log.WithFields(log.Fields{
		"id":      req.Id,
		"request": utils.BytesToHex(resp),
	}).Debug("[40] TransactionRecordConfirmed message sent")
	return nil
}

func SendRemoteRebootRequest(req *dtos.RemoteRebootRequestMessage) error {
	c, _, err := utils.GetClientByIPAddress(req.Id)
	if err != nil {
		return err
	}
	resp := protocols.PackRemoteRebootRequestMessage(req)
	_, _ = c.Write(resp)
	log.WithFields(log.Fields{
		"id":      req.Id,
		"request": utils.BytesToHex(resp),
	}).Debug("[92] RemoteRebootRequest message sent")
	return nil
}

func SendSetBillingModelRequestMessage(req *dtos.SetBillingModelRequestMessage) error {
	c, _, err := utils.GetClientByIPAddress(req.Id)
	if err != nil {
		return err
	}
	resp := protocols.PackSetBillingModelRequestMessage(req)
	_, _ = c.Write(resp)
	log.WithFields(log.Fields{
		"id":      req.Id,
		"request": utils.BytesToHex(resp),
	}).Debug("[58] SetBillingModelRequest message sent")
	return nil
}
