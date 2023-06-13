package messageFromClient

import (
	"encoding/json"

	abstractmessage "github.com/SugaoTT/back/message"
)

type L2TP_INFO_REQUEST struct {
	abstractmessage.AbstractMessage
	SessionID   string
	SrcTunnelID string
	DstTunnelID string
	SrcUUID     string `json:"srcUUID"`
	SrcEthName  string `json:"srcEthName"`
	DstUUID     string `json:"dstUUID"`
	DstEthName  string `json:"dstEthName"`
}

func NewL2TP_INFO_REQUEST(inputMsg []byte) *L2TP_INFO_REQUEST {
	msg := &L2TP_INFO_REQUEST{}
	msg.AbstractMessage.MessageType = "L2TP_INFO_REQUEST"

	//具象的パラメータをmsgに追加
	var ev L2TP_INFO_REQUEST
	json.Unmarshal(inputMsg, &ev)
	msg.SrcUUID = ev.SrcUUID
	msg.SrcEthName = ev.SrcEthName
	msg.DstUUID = ev.DstUUID
	msg.DstEthName = ev.DstEthName
	return msg
}
