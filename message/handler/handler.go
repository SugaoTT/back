package handler

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/websocket"

	k8s "github.com/SugaoTT/back/k8s"
	l2tpData "github.com/SugaoTT/back/manager"
	message "github.com/SugaoTT/back/message/concrete/toClient"
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
	inputMsg := []byte(msg)

	fmt.Println(msg)

	var messageType MessageType
	json.Unmarshal(inputMsg, &messageType)

	fmt.Println(messageType.MessageType)
	switch messageType.MessageType {
	case "LAUNCH_NETWORK":
		break
	case "L2TP_SESSION_ID_REQUEST":
		fmt.Println("L2TP_SESSION_ID_REQUESTが届きました")
		result := l2tpData.GenerateSessionID()
		session := message.NewL2TP_SESSION_ID(inputMsg)
		session.SetSessionID(strconv.Itoa(result))

		jsonData, err := json.Marshal(session)
		if err != nil {
			fmt.Println("JSON変換エラー:", err)
			return
		}

		// JSONデータを文字列として表示
		fmt.Println("sendToClient-JSON: " + string(jsonData))

		websocket.Message.Send(ws, string(jsonData))

		break
	case "L2TP_TUNNEL_ID_REQUEST":
		fmt.Println("L2TP_TUNNEL_ID_REQUESTが届きました")
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

		websocket.Message.Send(ws, string(jsonData))
		break
	case "console":
		var consoleMessage ConsoleMessage
		json.Unmarshal(inputMsg, &consoleMessage)
		fmt.Println(consoleMessage.TargetUUID)

		outputCommand := strings.Split(consoleMessage.Content, " ")
		//fmt.Println(len(outputCommand))
		k8s.Pod_exec(ws, outputCommand)

		break
	}
}
