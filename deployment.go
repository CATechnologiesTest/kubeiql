package main

import (
	"context"
)

type deployment struct {
	Id        string
	Metadata  metadata
	Owner     resource
	RootOwner resource
}

type deploymentResolver struct {
	ctx context.Context
	d   deployment
}

func mapToDeployment(
	ctx context.Context,
	jsonObj map[string]interface{}) deployment {
	owner := getOwner(ctx, jsonObj)
	rootOwner := getRootOwner(ctx, jsonObj)
	meta := mapToMetadata(mapItem(jsonObj, "metadata"))
	return deployment{(mapItem(jsonObj, "metadata")["uid"]).(string),
		meta,
		owner,
		rootOwner}
}

func (r *deploymentResolver) Id() string {
	return r.d.Id
}

func (r *deploymentResolver) Metadata() *metadataResolver {
	return &metadataResolver{r.ctx, r.d.Metadata}
}

func (r *deploymentResolver) Owner() *resourceResolver {
	owner := r.d.Owner
	if owner == nil {
		return &resourceResolver{
			r.ctx, getOwner(r.ctx, getK8sResource("Deployment", r.d.Id))}
	}

	return &resourceResolver{r.ctx, owner}
}

func (r *deploymentResolver) RootOwner() *resourceResolver {
	rootOwner := r.d.RootOwner
	if rootOwner == nil {
		return &resourceResolver{
			r.ctx, getRootOwner(r.ctx, getK8sResource("Deployment", r.d.Id))}
	}
	return &resourceResolver{r.ctx, rootOwner}
}
