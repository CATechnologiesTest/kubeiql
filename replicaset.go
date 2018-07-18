package main

import (
	"context"
	//	"fmt"
)

type replicaSet struct {
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
	meta :=
		mapToMetadata(ctx, getNamespace(jsonObj), mapItem(jsonObj, "metadata"))
	return replicaSet{meta, owner, rootOwner}
}

func (r *replicaSetResolver) Kind() string {
	return ReplicaSetKind
}

func (r *replicaSetResolver) Metadata() *metadataResolver {
	return &metadataResolver{r.ctx, r.r.Metadata}
}

func (r *replicaSetResolver) Owner() *resourceResolver {
	return &resourceResolver{r.ctx, r.r.Owner}
}

func (r *replicaSetResolver) RootOwner() *resourceResolver {
	return &resourceResolver{r.ctx, r.r.RootOwner}
}
