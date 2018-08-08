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
	"strings"
)

// ReplicaSets manage replicated pods
type replicaSet struct {
	Metadata  metadata
	Owner     resource
	RootOwner resource
	Pods      *[]pod
}

type replicaSetResolver struct {
	ctx context.Context
	r   replicaSet
}

// Translate unmarshalled json into a deployment object
func mapToReplicaSet(
	ctx context.Context,
	jsonObj JsonObject) replicaSet {
	placeholder := &ownerRef{ctx, jsonObj, nil}
	owner := placeholder
	rootOwner := placeholder
	meta :=
		mapToMetadata(ctx, getNamespace(jsonObj), mapItem(jsonObj, "metadata"))
	return replicaSet{meta, owner, rootOwner, nil}
}

// ReplicaSets have pods as children
func getReplicaSetPods(ctx context.Context, r replicaSet) *[]pod {
	rsName := *r.Metadata.Name
	rsNamePrefix := rsName + "-"
	rsNamespace := *r.Metadata.Namespace

	pset := getAllK8sObjsOfKindInNamespace(
		ctx,
		PodKind,
		rsNamespace,
		func(jobj JsonObject) bool {
			return (strings.HasPrefix(getName(jobj), rsNamePrefix) &&
				hasMatchingOwner(jobj, rsName, ReplicaSetKind))
		})

	results := make([]pod, len(pset))

	for idx, p := range pset {
		pr := p.(*podResolver)
		results[idx] = pr.p
	}

	return &results
}

// Resource method implementations
func (r *replicaSetResolver) Kind() string {
	return ReplicaSetKind
}

func (r *replicaSetResolver) Metadata() metadataResolver {
	return metadataResolver{r.ctx, r.r.Metadata}
}

func (r *replicaSetResolver) Owner() *resourceResolver {
	if oref, ok := r.r.Owner.(*ownerRef); ok {
		r.r.Owner = getOwner(oref.ctx, oref.ref)
	}
	return &resourceResolver{r.ctx, r.r.Owner}
}

func (r *replicaSetResolver) RootOwner() *resourceResolver {
	if oref, ok := r.r.Owner.(*ownerRef); ok {
		r.r.Owner = getOwner(oref.ctx, oref.ref)
	}
	return &resourceResolver{r.ctx, r.r.RootOwner}
}

// Resolve child Pods
func (r *replicaSetResolver) Pods() []*podResolver {
	if r.r.Pods == nil {
		r.r.Pods = getReplicaSetPods(r.ctx, r.r)
	}

	var res []*podResolver
	for _, p := range *r.r.Pods {
		res = append(res, &podResolver{r.ctx, p})
	}
	if res == nil {
		res = make([]*podResolver, 0)
	}
	return res
}
