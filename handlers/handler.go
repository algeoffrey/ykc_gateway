package handlers

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"ykc-proxy-server/dtos"
	"ykc-proxy-server/services"
	"ykc-proxy-server/utils"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func StartChargingHandler(c *gin.Context) {
	var req dtos.StartChargingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	if req.DeviceID == "" {
		c.JSON(400, gin.H{"error": "client ID is required"})
		return
	}

	if req.Port == 0 {
		c.JSON(400, gin.H{"error": "port is required"})
		return
	}

	if req.OrderNumber == "" {
		c.JSON(400, gin.H{"error": "order number is required"})
		return
	}
	if req.Value == 0 {
		c.JSON(400, gin.H{"error": "value is required"})
		return
	}

	err := services.StartCharging(req.DeviceID, req.Port, req.OrderNumber)

	utils.StartChargingSession(req.DeviceID, req.Port, req.Value)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"status": "message sent"})
}

func StopChargingHandler(c *gin.Context) {
	var req dtos.StopChargingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body"})
		return
	}

	if req.DeviceID == "" {
		c.JSON(400, gin.H{"error": "client ID is required"})
		return
	}

	if req.OrderNumber == "" {
		c.JSON(400, gin.H{"error": "order number is required"})
		return
	}

	if req.Port == 0 {
		c.JSON(400, gin.H{"error": "port is required"})
		return
	}

	err := services.StopCharging(req.DeviceID, req.Port, req.OrderNumber)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	session := utils.StopChargingSession(req.DeviceID, req.Port)
	if session == nil {
		log.Warn("No active charging session found to stop")
	}

	c.JSON(200, gin.H{"status": "message sent"})
}

func VerificationResponseRouter(c *gin.Context) {
	var req dtos.VerificationResponseMessage
	if c.ShouldBind(&req) == nil {
		err := services.ResponseToVerification(&req)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
	}
	c.JSON(200, gin.H{"message": "done"})
}

func VerificationHandler(opt *dtos.Options, buf []byte, hex []string, header *dtos.Header, conn net.Conn,
) {
	msg := services.Verification(opt, buf, hex, header, conn)

	if opt.AutoVerification {
		m := &dtos.VerificationResponseMessage{
			Header: &dtos.Header{
				Seq:       0,
				Encrypted: false,
			},
			Id:     msg.Id,
			Result: true,
		}
		_ = services.ResponseToVerification(m)
		return
	}

	//forward
	if opt.MessageForwarder != nil {
		//convert msg to json string bytes
		b, _ := json.Marshal(msg)
		_ = opt.MessageForwarder.Publish("01", b)
	}
}

func HeartbeatHandler(buf []byte, header *dtos.Header, conn net.Conn) {
	services.Hearthbeat(buf, header, conn)
	// Send Heartbeat Response
	// _ = services.SendHeartbeatResponse(conn, header)
}

func BillingModelVerificationHandler(opt *dtos.Options, hex []string, header *dtos.Header, conn net.Conn) {
	msg := services.BillingModelVerification(opt, hex, header, conn)
	//auto response
	if opt.AutoBillingModelVerify {
		m := &dtos.BillingModelVerificationResponseMessage{
			Header: &dtos.Header{
				Seq:       0,
				Encrypted: false,
			},
			Id:               msg.Id,
			BillingModelCode: msg.BillingModelCode,
			Result:           true,
		}
		_ = services.ResponseToBillingModelVerification(m)
		return
	}

	//forward
	if opt.MessageForwarder != nil {
		//convert msg to json string bytes
		b, _ := json.Marshal(msg)
		_ = opt.MessageForwarder.Publish("05", b)
	}
}

func BillingModelRequestMessageHandler(opt *dtos.Options, hex []string, header *dtos.Header, conn net.Conn) {
	msg := services.BillingModelRequestMessage(opt, hex, header, conn)
	//forward
	if opt.MessageForwarder != nil {
		//convert msg to json string bytes
		b, _ := json.Marshal(msg)
		_ = opt.MessageForwarder.Publish("09", b)
	}
}

func BillingModelResponseMessageHandler(c *gin.Context) {
	var req dtos.BillingModelResponseMessage
	if c.ShouldBind(&req) == nil {
		err := services.SendBillingModelResponseMessage(&req)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
	}
	c.JSON(200, gin.H{"message": "done"})
}

func BillingModelVerificationResponseHandler(c *gin.Context) {
	var req dtos.BillingModelVerificationResponseMessage
	if c.ShouldBind(&req) == nil {
		err := services.ResponseToBillingModelVerification(&req)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
	}
	c.JSON(200, gin.H{"message": "done"})
}

func RemoteBootstrapRequestHandler(c *gin.Context) {
	var req dtos.RemoteBootstrapRequestMessage
	if c.ShouldBind(&req) == nil {
		err := services.SendRemoteBootstrapRequest(&req)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
	}
	c.JSON(200, gin.H{"message": "done"})
}

func RemoteShutdownRequestHandler(c *gin.Context) {
	var req dtos.RemoteShutdownRequestMessage
	if c.ShouldBind(&req) == nil {
		err := services.SendRemoteShutdownRequest(&req)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
	}
	c.JSON(200, gin.H{"message": "done"})
}

