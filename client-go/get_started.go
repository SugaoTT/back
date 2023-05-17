package main

import (
    "context"
    "fmt"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "log"
    "os"
    "path/filepath"
)
 
func main() {
    // Kubeconfigのファイルパスを指定
    kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
    config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    if err != nil {
        log.Fatal(err)
    }
 
    // Kubeconfigを読み込む
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        log.Fatal(err)
    }
 
    // Pod一覧を呼び出す
    namespace := "kube-system"
    pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        log.Fatalln("failed to get pods:", err)
    }
 
    // print pods
    // pods.Items: []v1.Pod
    for i, pod := range pods.Items {
        fmt.Printf("[%d] %s\n", i, pod.GetName())
    }


	
}