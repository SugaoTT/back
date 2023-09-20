package main

import (
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

func main() {
	kubeconfig := filepath.Join("/Users", "sugaott", "school", "study", "code", "k8s", "kubectl", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	podName := "klish"     // 替换为你要与之交互的 Pod 名称
	namespace := "default" // 替换为 Pod 所在的 namespace 名称

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec")

	req.VersionedParams(&corev1.PodExecOptions{
		Command: []string{"/bin/sh"}, // Here we open a shell in the container
		//	Container: "my-container",
		Stdin:  true,
		Stdout: true,
		Stderr: true,
		TTY:    false,
	}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		panic(err.Error())
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin, // Now we pass the os.Stdin to interactively send inputs to the container.
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    false, // Enable TTY
	})
	if err != nil {
		panic(err.Error())
	}
}
