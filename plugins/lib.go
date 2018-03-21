package plugins

import (
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	FlagKubeConfig        = "kubeconfig"
	FlagKubeConfigContext = "context"
	FlagHost              = "host"
)

func BuildConfig(kubeconfigPath, contextName string) (*restclient.Config, error) {
	var loader clientcmd.ClientConfigLoader
	if kubeconfigPath == "" {
		return restclient.InClusterConfig()
	} else {
		loader = &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath}
	}
	overrides := &clientcmd.ConfigOverrides{
		CurrentContext:  contextName,
		ClusterDefaults: clientcmd.ClusterDefaults,
	}
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loader, overrides).ClientConfig()
}

func GetClient(kubeconfigPath, contextName string) (kubernetes.Interface, error) {
	cfg, err := BuildConfig(kubeconfigPath, contextName)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(cfg)
}
