package k8s

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"text/template"

	messageFromClient "github.com/SugaoTT/back/message/concrete/fromClient"
)

func Pod_apply(msgOf_LAUNCH_NETWORK_REQUEST *messageFromClient.LAUNCH_NETWORK_REQUEST) {

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
	json.Unmarshal([]byte(msgOf_LAUNCH_NETWORK_REQUEST.GetNetworkTopology()), &ev)

	// Pod 名とコンテナイメージを設定する変数
	//containerImage := "frrouting/frr:v8.1.0"

	fmt.Println("テンプレート作成しています")

	var yamlTemplate string

	i := 0
	j := 0
	k := 0
	interfaces := ev.Interface.Items
	uuidPrefix := ev.PodName[:8]
	callConnectCNI := "connect-" + uuidPrefix

	fmt.Println(uuidPrefix)

	switch ev.PodType {
	case "Router":

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

		for ; k < len(interfaces); k++ {
			yamlTemplate += fmt.Sprintf(`
          {
            "name": "%s",
            "target-pod-name": "%s",
            "target-pod-nic": "%s",
            "self-tunnel-id": "%s",
            "target-tunnel-id": "%s",
            "session-id": "%s"
          }`, ev.Interface.Items[k].Name, ev.Interface.Items[k].TargetPodName[:8], ev.Interface.Items[k].TargetPodNic, ev.Interface.Items[k].SelfTunnelID, ev.Interface.Items[k].TargetTunnelID, ev.Interface.Items[k].SessionID)

			if k != len(interfaces)-1 {
				yamlTemplate += `,`
			}
		}
		yamlTemplate += `
        ]
      }
    }`

		for ; j < len(interfaces); j++ {
			bridgeName := uuidPrefix + "-net" + strconv.Itoa(j+1)
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

		// yamlTemplate = `

		// `

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
		yamlTemplate += `
      {"name": "{{ .Connect }}"}
    ]'
spec: 
  containers:
  - name: {{ .Name }}
    image: frrouting/frr:v8.1.0
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
                - "ip link set eth0 down"`

		//       {"name": "r1c-net1"},
		//       {"name": "r1c-net2"},
		//       {"name": {{ .Connect }}}
		//     ]'
		// spec:
		//   nodeName: sugao-k8s-worker3
		//   containers:
		//   - name: {{ .Name }}
		//     image: frrouting/frr:v8.1.0
		//     command:
		//     - /sbin/init
		//     securityContext:
		//       privileged: true
		//     lifecycle:
		//           postStart:
		//             exec:
		//               command:
		//                 - sh
		//                 - -c
		//                 - "ip link set eth0 down"
		// `

		// 		yamlTemplate = `
		// apiVersion: v1
		// kind: Pod
		// metadata:
		//   name: {{ .Name }}
		//   labels:
		//     app: frr
		// spec:
		//   containers:
		//   - name: example-container
		//     image: {{ .Image }}
		//     ports:
		//     - containerPort: 80
		// `
	case "Switch":
	case "Host":
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

	fmt.Println(string(yamlBuffer.Bytes()))

	// // kubectl コマンドを実行する
	// cmd := exec.Command("kubectl", "apply", "-f", "-")
	// cmd.Stdin = &yamlBuffer
	// output, err := cmd.CombinedOutput()
	// if err != nil {
	// 	fmt.Println("Error deploying Pod: ", err)
	// 	return
	// }
	// fmt.Println(string(output))

	// time.Sleep(500 * time.Millisecond)

	// l := 0
	// // kubectl コマンドを実行する
	// output2, _ := exec.Command("kubectl", "delete", "pods", ev.PodName).Output()
	// fmt.Println(string(output2))
	// for ; l < len(interfaces); l++ {
	// 	output3, _ := exec.Command("kubectl", "delete", "network-attachment-definitions", ev.PodName+"-net"+strconv.Itoa(l+1)).Output()
	// 	fmt.Println(string(output3))
	// }
	// output4, _ := exec.Command("kubectl", "delete", "network-attachment-definitions", "connect-"+ev.PodName).Output()
	// fmt.Println(string(output4))

}
