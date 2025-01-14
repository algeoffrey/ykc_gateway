package services

import (
	"encoding/hex"
	"ykc-proxy-server/protocols"
	"ykc-proxy-server/utils"
)

func StartCharging(IPAddress string, port int, orderNumber string) error {
	conn, imei, err := utils.GetClientByIPAddress(IPAddress)
	if err != nil {
		return err
	}
	hexPort := []byte{byte(port)} // Convert port number directly to byte (1->0x01, 2->0x02, etc)
	orderNumberHex, _ := hex.DecodeString(orderNumber)
	utils.PrintHex(orderNumberHex)
	packet := protocols.ParseStartChargingRequest(imei, hexPort, orderNumberHex)
	err = utils.SendMessage(conn, packet)
	if err != nil {
		return err
	}
	return nil
}

func StopCharging(IPAddress string, orderNumber string) error {
	conn, imei, err := utils.GetClientByIPAddress(IPAddress)
	if err != nil {
		return err
	}
	orderNumberHex := utils.ASCIIToHex(orderNumber)
	packet := protocols.ParseStopChargingRequest(imei, orderNumberHex)
	err = utils.SendMessage(conn, packet)
	if err != nil {
		return err
	}
	return nil
}
