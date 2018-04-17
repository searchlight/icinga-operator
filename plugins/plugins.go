package plugins

import "github.com/appscode/searchlight/pkg/icinga"

const (
	FlagKubeConfig        = "kubeconfig"
	FlagKubeConfigContext = "context"
	FlagHost              = "host"
	FlagCheckInterval     = "icinga.checkInterval"
)

type PluginInterface interface {
	Check() (icinga.State, interface{})
}
