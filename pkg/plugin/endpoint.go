package plugin

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

func exploreEndpoints(configFlags *genericclioptions.ConfigFlags, clientset *kubernetes.Clientset, namespace string, name string, port int32) (trees []Tree, err error) {
	trees = []Tree{}
	endpoints, err := clientset.CoreV1().Endpoints(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return
	}
	for _, endpoint := range endpoints.Subsets {
		for _, address := range endpoint.Addresses {
			ref := address.TargetRef
			info := fmt.Sprintf("%s:%d", address.IP, port)
			tree := NewTree(ref.Kind, fmt.Sprintf("%s.%s", ref.Namespace, ref.Name), &info)
			tree.setStatus(true)
			if ref.Kind == "Pod" {
				if container, found, _ := explorePod(configFlags, clientset, namespace, ref.Name, port); found {
					tree.addChild(container)
				}
			} // what else?
			trees = append(trees, tree)
		}
		for _, address := range endpoint.NotReadyAddresses {
			ref := address.TargetRef
			tree := NewTree(ref.Kind, fmt.Sprintf("%s.%s", ref.Namespace, ref.Name), nil)
			tree.setStatus(false)
			trees = append(trees, tree)
		}
	}

	return
}
