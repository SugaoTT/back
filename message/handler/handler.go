package handler

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/websocket"

	k8s "github.com/SugaoTT/back/k8s"
	l2tpData "github.com/SugaoTT/back/manager"
	messageFromClient "github.com/SugaoTT/back/message/concrete/fromClient"
	messageToClient "github.com/SugaoTT/back/message/concrete/toClient"
	//abstractMessage "github.com/SugaoTT/back/message"
)

type MessageType struct {
	MessageType string `json:"messageType"`
}

type ConsoleMessage struct {
	TargetUUID     string `json:"target-uuid"`
	TargetNodeName string `json:"target-node-name"`
	Content        string `json:"content"`
}

/** 受け取ったメッセージに応じた処理を実施 */
func HandleMessage(ws *websocket.Conn, msg string) {
	//inputMsg := []byte(msg)

	fmt.Println(msg)

	var messageType MessageType
	json.Unmarshal([]byte(msg), &messageType)

	fmt.Println("メッセージのタイプ: " + messageType.MessageType)
	switch messageType.MessageType {
	case "LAUNCH_NETWORK_REQUEST":
		fmt.Println("LAUNCH_NETWORK_REQUESTが届きました")
		LAUNCH_NETWORK_REQUEST(ws, []byte(msg))
		break
	case "L2TP_INFO_REQUEST":
		fmt.Println("L2TP_INFO_REQUESTが届きました")
		L2TP_INFO_REQUEST(ws, []byte(msg))

		break
	case "L2TP_TUNNEL_ID_REQUEST":
		/*fmt.Println("L2TP_TUNNEL_ID_REQUESTが届きました")
		result := l2tpData.GenerateTunnelID()

		tunnel := message.NewL2TP_TUNNEL_ID()
		tunnel.SetTunnelID(strconv.Itoa(result))

		jsonData, err := json.Marshal(tunnel)
		if err != nil {
			fmt.Println("JSON変換エラー:", err)
			return
		}

		// JSONデータを文字列として表示
		fmt.Println(string(jsonData))

		websocket.Message.Send(ws, string(jsonData))*/
		break
	case "REMOVE_NETWORK_REQUEST":
		REMOVE_NETWORK_REQUEST(ws, []byte(msg))
		break
	case "console":
		var consoleMessage ConsoleMessage
		json.Unmarshal([]byte(msg), &consoleMessage)
		fmt.Println(consoleMessage.TargetUUID)

		//outputCommand := strings.Split(consoleMessage.Content, " ")

		//str := `vtysh -c "show ip route"`

		// Split by space but respect quoted strings as single element
		re := regexp.MustCompile(`[^\s"']+|"([^"]*)"|'([^']*)`)
		outputCommand := re.FindAllString(consoleMessage.Content, -1)

		// Remove quotes from quoted strings
		for i, part := range outputCommand {
			if strings.HasPrefix(part, `"`) && strings.HasSuffix(part, `"`) {
				outputCommand[i] = strings.Trim(part, `"`)
			}
		}

		fmt.Println(outputCommand)

		//fmt.Println(len(outputCommand))
		k8s.Pod_exec_new(ws, outputCommand, consoleMessage.TargetUUID)

		break
	}
}

func LAUNCH_NETWORK_REQUEST(ws *websocket.Conn, msg []byte) {

	/* ネットワークトポロジに関する構造体 */
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

	//fromClientのmsgを作成
	msgOf_LAUNCH_NETWORK_REQUEST := messageFromClient.NewLAUNCH_NETWORK_REQUEST(msg)
	fmt.Println(msgOf_LAUNCH_NETWORK_REQUEST)

	//fmt.Println("net: " + msgOf_LAUNCH_NETWORK_REQUEST.GetNetworkTopologies())

	var ev NetworkTopologies
	json.Unmarshal([]byte(msgOf_LAUNCH_NETWORK_REQUEST.GetNetworkTopologies()), &ev)

	fmt.Println(ev.PodList)

	for _, networkTopology := range ev.PodList {
		fmt.Println(networkTopology.PodType)
		var json, _ = json.Marshal(networkTopology)
		k8s.Pod_apply(string(json))
	}

	//ネットワークトポロジが受け取れたので，これを用いてk8sに発行するためのマニフェストを作成する

	//k8s.Pod_apply(msgOf_LAUNCH_NETWORK_REQUEST)

	time.Sleep(5000 * time.Millisecond)

	for _, networkTopology := range ev.PodList {
		status := k8s.Chack_Status(networkTopology.PodName[:8])
		if status {
			continue
		}

		for !status {
			time.Sleep(time.Millisecond * 500)
			status = k8s.Chack_Status(networkTopology.PodName[:8])
		}
	}

	//toClientのmsgを作成
	msgOf_LAUNCH_NETWORK_SUCCESS := messageToClient.NewLAUNCH_NETWORK_SUCCESS()

	jsonData, err := json.Marshal(msgOf_LAUNCH_NETWORK_SUCCESS)
	if err != nil {
		fmt.Println("JSON変換エラー:", err)
		return
	}

	fmt.Println(strings.Replace(string(jsonData), "\"", "", -1))

	// JSONデータを文字列として表示
	fmt.Println("sendToClient-JSON: " + string(jsonData))

	websocket.Message.Send(ws, string(jsonData))
}

