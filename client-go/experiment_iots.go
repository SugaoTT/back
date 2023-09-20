package main

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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

	const sets = 25 // 作成するDeploymentのセット数
	//const targetNode = "sugao-k8s-worker2" // Podをスケジュールするノードの名前

	for i := 20; i < sets; i++ {
		//createDeploymentIOTS_SINGLENODE(clientset, "routers", "frrouting/frr:v8.1.0", i, 4, targetNode)            // コンテナAのデプロイメント
		//createDeploymentIOTS_SINGLENODE(clientset, "switches", "openshift/openvswitch:v3.9.0", i, 4, targetNode)   // コンテナBのデプロイメント
		//createDeploymentIOTS_SINGLENODE(clientset, "hosts", "sugaott/sugaott-ubuntu-focal:1.4", i, 12, targetNode) // コンテナCのデプロイメント

		createDeploymentIOTS_MULTINODE(clientset, "routers", "frrouting/frr:v8.1.0", i, 4)            // コンテナAのデプロイメント
		createDeploymentIOTS_MULTINODE(clientset, "switches", "openshift/openvswitch:v3.9.0", i, 4)   // コンテナBのデプロイメント
		createDeploymentIOTS_MULTINODE(clientset, "hosts", "sugaott/sugaott-ubuntu-focal:1.4", i, 12) // コンテナCのデプロイメント

	}
}

func createDeploymentIOTS_SINGLENODE(clientset *kubernetes.Clientset, containerName, containerImage string, setNumber int, replicas int32, targetNode string) {
	deploymentName := fmt.Sprintf("%s-deployment-set%d", containerName, setNumber)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: "default",
			Labels: map[string]string{ // この部分を追加
				"app": containerName,
				"set": fmt.Sprintf("%d", setNumber),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": containerName,
					"set": fmt.Sprintf("%d", setNumber),
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": containerName,
						"set": fmt.Sprintf("%d", setNumber),
					},
				},
				Spec: corev1.PodSpec{
					HostNetwork: true,
					NodeSelector: map[string]string{
						"kubernetes.io/hostname": targetNode,
					},
					Containers: []corev1.Container{
						{
							Name:  containerName,
							Image: containerImage,
						},
					},
				},
			},
		},
	}

	// Deploymentの作成
	_, err := clientset.AppsV1().Deployments("default").Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(fmt.Errorf("failed to create deployment %s: %v", deploymentName, err))
	}

	fmt.Printf("Deployment %s created successfully!\n", deploymentName)

	time.Sleep(time.Millisecond * 5000)
}

func createDeploymentIOTS_MULTINODE(clientset *kubernetes.Clientset, containerName, containerImage string, setNumber int, replicas int32) {
	deploymentName := fmt.Sprintf("%s-deployment-set%d", containerName, setNumber)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: "default",
			Labels: map[string]string{ // この部分を追加
				"app": containerName,
				"set": fmt.Sprintf("%d", setNumber),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": containerName,
					"set": fmt.Sprintf("%d", setNumber),
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": containerName,
						"set": fmt.Sprintf("%d", setNumber),
					},
				},
				Spec: corev1.PodSpec{
					HostNetwork: true,
					Containers: []corev1.Container{
						{
							Name:  containerName,
							Image: containerImage,
						},
					},
				},
			},
		},
	}

	// Deploymentの作成
	_, err := clientset.AppsV1().Deployments("default").Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(fmt.Errorf("failed to create deployment %s: %v", deploymentName, err))
	}

	fmt.Printf("Deployment %s created successfully!\n", deploymentName)

	time.Sleep(time.Millisecond * 5000)
}
