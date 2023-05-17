package message

type L2TP_TUNNEL_ID struct {
	AbstractMessage
	TunnelID string
}

func NewL2TP_TUNNEL_ID() *L2TP_TUNNEL_ID {
	msg := &L2TP_TUNNEL_ID{}
	msg.AbstractMessage.MessageType = "L2TP_TUNNEL_ID"
	return msg
}

func (msg *L2TP_TUNNEL_ID) SetTunnelID(TunnelID string) {
	msg.TunnelID = TunnelID
}

func (msg *L2TP_TUNNEL_ID) GetTunnelID() string {
	return msg.TunnelID
}
