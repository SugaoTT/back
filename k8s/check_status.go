package k8s

import (
	"context"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func Chack_Status(uuid string) bool {
	kubeconfig := filepath.Join("/Users", "sugaott", "school", "study", "code", "k8s", "kubectl", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	//podName := "r1"        // replace with your pod name
	namespace := "default" // replace with your namespace

	//@TODO: podがnot foundの時になんかk8sの調子が悪くなる
	pod, err := clientset.CoreV1().Pods(namespace).Get(context.TODO(), uuid, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	return pod.Status.Phase == corev1.PodRunning
}
