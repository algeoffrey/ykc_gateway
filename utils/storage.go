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

var chargingSessions sync.Map

type ChargingSession struct {
	DeviceID         string
	Port             int
	StartTime        time.Time
	StopTime         *time.Time // Pointer to allow nil for active sessions
	ExpectedStopTime *time.Time // Pointer to allow nil for active sessions
	MaxWatt          int
	Value            int
	Duration         int
}

// StartChargingSession creates a new charging session for a device and port
func StartChargingSession(deviceID string, port int, value int) {
	key := getSessionKey(deviceID, port)
	session := ChargingSession{
		DeviceID:  deviceID,
		Port:      port,
		StartTime: time.Now(),
		StopTime:  nil,
		MaxWatt:   0, // Track maximum wattage during charging session
		Value:     value,
		Duration:  0,
	}
	chargingSessions.Store(key, &session)
}

func UpdateSessionMaxWatt(deviceID string, port int, watt int) *ChargingSession {
	key := getSessionKey(deviceID, port)
	if value, exists := chargingSessions.Load(key); exists {
		session := value.(*ChargingSession)
		if session.MaxWatt < int(watt) {
			session.MaxWatt = watt
		}
		chargingSessions.Store(key, session)
		return session
	}
	return nil
}

func UpdateSessionExpectedStopTime(deviceID string, port int) *ChargingSession {
	key := getSessionKey(deviceID, port)
	if value, exists := chargingSessions.Load(key); exists {
		session := value.(*ChargingSession)
		expectedStop := session.StartTime.Add(1000 * time.Second) // 1000 seconds
		session.ExpectedStopTime = &expectedStop
		chargingSessions.Store(key, session)
		return session
	}
	return nil
}

// StopChargingSession ends a charging session and records the stop time
func StopChargingSession(deviceID string, port int) *ChargingSession {
	key := getSessionKey(deviceID, port)
	if value, exists := chargingSessions.Load(key); exists {
		session := value.(*ChargingSession)
		now := time.Now()
		session.StopTime = &now
		duration := now.Sub(session.StartTime).Seconds()
		// Keep the session in storage rather than deleting it
		chargingSessions.Store(key, session)
		session.Duration = int(duration)
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
