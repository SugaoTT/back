package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

type myWriter struct {
	writer io.Writer
}

func (m *myWriter) Write(p []byte) (int, error) {
	processOutput(p) // Pass the data to your function here.
	return len(p), nil
}

func processOutput(data []byte) {
	fmt.Printf("%s", data)
}

func main() {

	kubeconfig := filepath.Join("/Users", "sugaott", "school", "study", "code", "k8s", "kubectl", "config")

	// Use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, _ := kubernetes.NewForConfig(config)

	// Specify the Pod and the command
	podName := "klish"
	command := []string{"bash", "-c", "CLISH_PATH=xml-examples/clish bin/clish --lockless"}

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace("default").
		SubResource("exec")
	option := &corev1.PodExecOptions{
		Command: command,
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}
	fmt.Println(command)
	req.VersionedParams(
		option,
		scheme.ParameterCodec,
	)

	// Execute the command
	exec, _ := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	exec.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: &myWriter{},
		//Stderr: os.Stderr,
		Tty: true,
	})
}
