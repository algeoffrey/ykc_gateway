package utils

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"ykc-proxy-server/dtos"

	log "github.com/sirupsen/logrus"
)

var clients sync.Map

func StoreClient(client dtos.ClientInfo, conn net.Conn) {
	clients.Store(client, conn)
}

func GetClient(info dtos.ClientInfo) (net.Conn, error) {
	value, ok := clients.Load(info)
	if ok {
		conn := value.(net.Conn)
		return conn, nil
	} else {
		return nil, errors.New("client does not exist")
	}
}

func GetClientByIPAddress(ipAddress string) (net.Conn, string, error) {
	var foundConn net.Conn
	var found bool
	var imei string

	clients.Range(func(key, value interface{}) bool {
		clientInfo := key.(dtos.ClientInfo) // Cast the key to ClientInfo
		fmt.Println(clientInfo)
		if clientInfo.IPAddress == ipAddress {
			foundConn = value.(net.Conn) // Cast the value to net.Conn
			imei = clientInfo.IMEI
			found = true
			return false // Stop iteration as we found the client
		}
		return true // Continue iteration
	})

	if found {
		return foundConn, imei, nil
	}
	return nil, "", errors.New("client does not exist")
}

func SendMessage(conn net.Conn, message []byte) error {
	// Convert message to bytes or proper format
	PrintHexAndByte(message)
	// Send the message to the device
	_, err := conn.Write(message)
	if err != nil {
		log.Error("Error sending message:", err)
		return err
	}

	log.Infof("Sent message to %s: %s", conn.RemoteAddr().String(), message)
	return nil
}
