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

// StatefulSets manage pods that have dependencies
type statefulSet struct {
	Metadata  metadata
	Owner     resource
	RootOwner resource
	Pods      *[]pod
}

type statefulSetResolver struct {
	ctx context.Context
	s   statefulSet
}

// Translate unmarshalled json into a deployment object
func mapToStatefulSet(
	ctx context.Context,
	jsonObj JsonObject) statefulSet {
	placeholder := &ownerRef{ctx, jsonObj, nil}
	owner := placeholder
	rootOwner := placeholder
	meta :=
		mapToMetadata(ctx, getNamespace(jsonObj), mapItem(jsonObj, "metadata"))
	return statefulSet{meta, owner, rootOwner, nil}
}

// StatefulSets have pods as children
func getStatefulSetPods(ctx context.Context, s statefulSet) *[]pod {
	ssName := *s.Metadata.Name
	ssNamePrefix := ssName + "-"
	ssNamespace := *s.Metadata.Namespace

	pset := getAllK8sObjsOfKindInNamespace(
		ctx,
		PodKind,
		ssNamespace,
		func(jobj JsonObject) bool {
			return (strings.HasPrefix(getName(jobj), ssNamePrefix) &&
				hasMatchingOwner(jobj, ssName, StatefulSetKind))
		})

	results := make([]pod, len(pset))

	for idx, p := range pset {
		pr := p.(*podResolver)
		results[idx] = pr.p
	}

	return &results
}

// Resource method implementations
func (r *statefulSetResolver) Kind() string {
	return StatefulSetKind
}

func (r *statefulSetResolver) Metadata() metadataResolver {
	return metadataResolver{r.ctx, r.s.Metadata}
}

func (r *statefulSetResolver) Owner() *resourceResolver {
	if oref, ok := r.s.Owner.(*ownerRef); ok {
		r.s.Owner = getOwner(oref.ctx, oref.ref)
	}
	return &resourceResolver{r.ctx, r.s.Owner}
}

func (r *statefulSetResolver) RootOwner() *resourceResolver {
	if oref, ok := r.s.Owner.(*ownerRef); ok {
		r.s.Owner = getOwner(oref.ctx, oref.ref)
	}
	return &resourceResolver{r.ctx, r.s.RootOwner}
}

// Resolve child Pods
func (r *statefulSetResolver) Pods() []*podResolver {
	if r.s.Pods == nil {
		r.s.Pods = getStatefulSetPods(r.ctx, r.s)
	}

	var res []*podResolver
	for _, p := range *r.s.Pods {
		res = append(res, &podResolver{r.ctx, p})
	}
	if res == nil {
		res = make([]*podResolver, 0)
	}
	return res
}
