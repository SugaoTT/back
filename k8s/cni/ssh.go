package main

import (
	//k8s "github.com/SugaoTT/back/k8s"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	ls, _ := exec.Command("ssh", "-i", "/home/ubuntu/.ssh/id_ed25519", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", "ubuntu@192.168.0.224", "kubectl", "get", "pods", "-o", "wide", "|", "grep", "-w", "h1", "|", "awk", "'{print $7}'").Output()
	ls = []byte(strings.TrimRight(string(ls), "\n"))
	fmt.Println(string(ls))

	max, err := calcL2tpEth()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	// 最大値を表示します。
	fmt.Printf("l2tpethの最大値は%dです。\n", max)

	selfNodeIP, _ := exec.Command("ssh", "-i", "/home/ubuntu/.ssh/id_ed25519", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", "ubuntu@192.168.0.224", "kubectl", "get", "node", "-o", "wide", "|", "grep", "-w", "sugao-k8s-worker1", "|", "awk", "'{print $6}'").Output()
	selfNodeIP = []byte(strings.TrimRight(string(selfNodeIP), "\n"))

	fmt.Println(string(selfNodeIP))

	str := "192.168.0.210"
	if string(selfNodeIP) == str {
		fmt.Println("OK")
	}
}

func calcL2tpEth() (int, error) {
	// l2tpethの最大値を検索するために、ip linkコマンドを実行します。
	output, err := exec.Command("ssh", "-i", "/home/ubuntu/.ssh/id_ed25519", "-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null", "ubuntu@192.168.0.210", "ip link show | grep l2tpeth | awk '{print $2}' | sed 's/l2tpeth//g'").Output()
	if err != nil {
		return 0, fmt.Errorf("failed to execute command: %s", err)
	}

	// 出力から最大値を抽出します。
	max := 0
	numbers := strings.Split(string(output), "\n")
	for _, number := range numbers {
		if len(number) > 0 {
			// インターフェース名の最後に":"が付いている場合、それを取り除きます。
			if strings.HasSuffix(number, ":") {
				number = number[:len(number)-1]
			}
			n, err := strconv.Atoi(number)
			if err == nil && n > max {
				max = n
			}
		}
	}

	return max, nil
}
