package plugin

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

func exploreService(configFlags *genericclioptions.ConfigFlags, clientset *kubernetes.Clientset, namespace string, name string, port int32) (tree Tree, err error) {
	service, err := clientset.CoreV1().Services(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return
	}
	for _, p := range service.Spec.Ports {
		if p.Port == port {
			info := fmt.Sprintf("%d", port)
			tree = NewTree("Service", fmt.Sprintf("%s.%s", service.Namespace, service.Name), &info)
			var endpoints []Tree
			endpoints, err = exploreEndpoints(configFlags, clientset, namespace, name, p.TargetPort.IntVal)
			for _, endpoint := range endpoints {
				tree.addChild(endpoint)
			}
		}
	}
	return
}
