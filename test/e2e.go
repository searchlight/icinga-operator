package e2e

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/appscode/searchlight/cmd/searchlight/app"
	"github.com/appscode/searchlight/pkg/client/icinga"
	"github.com/appscode/searchlight/pkg/client/k8s"
	acw "github.com/appscode/searchlight/pkg/watcher"
	"github.com/appscode/searchlight/test/mini"
	"github.com/appscode/searchlight/util"
	kapi "k8s.io/kubernetes/pkg/api"
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

type icingaClient struct {
	isIcingaClientSet bool
	client            *icinga.IcingaClient
	once              sync.Once
}

var e2eIcingaClient = icingaClient{isIcingaClientSet: false}

var (
	IcingaAddress string = os.Getenv("ICINGA_ADDRESS")
	IcingaAPIUser string = os.Getenv("ICINGA_API_USER")
	IcingaAPIPass string = os.Getenv("ICINGA_API_PASS")
)

func getIcingaClient() (icingaClient *icinga.IcingaClient, err error) {
	if e2eIcingaClient.isIcingaClientSet {
		icingaClient = e2eIcingaClient.client
		return
	}
	e2eIcingaClient.once.Do(
		func() {
			var kubeClient *k8s.KubeClient
			kubeClient, err = getKubeClient()
			if err != nil {
				return
			}

			// Secret will be created with information of Icinga2 running in Docker for travisCI test
			secretMap := map[string]string{
				icinga.IcingaAPIUser: IcingaAPIUser,
				icinga.IcingaAPIPass: IcingaAPIPass,
				icinga.IcingaAddress: IcingaAddress,
			}

			namespace := kapi.NamespaceSystem
			icingaSecretName, err := mini.CreateIcingaSecret(kubeClient, namespace, secretMap)

			icingaClient, err = icinga.NewIcingaClient(kubeClient.Client, icingaSecretName, namespace)
			if err != nil {
				return
			}

			e2eIcingaClient.client = icingaClient
			e2eIcingaClient.isIcingaClientSet = true
		},
	)
	return
}

type kubeWatcher struct {
	isWatcherSet     bool
	isIcingaIncluded bool
	watcher          *app.Watcher
	once             sync.Once
}

var e2eWatcher = kubeWatcher{isWatcherSet: false}

func runKubeD(setIcingaClient bool) (w *app.Watcher, err error) {
	// Watcher is already running.. Want to Add IcingaClient
	if e2eWatcher.isWatcherSet && setIcingaClient {
		// IcingaClient is already added in watcher
		if e2eWatcher.isIcingaIncluded {
			w = e2eWatcher.watcher
			return
		} else {
			// Added IcingaClient
			var icingaClient *icinga.IcingaClient
			icingaClient, err = getIcingaClient()
			if err != nil {
				return
			}
			e2eWatcher.watcher.IcingaClient = icingaClient
			e2eWatcher.isIcingaIncluded = true
		}
		w = e2eWatcher.watcher
		return
	}

	// Watcher is already running
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
					KubeClient: kubeClient.Client,
					ExtClient:  kubeClient.ExtClient,
					SyncPeriod: time.Minute * 2,
				},
			}

			// Set IcingaClient
			if setIcingaClient {
				// Added IcingaClient
				var icingaClient *icinga.IcingaClient
				icingaClient, err = getIcingaClient()
				if err != nil {
					return
				}
				w.IcingaClient = icingaClient
				e2eWatcher.isIcingaIncluded = true
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
