[![Go Report Card](https://goreportcard.com/badge/github.com/spaghettifunk/linkerd2-operator)](https://goreportcard.com/report/github.com/spaghettifunk/linkerd2-operator)
[![codecov](https://codecov.io/gh/spaghettifunk/linkerd2-operator/branch/main/graph/badge.svg)](https://codecov.io/gh/spaghettifunk/linkerd2-operator)

# Linkerd2 Operator

## :warning: Project status (pre-alpha)

The project is in active **development** and things are rapidly changing.

## Overview

This Linkerd2 operator is a simple operator that takes care of deploying all the components of Linkerd2

## Prerequisites

- [go][go_tool] version v1.13+.
- [docker][docker_tool] version 17.03+
- [kubectl][kubectl_tool] v1.14.1+
- [operator-sdk][operator_install]
- [kustomize][kustomize_tool]
- Access to a Kubernetes v1.14.1+ cluster

## Getting Started

### Cloning the repository

Checkout this repository

```
$ mkdir -p $GOPATH/src/github.com/spaghettifunk
$ cd $GOPATH/src/github.com/spaghettifunk
$ git clone https://github.com/spaghettifunk/linkerd2-operator.git
$ cd linkerd2-operator
```

### Installing

There is a `Makefile` for convenience. Run the following commands to start

```shell
export POD_NAMESPACE=linkerd
make install
make run
```

In a new shell, you can now create the Linkerd deployments. There is an example of usage in the `config/sample/linkerd.example.yaml`. Use that for testing. Assuming you want to use that file, run the following

```yaml
kubectl create ns linkerd
kubectl apply -f config/sample/linkerd.example.yaml
```

In the previous shell you should see that Kubernetes is trying to reconcile the object.

### Uninstalling

To uninstall all that was performed in the above step run `make uninstall`.

### Troubleshooting

Use the following command to check the operator logs.

```shell
kubectl logs deployment.apps/linkerd2-operator -n linkerd
```

### Running Tests

Run `make test-e2e` to run the integration e2e tests with different options. For
more information see the [writing e2e tests][golang-e2e-tests] guide.

[dep_tool]: https://golang.github.io/dep/docs/installation.html
[go_tool]: https://golang.org/dl/
[kubectl_tool]: https://kubernetes.io/docs/tasks/tools/install-kubectl/
[docker_tool]: https://docs.docker.com/install/
[operator_sdk]: https://github.com/operator-framework/operator-sdk
[operator_install]: https://sdk.operatorframework.io/docs/install-operator-sdk/
[golang-e2e-tests]: https://sdk.operatorframework.io/docs/golang/e2e-tests/
