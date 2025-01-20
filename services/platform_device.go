package services

import (
	"encoding/hex"
	"ykc-proxy-server/protocols"
	"ykc-proxy-server/utils"
)

func StartCharging(deviceID string, port int, orderNumber string) error {

	conn, _, err := utils.GetClientByDeviceID(deviceID)
	if err != nil {
		return err
	}
	hexPort := []byte{byte(port)} // Convert port number directly to byte (1->0x01, 2->0x02, etc)
	orderNumberHex, _ := hex.DecodeString(orderNumber)
	utils.PrintHex(orderNumberHex)
	packet := protocols.ParseStartChargingRequest(hexPort, orderNumberHex)
	err = utils.SendMessage(conn, packet)
	if err != nil {
		return err
	}
	return nil
}

func StopCharging(deviceID string, port int, orderNumber string) error {
	conn, _, err := utils.GetClientByDeviceID(deviceID)
	if err != nil {
		return err
	}

	hexPort := []byte{byte(port)} // Convert port number directly to byte (1->0x01, 2->0x02, etc)
	orderNumberHex, _ := hex.DecodeString(orderNumber)
	packet := protocols.ParseStopChargingRequest(orderNumberHex, hexPort)
	err = utils.SendMessage(conn, packet)
	if err != nil {
		return err
	}
	return nil
}
