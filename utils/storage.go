package utils

import (
	"sync"
	"ykc-proxy-server/dtos"
)

var heartbeats sync.Map

// StoreHeartbeat stores the heartbeat message for a given IMEI
func StoreHeartbeat(deviceID string, heartbeat *dtos.HeartbeatMessage) {
	heartbeats.Store(deviceID, heartbeat)
}

// GetHeartbeat retrieves the stored heartbeat message for a given IMEI
func GetHeartbeat(deviceID string) (*dtos.HeartbeatMessage, bool) {
	value, ok := heartbeats.Load(deviceID)
	if !ok {
		return nil, false
	}
	return value.(*dtos.HeartbeatMessage), true
}

var portReports sync.Map

// StorePortReport stores the charging port data report for a device ID and port number
func StorePortReport(deviceID string, portNumber int, report interface{}) {
	key := getPortReportKey(deviceID, portNumber)
	portReports.Store(key, report)
}

// GetPortReport retrieves the stored charging port data report for a device ID and port number
func GetPortReport(deviceID string, portNumber int) (interface{}, bool) {
	key := getPortReportKey(deviceID, portNumber)
	return portReports.Load(key)
}

// ClearPortReport removes the stored charging port data report for a device ID and port number
func ClearPortReport(deviceID string, portNumber int) {
	key := getPortReportKey(deviceID, portNumber)
	portReports.Delete(key)
}

// Helper function to generate storage key
func getPortReportKey(deviceID string, portNumber int) string {
	return deviceID + "_" + string(rune(portNumber))
}

var reports85 sync.Map

// Store85Report stores the 0x85 report for a device ID and port number
func Store85Report(deviceID string, portNumber int, report interface{}) {
	key := get85ReportKey(deviceID, portNumber)
	reports85.Store(key, report)
}

// Get85Report retrieves the stored 0x85 report for a device ID and port number
func Get85Report(deviceID string, portNumber int) (interface{}, bool) {
	key := get85ReportKey(deviceID, portNumber)
	return reports85.Load(key)
}

// Clear85Report removes the stored 0x85 report for a device ID and port number
func Clear85Report(deviceID string, portNumber int) {
	key := get85ReportKey(deviceID, portNumber)
	reports85.Delete(key)
}

// Helper function to generate storage key for 0x85 reports
func get85ReportKey(deviceID string, portNumber int) string {
	return deviceID + "_85_" + string(rune(portNumber))
}
