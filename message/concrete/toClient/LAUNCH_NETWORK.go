package messageToClient

import (
	abstractmessage "github.com/SugaoTT/back/message"
)

type LAUNCH_NETWORK struct {
	abstractmessage.AbstractMessage
	NetworkTopology string
}

func NewLAUNCH_NETWORK() *LAUNCH_NETWORK {
	msg := &LAUNCH_NETWORK{}
	msg.AbstractMessage.MessageType = "LAUNCH_NETWORK"
	return msg
}

func (msg *LAUNCH_NETWORK) SetNetworkTopology(networkTopology string) {
	msg.NetworkTopology = networkTopology
}

func (msg *LAUNCH_NETWORK) GetNetworkTopology() string {
	return msg.NetworkTopology
}
