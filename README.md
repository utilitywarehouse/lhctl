# lhctl

Attempt to create a simple longhorn CLI to automate upgrading storage nodes.

## Installation

Grab a binary from the available [releases](https://github.com/utilitywarehouse/lhctl/releases)

## Usage

The cli client needs an http endpoint for most of the available commands.
Currently that could be done with `--url=` flag. That could be easily obtained
by kubernetes, for example:

```
CONTEXT=<kube-context>
NAMESPACE=<longhorn-namespace>
LH_SVC=$(kubectl --context=${CONTEXT} --namespace=${NAMESPACE} get svc | grep longhorn-frontend | awk '{print $4}')
lhctl --url=http://${LH_SVC}/v1 get volume
```

Available commands can be see using help:

```
# lhctl --help

CLI to perform basic actions against longhorn api to help with
automated tasks.
For example:

# lhctl get volume

# lhctl get node

Usage:
  lhctl [command]

Available Commands:
  attach      Attach volume from node
  completion  Generates bash completion scripts
  delete      Delete resources resources
  detach      Detach volume from node
  disable     Disable scheduling on node
  enable      Enable scheduling on node
  get         List resources
  help        Help about any command
  update      Update resources

Flags:
  -h, --help         help for lhctl
  -t, --toggle       Help message for toggle
      --url string   longhorn manager url (example: http://10.88.1.3/v1)

Use "lhctl [command] --help" for more information about a command.
```

## Completion

For `linux` users:
```
lhctl  completion > /etc/bash_completion.d/lhctl 
```

and reload your shell.
