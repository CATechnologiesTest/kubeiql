package main

import (
	"context"
)

type metadata struct {
	CreationTimestamp *string
	GenerateName      string
	Labels            *[]label
	Name              string
	Namespace         string
	OwnerReferences   *[]resource
	ResourceVersion   string
	SelfLink          string
	Uid               string
}

type metadataResolver struct {
	ctx context.Context
	m   *metadata
}

func (r *metadataResolver) CreationTimestamp() *string {
	return r.m.CreationTimestamp
}

func (r *metadataResolver) GenerateName() string {
	return r.m.GenerateName
}

func (r *metadataResolver) Labels() []*labelResolver {
	var labelResolvers []*labelResolver
	labels := r.m.Labels
	if labels == nil {
		labels = getMetadataLabels(r.m)
	}
	for _, label := range *labels {
		lab := label
		labelResolvers = append(labelResolvers, &labelResolver{r.ctx, &lab})
	}
	return labelResolvers
}

func (r *metadataResolver) Name() string {
	return r.m.Name
}

func (r *metadataResolver) Namespace() string {
	return r.m.Namespace
}

func (r *metadataResolver) OwnerReferences() []*resourceResolver {
	var ownerResolvers []*resourceResolver
	owners := r.m.OwnerReferences
	if owners == nil {
		owners = getMetadataOwnerReferences(r.m)
	}
	for _, owner := range *owners {
		own := owner
		ownerResolvers = append(ownerResolvers, &resourceResolver{r.ctx, &own})
	}
	return ownerResolvers
}

func (r *metadataResolver) ResourceVersion() string {
	return r.m.ResourceVersion
}

func (r *metadataResolver) SelfLink() string {
	return r.m.SelfLink
}

func (r *metadataResolver) Uid() string {
	return r.m.Uid
}

func getMetadataLabels(m *metadata) *[]label {
	return nil
}

func getMetadataOwnerReferences(m *metadata) *[]resource {
	return nil
}
