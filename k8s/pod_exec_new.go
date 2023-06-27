package k8s

import (
	"fmt"
	"path/filepath"
	"strings"

	"golang.org/x/net/websocket"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

type myWriter struct {
	ws *websocket.Conn
}

func (mw *myWriter) Write(p []byte) (n int, err error) {
	myFunction(p) // pを自分の関数に渡す
	websocket.Message.Send(mw.ws, strings.TrimSpace(string(p)))
	return len(p), nil
}

func myFunction(p []byte) {
	// pを処理する
	fmt.Print(string(p))
}

func Pod_exec_new(ws *websocket.Conn, outputCommand []string, uuid string) {

	uuidPrefix := uuid[:8]

	defer websocket.Message.Send(ws, strings.TrimSpace("EXEC COMPLETE"))

	//kubeconfig := os.Getenv("KUBECONFIG")
	kubeconfig := filepath.Join("/Users", "sugaott", "school", "study", "code", "k8s", "kubectl", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		//	panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		//	panic(err)
	}

	//podName := "r1"
	namespace := "default"

	execReq := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(uuidPrefix).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Command: outputCommand,
			Stdout:  true,
			Stderr:  true,
		}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(config, "POST", execReq.URL())
	if err != nil {
		//	panic(err)
	}

	mw := &myWriter{ws: ws}
	err = executor.Stream(remotecommand.StreamOptions{
		Stdout: mw,
		Stderr: mw,
	})
	if err != nil {
		//	panic(err)
	}
}
