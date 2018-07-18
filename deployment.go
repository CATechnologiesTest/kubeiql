// Copyright 2018 Yipee.io
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

// Top level Kubernetes replicated controller. Deployments are built out
// of ReplicaSets.
type deployment struct {
	Metadata    metadata
	Owner       resource
	RootOwner   resource
	ReplicaSets *[]replicaSet
}

type deploymentResolver struct {
	ctx context.Context
	d   deployment
}

// Translate unmarshalled json into a deployment object
func mapToDeployment(
	ctx context.Context,
	jsonObj map[string]interface{}) deployment {
	meta :=
		mapToMetadata(ctx, getNamespace(jsonObj), mapItem(jsonObj, "metadata"))
	return deployment{meta, nil, nil, nil}
}

// Retrieve the ReplicaSets comprising the deployment
func getReplicaSets(ctx context.Context, d deployment) *[]replicaSet {
	depName := d.Metadata.Name
	depNamePrefix := depName + "-"
	depNamespace := d.Metadata.Namespace

	rsets := getAllK8sObjsOfKindInNamespace(
		ctx,
		"ReplicaSet",
		depNamespace,
		func(jobj map[string]interface{}) bool {
			return (strings.HasPrefix(getName(jobj), depNamePrefix) &&
				hasMatchingOwner(jobj, depName, DeploymentKind))
		})

	results := make([]replicaSet, len(rsets))

	for idx, rs := range rsets {
		rsr := rs.(*replicaSetResolver)
		results[idx] = rsr.r
	}

	return &results
}

// Resource method implementations
func (r *deploymentResolver) Kind() string {
	return DeploymentKind
}

func (r *deploymentResolver) Metadata() *metadataResolver {
	return &metadataResolver{r.ctx, r.d.Metadata}
}

func (r *deploymentResolver) Owner() *resourceResolver {
	return &resourceResolver{r.ctx, &deploymentResolver{r.ctx, r.d}}
}

func (r *deploymentResolver) RootOwner() *resourceResolver {
	return &resourceResolver{r.ctx, &deploymentResolver{r.ctx, r.d}}
}

// Resolve child ReplicaSets
func (r *deploymentResolver) ReplicaSets() []*replicaSetResolver {
	if r.d.ReplicaSets == nil {
		r.d.ReplicaSets = getReplicaSets(r.ctx, r.d)
	}

	var res []*replicaSetResolver
	for _, rs := range *r.d.ReplicaSets {
		res = append(res, &replicaSetResolver{r.ctx, rs})
	}
	if res == nil {
		res = make([]*replicaSetResolver, 0)
	}
	return res
}
