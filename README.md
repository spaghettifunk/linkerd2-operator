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
### Pulling the dependencies

Run the following command

```
$ go mod tidy
```

### Building the operator

Build the Linkerd2 operator image and push it to a public registry, such as quay.io:

```
$ export IMAGE=quay.io/{YOUR_REPOSITORY}/linkerd2-operator:v0.0.1
$ operator-sdk build $IMAGE
$ docker push $IMAGE
```

### Using the image

```
# Update the operator manifest to use the built image name (if you are performing these steps on OSX, see note below)
$ sed -i 's|REPLACE_IMAGE|quay.io/{YOUR_REPOSITORY}/linkerd2-operator:v0.0.1|g' deploy/operator.yaml
# On OSX use:
$ sed -i "" 's|REPLACE_IMAGE|quay.io/{YOUR_REPOSITORY}/linkerd2-operator:v0.0.1|g' deploy/operator.yaml
```

**NOTE** The `quay.io/{YOUR_REPOSITORY}/linkerd2-operator:v0.0.1` is an example. You should build and push the image for your repository.

### Installing

Run `make install` to install the operator. Check that the operator is running in the cluster, also check that the LINKERD service was deployed.

Following the expected result.

```shell
$ kubectl get all -n linkerd
NAME                                      READY   STATUS    RESTARTS   AGE
pod/example-linkerd-7c4df9b7b4-lzd6j      1/1     Running   0          64s
pod/example-linkerd-7c4df9b7b4-wbtkz      1/1     Running   0          64s
pod/example-linkerd-7c4df9b7b4-wt6jb      1/1     Running   0          64s
pod/linkerd2-operator-56f54d84bf-zrtfv    1/1     Running   0          69s

NAME                                 TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)             AGE
service/example-linkerd              ClusterIP   10.108.124.47   <none>        11211/TCP           63s

NAME                                 READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/example-linkerd      3/3     3            3           65s
deployment.apps/linkerd2-operator    1/1     1            1           70s

NAME                                            DESIRED   CURRENT   READY   AGE
replicaset.apps/example-linkerd-7c4df9b7b4      3         3         3       65s
replicaset.apps/linkerd2-operator-56f54d84bf    1         1         1       70s
```

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
