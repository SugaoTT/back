package main

import (
	"context"
	"fmt"
	"log"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	//"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func main() {
	// Kubeconfigのファイルパスを指定
	kubeconfig := filepath.Join("/Users", "sugaott", "school", "study", "code", "k8s", "kubectl", "config")
	config, _ := clientcmd.BuildConfigFromFlags("", kubeconfig)
	//client, _ := corev1client.NewForConfig(config)

	// Kubernetesクライアントを作成する
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	// Podの状態を確認
	for {
		pod, err := clientset.CoreV1().Pods("default").Get(context.Background(), "centos1", metav1.GetOptions{})
		if err != nil {
			panic(err)
		}
		if pod.Status.Phase == v1.PodRunning {
			fmt.Println("Pod is running")
			break
		}
		fmt.Println("Waiting for Pod to start")
		time.Sleep(2 * time.Second)
	}

}
