package plugin

import "k8s.io/cli-runtime/pkg/genericclioptions"

func currentNamespace(configFlags *genericclioptions.ConfigFlags) string {
	if v := *configFlags.Namespace; v != "" {
		return v
	}
	clientConfig := configFlags.ToRawKubeConfigLoader()
	defaultNamespace, _, err := clientConfig.Namespace()
	if err != nil {
		defaultNamespace = "default"
	}
	return defaultNamespace
}
