# OSB Starter Pack Operator

An operator, using the [Operator SDK](https://github.com/operator-framework/operator-sdk)
that can deploy an [OSB starter pack](https://github.com/pmorie/osb-starter-pack) based project
into a particular namespace.

## Who should use this project?

You should use this project if you're looking for a quick way to deploy an
Open Service Broker into a cluster.

## Prerequisites

You'll need:

- [`go`](https://golang.org/dl/)
- A running [Kubernetes](https://github.com/kubernetes/kubernetes)
- The [service-catalog](https://github.com/kubernetes-incubator/service-catalog)
  [installed](https://github.com/kubernetes-incubator/service-catalog/blob/master/docs/install.md)
  in that cluster

## Getting started

You can `git clone` it to start poking around right away.

### Get the project

```console
$ cd $GOPATH/src && mkdir -p github.com/shawn-hurley && cd github.com/shawn-hurley && git clone git://github.com/shawn-hurley/starter-pack-operator
```

Change into the project directory:

```console
$ cd $GOPATH/src/github.com/shawn-hurley/starter-pack-operator
```

### Deploying the operator

You'll need to run:

- `kubectl create ns test`
- `kubectl create -f deploy/rbac.yaml` -> sets up permissions for the operator to run.
- `kubectl create -f deploy/crd.yaml` -> installs the broker CRD.
- `kubectl create -f deploy/operator.yaml` -> will deploy the image to the test namespace.
- `kubectl create -f deploy/cr.yaml` -> will deploy the custom resource for your broker.


### Updating for our broker

If you want to deploy your broker, you will need to update the `deploy/cr.yaml` file. You can change the `spec.image` field to be your broker.
If you wanted you could edit the already running broker to change the image.

### Broker custom resource

```yaml
metadata:
  name: example
  namespace: test
  ...
spec:
  authenticateK8SToken: false
  image: quay.io/osb-starter-pack/servicebroker:latest
  port: 1338
  tlsSecretRef:
    name: tls-example
    namespace: test
status:
  phase: Running
```

**Note: that if you update the secret, you must have the cert, key and ca  cert.**
