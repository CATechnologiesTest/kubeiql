// Copyright 2018 Yipee.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
)

// Kubernetes metadata
type metadata struct {
	CreationTimestamp *string
	GenerateName      *string
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

// Translate unmarshalled json into a metadata object
func mapToMetadata(
	ctx context.Context, ns string, jsonObj map[string]interface{}) metadata {
	var m metadata
	var orefs []resource
	if ct, ok := jsonObj["creationTimestamp"].(string); ok {
		m.CreationTimestamp = &ct
	} else {
		m.CreationTimestamp = nil
	}
	if gn, ok := jsonObj["generateName"].(string); ok {
		m.CreationTimestamp = &gn
	} else {
		m.CreationTimestamp = nil
	}
	m.Labels = mapToLabels(mapItem(jsonObj, "labels"))
	m.Name = jsonObj["name"].(string)
	m.Namespace = jsonObj["namespace"].(string)
	m.ResourceVersion = jsonObj["resourceVersion"].(string)
	m.SelfLink = jsonObj["selfLink"].(string)
	m.Uid = jsonObj["uid"].(string)

	// Similar to getOwner
	if orArray := jsonObj["ownerReferences"]; orArray != nil {
		for _, oref := range orArray.([]interface{}) {
			ormap := oref.(map[string]interface{})
			orefs = append(
				orefs,
				getK8sResource(
					ctx,
					ormap["kind"].(string),
					ns,
					ormap["name"].(string)))
		}
	}

	m.OwnerReferences = &orefs
	return m
}

// Metadata methods
func (r *metadataResolver) CreationTimestamp() *string {
	return r.m.CreationTimestamp
}

func (r *metadataResolver) GenerateName() *string {
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
