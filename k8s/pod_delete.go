package k8s

import (
	"fmt"
	"os"
	"os/exec"
)

func Pod_delete(uuid string) {
	uuidPrefix := uuid[:8]
	yamlFilePath := uuidPrefix + ".yml"
	output, _ := exec.Command("kubectl", "delete", "-f", yamlFilePath).Output()
	fmt.Println(string(output))

	os.Remove(yamlFilePath)
}
