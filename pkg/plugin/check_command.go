package plugin

import (
	"fmt"
	"sort"
	"strings"

	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"github.com/appscode/searchlight/plugins/check_webhook"
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

	if plugin.Spec.Arguments.Vars != nil {
		for key := range plugin.Spec.Arguments.Vars.Fields {
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

	sort.Slice(args, func(i, j int) bool {
		return args[i].key < args[j].key
	})

	var command string
	flagList := make([]string, 0)

	if webhook == nil {
		// Command in CheckCommand
		parts := strings.Split(plugin.Spec.Command, " ")
		for i, part := range parts {
			if i == 0 {
				command = command + fmt.Sprintf(`"/%s"`, part)
			} else {
				command = command + fmt.Sprintf(`, "%s"`, part)
			}
		}

		// Arguments in CheckCommand
		for _, f := range args {
			flagList = append(flagList, fmt.Sprintf(`"--%s" = "%s"`, f.key, f.val))
		}
	} else {
		// Command in CheckCommand
		command = `"/hyperalert", "check_webhook"`

		// URL for webhook
		namespace := "default"
		if webhook.Namespace != "" {
			namespace = webhook.Namespace
		}
		url := fmt.Sprintf("http://%s.%s.svc/%s", webhook.Name, namespace, plugin.Name)
		flagList = append(flagList, fmt.Sprintf(`"--%s" = "%s"`, check_webhook.FlagWebhookURL, url))
		flagList = append(flagList, fmt.Sprintf(`"--%s" = "%s"`, check_webhook.FlagCheckCommand, plugin.Name))

		// Arguments in CheckCommand
		for i, f := range args {
			flagList = append(flagList, fmt.Sprintf(`"--key.%d" = "%s"`, i, f.key))
			flagList = append(flagList, fmt.Sprintf(`"--val.%d" = "%s"`, i, f.val))
		}
	}

	return fmt.Sprintf(checkCommandTemplate, plugin.Name, command, strings.Join(flagList, "\n\t"))
}
