package check_kube_exec

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/appscode/go/flags"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/pkg/util"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	apiv1 "k8s.io/client-go/pkg/api/v1"
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

func CheckKubeExec(req *Request) (icinga.State, interface{}) {
	kubeConfig, err := util.GetKubeConfig()
	if err != nil {
		return icinga.UNKNOWN, err
	}

	kubeClient, err := clientset.NewForConfig(kubeConfig)
	if err != nil {
		return icinga.UNKNOWN, err
	}

	pod, err := kubeClient.CoreV1().Pods(req.Namespace).Get(req.Pod, metav1.GetOptions{})
	if err != nil {
		return icinga.UNKNOWN, err
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
		return icinga.UNKNOWN, fmt.Sprintf(`Container "%v" not found`, req.Container)
	}

	execRequest := kubeClient.Core().RESTClient().Post().
		Resource("pods").
		Name(req.Pod).
		Namespace(req.Namespace).
		SubResource("exec").
		Param("container", req.Container)

	execRequest.VersionedParams(&apiv1.PodExecOptions{
		Container: req.Container,
		Command:   []string{req.Command},
		Stdin:     true,
		Stdout:    false,
		Stderr:    false,
		TTY:       false,
	}, internalversion.ParameterCodec)

	exec, err := remotecommand.NewExecutor(kubeConfig, "POST", execRequest.URL())
	if err != nil {
		return icinga.UNKNOWN, err
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
			return icinga.UNKNOWN, "Failed to find exit code."
		}
	}

	output := fmt.Sprintf("Exit Code: %v", exitCode)
	if exitCode != 0 {
		exitCode = 2
	}

	return icinga.State(exitCode), output
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
	var icingaHost string
	c := &cobra.Command{
		Use:     "check_kube_exec",
		Short:   "Check exit code of exec command on kubernetes container",
		Example: "",

		Run: func(cmd *cobra.Command, args []string) {
			flags.EnsureRequiredFlags(cmd, "host", "arg")

			host, err := icinga.ParseHost(icingaHost)
			if err != nil {
				fmt.Fprintln(os.Stdout, icinga.WARNING, "Invalid icinga host.name")
				os.Exit(3)
			}
			if host.Type != icinga.TypePod {
				fmt.Fprintln(os.Stdout, icinga.WARNING, "Invalid icinga host type")
				os.Exit(3)
			}
			req.Namespace = host.AlertNamespace
			req.Pod = host.ObjectName
			icinga.Output(CheckKubeExec(&req))
		},
	}

	c.Flags().StringVarP(&icingaHost, "host", "H", "", "Icinga host name")
	c.Flags().StringVarP(&req.Container, "container", "C", "", "Container name in specified pod")
	c.Flags().StringVarP(&req.Command, "cmd", "c", "/bin/sh", "Exec command. [Default: /bin/sh]")
	c.Flags().StringVarP(&req.Arg, "argv", "a", "", "Arguments for exec command. [Format: 'arg; arg; arg']")
	return c
}
