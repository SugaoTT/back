package messageToClient

import (
	abstractmessage "github.com/SugaoTT/back/message"
)

type L2TP_INFO struct {
	abstractmessage.AbstractMessage
	SessionID   string
	SrcTunnelID string
	DstTunnelID string
	SrcUUID     string `json:"srcUUID"`
	SrcEthName  string `json:"srcEthName"`
	DstUUID     string `json:"dstUUID"`
	DstEthName  string `json:"dstEthName"`
}

func NewL2TP_INFO() *L2TP_INFO {
	msg := &L2TP_INFO{}
	msg.AbstractMessage.MessageType = "L2TP_INFO"

	/*//具象的パラメータをmsgに追加
	var ev L2TP_INFO
	json.Unmarshal(inputMsg, &ev)
	msg.SrcUUID = ev.SrcUUID
	msg.SrcEthName = ev.SrcEthName
	msg.DstUUID = ev.DstUUID
	msg.DstEthName = ev.DstEthName*/
	return msg
}

func (msg *L2TP_INFO) SetSessionID(SessionID string) {
	msg.SessionID = SessionID
}

func (msg *L2TP_INFO) GetSessionID() string {
	return msg.SessionID
}

func (msg *L2TP_INFO) SetSrcTunnelID(SrcTunnelID string) {
	msg.SrcTunnelID = SrcTunnelID
}

func (msg *L2TP_INFO) GetSrcTunnelID() string {
	return msg.SrcTunnelID
}

func (msg *L2TP_INFO) SetDstTunnelID(DstTunnelID string) {
	msg.DstTunnelID = DstTunnelID
}

func (msg *L2TP_INFO) GetDstTunnelID() string {
	return msg.DstTunnelID
}

func (msg *L2TP_INFO) SetSrcUUID(SrcUUID string) {
	msg.SrcUUID = SrcUUID
}

func (msg *L2TP_INFO) GetSrcUUID() string {
	return msg.SrcUUID
}

func (msg *L2TP_INFO) SetSrcEthName(SrcEthName string) {
	msg.SrcEthName = SrcEthName
}

func (msg *L2TP_INFO) GetSrcEthName() string {
	return msg.SrcEthName
}

func (msg *L2TP_INFO) SetDstUUID(DstUUID string) {
	msg.DstUUID = DstUUID
}

func (msg *L2TP_INFO) GetDstUUID() string {
	return msg.DstUUID
}

func (msg *L2TP_INFO) SetDstEthName(DstEthName string) {
	msg.DstEthName = DstEthName
}

func (msg *L2TP_INFO) GetDstEthName() string {
	return msg.DstEthName
}
