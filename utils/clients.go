package utils

import (
	"errors"
	"net"
	"sync"

	log "github.com/sirupsen/logrus"
)

var clients sync.Map

func StoreClient(id string, conn net.Conn) {
	clients.Store(id, conn)
}

func GetClient(id string) (net.Conn, error) {
	value, ok := clients.Load(id)
	if ok {
		conn := value.(net.Conn)
		return conn, nil
	} else {
		return nil, errors.New("client does not exist")
	}
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
