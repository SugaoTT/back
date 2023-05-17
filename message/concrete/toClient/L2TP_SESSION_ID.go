package message

import (
	"encoding/json"
)

type L2TP_SESSION_ID struct {
	AbstractMessage
	SessionID  string
	SrcUUID    string `json:"srcUUID"`
	SrcEthName string `json:"srcEthName"`
	DstUUID    string `json:"dstUUID"`
	DstEthName string `json:"dstEthName"`
}

func NewL2TP_SESSION_ID(inputMsg []byte) *L2TP_SESSION_ID {
	msg := &L2TP_SESSION_ID{}
	msg.AbstractMessage.MessageType = "L2TP_SESSION_ID"

	//具象的パラメータをmsgに追加
	var ev L2TP_SESSION_ID
	json.Unmarshal(inputMsg, &ev)
	msg.SrcUUID = ev.SrcUUID
	msg.SrcEthName = ev.SrcEthName
	msg.DstUUID = ev.DstUUID
	msg.DstEthName = ev.DstEthName
	return msg
}

func (msg *L2TP_SESSION_ID) SetSessionID(SessionID string) {
	msg.SessionID = SessionID
}

func (msg *L2TP_SESSION_ID) GetSessionID() string {
	return msg.SessionID
}

func (msg *L2TP_SESSION_ID) SetSrcUUID(SrcUUID string) {
	msg.SrcUUID = SrcUUID
}

func (msg *L2TP_SESSION_ID) GetSrcUUID() string {
	return msg.SrcUUID
}

func (msg *L2TP_SESSION_ID) SetSrcEthName(SrcEthName string) {
	msg.SrcEthName = SrcEthName
}

func (msg *L2TP_SESSION_ID) GetSrcEthName() string {
	return msg.SrcEthName
}

func (msg *L2TP_SESSION_ID) SetDstUUID(DstUUID string) {
	msg.SrcUUID = DstUUID
}

func (msg *L2TP_SESSION_ID) GetDstUUID() string {
	return msg.DstUUID
}

func (msg *L2TP_SESSION_ID) SetDstEthName(DstEthName string) {
	msg.DstEthName = DstEthName
}

func (msg *L2TP_SESSION_ID) GetDstEthName() string {
	return msg.DstEthName
}
