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
--- apps.v1beta2.Deployment.default.example-my-chart (current)
+++ apps.v1beta2.Deployment.default.example-my-chart (target)
@@ -4,7 +4,7 @@
   annotations:
     deployment.kubernetes.io/revision: "1"
   creationTimestamp: "2019-09-04T17:56:35Z"
-  generation: 2
+  generation: 3
   labels:
     app.kubernetes.io/instance: example
     app.kubernetes.io/managed-by: Helm
@@ -13,12 +13,12 @@
     helm.sh/chart: my-chart-0.1.0
   name: example-my-chart
   namespace: default
-  resourceVersion: "1508776"
+  resourceVersion: "1508785"
   selfLink: /apis/apps/v1beta2/namespaces/default/deployments/example-my-chart
   uid: 56414197-cf3d-11e9-bdc6-049226c40916
 spec:
   progressDeadlineSeconds: 600
-  replicas: 3
+  replicas: 1
   revisionHistoryLimit: 10
   selector:
     matchLabels:

release "example" repaired

## Check that the number of replicas have been reset
$ kubectl get deployments.apps example-my-chart -o jsonpath={.spec.replicas}
1

## See that the release revision is still 1
$ helm list
NAME   	NAMESPACE	REVISION	UPDATED                             	STATUS  	CHART
example	default  	1       	2019-08-23 22:12:58.621967 -0400 EDT	deployed	my-chart-0.1.0
```