func L2TP_INFO_REQUEST(ws *websocket.Conn, msg []byte) {

	//fromClientのmsgを作成
	msgOf_L2TP_INFO_REQUEST := messageFromClient.NewL2TP_INFO_REQUEST(msg)

	//fmt.Println(msgOf_L2TP_INFO_REQUEST)

	//l2tpに関する情報を生成
	sessionID := l2tpData.GenerateSessionID()
	tunnelID := l2tpData.GenerateTunnelID()

	//toClientのmsgを作成
	msgOf_L2TP_INFO := messageToClient.NewL2TP_INFO()
	msgOf_L2TP_INFO.SetSrcUUID(msgOf_L2TP_INFO_REQUEST.SrcUUID)
	msgOf_L2TP_INFO.SetSrcEthName(msgOf_L2TP_INFO_REQUEST.SrcEthName)
	msgOf_L2TP_INFO.SetDstUUID(msgOf_L2TP_INFO_REQUEST.DstUUID)
	msgOf_L2TP_INFO.SetDstEthName(msgOf_L2TP_INFO_REQUEST.DstEthName)

	msgOf_L2TP_INFO.SetSessionID(strconv.Itoa(sessionID))
	msgOf_L2TP_INFO.SetSrcTunnelID(strconv.Itoa(tunnelID))
	msgOf_L2TP_INFO.SetDstTunnelID(strconv.Itoa(tunnelID + 1))

	jsonData, err := json.Marshal(msgOf_L2TP_INFO)
	if err != nil {
		fmt.Println("JSON変換エラー:", err)
		return
	}
	// JSONデータを文字列として表示
	fmt.Println("sendToClient-JSON: " + string(jsonData))

	/*

		//l2tpMsg := message.NewL2TP_INFO(msg)
		l2tpMsg.SetSessionID(strconv.Itoa(sessionID))
		l2tpMsg.SetSrcTunnelID(strconv.Itoa(tunnelID))
		l2tpMsg.SetDstTunnelID(strconv.Itoa(tunnelID + 1))

		jsonData, err := json.Marshal(l2tpMsg)
		if err != nil {
			fmt.Println("JSON変換エラー:", err)
			return
		}

		// JSONデータを文字列として表示
		fmt.Println("sendToClient-JSON: " + string(jsonData))
	*/

	websocket.Message.Send(ws, string(jsonData))
}

func REMOVE_NETWORK_REQUEST(ws *websocket.Conn, msg []byte) {
	fmt.Println("REMOVE_NETWORK_REQUESTが呼び出されました")
	/* ネットワークトポロジに関する構造体 */
	type Items struct {
		Name string `json:"name"`
	}

	type Pod struct {
		Items []Items `json:"items"`
	}

	type NetworkTopology struct {
		Interface Pod `json:"pod"`
	}

	//fromClientのmsgを作成
	msgOf_REMOVE_NETWORK_REQUEST := messageFromClient.NewREMOVE_NETWORK_REQUEST(msg)

	var ev NetworkTopology
	json.Unmarshal([]byte(msgOf_REMOVE_NETWORK_REQUEST.GetNetworkTopology()), &ev)

	fmt.Println(ev.Interface.Items)
	i := 0

	//ここから削除処理をする
	for ; i < len(ev.Interface.Items); i++ {
		//fmt.Println(ev.Interface.Items[i].Name)
		k8s.Pod_delete(ev.Interface.Items[i].Name)
	}
}
