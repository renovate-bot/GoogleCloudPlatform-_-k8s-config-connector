// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// *** DISCLAIMER ***
// Config Connector's go-client for CRDs is currently in ALPHA, which means
// that future versions of the go-client may include breaking changes.
// Please try it out and give us feedback!

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/GoogleCloudPlatform/k8s-config-connector/pkg/clients/generated/apis/workstations/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeWorkstationConfigs implements WorkstationConfigInterface
type FakeWorkstationConfigs struct {
	Fake *FakeWorkstationsV1alpha1
	ns   string
}

var workstationconfigsResource = v1alpha1.SchemeGroupVersion.WithResource("workstationconfigs")

var workstationconfigsKind = v1alpha1.SchemeGroupVersion.WithKind("WorkstationConfig")

// Get takes name of the workstationConfig, and returns the corresponding workstationConfig object, and an error if there is any.
func (c *FakeWorkstationConfigs) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.WorkstationConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(workstationconfigsResource, c.ns, name), &v1alpha1.WorkstationConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.WorkstationConfig), err
}

// List takes label and field selectors, and returns the list of WorkstationConfigs that match those selectors.
func (c *FakeWorkstationConfigs) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.WorkstationConfigList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(workstationconfigsResource, workstationconfigsKind, c.ns, opts), &v1alpha1.WorkstationConfigList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.WorkstationConfigList{ListMeta: obj.(*v1alpha1.WorkstationConfigList).ListMeta}
	for _, item := range obj.(*v1alpha1.WorkstationConfigList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested workstationConfigs.
func (c *FakeWorkstationConfigs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(workstationconfigsResource, c.ns, opts))

}

// Create takes the representation of a workstationConfig and creates it.  Returns the server's representation of the workstationConfig, and an error, if there is any.
func (c *FakeWorkstationConfigs) Create(ctx context.Context, workstationConfig *v1alpha1.WorkstationConfig, opts v1.CreateOptions) (result *v1alpha1.WorkstationConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(workstationconfigsResource, c.ns, workstationConfig), &v1alpha1.WorkstationConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.WorkstationConfig), err
}

// Update takes the representation of a workstationConfig and updates it. Returns the server's representation of the workstationConfig, and an error, if there is any.
func (c *FakeWorkstationConfigs) Update(ctx context.Context, workstationConfig *v1alpha1.WorkstationConfig, opts v1.UpdateOptions) (result *v1alpha1.WorkstationConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(workstationconfigsResource, c.ns, workstationConfig), &v1alpha1.WorkstationConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.WorkstationConfig), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeWorkstationConfigs) UpdateStatus(ctx context.Context, workstationConfig *v1alpha1.WorkstationConfig, opts v1.UpdateOptions) (*v1alpha1.WorkstationConfig, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(workstationconfigsResource, "status", c.ns, workstationConfig), &v1alpha1.WorkstationConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.WorkstationConfig), err
}

// Delete takes name of the workstationConfig and deletes it. Returns an error if one occurs.
func (c *FakeWorkstationConfigs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(workstationconfigsResource, c.ns, name, opts), &v1alpha1.WorkstationConfig{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeWorkstationConfigs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(workstationconfigsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.WorkstationConfigList{})
	return err
}

// Patch applies the patch and returns the patched workstationConfig.
func (c *FakeWorkstationConfigs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.WorkstationConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(workstationconfigsResource, c.ns, name, pt, data, subresources...), &v1alpha1.WorkstationConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.WorkstationConfig), err
}
