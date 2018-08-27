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
	"fmt"
)

// Represents all "active" components: (Pods, Deployments, DaemonSets,
// StatefulSets, ReplicaSets

type resource interface {
	Kind() string
	Metadata() metadataResolver
	Owner() *resourceResolver
	RootOwner() *resourceResolver
}

type resourceResolver struct {
	ctx context.Context
	r   resource
}

// Translate a map containing unmarshalled json into a resource instance.
func mapToResource(
	ctx context.Context,
	rMap JsonObject) resource {
	kind := getKind(rMap)

	switch kind {
	case DeploymentKind:
		mtd := mapToDeployment(ctx, rMap)
		return &deploymentResolver{ctx, mtd}
	case ReplicaSetKind:
		return &replicaSetResolver{ctx, mapToReplicaSet(ctx, rMap)}
	case DaemonSetKind:
		return &daemonSetResolver{ctx, mapToDaemonSet(ctx, rMap)}
	case StatefulSetKind:
		return &statefulSetResolver{ctx, mapToStatefulSet(ctx, rMap)}
	case PodKind:
		return &podResolver{ctx, mapToPod(ctx, rMap)}
	case ServiceKind:
		return &serviceResolver{ctx, mapToService(ctx, rMap)}
	}

	fmt.Printf("BAD KIND: %v\n", kind)
	return nil
}

// Turn an instance of a resource into one of its implementers
func (r *resourceResolver) ToPod() (*podResolver, bool) {
	c, ok := r.r.(*podResolver)
	return c, ok
}

func (r *resourceResolver) ToReplicaSet() (*replicaSetResolver, bool) {
	c, ok := r.r.(*replicaSetResolver)
	return c, ok
}

func (r *resourceResolver) ToDaemonSet() (*daemonSetResolver, bool) {
	c, ok := r.r.(*daemonSetResolver)
	return c, ok
}

func (r *resourceResolver) ToStatefulSet() (*statefulSetResolver, bool) {
	c, ok := r.r.(*statefulSetResolver)
	return c, ok
}

func (r *resourceResolver) ToDeployment() (*deploymentResolver, bool) {
	c, ok := r.r.(*deploymentResolver)
	return c, ok
}

func (r *resourceResolver) ToService() (*serviceResolver, bool) {
	c, ok := r.r.(*serviceResolver)
	return c, ok
}

// Implementations of the methods common to all resources
func (r *resourceResolver) Kind() string {
	return r.r.Kind()
}

func (r *resourceResolver) Metadata() metadataResolver {
	return r.r.Metadata()
}

func (r *resourceResolver) Owner() *resourceResolver {
	return r.r.Owner()
}

func (r *resourceResolver) RootOwner() *resourceResolver {
	return r.r.RootOwner()
}
