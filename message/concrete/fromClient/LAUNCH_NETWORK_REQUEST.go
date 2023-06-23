package messageFromClient

import (
	"encoding/json"
	"fmt"

	abstractmessage "github.com/SugaoTT/back/message"
)

type LAUNCH_NETWORK_REQUEST struct {
	abstractmessage.AbstractMessage
	NetworkTopologies string
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

// NetworkTopology のスライスを保持する新しい構造体
type NetworkTopologies struct {
	PodList []NetworkTopology `json:"podList"`
}

type OuterRequest struct {
	MessageType     string `json:"messageType"`
	NetworkTopology string `json:"networkTopology"` // 注意: これはJSON文字列として格納されています
}

func NewLAUNCH_NETWORK_REQUEST(inputMsg []byte) *LAUNCH_NETWORK_REQUEST {
	msg := &LAUNCH_NETWORK_REQUEST{}
	msg.AbstractMessage.MessageType = "LAUNCH_NETWORK_REQUEST"

	fmt.Println(inputMsg)

	var outerRequest OuterRequest
	json.Unmarshal([]byte(inputMsg), &outerRequest)

	var ev NetworkTopologies
	json.Unmarshal([]byte(outerRequest.NetworkTopology), &ev)

	fmt.Println("ev.PodList:")
	fmt.Println(ev.PodList)

	msg.NetworkTopologies = outerRequest.NetworkTopology

	// //具象的パラメータをmsgに追加
	// var ev LAUNCH_NETWORK_REQUEST
	// json.Unmarshal(inputMsg, &ev)
	// msg.NetworkTopologies = ev.NetworkTopologies

	// fmt.Println("ev.NetworkTopologies: " + ev.NetworkTopologies)

	//ここでデコードとかしてしまおう

	//data := "{\"networkTopology\":" + ev.NetworkTopology + "\"}"
	//fmt.Println("ここ: " + ev.NetworkTopology)

	//var networkTopology NetworkTopology
	//json.Unmarshal([]byte(ev.NetworkTopology), &networkTopology)

	return msg
}

func (msg *LAUNCH_NETWORK_REQUEST) SetNetworkTopologies(networkTopologies string) {
	msg.NetworkTopologies = networkTopologies
}

func (msg *LAUNCH_NETWORK_REQUEST) GetNetworkTopologies() string {
	return msg.NetworkTopologies
}
