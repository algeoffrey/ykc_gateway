package dtos

import "ykc-proxy-server/forwarder"

type Options struct {
	Host                         string
	TcpPort                      int
	HttpPort                     int
	AutoVerification             bool
	AutoHeartbeatResponse        bool
	AutoBillingModelVerify       bool
	AutoTransactionRecordConfirm bool
	MessagingServerType          string
	Servers                      []string
	Username                     string
	Password                     string
	MessageForwarder             forwarder.MessageForwarder
	PublishSubjectPrefix         string
}
