package main

import (
	"context"
	//	"fmt"
)

type pod struct {
	Metadata  metadata
	Owner     resource
	RootOwner resource
}

type podResolver struct {
	ctx context.Context
	p   pod
}

func mapToPod(ctx context.Context, jsonObj map[string]interface{}) pod {
	owner := getOwner(ctx, jsonObj)
	rootOwner := getRootOwner(ctx, jsonObj)
	meta :=
		mapToMetadata(ctx, getNamespace(jsonObj), mapItem(jsonObj, "metadata"))
	return pod{meta, owner, rootOwner}
}

func (r *podResolver) Kind() string {
	return PodKind
}

func (r *podResolver) Metadata() *metadataResolver {
	meta := r.p.Metadata
	return &metadataResolver{r.ctx, meta}
}

func (r *podResolver) Owner() *resourceResolver {
	return &resourceResolver{r.ctx, r.p.Owner}
}

func (r *podResolver) RootOwner() *resourceResolver {
	return &resourceResolver{r.ctx, r.p.RootOwner}
}
