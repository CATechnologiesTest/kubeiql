package main

import (
	"context"
	//	"fmt"
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

	if kind == DeploymentKind {
		return &deploymentResolver{ctx, mapToDeployment(ctx, rMap)}
	}

	if kind == ReplicaSetKind {
		return &replicaSetResolver{ctx, mapToReplicaSet(ctx, rMap)}
	}

	if kind == PodKind {
		return &podResolver{ctx, mapToPod(ctx, rMap)}
	}

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
