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

package notifier

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	api "go.searchlight.dev/icinga-operator/apis/monitoring/v1alpha1"
)

func (n *notifier) RenderSubject(receiver api.Receiver) string {
	opts := n.options
	switch api.AlertType(opts.notificationType) {
	case api.NotificationAcknowledgement:
		return fmt.Sprintf("Problem Acknowledged: Service [%s] for [%s] is in \"%s\" state", opts.alertName, opts.hostname, receiver.State)
	case api.NotificationRecovery:
		return fmt.Sprintf("Problem Recovered: Service [%s] for [%s] was in \"%s\" state.", opts.alertName, opts.hostname, receiver.State)
	case api.NotificationProblem:
		return fmt.Sprintf("Problem Detected: Service [%s] for [%s] is in \"%s\" state.", opts.alertName, opts.hostname, receiver.State)
	default:
		return fmt.Sprintf("Service [%s] for [%s] is in \"%s\" state.", opts.alertName, opts.hostname, receiver.State)
	}
}

type TemplateData struct {
	AlertNamespace     string
	AlertType          string
	AlertName          string
	ObjectName         string
	IcingaHostName     string
	IcingaServiceName  string
	IcingaCheckCommand string
	IcingaType         string
	IcingaState        string
	IcingaOutput       string
	Author             string
	Comment            string
	IcingaTime         time.Time
}

func (n *notifier) RenderMail(alert api.Alert) (string, error) {
	opts := n.options
	host := opts.host
	data := TemplateData{
		AlertName:          alert.GetName(),
		AlertNamespace:     host.AlertNamespace,
		AlertType:          host.Type,
		ObjectName:         host.ObjectName,
		IcingaHostName:     n.options.hostname,
		IcingaServiceName:  alert.GetName(),
		IcingaCheckCommand: alert.Command(),
		IcingaType:         opts.notificationType,
		IcingaState:        strings.ToUpper(opts.serviceState),
		IcingaOutput:       opts.serviceOutput,
		Author:             opts.author,
		Comment:            opts.comment,
		IcingaTime:         opts.time,
	}

	var buf bytes.Buffer
	if err := mailTemplate.Execute(&buf, data); err != nil {
		return "", err
	}
	config := buf.String()
	return config, nil
}
