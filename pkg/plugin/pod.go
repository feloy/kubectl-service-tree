package plugin

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

func explorePod(configFlags *genericclioptions.ConfigFlags, clientset *kubernetes.Clientset, namespace string, name string, port int32) (tree Tree, found bool, err error) {
	var pod *v1.Pod
	pod, err = clientset.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return
	}
	if len(pod.Spec.Containers) == 1 {
		return
	}

	for _, container := range pod.Spec.Containers {
		for _, p := range container.Ports {
			if p.ContainerPort == port {
				found = true
				tree = NewTree("Container", container.Name, nil)
				return
			}
		}
	}

	return
}
