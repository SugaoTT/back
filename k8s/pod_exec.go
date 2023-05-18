package k8s

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"

	"golang.org/x/net/websocket"
)

func Pod_exec(ws *websocket.Conn, outputCommand []string) {
	// Kubeconfigのファイルパスを指定
	kubeconfig := filepath.Join("/Users", "sugaott", "school", "study", "code", "k8s", "kubectl", "config")
	config, _ := clientcmd.BuildConfigFromFlags("", kubeconfig)
	//client, _ := corev1client.NewForConfig(config)

	var isCommandComplate bool = false

	var stdout, stderr bytes.Buffer
	var prevStdoutStr = ""

	//バッファ監視用のgoルーチン起動(並列処理)
	go func() {
		for {
			if stdout.String() != "" {
				if prevStdoutStr == "" || prevStdoutStr != stdout.String() {
					var output string = strings.Replace(stdout.String(), prevStdoutStr, "", -1)
					fmt.Print(output)
					err := websocket.Message.Send(ws, strings.TrimSpace(output))
					//err := websocket.Message.Send(ws, "    "+output)
					if err != nil {
					}
					prevStdoutStr = stdout.String()
					time.Sleep(time.Millisecond * 5)
				}
			}
		}
	}()

	//Podにexecしてコマンド発行処理を実施
	err := ExecInPod(config, &stdout, &stderr, "default", "h1", outputCommand, &isCommandComplate)
	if err != nil {
		//いい感じにエラー処理(print)を記述
	}

	//goルーチンの最後の処理が終わるまでwait
	time.Sleep(time.Millisecond * 200)
	websocket.Message.Send(ws, strings.TrimSpace("EXEC COMPLETE"))
	// if isCommandComplate {
	// 	websocket.Message.Send(ws, strings.TrimSpace("EXEC COMPLETE"))
	// }

}

func ExecInPod(config *rest.Config, stdout, stderr *bytes.Buffer, namespace, podName string, outputCommand []string, isCommandComplate *bool) error {
	k8sCli, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	//コンテナに発行するコマンド
	/*cmd := []string{
		"ping",
		"-c",
		"3",
		"127.0.0.1",
	}*/

	const tty = false
	req := k8sCli.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).SubResource("exec") //.Param("container", containerName)
	req.VersionedParams(
		&v1.PodExecOptions{
			Command: outputCommand,
			Stdin:   false,
			Stdout:  true,
			Stderr:  true,
			TTY:     tty,
		},
		scheme.ParameterCodec,
	)

	//コンテナにコマンド発行
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return err
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: stdout,
		Stderr: stderr,
	})
	if err != nil {
		return err
	}

	*isCommandComplate = true

	return err
}
