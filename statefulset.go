package main

import (
	"context"
	//	"fmt"
	"strings"
)

type statefulSet struct {
	Metadata  metadata
	Owner     resource
	RootOwner resource
	Pods      *[]pod
}

type statefulSetResolver struct {
	ctx context.Context
	s   statefulSet
}

func mapToStatefulSet(
	ctx context.Context,
	jsonObj map[string]interface{}) statefulSet {
	placeholder := &ownerRef{ctx, jsonObj, nil}
	owner := placeholder
	rootOwner := placeholder
	meta :=
		mapToMetadata(ctx, getNamespace(jsonObj), mapItem(jsonObj, "metadata"))
	return statefulSet{meta, owner, rootOwner, nil}
}

func getStatefulSetPods(ctx context.Context, s statefulSet) *[]pod {
	rsName := s.Metadata.Name
	rsNamePrefix := rsName + "-"
	rsNamespace := s.Metadata.Namespace

	pset := getAllK8sObjsOfKindInNamespace(
		ctx,
		PodKind,
		rsNamespace,
		func(jobj map[string]interface{}) bool {
			return (strings.HasPrefix(getName(jobj), rsNamePrefix) &&
				hasMatchingOwner(jobj, rsName, StatefulSetKind))
		})

	results := make([]pod, len(pset))

	for idx, p := range pset {
		pr := p.(*podResolver)
		results[idx] = pr.p
	}

	return &results
}

func (r *statefulSetResolver) Kind() string {
	return StatefulSetKind
}

func (r *statefulSetResolver) Metadata() *metadataResolver {
	return &metadataResolver{r.ctx, r.s.Metadata}
}

func (r *statefulSetResolver) Owner() *resourceResolver {
	if oref, ok := r.s.Owner.(*ownerRef); ok {
		r.s.Owner = getOwner(oref.ctx, oref.ref)
	}
	return &resourceResolver{r.ctx, r.s.Owner}
}

func (r *statefulSetResolver) RootOwner() *resourceResolver {
	if oref, ok := r.s.Owner.(*ownerRef); ok {
		r.s.Owner = getOwner(oref.ctx, oref.ref)
	}
	return &resourceResolver{r.ctx, r.s.RootOwner}
}

func (r *statefulSetResolver) Pods() []*podResolver {
	if r.s.Pods == nil {
		r.s.Pods = getStatefulSetPods(r.ctx, r.s)
	}

	var res []*podResolver
	for _, p := range *r.s.Pods {
		res = append(res, &podResolver{r.ctx, p})
	}
	if res == nil {
		res = make([]*podResolver, 0)
	}
	return res
}
