package main

import (
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func SendCustomMessage(c *gin.Context) {
	// Get the client ID from query parameter
	clientID := c.DefaultQuery("clientID", "")
	if clientID == "" {
		c.JSON(400, gin.H{"error": "client ID is required"})
		return
	}

	// Retrieve the client connection using the GetClient function
	conn, err := GetClient(clientID)
	if err != nil {
		c.JSON(404, gin.H{"error": "client not found"})
		return
	}

	// Define the message you want to send (in bytes format)
	message := []byte{0x5A, 0xA5, 0x11, 0x00, 0x82, 0x1F, 0x1E, 0x0A, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0xDC}

	// Send the message to the client
	err = sendMessage(conn, message)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to send message"})
		return
	}

	// Respond to HTTP request
	c.JSON(200, gin.H{"status": "message sent"})
}

func main() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	opt := parseOptions()

	//define message forwarder
	var f MessageForwarder
	switch opt.MessagingServerType {
	case "http":
		f := &HTTPForwarder{
			Endpoints: opt.Servers,
		}
		opt.MessageForwarder = f
		break
	case "nats":
		servers := strings.Join(opt.Servers, ",")
		f := &NatsForwarder{
			Servers:  servers,
			Username: opt.Username,
			Password: opt.Password,
		}
		f.Connect()
		opt.MessageForwarder = f
		break
	default:
		f = nil
		opt.MessageForwarder = f
		break
	}

	go enableTcpServer(opt)
	go enableHttpServer(opt)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	sig := <-sigChan
	log.Info("exit:", sig)
	os.Exit(0)
}

func enableTcpServer(opt *Options) {
	host := opt.Host
	port := strconv.Itoa(opt.TcpPort)
	addr, err := net.ResolveTCPAddr("tcp", host+":"+port)
	if err != nil {
		log.Error("error resolving address:", err)
		return
	}

	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Error("error listening:", err)
		return
	}
	defer ln.Close()
	log.Info("server listening on", addr.String())

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Error("error accepting connection:", err)
			continue
		}
		StoreClient(conn.RemoteAddr().String(), conn)
		go handleConnection(opt, conn)
	}
}

func enableHttpServer(opt *Options) {
	r := gin.Default()
	r.POST("/proxy/02", VerificationResponseRouter)
	r.POST("/proxy/06", BillingModelVerificationResponseRouter)
	r.POST("/proxy/0a", BillingModelResponseMessageRouter)
	r.POST("/proxy/34", RemoteBootstrapRequestRouter)
	r.POST("/proxy/36", RemoteShutdownRequestRouter)
	r.POST("/proxy/40", TransactionRecordConfirmedRouter)
	r.POST("/proxy/58", SetBillingModelRequestRouter)
	r.POST("/proxy/92", RemoteRebootRequestMessageRouter)
	r.POST("/send-custom-message", SendCustomMessage)
	host := opt.Host
	port := strconv.Itoa(opt.HttpPort)
	err := r.Run(host + ":" + port)
	if err != nil {
		panic(err)
	}
}

func handleConnection(opt *Options, conn net.Conn) {
	defer conn.Close()

	log.WithFields(log.Fields{
		"address": conn.RemoteAddr().String(),
	}).Info("new client connected")

	var connErr error
	for connErr == nil {
		connErr = drain(opt, conn)
		time.Sleep(time.Millisecond * 1)
	}

}

func drain(opt *Options, conn net.Conn) error {
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Error("Error reading: ", err)
		return err
	}

	hex := BytesToHex(buf[:n])

	encrypted := false
	if buf[4] == byte(0x01) {
		encrypted = true
	}

	length := buf[1]
	seq := buf[3]<<8 | buf[2]

	header := &Header{
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
	case Verification:
		VerificationRouter(opt, buf, hex, header, conn)
		break
	case Heartbeat:
		HeartbeatRouter(buf, header, conn)
		break
	case BillingModelVerification:
		BillingModelVerificationRouter(opt, hex, header, conn)
		break
	case BillingModelRequest:
		BillingModelRequestMessageRouter(opt, hex, header, conn)
		break
	case OfflineDataReport:
		OfflineDataReportMessageRouter(opt, buf, hex, header)
		break
	case ChargingFinished:
		ChargingFinishedMessageRouter(opt, hex, header)
		break
	case RemoteBootstrapResponse:
		RemoteBootstrapResponseRouter(opt, hex, header)
		break
	case RemoteShutdownResponse:
		RemoteShutdownResponseRouter(opt, hex, header)
		break
	case SetBillingModelResponse:
		SetBillingModelResponseMessageRouter(opt, hex, header)
		break
	case RemoteRebootResponse:
		RemoteRebootResponseMessageRouter(opt, hex, header)
		break
	case TransactionRecord:
		TransactionRecordMessageRouter(opt, buf, hex, header)
		break
	case DeviceLogin:
		log.Debug("Handling Device Login...")
		DeviceLoginRouter(opt, buf, header, conn)
		break
	case RemoteStart:
		RemoteStartRouter(buf, header, conn)
		break
	case RemoteStop:
		RemoteStopRouter(buf, header, conn)
		break
	case SubmitFinalStatus:
		SubmitFinalStatusRouter(opt, buf, header, conn)
	default:
		log.WithFields(log.Fields{
			"frame_id": int(buf[5]),
		}).Info("unsupported message")
		break

	}
	return nil
}
func sendMessage(conn net.Conn, message []byte) error {
	// Convert message to bytes or proper format
	msgBytes := []byte(message)

	// Send the message to the device
	_, err := conn.Write(msgBytes)
	if err != nil {
		log.Error("Error sending message:", err)
		return err
	}

	log.Infof("Sent message to %s: %s", conn.RemoteAddr().String(), message)
	return nil
}
