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
	m   metadata
}

func mapToMetadata(jsonObj map[string]interface{}) metadata {
	var m metadata
	if ct, ok := jsonObj["CreationTimestamp"].(string); ok {
		m.CreationTimestamp = &ct
	} else {
		m.CreationTimestamp = nil
	}
	m.GenerateName = jsonObj["GenerateName"].(string)
	m.Labels = mapToLabels(mapItem(jsonObj, "labels"))
	m.Name = jsonObj["Name"].(string)
	m.Namespace = jsonObj["Namespace"].(string)
	m.OwnerReferences = mapToOwnerReferences(mapItem(jsonObj, "OwnerReferences"))
	m.SelfLink = jsonObj["SelfLink"].(string)
	m.Uid = jsonObj["Uid"].(string)

	return m
}

func mapToOwnerReferences(orMap map[string]interface{}) *[]resource {
	return nil
}

func (r *metadataResolver) CreationTimestamp() *string {
	return r.m.CreationTimestamp
}

func (r *metadataResolver) GenerateName() string {
	return r.m.GenerateName
}

func (r *metadataResolver) Labels() []*labelResolver {
	var labelResolvers []*labelResolver
	for _, label := range *r.m.Labels {
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
	for _, owner := range *r.m.OwnerReferences {
		own := owner
		ownerResolvers = append(ownerResolvers, &resourceResolver{r.ctx, own})
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
