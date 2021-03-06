/*
Copyright 2018 The Kubernetes Authors.

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

package client

import (
	"context"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// client is a client.Client that reads and writes directly from/to an API server.  It lazily initializes
// new clients at the time they are used, and caches the client.
type unstructuredClient struct {
	cache      *clientCache
	paramCodec runtime.ParameterCodec
}

// Create implements client.Client
func (uc *unstructuredClient) Create(ctx context.Context, obj runtime.Object, opts ...CreateOption) error {
	u, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("unstructured client did not understand object: %T", obj)
	}

	gvk := u.GroupVersionKind()

	o, err := uc.cache.getObjMeta(obj)
	if err != nil {
		return err
	}

	createOpts := &CreateOptions{}
	createOpts.ApplyOptions(opts)
	result := o.Post().
		NamespaceIfScoped(o.GetNamespace(), o.isNamespaced()).
		Resource(o.resource()).
		Body(obj).
		VersionedParams(createOpts.AsCreateOptions(), uc.paramCodec).
		Context(ctx).
		Do().
		Into(obj)

	u.SetGroupVersionKind(gvk)
	return result
}

// Update implements client.Client
func (uc *unstructuredClient) Update(ctx context.Context, obj runtime.Object, opts ...UpdateOption) error {
	u, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("unstructured client did not understand object: %T", obj)
	}

	gvk := u.GroupVersionKind()

	o, err := uc.cache.getObjMeta(obj)
	if err != nil {
		return err
	}

	updateOpts := UpdateOptions{}
	updateOpts.ApplyOptions(opts)
	result := o.Put().
		NamespaceIfScoped(o.GetNamespace(), o.isNamespaced()).
		Resource(o.resource()).
		Name(o.GetName()).
		Body(obj).
		VersionedParams(updateOpts.AsUpdateOptions(), uc.paramCodec).
		Context(ctx).
		Do().
		Into(obj)

	u.SetGroupVersionKind(gvk)
	return result
}

// Delete implements client.Client
func (uc *unstructuredClient) Delete(ctx context.Context, obj runtime.Object, opts ...DeleteOption) error {
	_, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("unstructured client did not understand object: %T", obj)
	}

	o, err := uc.cache.getObjMeta(obj)
	if err != nil {
		return err
	}

	deleteOpts := DeleteOptions{}
	deleteOpts.ApplyOptions(opts)
	return o.Delete().
		NamespaceIfScoped(o.GetNamespace(), o.isNamespaced()).
		Resource(o.resource()).
		Name(o.GetName()).
		Body(deleteOpts.AsDeleteOptions()).
		Context(ctx).
		Do().
		Error()
}

// DeleteAllOf implements client.Client
func (uc *unstructuredClient) DeleteAllOf(ctx context.Context, obj runtime.Object, opts ...DeleteAllOfOption) error {
	_, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("unstructured client did not understand object: %T", obj)
	}

	o, err := uc.cache.getObjMeta(obj)
	if err != nil {
		return err
	}

	deleteAllOfOpts := DeleteAllOfOptions{}
	deleteAllOfOpts.ApplyOptions(opts)
	return o.Delete().
		NamespaceIfScoped(deleteAllOfOpts.ListOptions.Namespace, o.isNamespaced()).
		Resource(o.resource()).
		VersionedParams(deleteAllOfOpts.AsListOptions(), uc.paramCodec).
		Body(deleteAllOfOpts.AsDeleteOptions()).
		Context(ctx).
		Do().
		Error()
}

// Patch implements client.Client
func (uc *unstructuredClient) Patch(ctx context.Context, obj runtime.Object, patch Patch, opts ...PatchOption) error {
	_, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("unstructured client did not understand object: %T", obj)
	}

	o, err := uc.cache.getObjMeta(obj)
	if err != nil {
		return err
	}

	data, err := patch.Data(obj)
	if err != nil {
		return err
	}

	patchOpts := &PatchOptions{}
	return o.Patch(patch.Type()).
		NamespaceIfScoped(o.GetNamespace(), o.isNamespaced()).
		Resource(o.resource()).
		Name(o.GetName()).
		VersionedParams(patchOpts.ApplyOptions(opts).AsPatchOptions(), uc.paramCodec).
		Body(data).
		Context(ctx).
		Do().
		Into(obj)
}

// Get implements client.Client
func (uc *unstructuredClient) Get(ctx context.Context, key ObjectKey, obj runtime.Object) error {
	u, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("unstructured client did not understand object: %T", obj)
	}

	gvk := u.GroupVersionKind()

	r, err := uc.cache.getResource(obj)
	if err != nil {
		return err
	}

	result := r.Get().
		NamespaceIfScoped(key.Namespace, r.isNamespaced()).
		Resource(r.resource()).
		Context(ctx).
		Name(key.Name).
		Do().
		Into(obj)

	u.SetGroupVersionKind(gvk)

	return result
}

// List implements client.Client
func (uc *unstructuredClient) List(ctx context.Context, obj runtime.Object, opts ...ListOption) error {
	u, ok := obj.(*unstructured.UnstructuredList)
	if !ok {
		return fmt.Errorf("unstructured client did not understand object: %T", obj)
	}

	gvk := u.GroupVersionKind()
	if strings.HasSuffix(gvk.Kind, "List") {
		gvk.Kind = gvk.Kind[:len(gvk.Kind)-4]
	}

	listOpts := ListOptions{}
	listOpts.ApplyOptions(opts)

	r, err := uc.cache.getResource(obj)
	if err != nil {
		return err
	}

	return r.Get().
		NamespaceIfScoped(listOpts.Namespace, r.isNamespaced()).
		Resource(r.resource()).
		VersionedParams(listOpts.AsListOptions(), uc.paramCodec).
		Context(ctx).
		Do().
		Into(obj)
}

func (uc *unstructuredClient) UpdateStatus(ctx context.Context, obj runtime.Object, opts ...UpdateOption) error {
	_, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("unstructured client did not understand object: %T", obj)
	}

	o, err := uc.cache.getObjMeta(obj)
	if err != nil {
		return err
	}

	return o.Put().
		NamespaceIfScoped(o.GetNamespace(), o.isNamespaced()).
		Resource(o.resource()).
		Name(o.GetName()).
		SubResource("status").
		Body(obj).
		VersionedParams((&UpdateOptions{}).ApplyOptions(opts).AsUpdateOptions(), uc.paramCodec).
		Context(ctx).
		Do().
		Into(obj)
}

func (uc *unstructuredClient) PatchStatus(ctx context.Context, obj runtime.Object, patch Patch, opts ...PatchOption) error {
	u, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("unstructured client did not understand object: %T", obj)
	}

	gvk := u.GroupVersionKind()

	o, err := uc.cache.getObjMeta(obj)
	if err != nil {
		return err
	}

	data, err := patch.Data(obj)
	if err != nil {
		return err
	}

	patchOpts := &PatchOptions{}
	result := o.Patch(patch.Type()).
		NamespaceIfScoped(o.GetNamespace(), o.isNamespaced()).
		Resource(o.resource()).
		Name(o.GetName()).
		SubResource("status").
		Body(data).
		VersionedParams(patchOpts.ApplyOptions(opts).AsPatchOptions(), uc.paramCodec).
		Context(ctx).
		Do().
		Into(u)

	u.SetGroupVersionKind(gvk)
	return result
}
