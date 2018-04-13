package plugins

import "github.com/appscode/searchlight/pkg/icinga"

const (
	FlagKubeConfig        = "iv-kubeconfig"
	FlagKubeConfigContext = "iv-context"
	FlagHost              = "iv-host"
	FlagCheckInterval     = "iv-checkInterval"
)

type PluginInterface interface {
	Check() (icinga.State, interface{})
}
