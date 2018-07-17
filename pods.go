package main

import (
	"context"
)

type pod struct {
	Id        string
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
	meta := mapToMetadata(mapItem(jsonObj, "metadata"))
	return pod{(mapItem(jsonObj, "metadata")["uid"]).(string),
		meta,
		owner,
		rootOwner}
}

func (r *podResolver) Id() string {
	return r.p.Id
}

func (r *podResolver) Metadata() *metadataResolver {
	meta := r.p.Metadata
	return &metadataResolver{r.ctx, meta}
}

func (r *podResolver) Owner() *resourceResolver {
	owner := r.p.Owner
	if owner == nil {
		return &resourceResolver{
			r.ctx, getOwner(r.ctx, getK8sResource("Pod", r.p.Id))}
	}
	return &resourceResolver{r.ctx, owner}
}

func (r *podResolver) RootOwner() *resourceResolver {
	rootOwner := r.p.RootOwner
	if rootOwner == nil {
		return &resourceResolver{
			r.ctx, getRootOwner(r.ctx, getK8sResource("Pod", r.p.Id))}
	}
	return &resourceResolver{r.ctx, rootOwner}
}

func getPodMetadata(p *pod) metadata {
	meta := mapToMetadata(mapItem(getK8sResource("Pod", p.Id), "Metadata"))
	return meta
}
