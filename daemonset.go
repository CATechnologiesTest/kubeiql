// Copyright (c) 2018 CA. All rights reserved.
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

package main

import (
	"context"
	//	"fmt"
	"sort"
	"strings"
)

// DaemonSets place a pod on each server
type daemonSet struct {
	Metadata  metadata
	Owner     resource
	RootOwner resource
	Pods      *[]pod
}

type daemonSetResolver struct {
	ctx context.Context
	d   daemonSet
}

// Translate unmarshalled json into a deployment object
func mapToDaemonSet(
	ctx context.Context,
	jsonObj JsonObject) daemonSet {
	placeholder := &ownerRef{ctx, jsonObj, nil}
	owner := placeholder
	rootOwner := placeholder
	meta :=
		mapToMetadata(ctx, getNamespace(jsonObj), mapItem(jsonObj, "metadata"))
	return daemonSet{meta, owner, rootOwner, nil}
}

// DaemonSets have pods as children
func getDaemonSetPods(ctx context.Context, d daemonSet) *[]pod {
	dsName := *d.Metadata.Name
	dsNamePrefix := dsName + "-"
	dsNamespace := *d.Metadata.Namespace

	pset := getAllK8sObjsOfKindInNamespace(
		ctx,
		PodKind,
		dsNamespace,
		func(jobj JsonObject) bool {
			return (strings.HasPrefix(getName(jobj), dsNamePrefix) &&
				hasMatchingOwner(jobj, dsName, DaemonSetKind))
		})

	results := make([]pod, len(pset))

	for idx, p := range pset {
		pr := p.(*podResolver)
		results[idx] = pr.p
	}

	sort.Slice(
		results,
		func(i, j int) bool {
			return *results[i].Metadata.Name < *results[j].Metadata.Name
		})

	return &results
}

// Resource method implementations
func (r *daemonSetResolver) Kind() string {
	return DaemonSetKind
}

func (r *daemonSetResolver) Metadata() metadataResolver {
	return metadataResolver{r.ctx, r.d.Metadata}
}

func (r *daemonSetResolver) Owner() *resourceResolver {
	if oref, ok := r.d.Owner.(*ownerRef); ok {
		r.d.Owner = getOwner(oref.ctx, oref.ref)
	}
	return &resourceResolver{r.ctx, r.d.Owner}
}

func (r *daemonSetResolver) RootOwner() *resourceResolver {
	if oref, ok := r.d.Owner.(*ownerRef); ok {
		r.d.Owner = getOwner(oref.ctx, oref.ref)
	}
	return &resourceResolver{r.ctx, r.d.RootOwner}
}

// Resolve child Pods
func (r *daemonSetResolver) Pods() []*podResolver {
	if r.d.Pods == nil {
		r.d.Pods = getDaemonSetPods(r.ctx, r.d)
	}

	var res []*podResolver
	for _, p := range *r.d.Pods {
		res = append(res, &podResolver{r.ctx, p})
	}
	if res == nil {
		res = make([]*podResolver, 0)
	}
	return res
}
