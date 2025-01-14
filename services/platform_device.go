package services

import (
	"ykc-proxy-server/protocols"
	"ykc-proxy-server/utils"
)

func StartCharging(IPAddress string) error {
	conn, imei, err := utils.GetClientByIPAddress(IPAddress)
	if err != nil {
		return err
	}
	packet := protocols.ParseStartChargingRequest(imei)
	err = utils.SendMessage(conn, packet)
	if err != nil {
		return err
	}
	return nil
}

func StopCharging(IPAddress string) error {
	conn, imei, err := utils.GetClientByIPAddress(IPAddress)
	if err != nil {
		return err
	}
	packet := protocols.ParseStopChargingRequest(imei)
	err = utils.SendMessage(conn, packet)
	if err != nil {
		return err
	}
	return nil
}
