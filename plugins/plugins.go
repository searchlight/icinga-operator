package plugins

import "github.com/appscode/searchlight/pkg/icinga"

type PluginInterface interface {
	Check() (icinga.State, interface{})
}
