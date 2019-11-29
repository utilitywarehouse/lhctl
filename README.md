# lhctl

Attempt to create a simple longhorn CLI to automate upgrading storage nodes.

## Usage

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
  completion  Generates bash completion scripts
  delete      Delete resources resources
  disable     Disable scheduling on node
  enable      Enable scheduling on node
  get         List resources
  help        Help about any command

Flags:
  -h, --help         help for lhctl
  -t, --toggle       Help message for toggle
      --url string   longhorn manager url (example: http://10.88.1.3/v1)

Use "lhctl [command] --help" for more information about a command.

```

Currently you need to procide `--url=` flag for all the commands that interact
with longhorn api.

## Completion

For `linux` users:
```
lhctl  completion > /etc/bash_completion.d/lhctl 
```

and reload your shell.
