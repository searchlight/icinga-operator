package check_kube_exec

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/appscode/go/flags"
	"github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/util"
	"github.com/spf13/cobra"
	kapi "k8s.io/kubernetes/pkg/api"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/client/unversioned/remotecommand"
	remotecommandserver "k8s.io/kubernetes/pkg/kubelet/server/remotecommand"
	utilexec "k8s.io/kubernetes/pkg/util/exec"
)

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

func CheckKubeExec(req *Request) (util.IcingaState, interface{}) {
	kubeConfig, err := k8s.GetKubeConfig()
	if err != nil {
		return util.Unknown, err
	}

	kubeClient, err := clientset.NewForConfig(kubeConfig)
	if err != nil {
		return util.Unknown, err
	}

	pod, err := kubeClient.Core().Pods(req.Namespace).Get(req.Pod)
	if err != nil {
		return util.Unknown, err
	}

	foundContainer := false
	if req.Container != "" {
		for _, container := range pod.Spec.Containers {
			if container.Name == req.Container {
				foundContainer = true
				break
			}
		}
	} else {
		foundContainer = true
	}

	if !foundContainer {
		return util.Unknown, fmt.Sprintf(`Container "%v" not found`, req.Container)
	}

	execRequest := kubeClient.Core().RESTClient().Post().
		Resource("pods").
		Name(req.Pod).
		Namespace(req.Namespace).
		SubResource("exec").
		Param("container", req.Container)

	execRequest.VersionedParams(&kapi.PodExecOptions{
		Container: req.Container,
		Command:   []string{req.Command},
		Stdin:     true,
		Stdout:    false,
		Stderr:    false,
		TTY:       false,
	}, kapi.ParameterCodec)

	exec, err := remotecommand.NewExecutor(kubeConfig, "POST", execRequest.URL())
	if err != nil {
		return util.Unknown, err
	}

	stdIn := newStringReader([]string{"-c", req.Arg})
	stdOut := new(Writer)
	stdErr := new(Writer)

	err = exec.Stream(remotecommand.StreamOptions{
		SupportedProtocols: remotecommandserver.SupportedStreamingProtocols,
		Stdin:              stdIn,
		Stdout:             stdOut,
		Stderr:             stdErr,
		Tty:                false,
	})

	var exitCode int
	if err == nil {
		exitCode = 0
	} else {
		if exitErr, ok := err.(utilexec.ExitError); ok && exitErr.Exited() {
			exitCode = exitErr.ExitStatus()
		} else {
			return util.Unknown, "Failed to find exit code."
		}
	}

	output := fmt.Sprintf("Exit Code: %v", exitCode)
	if exitCode != 0 {
		exitCode = 2
	}

	return util.IcingaState(exitCode), output
}

type Request struct {
	Pod       string
	Container string
	Namespace string
	Command   string
	Arg       string
}

func NewCmd() *cobra.Command {
	var req Request
	var host string
	c := &cobra.Command{
		Use:     "check_kube_exec",
		Short:   "Check exit code of exec command on kubernetes container",
		Example: "",

		Run: func(cmd *cobra.Command, args []string) {
			flags.EnsureRequiredFlags(cmd, "host", "arg")

			parts := strings.Split(host, "@")
			if len(parts) != 2 {
				fmt.Fprintln(os.Stdout, util.State[3], "Invalid icinga host.name")
				os.Exit(3)
			}
			req.Pod = parts[0]
			req.Namespace = parts[1]
			util.Output(CheckKubeExec(&req))
		},
	}

	c.Flags().StringVarP(&host, "host", "H", "", "Icinga host name")
	c.Flags().StringVarP(&req.Container, "container", "C", "", "Container name in specified pod")
	c.Flags().StringVarP(&req.Command, "cmd", "c", "/bin/sh", "Exec command. [Default: /bin/sh]")
	c.Flags().StringVarP(&req.Arg, "argv", "a", "", "Arguments for exec command. [Format: 'arg; arg; arg']")
	return c
}
