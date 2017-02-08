package e2e

import (
	"fmt"
	"sync"
	"time"

	acw "github.com/appscode/k8s-addons/pkg/watcher"
	"github.com/appscode/searchlight/cmd/searchlight/app"
	"github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/util"
)

type testData struct {
	Data                map[string]interface{}
	ExpectedIcingaState util.IcingaState
}

type dataConfig struct {
	ObjectType   string
	CheckCommand string
	Namespace    string
}

type kubeClient struct {
	isClientSet bool
	client      *k8s.KubeClient
	once        sync.Once
}

var e2eClient = kubeClient{isClientSet: false}

func getKubeClient() (kubeClient *k8s.KubeClient, err error) {
	if e2eClient.isClientSet {
		kubeClient = e2eClient.client
		return
	}
	e2eClient.once.Do(
		func() {
			kubeClient, err = k8s.NewClient()
			e2eClient.client = kubeClient
			e2eClient.isClientSet = true
		},
	)
	return
}

type kubeWatcher struct {
	isWatcherSet bool
	watcher      *app.Watcher
	once         sync.Once
}

var e2eWatcher = kubeWatcher{isWatcherSet: false}

func runKubeD() (w *app.Watcher, err error) {
	if e2eWatcher.isWatcherSet {
		w = e2eWatcher.watcher
		return
	}

	e2eWatcher.once.Do(
		func() {
			fmt.Println("-- TestE2E: Waiting for kubed")
			var kubeClient *k8s.KubeClient
			if kubeClient, err = getKubeClient(); err != nil {
				return
			}

			w = &app.Watcher{
				Watcher: acw.Watcher{
					Client:                  kubeClient.Client,
					AppsCodeExtensionClient: kubeClient.AppscodeExtensionClient,
					SyncPeriod:              time.Minute * 2,
				},
			}

			w.Watcher.Dispatch = w.Dispatch
			go w.Run()
			time.Sleep(time.Second * 10)
			e2eWatcher.watcher = w
			e2eWatcher.isWatcherSet = true
		},
	)
	return
}
