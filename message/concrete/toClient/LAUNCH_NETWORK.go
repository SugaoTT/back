package messageToClient

import (
	abstractmessage "github.com/SugaoTT/back/message"
)

type LAUNCH_NETWORK_SUCCESS struct {
	abstractmessage.AbstractMessage
	NetworkTopology string
}

func NewLAUNCH_NETWORK_SUCCESS() *LAUNCH_NETWORK_SUCCESS {
	msg := &LAUNCH_NETWORK_SUCCESS{}
	msg.AbstractMessage.MessageType = "LAUNCH_NETWORK_SUCCESS"
	return msg
}

func (msg *LAUNCH_NETWORK_SUCCESS) SetNetworkTopology(networkTopology string) {
	msg.NetworkTopology = networkTopology
}

func (msg *LAUNCH_NETWORK_SUCCESS) GetNetworkTopology() string {
	return msg.NetworkTopology
}
