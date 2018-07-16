package main

import (
	"context"
)

type pod struct {
	ID        string
	Metadata  *metadata
	Owner     *resource
	RootOwner *resource
}

type podResolver struct {
	ctx context.Context
	p   *pod
}

func (r *podResolver) ID() string {
	return r.p.ID
}

func (r *podResolver) Metadata() *metadataResolver {
	meta := r.p.Metadata
	if meta == nil {
		meta = getPodMetadata(r.p)
	}
	return &metadataResolver{r.ctx, meta}
}

func (r *podResolver) Owner() *resourceResolver {
	owner := r.p.Owner
	if owner == nil {
		owner = getPodOwner(r.p)
	}
	return &resourceResolver{r.ctx, owner}
}

func (r *podResolver) RootOwner() *resourceResolver {
	rootOwner := r.p.RootOwner
	if rootOwner == nil {
		rootOwner = getPodRootOwner(r.p)
	}
	return &resourceResolver{r.ctx, rootOwner}
}

func getPodMetadata(p *pod) *metadata {
	return nil
}

func getPodOwner(p *pod) *resource {
	return nil
}

func getPodRootOwner(p *pod) *resource {
	return nil
}
