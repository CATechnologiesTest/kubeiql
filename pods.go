package main

import (
	"context"
)

type pod struct {
	Id        string
	Metadata  *metadata
	Owner     *resource
	RootOwner *resource
}

type podResolver struct {
	ctx context.Context
	p   *pod
}

func mapToPod(jsonObj map[string]interface{}) pod {
	owner, rootOwner := getOwners(jsonObj)
	meta := mapToMetadata(mapItem(jsonObj, "metadata"))
	return pod{(mapItem(jsonObj, "metadata")["uid"]).(string),
		&meta,
		owner,
		rootOwner}
}

func (r *podResolver) Id() string {
	return r.p.Id
}

func (r *podResolver) Metadata() *metadataResolver {
	meta := r.p.Metadata
	if meta == nil {
		meta = getPodMetadata(r.p)
	}
	return &metadataResolver{r.ctx, meta}
}

func (r *podResolver) Owner() *resourceResolver {
	owner := r.p.Owner
	if owner == nil {
		owner = getPodOwner(r.p.Id)
	}
	return &resourceResolver{r.ctx, owner}
}

func (r *podResolver) RootOwner() *resourceResolver {
	rootOwner := r.p.RootOwner
	if rootOwner == nil {
		rootOwner = getPodRootOwner(r.p.Id)
	}
	return &resourceResolver{r.ctx, rootOwner}
}

func getPodMetadata(p *pod) *metadata {
	meta := mapToMetadata(mapItem(getK8sResource(p.Id), "Metadata"))
	return &meta
}

func getPodOwner(pid string) *resource {
	if podval := getK8sResource(pid); podval != nil {
		if orefs := podval["OwnerReferences"]; orefs != nil {
			orefArray := orefs.([]map[string]interface{})
			if len(orefArray) > 0 {
				if res := getK8sResource(
					orefArray[0]["uid"].(string)); res != nil {
					return mapToResource(res)
				}
			}
		} else {
			return mapToResource(podval)
		}
	}

	return nil
}

func getPodRootOwner(pid string) *resource {
	result := getPodOwner(pid)

	if (*result).Id() == pid {
		return result
	}

	return getPodRootOwner((*getPodOwner((*result).Id())).Id())
}
