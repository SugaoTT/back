// package main

// import (
// 	"fmt"
// 	"os"
// 	"path/filepath"
// 	"time"

// 	corev1 "k8s.io/api/core/v1"
// 	"k8s.io/client-go/kubernetes"
// 	"k8s.io/client-go/kubernetes/scheme"
// 	"k8s.io/client-go/rest"
// 	"k8s.io/client-go/tools/clientcmd"
// 	"k8s.io/client-go/tools/remotecommand"
// )

// type myWriter struct{}

// func (mw *myWriter) Write(p []byte) (n int, err error) {
// 	myFunction(p) // pを自分の関数に渡す
// 	return len(p), nil
// }

// func myFunction(p []byte) {
// 	// pを処理する
// 	fmt.Print(string(p))
// }

// func main() {
// 	//kubeconfig := os.Getenv("KUBECONFIG")
// 	kubeconfig := filepath.Join("/Users", "sugaott", "school", "study", "code", "k8s", "kubectl", "config")
// 	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
// 	if err != nil {
// 		panic(err)
// 	}

// 	clientset, err := kubernetes.NewForConfig(config)
// 	if err != nil {
// 		panic(err)
// 	}

// 	podName := "klish"
// 	namespace := "default"

// 	execCommand(clientset, config, podName, namespace, []string{"bash", "-c", "CLISH_PATH=xml-examples/clish bin/clish --lockless"})

// 	time.Sleep(5 * time.Second)

// 	execCommand(clientset, config, podName, namespace, []string{"ping", "-c", "4", "127.0.0.1"})
// }

// func execCommand(clientset *kubernetes.Clientset, config *rest.Config, podName string, namespace string, command []string) {
// 	execReq := clientset.CoreV1().RESTClient().Post().
// 		Resource("pods").
// 		Name(podName).
// 		Namespace(namespace).
// 		SubResource("exec").
// 		VersionedParams(&corev1.PodExecOptions{
// 			Command: command,
// 			Stdout:  true,
// 			Stderr:  true,
// 			TTY:     true,
// 		}, scheme.ParameterCodec)

// 	executor, err := remotecommand.NewSPDYExecutor(config, "POST", execReq.URL())
// 	if err != nil {
// 		panic(err)
// 	}

// 	mw := &myWriter{}
// 	err = executor.Stream(remotecommand.StreamOptions{
// 		Stdout: os.Stdout,
// 		Stderr: mw,
// 	})
// 	if err != nil {
// 		panic(err)
// 	}
// }
