package app

import (
	"fmt"
	"os"
	"time"

	"github.com/appscode/go/runtime"
	"github.com/appscode/log"
	_ "github.com/appscode/searchlight/api/install"
	acs "github.com/appscode/searchlight/client/clientset"
	"github.com/appscode/searchlight/cmd/searchlight/app/options"
	"github.com/appscode/searchlight/pkg/client/icinga"
	acw "github.com/appscode/searchlight/pkg/watcher"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
)

func Run(config *options.Config) {
	log.Infoln("Configuration:", config)
	defer runtime.HandleCrash()

	c, err := clientcmd.BuildConfigFromFlags(config.Master, config.KubeConfig)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	w := &Watcher{
		Watcher: acw.Watcher{
			Client:     clientset.NewForConfigOrDie(c),
			ExtClient:  acs.NewForConfigOrDie(c),
			SyncPeriod: time.Minute * 2,
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
