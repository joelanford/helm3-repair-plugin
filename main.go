/*
Copyright The Helm Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"helm.sh/helm/cmd/helm/require"
	helmaction "helm.sh/helm/pkg/action"
	"helm.sh/helm/pkg/cli"
	"helm.sh/helm/pkg/kube"
	"helm.sh/helm/pkg/storage"
	"helm.sh/helm/pkg/storage/driver"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/klog"

	// Import to initialize client auth plugins.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/joelanford/helm3-repair-plugin/pkg/action"
)

var (
	settings   cli.EnvSettings
	config     genericclioptions.RESTClientGetter
	configOnce sync.Once
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func debug(format string, v ...interface{}) {
	if settings.Debug {
		format = fmt.Sprintf("[debug] %s\n", format)
		log.Output(2, fmt.Sprintf(format, v...))
	}
}

func initKubeLogs() {
	pflag.CommandLine.SetNormalizeFunc(wordSepNormalizeFunc)
	gofs := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(gofs)
	pflag.CommandLine.AddGoFlagSet(gofs)
	pflag.CommandLine.Set("logtostderr", "true")
}

func newRootCmd(cfg *helmaction.Configuration, w io.Writer, args []string) *cobra.Command {
	client := action.NewRepair(cfg)

	cmd := &cobra.Command{
		Use:   "helm repair [NAME]",
		Short: "Repair a release that has been modified outside of helm",
		Long:  "Repair a release that has been modified outside of helm",
		Args:  require.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			res, didRepair, err := client.Run(args[0])
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				os.Exit(1)
			}
			prefix := ""
			if client.DryRun {
				prefix = "DRY RUN: "
			}
			if didRepair {
				fmt.Printf("%srelease %q repaired\n", prefix, res.Name)
			} else {
				fmt.Printf("%srelease %q already up-to-date\n", prefix, res.Name)
			}
		},
	}

	cmd.Flags().BoolVar(&client.DryRun, "dry-run", false, "simulate a repair")

	return cmd
}

func main() {
	initKubeLogs()

	actionConfig := new(helmaction.Configuration)
	cmd := newRootCmd(actionConfig, os.Stdout, os.Args[1:])

	// Initialize the rest of the actionConfig
	initActionConfig(actionConfig, false)

	if err := cmd.Execute(); err != nil {
		debug("%+v", err)
		os.Exit(1)
	}
}

func initActionConfig(actionConfig *helmaction.Configuration, allNamespaces bool) {
	kc := kube.New(kubeConfig())
	kc.Log = debug

	clientset, err := kc.Factory.KubernetesClientSet()
	if err != nil {
		// TODO return error
		log.Fatal(err)
	}
	var namespace string
	if !allNamespaces {
		namespace = getNamespace()
	}

	var store *storage.Storage
	switch os.Getenv("HELM_DRIVER") {
	case "secret", "secrets", "":
		d := driver.NewSecrets(clientset.CoreV1().Secrets(namespace))
		d.Log = debug
		store = storage.Init(d)
	case "configmap", "configmaps":
		d := driver.NewConfigMaps(clientset.CoreV1().ConfigMaps(namespace))
		d.Log = debug
		store = storage.Init(d)
	case "memory":
		d := driver.NewMemory()
		store = storage.Init(d)
	default:
		// Not sure what to do here.
		panic("Unknown driver in HELM_DRIVER: " + os.Getenv("HELM_DRIVER"))
	}

	actionConfig.RESTClientGetter = kubeConfig()
	actionConfig.KubeClient = kc
	actionConfig.Releases = store
	actionConfig.Log = debug
}

func kubeConfig() genericclioptions.RESTClientGetter {
	configOnce.Do(func() {
		config = kube.GetConfig(settings.KubeConfig, settings.KubeContext, settings.Namespace)
	})
	return config
}

func getNamespace() string {
	if ns, _, err := kubeConfig().ToRawKubeConfigLoader().Namespace(); err == nil {
		return ns
	}
	return "default"
}

// wordSepNormalizeFunc changes all flags that contain "_" separators
func wordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	return pflag.NormalizedName(strings.ReplaceAll(name, "_", "-"))
}
