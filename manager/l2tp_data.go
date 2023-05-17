package manager

import (
	"sync"
)

type Data struct {
	sessionID int
	tunnelID  int
}

var (
	data  Data
	mutex sync.Mutex
)

func GenerateSessionID() int {
	mutex.Lock()
	defer mutex.Unlock()

	data.sessionID++
	//data.tunnelID++

	return data.sessionID
}

func GenerateTunnelID() int {
	mutex.Lock()
	defer mutex.Unlock()

	data.tunnelID += 2
	//data.tunnelID++

	return data.tunnelID
}
