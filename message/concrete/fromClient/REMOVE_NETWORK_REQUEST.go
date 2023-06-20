package messageFromClient

import (
	"encoding/json"
	"fmt"

	abstractmessage "github.com/SugaoTT/back/message"
)

type REMOVE_NETWORK_REQUEST struct {
	abstractmessage.AbstractMessage
	NetworkTopology string
}

type RemoveItems struct {
	Name string `json:"name"`
}

type Pod struct {
	Items []Items `json:"items"`
}

type RemoveNetworkTopology struct {
	Interface Pod `json:"pod"`
}

func NewREMOVE_NETWORK_REQUEST(inputMsg []byte) *REMOVE_NETWORK_REQUEST {
	msg := &REMOVE_NETWORK_REQUEST{}
	msg.AbstractMessage.MessageType = "REMOVE_NETWORK_REQUEST"

	//具象的パラメータをmsgに追加
	var ev REMOVE_NETWORK_REQUEST
	json.Unmarshal(inputMsg, &ev)
	msg.NetworkTopology = ev.NetworkTopology

	//ここでデコードとかしてしまおう

	//data := "{\"networkTopology\":" + ev.NetworkTopology + "\"}"
	fmt.Println("ここ: " + ev.NetworkTopology)

	//var networkTopology NetworkTopology
	//json.Unmarshal([]byte(ev.NetworkTopology), &networkTopology)

	return msg
}

func (msg *REMOVE_NETWORK_REQUEST) SetNetworkTopology(networkTopology string) {
	msg.NetworkTopology = networkTopology
}

func (msg *REMOVE_NETWORK_REQUEST) GetNetworkTopology() string {
	return msg.NetworkTopology
}
