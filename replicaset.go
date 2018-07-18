package main

import (
	"context"
	//	"fmt"
	"strings"
)

type replicaSet struct {
	Metadata  metadata
	Owner     resource
	RootOwner resource
	Pods      *[]pod
}

type replicaSetResolver struct {
	ctx context.Context
	r   replicaSet
}

func mapToReplicaSet(
	ctx context.Context,
	jsonObj map[string]interface{}) replicaSet {
	placeholder := &ownerRef{ctx, jsonObj, nil}
	owner := placeholder
	rootOwner := placeholder
	meta :=
		mapToMetadata(ctx, getNamespace(jsonObj), mapItem(jsonObj, "metadata"))
	return replicaSet{meta, owner, rootOwner, nil}
}

func getReplicaSetPods(ctx context.Context, r replicaSet) *[]pod {
	rsName := r.Metadata.Name
	rsNamePrefix := rsName + "-"
	rsNamespace := r.Metadata.Namespace

	pset := getAllK8sObjsOfKindInNamespace(
		ctx,
		PodKind,
		rsNamespace,
		func(jobj map[string]interface{}) bool {
			return (strings.HasPrefix(getName(jobj), rsNamePrefix) &&
				hasMatchingOwner(jobj, rsName, ReplicaSetKind))
		})

	results := make([]pod, len(pset))

	for idx, p := range pset {
		pr := p.(*podResolver)
		results[idx] = pr.p
	}

	return &results
}

func (r *replicaSetResolver) Kind() string {
	return ReplicaSetKind
}

func (r *replicaSetResolver) Metadata() *metadataResolver {
	return &metadataResolver{r.ctx, r.r.Metadata}
}

func (r *replicaSetResolver) Owner() *resourceResolver {
	if oref, ok := r.r.Owner.(*ownerRef); ok {
		r.r.Owner = getOwner(oref.ctx, oref.ref)
	}
	return &resourceResolver{r.ctx, r.r.Owner}
}

func (r *replicaSetResolver) RootOwner() *resourceResolver {
	if oref, ok := r.r.Owner.(*ownerRef); ok {
		r.r.Owner = getOwner(oref.ctx, oref.ref)
	}
	return &resourceResolver{r.ctx, r.r.RootOwner}
}

func (r *replicaSetResolver) Pods() []*podResolver {
	if r.r.Pods == nil {
		r.r.Pods = getReplicaSetPods(r.ctx, r.r)
	}

	var res []*podResolver
	for _, p := range *r.r.Pods {
		res = append(res, &podResolver{r.ctx, p})
	}
	if res == nil {
		res = make([]*podResolver, 0)
	}
	return res
}
