/*
Copyright 2022 Contributors to The HwameiStor project.

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/hwameistor/local-disk-manager/pkg/apis/hwameistor/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeLocalDiskClaims implements LocalDiskClaimInterface
type FakeLocalDiskClaims struct {
	Fake *FakeHwameistorV1alpha1
}

var localdiskclaimsResource = schema.GroupVersionResource{Group: "hwameistor.io", Version: "v1alpha1", Resource: "localdiskclaims"}

var localdiskclaimsKind = schema.GroupVersionKind{Group: "hwameistor.io", Version: "v1alpha1", Kind: "LocalDiskClaim"}

// Get takes name of the localDiskClaim, and returns the corresponding localDiskClaim object, and an error if there is any.
func (c *FakeLocalDiskClaims) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.LocalDiskClaim, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(localdiskclaimsResource, name), &v1alpha1.LocalDiskClaim{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.LocalDiskClaim), err
}

// List takes label and field selectors, and returns the list of LocalDiskClaims that match those selectors.
func (c *FakeLocalDiskClaims) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.LocalDiskClaimList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(localdiskclaimsResource, localdiskclaimsKind, opts), &v1alpha1.LocalDiskClaimList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.LocalDiskClaimList{ListMeta: obj.(*v1alpha1.LocalDiskClaimList).ListMeta}
	for _, item := range obj.(*v1alpha1.LocalDiskClaimList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested localDiskClaims.
func (c *FakeLocalDiskClaims) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(localdiskclaimsResource, opts))
}

// Create takes the representation of a localDiskClaim and creates it.  Returns the server's representation of the localDiskClaim, and an error, if there is any.
func (c *FakeLocalDiskClaims) Create(ctx context.Context, localDiskClaim *v1alpha1.LocalDiskClaim, opts v1.CreateOptions) (result *v1alpha1.LocalDiskClaim, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(localdiskclaimsResource, localDiskClaim), &v1alpha1.LocalDiskClaim{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.LocalDiskClaim), err
}

// Update takes the representation of a localDiskClaim and updates it. Returns the server's representation of the localDiskClaim, and an error, if there is any.
func (c *FakeLocalDiskClaims) Update(ctx context.Context, localDiskClaim *v1alpha1.LocalDiskClaim, opts v1.UpdateOptions) (result *v1alpha1.LocalDiskClaim, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(localdiskclaimsResource, localDiskClaim), &v1alpha1.LocalDiskClaim{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.LocalDiskClaim), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeLocalDiskClaims) UpdateStatus(ctx context.Context, localDiskClaim *v1alpha1.LocalDiskClaim, opts v1.UpdateOptions) (*v1alpha1.LocalDiskClaim, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(localdiskclaimsResource, "status", localDiskClaim), &v1alpha1.LocalDiskClaim{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.LocalDiskClaim), err
}

// Delete takes name of the localDiskClaim and deletes it. Returns an error if one occurs.
func (c *FakeLocalDiskClaims) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(localdiskclaimsResource, name), &v1alpha1.LocalDiskClaim{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeLocalDiskClaims) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(localdiskclaimsResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.LocalDiskClaimList{})
	return err
}

// Patch applies the patch and returns the patched localDiskClaim.
func (c *FakeLocalDiskClaims) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.LocalDiskClaim, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(localdiskclaimsResource, name, pt, data, subresources...), &v1alpha1.LocalDiskClaim{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.LocalDiskClaim), err
}
