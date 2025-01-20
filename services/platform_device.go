package services

import (
	"encoding/hex"
	"ykc-proxy-server/protocols"
	"ykc-proxy-server/utils"
)

func StartCharging(IMEI string, port int, orderNumber string) error {

	conn, _, err := utils.GetClientByIMEI(IMEI)
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

func StopCharging(IMEI string, port int, orderNumber string) error {
	conn, _, err := utils.GetClientByIMEI(IMEI)
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
