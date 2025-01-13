package services

import (
	"fmt"
	"ykc-proxy-server/protocols"
	"ykc-proxy-server/utils"
)

func StartCharging(IPAddress string) error {

	fmt.Println(IPAddress)
	// // Retrieve the client connection using the GetClient function
	conn, imei, err := utils.GetClientByIPAddress(IPAddress)

	if err != nil {
		return err
	}
	packet := protocols.ParseStartChargingRequest(imei)

	// packet := []byte{
	// 	0x5A, 0xA5, // header
	// 	0x16, 0x00, // length
	// 	0x83, 0x00, // cmd
	// 	0x38, 0x36, 0x31, 0x34, 0x33, 0x35, 0x30, 0x37, 0x33, 0x39, 0x30, 0x30, 0x38, 0x34, 0x33, //emai,
	// 	0x01,                   //port
	// 	0x01, 0x01, 0x05, 0x01, // order
	// 	0x03,                   // payment mode
	// 	0x00, 0x00, 0x00, 0x00, // card number
	// 	0x05,                   // charging mode
	// 	0xE8, 0x03, 0x00, 0x00, // time number
	// 	0x64, 0x00, 0x00, 0x00, //balance
	// 	0xED, // checksum
	// }

	// // Send the message to the client
	err = utils.SendMessage(conn, packet)
	if err != nil {
		return err
	}

	return nil
}

func StopCharging(IPAddress string) error {
	fmt.Println(IPAddress)
	conn, _, err := utils.GetClientByIPAddress(IPAddress)
	if err != nil {
		return err
	}
	packet := []byte{
		0x5A, 0xA5,
		0x08, 0x00,
		0x84, 0x00,
		0x01, 0x01, 0x05, 0x01,
		0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01,
		0x01, 0x01, 0x00, 0x00, 0x00, 0x8F}
	err = utils.SendMessage(conn, packet)
	if err != nil {
		return err
	}
	return nil
}
