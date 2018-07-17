package main

import (
	"context"
)

type replicaSet struct {
	Id        string
	Metadata  metadata
	Owner     resource
	RootOwner resource
}

type replicaSetResolver struct {
	ctx context.Context
	r   replicaSet
}

func mapToReplicaSet(
	ctx context.Context,
	jsonObj map[string]interface{}) replicaSet {
	owner := getOwner(ctx, jsonObj)
	rootOwner := getRootOwner(ctx, jsonObj)
	meta := mapToMetadata(mapItem(jsonObj, "metadata"))
	return replicaSet{(mapItem(jsonObj, "metadata")["uid"]).(string),
		meta,
		owner,
		rootOwner}
}

func (r *replicaSetResolver) Id() string {
	return r.r.Id
}

func (r *replicaSetResolver) Metadata() *metadataResolver {
	return &metadataResolver{r.ctx, r.r.Metadata}
}

func (r *replicaSetResolver) Owner() *resourceResolver {
	owner := r.r.Owner
	if owner == nil {
		return &resourceResolver{
			r.ctx, getOwner(r.ctx, getK8sResource("ReplicaSet", r.r.Id))}
	}

	return &resourceResolver{r.ctx, owner}
}

func (r *replicaSetResolver) RootOwner() *resourceResolver {
	rootOwner := r.r.RootOwner
	if rootOwner == nil {
		return &resourceResolver{
			r.ctx, getRootOwner(r.ctx, getK8sResource("ReplicaSet", r.r.Id))}
	}
	return &resourceResolver{r.ctx, rootOwner}
}
