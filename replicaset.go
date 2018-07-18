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
	owner := getOwner(ctx, jsonObj)
	rootOwner := getRootOwner(ctx, jsonObj)
	meta :=
		mapToMetadata(ctx, getNamespace(jsonObj), mapItem(jsonObj, "metadata"))
	return replicaSet{meta, owner, rootOwner, nil}
}

func getPods(ctx context.Context, r replicaSet) *[]pod {
	rsName := r.Metadata.Name
	rsNamePrefix := rsName + "-"
	rsNamespace := r.Metadata.Namespace

	psets := getAllK8sObjsOfKindInNamespace(
		ctx,
		"Pod",
		rsNamespace,
		func(jobj map[string]interface{}) bool {
			return (strings.HasPrefix(getName(jobj), rsNamePrefix) &&
				hasMatchingOwner(jobj, rsName, ReplicaSetKind))
		})

	results := make([]pod, len(psets))

	for idx, p := range psets {
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
	return &resourceResolver{r.ctx, r.r.Owner}
}

func (r *replicaSetResolver) RootOwner() *resourceResolver {
	return &resourceResolver{r.ctx, r.r.RootOwner}
}

func (r *replicaSetResolver) Pods() []*podResolver {
	if r.r.Pods == nil {
		r.r.Pods = getPods(r.ctx, r.r)
	}

	var res []*podResolver
	for _, p := range *r.r.Pods {
		res = append(res, &podResolver{r.ctx, p})
	}
	return res
}
