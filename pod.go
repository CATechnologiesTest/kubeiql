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
)

// The base Kubernetes component
type pod struct {
	Metadata  metadata
	Owner     resource
	RootOwner resource
}

type podResolver struct {
	ctx context.Context
	p   pod
}

// Translate unmarshalled json into a metadata object
func mapToPod(ctx context.Context, jsonObj map[string]interface{}) pod {
	placeholder := &ownerRef{ctx, jsonObj, nil}
	owner := placeholder
	rootOwner := placeholder
	meta :=
		mapToMetadata(ctx, getNamespace(jsonObj), mapItem(jsonObj, "metadata"))
	return pod{meta, owner, rootOwner}
}

// Resource method implementations
func (r *podResolver) Kind() string {
	return PodKind
}

func (r *podResolver) Metadata() *metadataResolver {
	meta := r.p.Metadata
	return &metadataResolver{r.ctx, meta}
}

func (r *podResolver) Owner() *resourceResolver {
	if oref, ok := r.p.Owner.(*ownerRef); ok {
		r.p.Owner = getOwner(oref.ctx, oref.ref)
	}
	return &resourceResolver{r.ctx, r.p.Owner}
}

func (r *podResolver) RootOwner() *resourceResolver {
	if oref, ok := r.p.RootOwner.(*ownerRef); ok {
		r.p.RootOwner = getRootOwner(oref.ctx, oref.ref)
	}
	return &resourceResolver{r.ctx, r.p.RootOwner}
}
