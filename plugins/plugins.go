package plugins

import "github.com/appscode/searchlight/pkg/icinga"

const (
	FlagKubeConfig        = "kubeconfig"
	FlagKubeConfigContext = "context"
	FlagHost              = "host"
)

type PluginInterface interface {
	Check() (icinga.State, interface{})
}
