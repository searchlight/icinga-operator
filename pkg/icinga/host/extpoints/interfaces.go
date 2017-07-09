package extpoints

import "github.com/appscode/searchlight/pkg/controller/types"

type IcingaHostType interface {
	CreateAlert(ctx *types.Option, specificObject string) error
	UpdateAlert(ctx *types.Option) error
	DeleteAlert(ctx *types.Option, specificObject string) error
}
