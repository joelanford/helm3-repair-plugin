module github.com/joelanford/helm3-repair-plugin

go 1.12

require (
	contrib.go.opencensus.io/exporter/ocagent v0.2.0 // indirect
	github.com/Azure/go-autorest v11.2.8+incompatible // indirect
	github.com/MakeNowJust/heredoc v0.0.0-20171113091838-e9091a26100e // indirect
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/semver v1.4.2 // indirect
	github.com/Masterminds/sprig v2.20.0+incompatible // indirect
	github.com/Microsoft/go-winio v0.4.12 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/census-instrumentation/opencensus-proto v0.1.0 // indirect
	github.com/deislabs/oras v0.6.0 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/docker/spdystream v0.0.0-20181023171402-6480d4af844c // indirect
	github.com/emicklei/go-restful v2.8.0+incompatible // indirect
	github.com/evanphx/json-patch v3.0.0+incompatible // indirect
	github.com/ghodss/yaml v1.0.0
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/golang/groupcache v0.0.0-20181024230925-c65c006176ff // indirect
	github.com/google/btree v1.0.0 // indirect
	github.com/google/uuid v1.1.0 // indirect
	github.com/googleapis/gnostic v0.2.0 // indirect
	github.com/gosuri/uitable v0.0.1 // indirect
	github.com/gotestyourself/gotestyourself v2.2.0+incompatible // indirect
	github.com/gregjones/httpcache v0.0.0-20181110185634-c63ab54fda8f // indirect
	github.com/huandu/xstrings v1.2.0 // indirect
	github.com/imdario/mergo v0.3.6 // indirect
	github.com/json-iterator/go v1.1.5 // indirect
	github.com/mattn/go-runewidth v0.0.4 // indirect
	github.com/mitchellh/go-wordwrap v1.0.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/pmezard/go-difflib v1.0.0
	github.com/sirupsen/logrus v1.4.1 // indirect
	github.com/spf13/cobra v0.0.4
	github.com/spf13/pflag v1.0.3
	github.com/xeipuuv/gojsonpointer v0.0.0-20180127040702-4e3ac2762d5f // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.1.0 // indirect
	go.opencensus.io v0.18.0 // indirect
	google.golang.org/api v0.0.0-20181221000618-65a46cafb132 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/square/go-jose.v2 v2.3.0 // indirect
	helm.sh/helm v0.0.0-20190820153828-fba311ba2362
	k8s.io/apimachinery v0.0.0
	k8s.io/cli-runtime v0.0.0
	k8s.io/client-go v0.0.0
	k8s.io/klog v0.3.1
	k8s.io/kube-openapi v0.0.0-20190816220812-743ec37842bf // indirect
	k8s.io/kubernetes v1.15.1 // indirect
	k8s.io/utils v0.0.0-20190809000727-6c36bc71fc4a // indirect
)

replace (
	github.com/docker/docker => github.com/docker/docker v0.0.0-20180612054059-a9fbbdc8dd87 // indirect
	k8s.io/api => k8s.io/api v0.0.0-20190718183219-b59d8169aab5
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190718185103-d1ef975d28ce
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190612205821-1799e75a0719
	k8s.io/apiserver => k8s.io/apiserver v0.0.0-20190718184206-a1aa83af71a7
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20190718185405-0ce9869d0015
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190718183610-8e956561bbf5
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20190718190308-f8e43aa19282
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.0.0-20190718190146-f7b0473036f9
	k8s.io/code-generator => k8s.io/code-generator v0.0.0-20190612205613-18da4a14b22b
	k8s.io/component-base => k8s.io/component-base v0.0.0-20190718183727-0ececfbe9772
	k8s.io/cri-api => k8s.io/cri-api v0.0.0-20190531030430-6117653b35f1
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.0.0-20190718190424-bef8d46b95de
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20190718184434-a064d4d1ed7a
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.0.0-20190718190030-ea930fedc880
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.0.0-20190718185641-5233cb7cb41e
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.0.0-20190718185913-d5429d807831
	k8s.io/kubelet => k8s.io/kubelet v0.0.0-20190718185757-9b45f80d5747
	k8s.io/kubernetes => k8s.io/kubernetes v1.15.1
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.0.0-20190718190548-039b99e58dbd
	k8s.io/metrics => k8s.io/metrics v0.0.0-20190718185242-1e1642704fe6
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.0.0-20190718184639-baafa86838c0
)
