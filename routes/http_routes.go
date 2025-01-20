package routes

import (
	"ykc-proxy-server/handlers"

	"github.com/gin-gonic/gin"
)

func SetupHttpRoutes(r *gin.Engine) {
	r.POST("/start", handlers.StartChargingHandler)
	r.POST("/stop", handlers.StopChargingHandler)
	r.GET("/heartbeat", handlers.GetHeartbeatHandler)
	r.GET("/port-report", handlers.GetLatestPortReportHandler)
	r.POST("/proxy/02", handlers.VerificationResponseRouter)
	r.POST("/proxy/06", handlers.BillingModelVerificationResponseHandler)
	r.POST("/proxy/0a", handlers.BillingModelResponseMessageHandler)
	r.POST("/proxy/34", handlers.RemoteBootstrapRequestHandler)
	r.POST("/proxy/36", handlers.RemoteShutdownRequestHandler)
	r.POST("/proxy/40", handlers.TransactionRecordConfirmedHandler)
	r.POST("/proxy/58", handlers.SetBillingModelRequestHandler)
	r.POST("/proxy/92", handlers.RemoteRebootRequestMessageHandler)
}
