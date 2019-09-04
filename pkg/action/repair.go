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

package action

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/pmezard/go-difflib/difflib"
	"helm.sh/helm/pkg/action"
	"helm.sh/helm/pkg/release"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	apitypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/kubernetes/scheme"
)

// Repair is the action for repairing a given release's deployed resources.
//
// It provides the implementation of 'helm repair'.
type Repair struct {
	cfg *action.Configuration

	DryRun bool
}

// NewRepair creates a new Repair object with the given configuration.
func NewRepair(cfg *action.Configuration) *Repair {
	return &Repair{
		cfg: cfg,
	}
}

// Run executes 'helm repair' against the given release.
func (r *Repair) Run(name string) (*release.Release, bool, error) {
	get := action.NewGet(r.cfg)
	rel, err := get.Run(name)
	if err != nil {
		return nil, false, err
	}

	resources, err := r.cfg.KubeClient.Build(bytes.NewBufferString(rel.Manifest))
	if err != nil {
		return nil, false, errors.Wrap(err, "unable to build kubernetes objects from release manifest")
	}

	didRepair := false
	err = resources.Visit(func(target *resource.Info, err error) error {
		if err != nil {
			return err
		}

		helper := resource.NewHelper(target.Client, target.Mapping)
		current, err := helper.Get(target.Namespace, target.Name, target.Export)
		if apierrors.IsNotFound(err) {
			if _, err := helper.Create(target.Namespace, true, target.Object, nil); err != nil {
				return errors.Wrapf(err, "unable to recreate release resource %q", target.Name)
			}
			didRepair = true
			return nil
		} else if err != nil {
			return errors.Wrapf(err, "unable to get release resource %q", target.Name)
		}

		patch, err := generateStrategicMergePatch(current, target)
		if err != nil {
			return errors.Wrapf(err, "unable to generate strategic merge patch for release resource %q", target.Name)
		}

		if string(patch) == "{}" {
			return nil
		}

		patchOpts := &metav1.PatchOptions{}
		if r.DryRun {
			patchOpts.DryRun = []string{metav1.DryRunAll}
		}
		repaired, err := helper.Patch(target.Namespace, target.Name, apitypes.StrategicMergePatchType, patch, patchOpts)
		if err != nil {
			return errors.Wrapf(err, "unable to patch release resource %q", target.Name)
		}
		didRepair = true
		if err := printDiff(current, repaired); err != nil {
			return errors.Wrapf(err, "unable to print diff for resource %q", target.Name)
		}
		return nil
	})
	return rel, didRepair, err
}

func generateStrategicMergePatch(current runtime.Object, target *resource.Info) ([]byte, error) {
	currentJSON, err := json.Marshal(current)
	if err != nil {
		return nil, err
	}
	targetJSON, err := json.Marshal(target.Object)
	if err != nil {
		return nil, err
	}
	targetVersioned := asVersioned(target)
	patchMeta, err := strategicpatch.NewPatchMetaFromStruct(targetVersioned)
	if err != nil {
		return nil, err
	}
	return strategicpatch.CreateThreeWayMergePatch(targetJSON, targetJSON, currentJSON, patchMeta, true)
}

func asVersioned(info *resource.Info) runtime.Object {
	var gv = runtime.GroupVersioner(schema.GroupVersions(scheme.Scheme.PrioritizedVersionsAllGroups()))
	if info.Mapping != nil {
		gv = info.Mapping.GroupVersionKind.GroupVersion()
	}
	obj, _ := runtime.ObjectConvertor(scheme.Scheme).ConvertToVersion(info.Object, gv)
	return obj
}

func printDiff(a, b runtime.Object) (err error) {
	yamlA, err := yaml.Marshal(a)
	if err != nil {
		return err
	}
	ua := unstructured.Unstructured{}
	ua.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(a)
	if err != nil {
		return err
	}

	yamlB, err := yaml.Marshal(b)
	if err != nil {
		return err
	}
	ub := unstructured.Unstructured{}
	ub.Object, err = runtime.DefaultUnstructuredConverter.ToUnstructured(b)
	if err != nil {
		return err
	}

	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(yamlA)),
		B:        difflib.SplitLines(string(yamlB)),
		FromFile: fmt.Sprintf("%s.%s.%s.%s (current)", strings.ReplaceAll(ua.GroupVersionKind().GroupVersion().String(), "/", "."), ua.GroupVersionKind().Kind, ua.GetNamespace(), ua.GetName()),
		ToFile:   fmt.Sprintf("%s.%s.%s.%s (target)", strings.ReplaceAll(ub.GroupVersionKind().GroupVersion().String(), "/", "."), ub.GroupVersionKind().Kind, ub.GetNamespace(), ub.GetName()),
		Context:  3,
	}
	text, err := difflib.GetUnifiedDiffString(diff)
	if err != nil {
		return err
	}
	fmt.Println(text)
	return nil
}
