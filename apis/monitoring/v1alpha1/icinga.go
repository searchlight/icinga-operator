package v1alpha1

import (
	"strings"
	"sync"

	"github.com/appscode/go-notify/unified"
	"github.com/appscode/go/log"
	"github.com/appscode/searchlight/data"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	CheckPodInfluxQuery = "influx-query"
	CheckPodStatus      = "pod-status"
	CheckPodVolume      = "pod-volume"
	CheckPodExec        = "pod-exec"
)

const (
	CheckNodeInfluxQuery = "influx-query"
	CheckNodeVolume      = "node-volume"
	CheckNodeStatus      = "node-status"
)

const (
	CheckComponentStatus = "component-status"
	CheckJsonPath        = "json-path"
	CheckNodeExists      = "node-exists"
	CheckPodExists       = "pod-exists"
	CheckEvent           = "event"
	CheckCACert          = "ca-cert"
	CheckHttp            = "any-http"
	CheckEnv             = "env"
	CheckDummy           = "dummy"
	//CheckICMP           = "icmp"
	//CheckDIG            = "dig"
	//CheckDNS            = "dns"
)

// +k8s:deepcopy-gen=false
type Registry struct {
	reg map[string]IcingaCommand
	mu  sync.RWMutex
}

func (c *Registry) Get(cmd string) (IcingaCommand, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.reg[cmd]
	return v, ok
}

func (c *Registry) Insert(cmd string, v IcingaCommand) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.reg[cmd] = v
}

func (c *Registry) Delete(cmd string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.reg, cmd)
}

// +k8s:deepcopy-gen=false
type IcingaCommand struct {
	Name   string
	Vars   map[string]data.CommandVar
	States []string
}

var (
	PodCommands     = &Registry{reg: map[string]IcingaCommand{}}
	NodeCommands    = &Registry{reg: map[string]IcingaCommand{}}
	ClusterCommands = &Registry{reg: map[string]IcingaCommand{}}
)

func init() {
	clusterChecks, err := data.LoadClusterChecks()
	if err != nil {
		log.Fatal(err)
	}
	for _, cmd := range clusterChecks.Command {
		vars := make(map[string]data.CommandVar)
		for _, v := range cmd.Vars {
			vars[v.Name] = v
		}
		ClusterCommands.Insert(cmd.Name, IcingaCommand{
			Name:   cmd.Name,
			Vars:   vars,
			States: cmd.States,
		})
	}

	nodeChecks, err := data.LoadNodeChecks()
	if err != nil {
		log.Fatal(err)
	}
	for _, cmd := range nodeChecks.Command {
		vars := make(map[string]data.CommandVar)
		for _, v := range cmd.Vars {
			vars[v.Name] = v
		}
		NodeCommands.Insert(cmd.Name, IcingaCommand{
			Name:   cmd.Name,
			Vars:   vars,
			States: cmd.States,
		})
	}

	podChecks, err := data.LoadPodChecks()
	if err != nil {
		log.Fatal(err)
	}
	for _, cmd := range podChecks.Command {
		vars := make(map[string]data.CommandVar)
		for _, v := range cmd.Vars {
			vars[v.Name] = v
		}
		PodCommands.Insert(cmd.Name, IcingaCommand{
			Name:   cmd.Name,
			Vars:   vars,
			States: cmd.States,
		})
	}
}

func checkNotifiers(kc kubernetes.Interface, alert Alert) error {
	if alert.GetNotifierSecretName() == "" && len(alert.GetReceivers()) == 0 {
		return nil
	}
	secret, err := kc.CoreV1().Secrets(alert.GetNamespace()).Get(alert.GetNotifierSecretName(), metav1.GetOptions{})
	if err != nil {
		return err
	}
	for _, r := range alert.GetReceivers() {
		_, err = unified.LoadVia(r.Notifier, func(key string) (value string, found bool) {
			var bytes []byte
			bytes, found = secret.Data[key]
			value = string(bytes)
			return
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func AlertType(t string) IncidentNotificationType {
	switch strings.ToUpper(t) {
	case "PROBLEM":
		return NotificationProblem
	case "ACKNOWLEDGEMENT":
		return NotificationAcknowledgement
	case "RECOVERY":
		return NotificationRecovery
	default:
		return NotificationCustom
	}
}
