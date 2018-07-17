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
	r   *replicaSet
}

func (r *replicaSetResolver) Id() string {
	return r.r.Id
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
