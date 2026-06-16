package kubecontexts

import (
	"sort"

	"k8s.io/client-go/tools/clientcmd"
)

// Load returns the Kubernetes context names from the configured kubeconfig.
func Load(kubeconfig string) ([]string, error) {
	rules := clientcmd.NewDefaultClientConfigLoadingRules()
	if kubeconfig != "" {
		rules.ExplicitPath = kubeconfig
	}

	config, err := rules.Load()
	if err != nil {
		return nil, err
	}

	contexts := make([]string, 0, len(config.Contexts))
	for name := range config.Contexts {
		contexts = append(contexts, name)
	}

	sort.Strings(contexts)
	return contexts, nil
}
