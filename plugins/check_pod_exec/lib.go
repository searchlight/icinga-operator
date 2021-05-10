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

package check_pod_exec

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"go.searchlight.dev/icinga-operator/pkg/icinga"
	"go.searchlight.dev/icinga-operator/plugins"

	"github.com/spf13/cobra"
	"gomodules.xyz/x/flags"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	utilexec "k8s.io/client-go/util/exec"
	"kmodules.xyz/client-go/tools/clientcmd"
)

type plugin struct {
	config  *restclient.Config
	client  corev1.CoreV1Interface
	options options
}

var _ plugins.PluginInterface = &plugin{}

func newPluginFromConfig(opts options) (*plugin, error) {
	config, err := clientcmd.BuildConfigFromContext(opts.kubeconfigPath, opts.contextName)
	if err != nil {
		return nil, err
	}
	return &plugin{config, kubernetes.NewForConfigOrDie(config).CoreV1(), opts}, nil
}

type options struct {
	kubeconfigPath string
	contextName    string
	// options
	podName   string
	container string
	namespace string
	command   string
	arg       string
	// IcingaHost
	host *icinga.IcingaHost
}

func (o *options) complete(cmd *cobra.Command) error {
	hostname, err := cmd.Flags().GetString(plugins.FlagHost)
	if err != nil {
		return err
	}
	o.host, err = icinga.ParseHost(hostname)
	if err != nil {
		return errors.New("invalid icinga host.name")
	}
	o.podName = o.host.ObjectName
	o.namespace = o.host.AlertNamespace

	o.kubeconfigPath, err = cmd.Flags().GetString(plugins.FlagKubeConfig)
	if err != nil {
		return err
	}
	o.contextName, err = cmd.Flags().GetString(plugins.FlagKubeConfigContext)
	if err != nil {
		return err
	}
	return nil
}

func (o *options) validate() error {
	if o.host.Type != icinga.TypePod {
		return errors.New("invalid icinga host type")
	}
	return nil
}

type Writer struct {
	Str []string
}

func (w *Writer) Write(p []byte) (n int, err error) {
	str := string(p)
	if len(str) > 0 {
		w.Str = append(w.Str, str)
	}
	return len(str), nil
}

func newStringReader(ss []string) io.Reader {
	formattedString := strings.Join(ss, "\n")
	reader := strings.NewReader(formattedString)
	return reader
}

func (p *plugin) Check() (icinga.State, interface{}) {
	opts := p.options
	pod, err := p.client.Pods(opts.namespace).Get(context.TODO(), opts.podName, metav1.GetOptions{})
	if err != nil {
		return icinga.Unknown, err
	}

	if opts.container != "" {
		notFound := true
		for _, container := range pod.Spec.Containers {
			if container.Name == opts.container {
				notFound = false
				break
			}
		}
		if notFound {
			return icinga.Unknown, fmt.Sprintf(`Container "%v" not found`, opts.container)
		}
	}

	execRequest := p.client.RESTClient().Post().
		Resource("pods").
		Name(opts.podName).
		Namespace(opts.namespace).
		SubResource("exec").
		Param("container", opts.container).
		Param("command", opts.command).
		Param("stdin", "true").
		Param("stdout", "true").
		Param("stderr", "true").
		Param("tty", "false")

	exec, err := remotecommand.NewSPDYExecutor(p.config, "POST", execRequest.URL())
	if err != nil {
		return icinga.Unknown, err
	}

	stdIn := newStringReader([]string{opts.arg})
	stdOut := new(Writer)
	stdErr := new(Writer)

	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdIn,
		Stdout: stdOut,
		Stderr: stdErr,
		Tty:    false,
	})

	var exitCode int
	if err == nil {
		exitCode = 0
	} else {
		if exitErr, ok := err.(utilexec.ExitError); ok && exitErr.Exited() {
			exitCode = exitErr.ExitStatus()
		} else {
			return icinga.Unknown, "Failed to find exit code."
		}
	}

	output := fmt.Sprintf("Exit Code: %v", exitCode)
	if exitCode != 0 {
		exitCode = 2
	}

	return icinga.State(exitCode), output
}

const (
	flagArgv = "argv"
)

func NewCmd() *cobra.Command {
	var opts options

	c := &cobra.Command{
		Use:   "check_pod_exec",
		Short: "Check exit code of exec command on Kubernetes container",

		Run: func(cmd *cobra.Command, args []string) {
			flags.EnsureRequiredFlags(cmd, plugins.FlagHost, flagArgv)

			if err := opts.complete(cmd); err != nil {
				icinga.Output(icinga.Unknown, err)
			}
			if err := opts.validate(); err != nil {
				icinga.Output(icinga.Unknown, err)
			}
			plugin, err := newPluginFromConfig(opts)
			if err != nil {
				icinga.Output(icinga.Unknown, err)
			}
			icinga.Output(plugin.Check())
		},
	}

	c.Flags().StringP(plugins.FlagHost, "H", "", "Icinga host name")
	c.Flags().StringVarP(&opts.container, "container", "C", "", "Container name in specified pod")
	c.Flags().StringVarP(&opts.command, "cmd", "c", "/bin/sh", "Exec command. [Default: /bin/sh]")
	c.Flags().StringVarP(&opts.arg, flagArgv, "a", "", "Arguments for exec command. [Format: 'arg; arg; arg']")
	return c
}
