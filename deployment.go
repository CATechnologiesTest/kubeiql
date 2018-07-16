package main

import (
	"context"
)

type deployment struct {
	Id        string
	Metadata  *metadata
	Owner     *resource
	RootOwner *resource
}

type deploymentResolver struct {
	ctx context.Context
	d   *deployment
}

func (r *deploymentResolver) Id() string {
	return r.d.Id
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
