package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
	"ykc-proxy-server/dtos"
	"ykc-proxy-server/forwarder"
	"ykc-proxy-server/handlers"
	"ykc-proxy-server/routes"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	opt := parseOptions()

	//define message forwarder
	var f forwarder.MessageForwarder
	switch opt.MessagingServerType {
	case "http":
		f := &forwarder.HTTPForwarder{
			Endpoints: opt.Servers,
		}
		opt.MessageForwarder = f
	case "nats":
		servers := strings.Join(opt.Servers, ",")
		f := &forwarder.NatsForwarder{
			Servers:  servers,
			Username: opt.Username,
			Password: opt.Password,
		}
		f.Connect()
		opt.MessageForwarder = f
	default:
		f = nil
		opt.MessageForwarder = f
	}

	go enableTcpServer(opt)
	go enableHttpServer(opt)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	sig := <-sigChan
	log.Info("exit:", sig)
	os.Exit(0)
}

func enableTcpServer(opt *dtos.Options) {
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

		fmt.Println(conn.RemoteAddr().String())
		go handleConnection(opt, conn)
	}
}

func enableHttpServer(opt *dtos.Options) {
	fmt.Println("Http ok")
	r := gin.Default()
	r.POST("/start", handlers.StartChargingHandler)
	r.POST("/stop", handlers.StopChargingHandler)
	r.POST("/proxy/02", handlers.VerificationResponseRouter)
	r.POST("/proxy/06", handlers.BillingModelVerificationResponseHandler)
	r.POST("/proxy/0a", handlers.BillingModelResponseMessageHandler)
	r.POST("/proxy/34", handlers.RemoteBootstrapRequestHandler)
	r.POST("/proxy/36", handlers.RemoteShutdownRequestHandler)
	r.POST("/proxy/40", handlers.TransactionRecordConfirmedHandler)
	r.POST("/proxy/58", handlers.SetBillingModelRequestHandler)
	r.POST("/proxy/92", handlers.RemoteRebootRequestMessageHandler)
	host := opt.Host

	port := strconv.Itoa(opt.HttpPort)
	fmt.Println(port)
	fmt.Println(host)
	err := r.Run(host + ":" + port)
	if err != nil {
		panic(err)
	}
}

func handleConnection(opt *dtos.Options, conn net.Conn) {
	defer conn.Close()

	log.WithFields(log.Fields{
		"address": conn.RemoteAddr().String(),
	}).Info("new client connected")

	var connErr error
	for connErr == nil {
		connErr = routes.HandleChargingProtocol(opt, conn)
		time.Sleep(time.Millisecond * 1)
	}

}
