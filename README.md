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
  delete      Delete resources
  detach      Detach volume from node
  disable     Disable scheduling on node
  enable      Enable scheduling on node
  get         List resources
  help        Help about any command
  update      Update resources

Flags:
      --config string    config file (default is $HOME/.lhctl.yaml)
      --context string   config file context
  -h, --help             help for lhctl
      --pass string      password for http request
      --url string       longhorn manager url (example: http://10.88.1.3/v1)
      --user string      user for http request

Use "lhctl [command] --help" for more information about a command.
```

## Completion

For `linux` users:
```
lhctl  completion > /etc/bash_completion.d/lhctl 
```

and reload your shell.

## Auth

Rancher client supports basic auth on http requests. Since longhorn does not
implement auth, a reverse proxy that supports http basic auth can be used in
front of the longhorn ui to implement this feature. Example on kubernetes using
`traefik` ingresses:

```
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  labels:
    kubernetes.io/ingress.class: traefik
  annotations:
    ingress.kubernetes.io/auth-type: basic
    ingress.kubernetes.io/auth-secret: longhorn-admins
  name: longhorn
spec:
  rules:
  - host: longhorn.test.example.com
    http:
      paths:
      - backend:
          serviceName: longhorn-frontend
          servicePort: 80
        path: /
---
apiVersion: v1
kind: Secret
metadata:
  name: longhorn-admins
data:
  users: |
    dGVzdDokYXByMSRINnVza2trVyRJZ1hMUDZld1RyU3VCa1RycUU4d2ovCnRlc3QyOiRhcHIxJGQ5aHI5SEJCJDRIeHdnVWlyM0hQNEVzZ2dQL1FObzAK
```

htpasswd generation example:
```
$ pwgen 20 1
etee0aenah4uYarie2Ya
$ htpasswd -nb test etee0aenah4uYarie2Ya
test:$apr1$FiK5nZd8$n/MpMQJVrssqejCFpAOMF/
```

## Config File - Contexts

The client supports reading config from a `yaml` file, where someone can define
many contexts in case more than one longhorn endpoints are available. The
default location is `${HOME}/.lhctl.yaml` and an alternative one can be
specified using `--config` flag.

Config file example:

```
contexts:
- name: context-1
  url: https://longhorn.context.one
  user: test
  pass: test
- name: context-2
  url: https://longhorn.context.two
default: context-1
```

The desired context can be specified via `--context` flag, otherwise the default
context will be assumed.
User and pass are of course optional and rely on the user's ui setup as
described above.
