package main

import (
	"context"
)

type resource interface {
	Id() string
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

	if kind == "Deployment" {
		return &deploymentResolver{ctx, mapToDeployment(ctx, rMap)}
	}

	if kind == "ReplicaSet" {
		return &replicaSetResolver{ctx, mapToReplicaSet(ctx, rMap)}
	}

	if kind == "Pod" {
		return &podResolver{ctx, mapToPod(ctx, rMap)}
	}

	return nil
}

func (r *resourceResolver) Id() string {
	return r.r.Id()
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
