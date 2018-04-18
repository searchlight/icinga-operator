package plugin

import (
	"fmt"
	"strings"

	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
)

var checkCommandTemplate = `object CheckCommand "%s" {
  import "plugin-check-command"
  command = [ PluginDir + %s]

  arguments = {
	%s
  }
}`

func GenerateCheckCommand(plugin *api.SearchlightPlugin) string {
	type arg struct {
		key string
		val string
	}
	args := make([]arg, 0)

	args = append(args, arg{
		key: "icinga.checkInterval",
		val: "$service.check_interval$",
	})

	webhook := plugin.Spec.Webhook
	if webhook != nil {
		namespace := "default"

		if webhook.Namespace != "" {
			namespace = webhook.Namespace
		}

		args = append(args, arg{
			key: "url",
			val: fmt.Sprintf("http://%s.%s.svc", webhook.Name, namespace),
		})
	}

	if plugin.Spec.Arguments.Vars != nil {
		for key := range plugin.Spec.Arguments.Vars.Items {
			args = append(args, arg{
				key: key,
				val: fmt.Sprintf("$%s$", key),
			})
		}
	}

	for key, val := range plugin.Spec.Arguments.Host {
		args = append(args, arg{
			key: key,
			val: fmt.Sprintf("$host.%s$", val),
		})
	}

	flagList := make([]string, 0)

	if webhook == nil {
		for _, f := range args {
			flagList = append(flagList, fmt.Sprintf(`"--%s" = "%s"`, f.key, f.val))
		}
	} else {
		for i, f := range args {
			flagList = append(flagList, fmt.Sprintf(`"--key.%d" = "%s"`, i, f.key))
			flagList = append(flagList, fmt.Sprintf(`"--val.%d" = "%s"`, i, f.val))
		}
	}

	var command string
	parts := strings.Split(plugin.Spec.Command, " ")
	for i, part := range parts {
		if i == 0 {
			command = command + fmt.Sprintf(`"/%s"`, part)
		} else {
			command = command + fmt.Sprintf(`, "%s"`, part)
		}
	}

	return fmt.Sprintf(checkCommandTemplate, plugin.Name, command, strings.Join(flagList, "\n\t"))
}
