package message

type LAUNCH_NETWORK struct {
	AbstractMessage
	NetworkTopology string
}

func NewLAUNCH_NETWORK() *LAUNCH_NETWORK {
	msg := &LAUNCH_NETWORK{}
	//msg.SetMessageType("LAUNCH_NETWORK")
	return msg
}

func (msg *LAUNCH_NETWORK) SetNetworkTopology(networkTopology string) {
	msg.NetworkTopology = networkTopology
}

func (msg *LAUNCH_NETWORK) GetNetworkTopology() string {
	return msg.NetworkTopology
}
