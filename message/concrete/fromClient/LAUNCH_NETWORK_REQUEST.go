package messageFromClient

import (
	"encoding/json"

	abstractmessage "github.com/SugaoTT/back/message"
)

type LAUNCH_NETWORK_REQUEST struct {
	abstractmessage.AbstractMessage
	NetworkTopology string
}

type Items struct {
	Name           string `json:"name"`
	TargetPodName  string `json:"target-pod-name"`
	TargetPodNic   string `json:"target-pod-nic"`
	SelfTunnelID   string `json:"self-tunnel-id"`
	TargetTunnelID string `json:"target-tunnel-id"`
	SessionID      string `json:"session-id"`
}

type Interface struct {
	Items []Items `json:"items"`
}

type NetworkTopology struct {
	PodName   string    `json:"pod-name"`
	PodType   string    `json:"pod-type"`
	Interface Interface `json:"interface"`
}

func NewLAUNCH_NETWORK_REQUEST(inputMsg []byte) *LAUNCH_NETWORK_REQUEST {
	msg := &LAUNCH_NETWORK_REQUEST{}
	msg.AbstractMessage.MessageType = "LAUNCH_NETWORK_REQUEST"

	//具象的パラメータをmsgに追加
	var ev LAUNCH_NETWORK_REQUEST
	json.Unmarshal(inputMsg, &ev)
	msg.NetworkTopology = ev.NetworkTopology

	//ここでデコードとかしてしまおう

	//data := "{\"networkTopology\":" + ev.NetworkTopology + "\"}"
	//fmt.Println("ここ: " + ev.NetworkTopology)

	//var networkTopology NetworkTopology
	//json.Unmarshal([]byte(ev.NetworkTopology), &networkTopology)

	return msg
}

func (msg *LAUNCH_NETWORK_REQUEST) SetNetworkTopology(networkTopology string) {
	msg.NetworkTopology = networkTopology
}

func (msg *LAUNCH_NETWORK_REQUEST) GetNetworkTopology() string {
	return msg.NetworkTopology
}
