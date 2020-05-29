package plugin

import (
	"fmt"

	"k8s.io/api/extensions/v1beta1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

func exploreIngress(configFlags *genericclioptions.ConfigFlags, clientset *kubernetes.Clientset, ingress v1beta1.Ingress) (trees []Tree, err error) {
	trees = []Tree{}
	for _, rule := range ingress.Spec.Rules {
		for _, path := range rule.HTTP.Paths {
			host := "*"
			if rule.Host != "" {
				host = rule.Host
			}
			info := fmt.Sprintf("%s%s", host, path.Path)
			tree := NewTree(
				"Ingress",
				fmt.Sprintf("%s.%s", ingress.Namespace, ingress.Name),
				&info)
			var service Tree
			service, err = exploreService(configFlags, clientset, ingress.Namespace, path.Backend.ServiceName, path.Backend.ServicePort.IntVal)
			if err != nil {
				return
			}
			tree.addChild(service)
			trees = append(trees, tree)
		}
	}
	return
}
