package routes

import (
	"net"
	"strconv"
	"ykc-proxy-server/dtos"
	"ykc-proxy-server/handlers"
	"ykc-proxy-server/protocols"
	"ykc-proxy-server/services"
	"ykc-proxy-server/utils"

	log "github.com/sirupsen/logrus"
)

func Drain(opt *dtos.Options, conn net.Conn) error {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Error("Error reading: ", err)
		return err
	}

	hex := utils.BytesToHex(buf[:n])

	encrypted := false
	if buf[4] == byte(0x01) {
		encrypted = true
	}

	length := buf[1]
	seq := buf[3]<<8 | buf[2]

	header := &dtos.Header{
		Length:    int(length),
		Seq:       int(seq),
		Encrypted: encrypted,
		FrameId:   strconv.Itoa(int(buf[4])),
	}

	log.WithFields(log.Fields{
		"hex":       hex,
		"encrypted": encrypted,
		"length":    length,
		"seq":       seq,
		"frame_id":  int(buf[4]),
	}).Info("Received message")

	log.Debugf("buf[4] (frame_id in hex): %X", buf[4]) // Added for clarity

	switch buf[4] {
	case protocols.Verification:
		handlers.VerificationHandler(opt, buf, hex, header, conn)
	case protocols.Heartbeat:
		handlers.HeartbeatHandler(buf, header, conn)
	case protocols.BillingModelVerification:
		handlers.BillingModelVerificationHandler(opt, hex, header, conn)
	case protocols.BillingModelRequest:
		handlers.BillingModelRequestMessageHandler(opt, hex, header, conn)
	case protocols.OfflineDataReport:
		services.OfflineDataReportMessageRouter(opt, buf, hex, header)
	case protocols.ChargingFinished:
		services.ChargingFinishedMessageRouter(opt, hex, header)
	case protocols.RemoteBootstrapResponse:
		services.RemoteBootstrapResponseRouter(opt, hex, header)
	case protocols.RemoteShutdownResponse:
		services.RemoteShutdownResponseRouter(opt, hex, header)
	case protocols.SetBillingModelResponse:
		services.SetBillingModelResponseMessageRouter(opt, hex, header)
	case protocols.RemoteRebootResponse:
		services.RemoteRebootResponseMessageRouter(opt, hex, header)
	case protocols.TransactionRecord:
		services.TransactionRecordMessageRouter(opt, buf, hex, header)
	case protocols.DeviceLogin:
		services.DeviceLoginRouter(opt, buf, header, conn)
	case protocols.RemoteStart:
		services.RemoteStartRouter(buf, header, conn)
	case protocols.RemoteStop:
		services.RemoteStopRouter(buf, header, conn)
	case protocols.SubmitFinalStatus:
		services.SubmitFinalStatusRouter(opt, buf, header, conn)
	default:
		log.WithFields(log.Fields{
			"frame_id": int(buf[5]),
		}).Info("unsupported message")
	}
	return nil
}
