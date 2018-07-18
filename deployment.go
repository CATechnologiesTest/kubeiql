package main

import (
	"context"
	//	"fmt"
	"strings"
)

type deployment struct {
	Metadata    metadata
	Owner       resource
	RootOwner   resource
	ReplicaSets *[]replicaSet
}

type deploymentResolver struct {
	ctx context.Context
	d   deployment
}

func mapToDeployment(
	ctx context.Context,
	jsonObj map[string]interface{}) deployment {
	meta :=
		mapToMetadata(ctx, getNamespace(jsonObj), mapItem(jsonObj, "metadata"))
	return deployment{meta, nil, nil, nil}
}

func getReplicaSets(ctx context.Context, d deployment) *[]replicaSet {
	depName := d.Metadata.Name
	depNamePrefix := depName + "-"
	depNamespace := d.Metadata.Namespace

	rsets := getAllK8sObjsOfKindInNamespace(
		ctx,
		"ReplicaSet",
		depNamespace,
		func(jobj map[string]interface{}) bool {
			return (strings.HasPrefix(getName(jobj), depNamePrefix) &&
				hasMatchingOwner(jobj, depName, DeploymentKind))
		})

	results := make([]replicaSet, len(rsets))

	for idx, rs := range rsets {
		rsr := rs.(*replicaSetResolver)
		results[idx] = rsr.r
	}

	return &results
}

func (r *deploymentResolver) Kind() string {
	return DeploymentKind
}

func (r *deploymentResolver) Metadata() *metadataResolver {
	return &metadataResolver{r.ctx, r.d.Metadata}
}

func (r *deploymentResolver) Owner() *resourceResolver {
	return &resourceResolver{r.ctx, r.d.Owner}
}

func (r *deploymentResolver) RootOwner() *resourceResolver {
	return &resourceResolver{r.ctx, r.d.RootOwner}
}

func (r *deploymentResolver) ReplicaSets() []*replicaSetResolver {
	if r.d.ReplicaSets == nil {
		r.d.ReplicaSets = getReplicaSets(r.ctx, r.d)
	}

	var res []*replicaSetResolver
	for _, rs := range *r.d.ReplicaSets {
		res = append(res, &replicaSetResolver{r.ctx, rs})
	}
	return res
}
