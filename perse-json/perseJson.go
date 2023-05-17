package persejson

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type MessageType struct {
	MessageType string `json:"message-type"`
}

type ConsoleMessage struct {
	TargetUUID     string `json:"target-uuid"`
	TargetNodeName string `json:"target-node-name"`
	Content        string `json:"content"`
}

func main() {
	file, err := ioutil.ReadFile("./consoleMsg.json")
	if err != nil {
		log.Println("ReadError: ", err)
		os.Exit(1)
	}
	var messageType MessageType
	json.Unmarshal(file, &messageType)

	fmt.Println(messageType.MessageType)
	switch messageType.MessageType {
	case "launch-network":
		break
	case "console":
		var consoleMessage ConsoleMessage
		json.Unmarshal(file, &consoleMessage)
		fmt.Println(consoleMessage.TargetUUID)
		break
	}
}
