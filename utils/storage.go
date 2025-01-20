package utils

import (
	"fmt"
	"sync"
	"time"
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

var chargingSessions sync.Map

type ChargingSession struct {
	DeviceID  string
	Port      int
	StartTime time.Time
	StopTime  *time.Time // Pointer to allow nil for active sessions
}

// StartChargingSession creates a new charging session for a device and port
func StartChargingSession(deviceID string, port int) {
	key := getSessionKey(deviceID, port)
	session := ChargingSession{
		DeviceID:  deviceID,
		Port:      port,
		StartTime: time.Now(),
	}
	chargingSessions.Store(key, &session)
}

// StopChargingSession ends a charging session and records the stop time
func StopChargingSession(deviceID string, port int) *ChargingSession {
	key := getSessionKey(deviceID, port)
	if value, exists := chargingSessions.Load(key); exists {
		session := value.(*ChargingSession)
		now := time.Now()
		session.StopTime = &now
		chargingSessions.Delete(key) // Remove from active sessions
		return session
	}
	return nil
}

// GetChargingSession retrieves an active charging session if it exists
func GetChargingSession(deviceID string, port int) (*ChargingSession, bool) {
	key := getSessionKey(deviceID, port)
	if value, exists := chargingSessions.Load(key); exists {
		return value.(*ChargingSession), true
	}
	return nil, false
}

// Helper function to generate session storage key
func getSessionKey(deviceID string, port int) string {
	return fmt.Sprintf("%s_session_%d", deviceID, port)
}
