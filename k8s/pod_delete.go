package k8s

import (
	"fmt"
	"os/exec"
)

func Pod_delete(deleteCommand string) {

	baseCommand := "kubectl delete"
	fullCommand := baseCommand + " " + deleteCommand

	cmd, _ := exec.Command("bash", "-c", fullCommand).Output()
	fmt.Println(string(cmd))

	//uuidPrefix := uuid[:8]
	//yamlFilePath := uuidPrefix + ".yml"
	// fmt.Println("出力:"+"kubectl", "delete", deleteCommand)
	// output, _ := exec.Command("kubectl", "delete", deleteCommand).Output()
	// fmt.Println(string(output))

	//os.Remove(yamlFilePath)
}
