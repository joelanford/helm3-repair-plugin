# Helm 3 Repair Plugin

The `repair` plugin repairs Helm releases whose resources have been changed outside of Helm. It performs a three-way merge to restore Helm-defined values.

## Installation

```console
go get -d https://github.com/joelanford/helm3-repair-plugin
cd $GOPATH/src/github.com/joelanford/helm3-repair-plugin
make
helm plugin install .
```

## Example

```console
## Create an example chart
$ helm create my-chart
Creating my-chart

## Install a release
$ helm install example my-chart
NAME: example
LAST DEPLOYED: 2019-08-23 22:12:58.621967 -0400 EDT m=+0.123583973
NAMESPACE: default
STATUS: deployed

NOTES:
1. Get the application URL by running these commands:
  export POD_NAME=$(kubectl get pods -l "app=my-chart,release=example" -o jsonpath="{.items[0].metadata.name}")
  echo "Visit http://127.0.0.1:8080 to use your application"
  kubectl port-forward $POD_NAME 8080:80

## Check the number of replicas in the release deployment
$ kubectl get deployments.apps example-my-chart -o jsonpath={.spec.replicas}
1

## Change the number of replicas in the release deployment
$ kubectl patch deployments.apps example-my-chart -p '{"spec":{"replicas":3}}'
deployment.apps/example-my-chart patched

## Repair the release
$ helm repair example
release "example" repaired

## Check that the number of replicas have been reset
$ kubectl get deployments.apps example-my-chart -o jsonpath={.spec.replicas}
1

## See that the release revision is still 1
$ helm list
NAME   	NAMESPACE	REVISION	UPDATED                             	STATUS  	CHART
example	default  	1       	2019-08-23 22:12:58.621967 -0400 EDT	deployed	my-chart-0.1.0
```

