package handlers

import (
	"encoding/json"
	"fmt"
	"net"
	"ykc-proxy-server/dtos"
	"ykc-proxy-server/services"

	"github.com/gin-gonic/gin"
)

func StartChargingRouter(c *gin.Context) {
	fmt.Println("Charging start")
	clientID := c.DefaultQuery("clientID", "")
	if clientID == "" {
		c.JSON(400, gin.H{"error": "client ID is required"})
		return
	}
	err := services.StartCharging(clientID)
	if err != nil {
		c.JSON(400, gin.H{"error": err})
		return
	}
	c.JSON(200, gin.H{"status": "message sent"})
}

func StopChargingRouter(c *gin.Context) {

	fmt.Println("Charging stop")
	clientID := c.DefaultQuery("clientID", "")
	if clientID == "" {
		c.JSON(400, gin.H{"error": "client ID is required"})
		return
	}
	err := services.StopCharging(clientID)
	if err != nil {
		c.JSON(400, gin.H{"error": err})
		return
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
	_ = services.Hearthbeat(buf, header, conn)
	// Send Heartbeat Response
	_ = services.SendHeartbeatResponse(conn, header)
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
