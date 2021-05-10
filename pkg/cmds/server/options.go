/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server

import (
	"context"
	"flag"
	"time"

	api "go.searchlight.dev/icinga-operator/apis/monitoring/v1alpha1"
	cs "go.searchlight.dev/icinga-operator/client/clientset/versioned"
	"go.searchlight.dev/icinga-operator/pkg/admission/plugin"
	"go.searchlight.dev/icinga-operator/pkg/icinga"
	"go.searchlight.dev/icinga-operator/pkg/operator"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	crd_cs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	"kmodules.xyz/client-go/meta"
	hooks "kmodules.xyz/webhook-runtime/admission/v1beta1"
)

type OperatorOptions struct {
	ConfigRoot       string
	ConfigSecretName string
	ResyncPeriod     time.Duration
	MaxNumRequeues   int
	NumThreads       int
	IncidentTTL      time.Duration
	// V logging level, the value of the -v flag
	verbosity string
}

func NewOperatorOptions() *OperatorOptions {
	return &OperatorOptions{
		ConfigRoot:       "/srv",
		ConfigSecretName: "searchlight-operator",
		ResyncPeriod:     5 * time.Minute,
		MaxNumRequeues:   5,
		NumThreads:       1,
		IncidentTTL:      90 * 24 * time.Hour,
		verbosity:        "3",
	}
}

func (s *OperatorOptions) AddGoFlags(fs *flag.FlagSet) {
	fs.StringVar(&s.ConfigRoot, "config-dir", s.ConfigRoot, "Path to directory containing icinga2 config. This should be an emptyDir inside Kubernetes.")
	fs.StringVar(&s.ConfigSecretName, "config-secret-name", s.ConfigSecretName, "Name of Kubernetes secret used to pass icinga credentials.")
	fs.DurationVar(&s.ResyncPeriod, "resync-period", s.ResyncPeriod, "If non-zero, will re-list this often. Otherwise, re-list will be delayed aslong as possible (until the upstream source closes the watch or times out.")
	fs.DurationVar(&s.IncidentTTL, "incident-ttl", s.IncidentTTL, "Garbage collects incidents older than this duration. Set to 0 to disable garbage collection.")

	fs.BoolVar(&api.EnableStatusSubresource, "enable-status-subresource", api.EnableStatusSubresource, "If true, uses sub resource for Voyager crds.")
}

func (s *OperatorOptions) AddFlags(fs *pflag.FlagSet) {
	pfs := flag.NewFlagSet("searchlight", flag.ExitOnError)
	s.AddGoFlags(pfs)
	fs.AddGoFlagSet(pfs)
}

func (s *OperatorOptions) ApplyTo(cfg *operator.OperatorConfig) error {
	var err error

	cfg.ConfigRoot = s.ConfigRoot
	cfg.ConfigSecretName = s.ConfigSecretName
	cfg.ResyncPeriod = s.ResyncPeriod
	cfg.MaxNumRequeues = s.MaxNumRequeues
	cfg.NumThreads = s.NumThreads
	cfg.IncidentTTL = s.IncidentTTL
	cfg.Verbosity = s.verbosity

	if cfg.KubeClient, err = kubernetes.NewForConfig(cfg.ClientConfig); err != nil {
		return err
	}
	if cfg.ExtClient, err = cs.NewForConfig(cfg.ClientConfig); err != nil {
		return err
	}
	if cfg.CRDClient, err = crd_cs.NewForConfig(cfg.ClientConfig); err != nil {
		return err
	}
	cfg.AdmissionHooks = []hooks.AdmissionHook{&plugin.CRDValidator{}}

	secret, err := cfg.KubeClient.CoreV1().Secrets(meta.Namespace()).Get(context.TODO(), s.ConfigSecretName, metav1.GetOptions{})
	if err != nil {
		return errors.Wrapf(err, "failed to load secret: %s", s.ConfigSecretName)
	}

	mgr := &icinga.Configurator{
		ConfigRoot:       s.ConfigRoot,
		IcingaSecretName: s.ConfigSecretName,
		Expiry:           10 * 365 * 24 * time.Hour,
	}
	data, err := mgr.LoadConfig(func(key string) (value string, found bool) {
		var bytes []byte
		bytes, found = secret.Data[key]
		value = string(bytes)
		return
	})
	if err != nil {
		klog.Fatalln(err)
	}

	cfg.IcingaClient = icinga.NewClient(*data)
	for {
		if cfg.IcingaClient.Check().Get(nil).Do().Status == 200 {
			klog.Infoln("connected to icinga api")
			break
		}
		klog.Infoln("Waiting for icinga to start")
		time.Sleep(2 * time.Second)
	}

	return nil
}
