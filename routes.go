package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func PrintHexAndByte(data []byte) {
	// Convert the byte slice to a hexadecimal string
	hexString := hex.EncodeToString(data)

	// Split the hex string into groups of 4 characters (2 bytes per group)
	var groupedHex []string
	for i := 0; i < len(hexString); i += 4 {
		end := i + 4
		if end > len(hexString) {
			end = len(hexString)
		}
		groupedHex = append(groupedHex, hexString[i:end])
	}

	// Join the groups with a space separator
	formattedHex := strings.Join(groupedHex, " ")

	// Print the hexadecimal string in groups of 2 bytes
	fmt.Printf("Hex: %s\n", formattedHex)

	// Print the byte slice
	fmt.Printf("Byte: %v\n", data)
}
func VerificationRouter(opt *Options, buf []byte, hex []string, header *Header, conn net.Conn) {
	msg := PackVerificationMessage(buf, hex, header)

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
	StoreClient(msg.Id, conn)

	//auto response
	if opt.AutoVerification {
		m := &VerificationResponseMessage{
			Header: &Header{
				Seq:       0,
				Encrypted: false,
			},
			Id:     msg.Id,
			Result: true,
		}
		_ = ResponseToVerification(m)
		return
	}

	//forward
	if opt.MessageForwarder != nil {
		//convert msg to json string bytes
		b, _ := json.Marshal(msg)
		_ = opt.MessageForwarder.Publish("01", b)
	}

}

func VerificationResponseRouter(c *gin.Context) {
	var req VerificationResponseMessage
	if c.ShouldBind(&req) == nil {
		err := ResponseToVerification(&req)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
	}
	c.JSON(200, gin.H{"message": "done"})
}

func SendHeartbeatResponse(conn net.Conn, header *Header) error {
	resp := &bytes.Buffer{}

	// Frame Header
	resp.Write(HexToBytes("5AA5"))

	// Data Length
	resp.Write(HexToBytes("0400"))

	// Command
	resp.Write([]byte{0x82})

	// Reserved Field
	resp.Write([]byte{0x00})

	// Checksum
	checksum := CalculateChecksum(resp.Bytes()[2:])
	resp.Write([]byte{checksum})

	_, err := conn.Write(resp.Bytes())
	if err != nil {
		log.Errorf("Failed to send Heartbeat Response: %v", err)
		return err
	}
	log.Debug("Sent Heartbeat Response successfully")
	return nil
}

func HeartbeatRouter(buf []byte, header *Header, conn net.Conn) {
	msg := PackHeartbeatMessage(buf, header)
	if msg == nil {
		log.Error("Failed to parse Heartbeat message")
		return
	}

	log.WithFields(log.Fields{
		"header":         msg.Header,
		"signalValue":    msg.SignalValue,
		"temperature":    msg.Temperature,
		"totalPortCount": msg.TotalPortCount,
		"portStatus":     msg.PortStatus,
	}).Debug("[82] Heartbeat message")

	// Send Heartbeat Response
	_ = SendHeartbeatResponse(conn, header)
}

func BillingModelVerificationRouter(opt *Options, hex []string, header *Header, conn net.Conn) {
	msg := PackBillingModelVerificationMessage(hex, header)
	log.WithFields(log.Fields{
		"id":                 msg.Id,
		"billing_model_code": msg.BillingModelCode,
	}).Debug("[05] BillingModelVerification message")

	//auto response
	if opt.AutoBillingModelVerify {
		m := &BillingModelVerificationResponseMessage{
			Header: &Header{
				Seq:       0,
				Encrypted: false,
			},
			Id:               msg.Id,
			BillingModelCode: msg.BillingModelCode,
			Result:           true,
		}
		_ = ResponseToBillingModelVerification(m)
		return
	}

	//forward
	if opt.MessageForwarder != nil {
		//convert msg to json string bytes
		b, _ := json.Marshal(msg)
		_ = opt.MessageForwarder.Publish("05", b)
	}
}

func BillingModelRequestMessageRouter(opt *Options, hex []string, header *Header, conn net.Conn) {
	msg := PackBillingModelRequestMessage(hex, header)
	log.WithFields(log.Fields{
		"id": msg.Id,
	}).Debug("[09] BillingModelRequest message")

	//forward
	if opt.MessageForwarder != nil {
		//convert msg to json string bytes
		b, _ := json.Marshal(msg)
		_ = opt.MessageForwarder.Publish("09", b)
	}
}

