package main

import (
	"context"
	"fmt"
)

type resource interface {
	Kind() string
	Metadata() *metadataResolver
	Owner() *resourceResolver
	RootOwner() *resourceResolver
}

type resourceResolver struct {
	ctx context.Context
	r   resource
}

func mapToResource(
	ctx context.Context,
	rMap map[string]interface{}) resource {
	kind := getKind(rMap)

	switch kind {
	case DeploymentKind:
		return &deploymentResolver{ctx, mapToDeployment(ctx, rMap)}
	case ReplicaSetKind:
		return &replicaSetResolver{ctx, mapToReplicaSet(ctx, rMap)}
	case DaemonSetKind:
		return &daemonSetResolver{ctx, mapToDaemonSet(ctx, rMap)}
	case StatefulSetKind:
		return &statefulSetResolver{ctx, mapToStatefulSet(ctx, rMap)}
	case PodKind:
		return &podResolver{ctx, mapToPod(ctx, rMap)}
	}

	fmt.Printf("BAD KIND: %v\n", kind)
	return nil
}

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

func (r *resourceResolver) Kind() string {
	return r.r.Kind()
}

func (r *resourceResolver) Metadata() *metadataResolver {
	return r.r.Metadata()
}

func (r *resourceResolver) Owner() *resourceResolver {
	return r.r.Owner()
}

func (r *resourceResolver) RootOwner() *resourceResolver {
	return r.r.RootOwner()
}