func TransactionRecordConfirmedHandler(c *gin.Context) {
	var req dtos.TransactionRecordConfirmedMessage
	if c.ShouldBind(&req) == nil {
		err := services.SendTransactionRecordConfirmed(&req)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
	}
	c.JSON(200, gin.H{"message": "done"})
}

func SetBillingModelRequestHandler(c *gin.Context) {
	var req dtos.SetBillingModelRequestMessage
	if c.ShouldBind(&req) == nil {
		err := services.SendSetBillingModelRequestMessage(&req)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
	}
	c.JSON(200, gin.H{"message": "done"})
}

func RemoteRebootRequestMessageHandler(c *gin.Context) {
	var req dtos.RemoteRebootRequestMessage
	if c.ShouldBind(&req) == nil {
		err := services.SendRemoteRebootRequest(&req)
		if err != nil {
			c.JSON(500, gin.H{"message": err})
			return
		}
	}
	c.JSON(200, gin.H{"message": "done"})
}

func SubmitFinalStatusHandler(opt *dtos.Options, buf []byte, header *dtos.Header, conn net.Conn,
) {
	data := services.SubmitFinalStatus(opt, buf, header, conn)
	utils.PrintHex(data)
	err := utils.SendMessage(conn, data)
	if err != nil {
		log.Errorf("Failed to send Submit Final Status response: %v", err)
	} else {
		IPAddress := conn.RemoteAddr().String()
		fmt.Println(IPAddress)
		log.Debug("Sent Submit Final Status response successfully")
	}
}

func DeviceLoginHandler(opt *dtos.Options, buf []byte, header *dtos.Header, conn net.Conn,
) {
	msg, resMessage := services.DeviceLogin(opt, buf, header, conn)
	// Forward the Device Login message to an external system (optional)
	if msg != nil {
		if opt.MessageForwarder != nil {
			jsonMsg, err := json.Marshal(msg)
			if err != nil {
				log.Errorf("Failed to marshal Device Login message: %v", err)
				return
			}
			err = opt.MessageForwarder.Publish("81", jsonMsg)
			if err != nil {
				log.Errorf("Failed to publish Device Login message: %v", err)
			}
		}
	}
	if resMessage != nil {
		utils.SendMessage(conn, resMessage)
	}
}

func RemoteStartHandler(buf []byte, header *dtos.Header, conn net.Conn) {
	services.RemoteStart(buf, header, conn)

}

func RemoteStopHandler(buf []byte, header *dtos.Header, conn net.Conn) {
	services.RemoteStop(buf, header, conn)

}

func ChargingPortDataHandler(opt *dtos.Options, buf []byte, header *dtos.Header, conn net.Conn) {
	// Parse the CMD088 message
	services.ChargingPortData(opt, buf, header, conn)

	// Forward the Charging Port Data message to an external system (optional)
	// if msg != nil {
	// 	if opt.MessageForwarder != nil {
	// 		jsonMsg, err := json.Marshal(msg)
	// 		if err != nil {
	// 			log.Errorf("Failed to marshal Charging Port Data message: %v", err)
	// 			return
	// 		}
	// 		err = opt.MessageForwarder.Publish("88", jsonMsg)
	// 		if err != nil {
	// 			log.Errorf("Failed to publish Charging Port Data message: %v", err)
	// 		}
	// 	}
	// }

}

func GetHeartbeatHandler(c *gin.Context) {
	deviceID := c.Query("deviceId")
	if deviceID == "" {
		c.JSON(400, gin.H{"error": "deviceId is required"})
		return
	}

	heartbeat, found := utils.GetHeartbeat(deviceID)
	if !found {
		c.JSON(404, gin.H{"error": "no heartbeat found for device"})
		return
	}

	c.JSON(200, heartbeat)
}

func GetLatestPortReportHandler(c *gin.Context) {
	deviceID := c.Query("deviceId")
	if deviceID == "" {
		c.JSON(400, gin.H{"error": "deviceId is required"})
		return
	}

	port := c.Query("port")
	if port == "" {
		c.JSON(400, gin.H{"error": "port is required"})
		return
	}

	portNum, err := strconv.Atoi(port)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid port number"})
		return
	}

	portReport, found := utils.GetPortReport(deviceID, portNum)
	if !found {
		c.JSON(404, gin.H{"error": "no port report found for device"})
		return
	}

	c.JSON(200, portReport)
}

func GetChargingSessionHandler(c *gin.Context) {
	deviceID := c.Query("deviceId")
	if deviceID == "" {
		c.JSON(400, gin.H{"error": "deviceId is required"})
		return
	}

	port := c.Query("port")
	if port == "" {
		c.JSON(400, gin.H{"error": "port is required"})
		return
	}
	intPort, err := strconv.Atoi(port)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid port number"})
		return
	}
	session, found := utils.GetChargingSession(deviceID, intPort)
	if !found {
		c.JSON(404, gin.H{"error": "no charging session found for device"})
		return
	}

	c.JSON(200, session)
}
