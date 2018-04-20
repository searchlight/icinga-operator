---
title: Check command | Icinga2
description: How to add support of additional Check command in Searchlight
menu:
  product_searchlight_6.0.0-rc.0:
    identifier: add-check-command
    name: CheckCommand
    parent: developer-guide
    weight: 15
product_name: searchlight
menu_name: product_searchlight_6.0.0-rc.0
section_menu_id: setup
---

# Add Check command in Searchlight

This document will show you how to add additional Check command support in Searchlight.

You need to follow these steps:

* Write Plugin for your Check command.
* Configure *CheckCommand* object in Icinga2.
* Add *check_command* information in `data/files`.

## Write Plugin

First of all, you need to write a plugin that will be called by Icinga. This plugin will do its job and will exit with specific code.

|  Service State | Exit Code |
|----------------|-----------|
| OK             | 0         |
| Warning        | 1         |
| Critical       | 2         |
| Unknown        | 3         |

That means, if plugin exits with code `0`, Icinga Service State will be `OK`. Similarly if it exits with code `2`,
State will be `Critical`. Your plugin should determine the correct State of this Check command.

All standard output will be shown as *Plugin Output* in Icingaweb2.

Check some existing plugins [here](https://github.com/appscode/searchlight/tree/master/plugins)

You can add your plugin as sub-command of `hyperalert`.

## Configure *CheckCommand* object.

When you create an Alert using Searchlight, you need to provide *check_command* name in `spec.check`.
Operator uses this name as Icinga2 Service attribute `check_command`. And this *check_command* should be configured in Icinga2.

Lets see an example CheckCommand

```text
object CheckCommand "component-status" {
  import "plugin-check-command"
  command = [ PluginDir + "/hyperalert", "check_component_status"]

  arguments = {
    "--selector" = "$selector$"
    "--componentName" = "$componentName$"
    "--v" = "$host.vars.verbosity$"
  }
}
```

Here, `component-status` is the name of the *check_command* provided in `spec.check`. And when Service checks its State, it executes a plugin
defined as `command` in *CheckCommand* configuration.

In this example, `hyperalert` plugin is called with command `check_component_status`. This plugin is called with parameters defined in arguments.

You can pass some custom variables to your plugins. In your Alert object, you can add variables in `spec.Vars`.
These variables will be added as custom variables in Icinga Service object. You need to configure CheckCommand to forward these variables to plugin.

Here followings are custom variables of Service object provided by user in `spec.Vars`.

```text
    "--selector" = "$selector$"
    "--componentName" = "$componentName$"
```

You can also forward data from Host object to plugin if necessary.

```text
 "--host" = "$host.name$"
 "--v" = "$host.vars.verbosity$"
```

Here, `verbosity` is custom variable set in Host object by Searchlight.

## Add *check_command* information in data files

Searchlight operator validates provided information for an Alert. You need to add information of your Check command in data files
otherwise, operator will not create Icinga objects for your Alert.

```json
{
      "name": "component-status",
      "vars": [
        {
          "flag": {
            "short": "l",
            "long": "selector"
          },
          "name": "selector",
          "type": "STRING"
        },
        {
          "flag": {
            "long": "componentName"
          },
          "name": "componentName",
          "type": "STRING"
        }
      ],
       "states": [
         "OK",
         "Critical",
         "Unknown"
       ]
     }
```

Check command `component-status` has two custom variables and supported States are `OK`, `Critical` and `Unknown`.


# Build

Now run following command to build hyperalert.

```bash
$ ./hack/make.py build hyperalert
```

Then build icinga2 docker images and push to your registry

```bash
$ ./hack/docker/icinga/alpine/build.sh
$ ./hack/docker/icinga/alpine/build.sh push
```

Change docker registry and image tag.

Finally build searchlight docker image and push to your registry

```bash
$ ./hack/docker/searchlight/setup.sh
$ ./hack/docker/searchlight/setup.sh push
```

Change docker registry and image tag.
