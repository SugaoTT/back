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

	// Podの定義
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "example-pod",
			Labels: map[string]string{
				"app": "example",
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "example-container",
					Image: "nginx",
				},
			},
		},
	}

	// Podの作成
	createdPod, err := clientset.CoreV1().Pods("default").Create(context.Background(), pod, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created Pod %q\n", createdPod.GetName())

	// Podの状態を確認
	for {
		pod, err := clientset.CoreV1().Pods("default").Get(context.Background(), createdPod.GetName(), metav1.GetOptions{})
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

	// Podの削除
	err = clientset.CoreV1().Pods("default").Delete(context.Background(), createdPod.GetName(), metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deleted Pod %q\n", createdPod.GetName())

}