func BillingModelResponseMessageRouter(c *gin.Context) {
	var req BillingModelResponseMessage
	if c.ShouldBind(&req) == nil {
		err := SendBillingModelResponseMessage(&req)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
	}
	c.JSON(200, gin.H{"message": "done"})
}

func BillingModelVerificationResponseRouter(c *gin.Context) {
	var req BillingModelVerificationResponseMessage
	if c.ShouldBind(&req) == nil {
		err := ResponseToBillingModelVerification(&req)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
	}
	c.JSON(200, gin.H{"message": "done"})
}

func RemoteBootstrapRequestRouter(c *gin.Context) {
	var req RemoteBootstrapRequestMessage
	if c.ShouldBind(&req) == nil {
		err := SendRemoteBootstrapRequest(&req)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
	}
	c.JSON(200, gin.H{"message": "done"})
}

func RemoteBootstrapResponseRouter(opt *Options, hex []string, header *Header) {
	msg := PackRemoteBootstrapResponseMessage(hex, header)
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

func OfflineDataReportMessageRouter(opt *Options, raw []byte, hex []string, header *Header) {
	msg := PackOfflineDataReportMessage(hex, raw, header)
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

func RemoteShutdownResponseRouter(opt *Options, hex []string, header *Header) {
	msg := PackRemoteShutdownResponseMessage(hex, header)
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

func RemoteShutdownRequestRouter(c *gin.Context) {
	var req RemoteShutdownRequestMessage
	if c.ShouldBind(&req) == nil {
		err := SendRemoteShutdownRequest(&req)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
	}
	c.JSON(200, gin.H{"message": "done"})
}

func TransactionRecordConfirmedRouter(c *gin.Context) {
	var req TransactionRecordConfirmedMessage
	if c.ShouldBind(&req) == nil {
		err := SendTransactionRecordConfirmed(&req)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
	}
	c.JSON(200, gin.H{"message": "done"})
}

func TransactionRecordMessageRouter(opt *Options, raw []byte, hex []string, header *Header) {
	msg := PackTransactionRecordMessage(raw, hex, header)
	msgJson, _ := json.Marshal(msg)
	log.WithFields(log.Fields{
		"msg": string(msgJson),
	}).Debug("[3b] TransactionRecord message")

	if opt.AutoTransactionRecordConfirm {
		m := &TransactionRecordConfirmedMessage{
			Header: &Header{
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

func RemoteRebootResponseMessageRouter(opt *Options, hex []string, header *Header) {
	msg := PackRemoteRebootResponseMessage(hex, header)
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

func RemoteRebootRequestMessageRouter(c *gin.Context) {
	var req RemoteRebootRequestMessage
	if c.ShouldBind(&req) == nil {
		err := SendRemoteRebootRequest(&req)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
	}
	c.JSON(200, gin.H{"message": "done"})
}

func SetBillingModelRequestRouter(c *gin.Context) {
	var req SetBillingModelRequestMessage
	if c.ShouldBind(&req) == nil {
		err := SendSetBillingModelRequestMessage(&req)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
	}
	c.JSON(200, gin.H{"message": "done"})
}

func SetBillingModelResponseMessageRouter(opt *Options, hex []string, header *Header) {
	msg := PackSetBillingModelResponseMessage(hex, header)
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

func ChargingFinishedMessageRouter(opt *Options, hex []string, header *Header) {
	msg := PackChargingFinishedMessage(hex, header)
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

func DeviceLoginRouter(opt *Options, buf []byte, header *Header, conn net.Conn) {
	// Unpack Device Login Message
	msg := PackDeviceLoginMessage(buf, header)

	if msg == nil {
		log.Error("Failed to parse Device Login message due to checksum mismatch or invalid buffer")
		return
	}

	log.Debugf("Raw Buffer: %x", buf)

	// Debug: Log the parsed message
	log.WithFields(log.Fields{
		"imei":            msg.IMEI,
		"devicePortCount": msg.DevicePortCount,
		"hardwareVersion": msg.HardwareVersion,
		"softwareVersion": msg.SoftwareVersion,
		"ccid":            msg.CCID,
		"signalValue":     msg.SignalValue,
		"loginReason":     msg.LoginReason,
	}).Info("Parsed Device Login Message")

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
	// heartbeatPeriod := 30 // Default heartbeat interval (30 seconds)
	// if heartbeatPeriod < 10 || heartbeatPeriod > 250 {
	// 	heartbeatPeriod = 30 // Enforce valid range (10-250 seconds)
	// }

	// resp := &DeviceLoginResponseMessage{
	// 	Header: &Header{
	// 		Seq:       header.Seq,
	// 		Encrypted: false,
	// 	},
	// 	Time:            "00000000000000", // Reserved Time (BCD format)
	// 	HeartbeatPeriod: heartbeatPeriod,  // Valid interval
	// 	Result:          0x00,             // Login successful
	// }

	// // Pack the response message
	// data := PackDeviceLoginResponseMessage(resp)
	// PrintHexAndByte(data)
	// // Send the response back to the device
	// _, err := conn.Write(data)
	// if err != nil {
	// 	log.Errorf("Failed to send Device Login response: %v", err)
	// 	return
	// }

	message := []byte{
		0x5A, 0xA5, // Frame Header
		0x0C, 0x00, // Data Length
		0x81,                                                 // Command
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Data
		0xF0, 0x7D, // Footer
	}
	sendMessage(conn, message)
	log.Debug("Sent Device Login response successfully")

	// Forward the Device Login message to an external system (optional)
	if opt.MessageForwarder != nil {
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			log.Errorf("Failed to marshal Device Login message: %v", err)
			return
		}
		err = opt.MessageForwarder.Publish("81", jsonMsg)
		if err != nil {
			log.Errorf("Failed to publish Device Login message: %v", err)
		}
	}
}

func RemoteStartRouter(buf []byte, header *Header, conn net.Conn) {
	msg := PackRemoteStartMessage(buf, header)
	if msg == nil {
		log.Error("Failed to parse Remote Start message")
		return
	}

	log.WithFields(log.Fields{
		"port":            msg.Port,
		"orderNumber":     msg.OrderNumber,
		"startMethod":     msg.StartMethod,
		"cardNumber":      msg.CardNumber,
		"chargingMethod":  msg.ChargingMethod,
		"chargingParam":   msg.ChargingParam,
		"availableAmount": msg.AvailableAmount,
	}).Debug("[83] Remote Start message")

	// Auto Response
	response := &RemoteStartResponseMessage{
		Header:      header,
		Port:        msg.Port,
		OrderNumber: msg.OrderNumber,
		StartMethod: msg.StartMethod,
		Result:      0x00, // 0x00 for success
	}

	data := PackRemoteStartResponseMessage(response)
	_, err := conn.Write(data)
	if err != nil {
		log.Errorf("Failed to send Remote Start response: %v", err)
	} else {
		log.Debug("Sent Remote Start response successfully")
	}
}

func RemoteStopRouter(buf []byte, header *Header, conn net.Conn) {
	msg := PackRemoteStopMessage(buf, header)
	if msg == nil {
		log.Error("Failed to parse Remote Stop message")
		return
	}

	log.WithFields(log.Fields{
		"port":        msg.Port,
		"orderNumber": msg.OrderNumber,
	}).Debug("[84] Remote Stop message")

	// Auto Response
	response := &RemoteStopResponseMessage{
		Header:      header,
		Port:        msg.Port,
		OrderNumber: msg.OrderNumber,
		Result:      0x00, // 0x00 for success, other values for specific errors
	}

	data := PackRemoteStopResponseMessage(response)
	_, err := conn.Write(data)
	if err != nil {
		log.Errorf("Failed to send Remote Stop response: %v", err)
	} else {
		log.Debug("Sent Remote Stop response successfully")
	}
}

func SubmitFinalStatusRouter(opt *Options, buf []byte, header *Header, conn net.Conn) {
	msg := PackSubmitFinalStatusMessage(buf, header)
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

	// Auto Response
	response := &SubmitFinalStatusResponse{
		Header: header,
		Result: 0x00, // Success
	}

	data := PackSubmitFinalStatusResponse(response)
	_, err := conn.Write(data)
	if err != nil {
		log.Errorf("Failed to send Submit Final Status response: %v", err)
	} else {
		log.Debug("Sent Submit Final Status response successfully")
	}
}
