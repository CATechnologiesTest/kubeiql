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
	placeholder := &ownerRef{ctx, jsonObj, nil}
	owner := placeholder
	rootOwner := placeholder
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
	if oref, ok := r.p.Owner.(*ownerRef); ok {
		r.p.Owner = getOwner(oref.ctx, oref.ref)
	}
	return &resourceResolver{r.ctx, r.p.Owner}
}

func (r *podResolver) RootOwner() *resourceResolver {
	if oref, ok := r.p.RootOwner.(*ownerRef); ok {
		r.p.RootOwner = getRootOwner(oref.ctx, oref.ref)
	}
	return &resourceResolver{r.ctx, r.p.RootOwner}
}
