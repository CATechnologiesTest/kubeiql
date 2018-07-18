package main

import (
	"context"
	//	"fmt"
	"strings"
)

type daemonSet struct {
	Metadata  metadata
	Owner     resource
	RootOwner resource
	Pods      *[]pod
}

type daemonSetResolver struct {
	ctx context.Context
	d   daemonSet
}

func mapToDaemonSet(
	ctx context.Context,
	jsonObj map[string]interface{}) daemonSet {
	placeholder := &ownerRef{ctx, jsonObj, nil}
	owner := placeholder
	rootOwner := placeholder
	meta :=
		mapToMetadata(ctx, getNamespace(jsonObj), mapItem(jsonObj, "metadata"))
	return daemonSet{meta, owner, rootOwner, nil}
}

func getDaemonSetPods(ctx context.Context, d daemonSet) *[]pod {
	dsName := d.Metadata.Name
	dsNamePrefix := dsName + "-"
	dsNamespace := d.Metadata.Namespace

	pset := getAllK8sObjsOfKindInNamespace(
		ctx,
		PodKind,
		dsNamespace,
		func(jobj map[string]interface{}) bool {
			return (strings.HasPrefix(getName(jobj), dsNamePrefix) &&
				hasMatchingOwner(jobj, dsName, DaemonSetKind))
		})

	results := make([]pod, len(pset))

	for idx, p := range pset {
		pr := p.(*podResolver)
		results[idx] = pr.p
	}

	return &results
}

func (r *daemonSetResolver) Kind() string {
	return DaemonSetKind
}

func (r *daemonSetResolver) Metadata() *metadataResolver {
	return &metadataResolver{r.ctx, r.d.Metadata}
}

func (r *daemonSetResolver) Owner() *resourceResolver {
	if oref, ok := r.d.Owner.(*ownerRef); ok {
		r.d.Owner = getOwner(oref.ctx, oref.ref)
	}
	return &resourceResolver{r.ctx, r.d.Owner}
}

func (r *daemonSetResolver) RootOwner() *resourceResolver {
	if oref, ok := r.d.Owner.(*ownerRef); ok {
		r.d.Owner = getOwner(oref.ctx, oref.ref)
	}
	return &resourceResolver{r.ctx, r.d.RootOwner}
}

func (r *daemonSetResolver) Pods() []*podResolver {
	if r.d.Pods == nil {
		r.d.Pods = getDaemonSetPods(r.ctx, r.d)
	}

	var res []*podResolver
	for _, p := range *r.d.Pods {
		res = append(res, &podResolver{r.ctx, p})
	}
	if res == nil {
		res = make([]*podResolver, 0)
	}
	return res
}
