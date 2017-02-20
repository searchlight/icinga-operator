package app

import (
	"time"

	"github.com/appscode/errors"
	"github.com/appscode/go/runtime"
	_ "github.com/appscode/k8s-addons/api/install"
	acs "github.com/appscode/k8s-addons/client/clientset"
	acw "github.com/appscode/k8s-addons/pkg/watcher"
	"github.com/appscode/log"
	"github.com/appscode/searchlight/cmd/searchlight/app/options"
	"github.com/appscode/searchlight/pkg/client/icinga"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
)

func Run(config *options.Config) {
	log.Infoln("Configuration:", config)
	defer runtime.HandleCrash()

	c, err := clientcmd.BuildConfigFromFlags(config.Master, config.KubeConfig)
	if err != nil {
		errors.Exit(err)
	}

	w := &Watcher{
		Watcher: acw.Watcher{
			Client:                  clientset.NewForConfigOrDie(c),
			AppsCodeExtensionClient: acs.NewACExtensionsForConfigOrDie(c),
			SyncPeriod:              time.Minute * 2,
		},
	}
	if config.IcingaSecretName != "" {
		icingaClient, err := icinga.NewIcingaClient(w.Client, config.IcingaSecretName, config.IcingaSecretNamespace)
		if err != nil {
			log.Fatalln(err)
		}
		w.IcingaClient = icingaClient
	} else {
		log.Fatalln("Missing icinga secret")
	}

	log.Infoln("configuration loadded, running watcher")
	go w.Run()
}
