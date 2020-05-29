package plugin

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/replicatedhq/krew-plugin-template/pkg/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
)

func RunPlugin(configFlags *genericclioptions.ConfigFlags) error {
	log := logger.NewLogger()

	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return errors.Wrap(err, "failed to read kubeconfig")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "failed to create clientset")
	}

	ingresses, err := clientset.ExtensionsV1beta1().Ingresses(currentNamespace(configFlags)).List(metav1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to list ingresses")
	}

	trees := []Tree{}
	for _, ingress := range ingresses.Items {
		trees, err = exploreIngress(configFlags, clientset, ingress)
		if err != nil {
			return errors.Wrap(err, "failed to explore ingress")
		}
	}

	exteranalServices, err := clientset.CoreV1().Services(currentNamespace(configFlags)).List(metav1.ListOptions{}) // TODO filter by type
	if err != nil {
		return errors.Wrap(err, "failed to list services")
	}

	for _, service := range exteranalServices.Items {
		if service.Spec.Type == "ClusterIP" {
			continue
		}
		for _, port := range service.Spec.Ports {
			var tree Tree
			tree, err = exploreService(configFlags, clientset, service.Namespace, service.Name, port.Port)
			trees = append(trees, tree)
		}
	}

	if len(trees) == 0 {
		log.Info("no service found in namespace '%s'", currentNamespace(configFlags))
		return nil
	}

	first := true
	for _, tree := range trees {
		tree.display(os.Stdout)
		if first {
			fmt.Fprintf(os.Stdout, "\n")
		}
	}

	return nil
}
