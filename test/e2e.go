package e2e

import (
	"fmt"
	"time"

	acw "github.com/appscode/k8s-addons/pkg/watcher"
	"github.com/appscode/searchlight/cmd/searchlight/app"
	"github.com/appscode/searchlight/pkg/client"
)

type dataConfig struct {
	ObjectType   string
	CheckCommand string
	Namespace    string
}

func runKubeD(context *client.Context) *app.Watcher {
	fmt.Println("-- TestE2E: Waiting for kubed")
	w := &app.Watcher{
		Watcher: acw.Watcher{
			Client:                  context.KubeClient.Client,
			AppsCodeExtensionClient: context.KubeClient.AppscodeExtensionClient,
			SyncPeriod:              time.Minute * 2,
		},
		IcingaClient: context.IcingaClient,
	}

	w.Watcher.Dispatch = w.Dispatch
	go w.Run()
	time.Sleep(time.Second * 10)
	return w
}
