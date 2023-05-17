package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var labelSelector, kubeconfigPath string
	flag.StringVar(&labelSelector, "label-selector", "", "Label selector for the pods")
	flag.StringVar(&kubeconfigPath, "kubeconfig", "", "Path to the kubeconfig file")
	flag.Parse()

	if labelSelector == "" {
		fmt.Println("Label selector must be provided with the -label-selector flag")
		return
	}

	kubeconfig := filepath.Join("/Users", "sugaott", "school", "study", "code", "k8s", "kubectl", "config")
	config, _ := clientcmd.BuildConfigFromFlags("", kubeconfig)
	// Kubernetesクライアントを作成する
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	// Podのレディ状態をポーリングし、すべてがレディ状態になるのを待つ
	err = waitUntilAllPodsReady(clientset, labelSelector)
	if err != nil {
		fmt.Printf("Error waiting for all pods to be ready: %v\n", err)
		return
	}

	// すべてのPodがレディ状態になったことを示すログメッセージを出力
	fmt.Println("All pods are ready")

}

func waitUntilAllPodsReady(clientset *kubernetes.Clientset, labelSelector string) error {
	var ready bool

	for !ready {
		ready = true

		// Label Selectorを持つすべてのPodを取得
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{
			LabelSelector: labelSelector,
		})
		fmt.Println(pods)
		if err != nil {
			return err
		}
		// 各Podのレディ状態を確認し、すべてがレディ状態になるのを待つ
		for _, pod := range pods.Items {
			if !isPodReady(&pod) {
				ready = false
				break
			}
		}

		if !ready {
			// すべてのPodがレディ状態になっていない場合は、一定時間待機して再試行します
			time.Sleep(5 * time.Second)
		}
	}

	return nil
}

func isPodReady(pod *corev1.Pod) bool {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == corev1.PodReady && condition.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}
