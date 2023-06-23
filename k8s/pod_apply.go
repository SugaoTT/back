package k8s

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"text/template"

	"io/ioutil"
	//handler "github.com/SugaoTT/back/message/handler"
)

func Pod_apply(inputJson string) {

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

	var ev NetworkTopology
	json.Unmarshal([]byte(inputJson), &ev)

	// Pod 名とコンテナイメージを設定する変数
	//containerImage := "frrouting/frr:v8.1.0"

	fmt.Println("テンプレート作成しています")

	var yamlTemplate string

	i := 0
	interfaces := ev.Interface.Items
	uuidPrefix := ev.PodName[:8]
	callConnectCNI := "connect-" + uuidPrefix

	fmt.Println(uuidPrefix)

	yamlTemplate += `
apiVersion: "k8s.cni.cncf.io/v1"
kind: NetworkAttachmentDefinition
metadata:
  name: {{ .Connect }}
spec:
  config: >
    {
      "cniVersion": "0.3.0",
      "type": "connect",
      "pod-name": "{{ .Name }}",
      "interface": {
        "items": [`

	for ; i < len(interfaces); i++ {
		yamlTemplate += fmt.Sprintf(`
          {
            "name": "%s",
            "target-pod-name": "%s",
            "target-pod-nic": "%s",
            "self-tunnel-id": "%s",
            "target-tunnel-id": "%s",
            "session-id": "%s"
          }`, ev.Interface.Items[i].Name, ev.Interface.Items[i].TargetPodName[:8], ev.Interface.Items[i].TargetPodNic, ev.Interface.Items[i].SelfTunnelID, ev.Interface.Items[i].TargetTunnelID, ev.Interface.Items[i].SessionID)

		if i != len(interfaces)-1 {
			yamlTemplate += `,`
		}
	}
	i = 0
	yamlTemplate += `
        ]
      }
    }`

	for ; i < len(interfaces); i++ {
		bridgeName := uuidPrefix + "-net" + strconv.Itoa(i+1)
		yamlTemplate += fmt.Sprintf(`
---
apiVersion: "k8s.cni.cncf.io/v1"
kind: NetworkAttachmentDefinition
metadata:
  name: %s
spec:
  config: >
    {
      "cniVersion": "0.3.0",
      "type": "bridge",
      "bridge": "%s",
      "ipam": {
      }
    }`, bridgeName, bridgeName)
	}
	i = 0

	yamlTemplate += `
---
apiVersion: v1
kind: Pod
metadata: 
  name: {{ .Name }}
  annotations:
    k8s.v1.cni.cncf.io/networks: '[`

	fmt.Println(len(interfaces))

	for ; i < len(interfaces); i++ {
		interfaceName := uuidPrefix + "-net" + strconv.Itoa(i+1)
		yamlTemplate += fmt.Sprintf(`
      {"name": "%s"},`, interfaceName)
	}
	i = 0

	yamlTemplate += `
      {"name": "{{ .Connect }}"}
    ]'`

	switch ev.PodType {
	case "Router":

		yamlTemplate += `
spec: 
  containers:
  - name: {{ .Name }}
    image: frrouting/frr:v8.1.0
    securityContext:
      privileged: true
    lifecycle:
          postStart:
            exec:
              command:
                - sh
                - -c
                - "ip link set eth0 down"`

	case "Switch":
		yamlTemplate += `
spec: 
  containers:
  - name: {{ .Name }}
    image: openshift/openvswitch:v3.9.0
    command:
    - /sbin/init
    securityContext:
      privileged: true
    lifecycle:
          postStart:
            exec:
              command:
                - sh
                - -c
                - "/usr/share/openvswitch/scripts/ovs-ctl start;ip link set eth0 down`

		yamlTemplate += fmt.Sprintf(`;ovs-vsctl add-br %s`, uuidPrefix)

		for ; i < len(interfaces); i++ {
			netName := "net" + strconv.Itoa(i+1)
			yamlTemplate += fmt.Sprintf(`;ovs-vsctl add-port %s %s`, uuidPrefix, netName)
		}
		i = 0

		yamlTemplate += `"`
	case "Host":
		yamlTemplate += `
spec: 
  containers:
  - name: {{ .Name }}
    image: sugaott/sugaott-ubuntu-focal:1.4
    command: ["bash", "-c", "sleep infinity"]
    securityContext:
      privileged: true
    lifecycle:
          postStart:
            exec:
              command:
                - sh
                - -c
                - "ip link set eth0 down"`
	}

	// テンプレートを定義する

	fmt.Println(yamlTemplate)
	// テンプレートを適用する
	tmpl, err := template.New("yaml").Parse(yamlTemplate)
	if err != nil {
		fmt.Println("Error parsing YAML template: ", err)
		return
	}

	var yamlBuffer bytes.Buffer
	err = tmpl.Execute(&yamlBuffer, struct{ Name, Connect string }{uuidPrefix, callConnectCNI})
	if err != nil {
		fmt.Println("Error executing YAML template: ", err)
		return
	}

	//sampleString := []byte("This is a string")
	//ioutil.WriteFile("/Users/sugaott/school/study/code/back/tmp/"+uuidPrefix+".yml", yamlBuffer.Bytes(), 0644)
	ioutil.WriteFile(uuidPrefix+".yml", yamlBuffer.Bytes(), 0644)

	fmt.Println(string(yamlBuffer.Bytes()))

	//kubectl コマンドを実行する
	cmd := exec.Command("kubectl", "apply", "-f", "-")
	cmd.Stdin = &yamlBuffer
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error deploying Pod: ", err)
		return
	}
	fmt.Println(string(output))

	// time.Sleep(500 * time.Millisecond)

	// i := 0
	// // kubectl コマンドを実行する
	// output2, _ := exec.Command("kubectl", "delete", "pods", ev.PodName).Output()
	// fmt.Println(string(output2))
	// for ; i < len(interfaces); i++ {
	// 	output3, _ := exec.Command("kubectl", "delete", "network-attachment-definitions", ev.PodName+"-net"+strconv.Itoa(i+1)).Output()
	// 	fmt.Println(string(output3))
	// }
	// output4, _ := exec.Command("kubectl", "delete", "network-attachment-definitions", "connect-"+ev.PodName).Output()
	// fmt.Println(string(output4))

}
