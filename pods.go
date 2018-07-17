package main

import (
	"context"
)

type pod struct {
	Id        string
	Metadata  metadata
	Owner     resource
	RootOwner resource
}

type podResolver struct {
	ctx context.Context
	p   *pod
}

func mapToPod(jsonObj map[string]interface{}) pod {
	owner, rootOwner := getOwners(jsonObj)
	meta := mapToMetadata(mapItem(jsonObj, "metadata"))
	return pod{(mapItem(jsonObj, "metadata")["uid"]).(string),
		meta,
		owner,
		rootOwner}
}

func (r *podResolver) Id() string {
	return r.p.Id
}

func (r *podResolver) Metadata() *metadataResolver {
	meta := r.p.Metadata
	return &metadataResolver{r.ctx, meta}
}

func (r *podResolver) Owner() *resourceResolver {
	owner := r.p.Owner
	if owner == nil {
		return &resourceResolver{r.ctx, getPodOwner(r.ctx, r.p.Id)}
	}
	return &resourceResolver{r.ctx, owner}
}

func (r *podResolver) RootOwner() *resourceResolver {
	rootOwner := r.p.RootOwner
	if rootOwner == nil {
		return &resourceResolver{r.ctx, getPodRootOwner(r.ctx, r.p.Id)}
	}
	return &resourceResolver{r.ctx, rootOwner}
}

func getPodMetadata(p *pod) metadata {
	meta := mapToMetadata(mapItem(getK8sResource(p.Id), "Metadata"))
	return meta
}

func getPodOwner(ctx context.Context, pid string) resource {
	if podval := getK8sResource(pid); podval != nil {
		if orefs := podval["OwnerReferences"]; orefs != nil {
			orefArray := orefs.([]map[string]interface{})
			if len(orefArray) > 0 {
				if res := getK8sResource(
					orefArray[0]["uid"].(string)); res != nil {
					return mapToResource(ctx, res)
				}
			}
		} else {
			return mapToResource(ctx, podval)
		}
	}

	return nil
}

func getPodRootOwner(ctx context.Context, pid string) resource {
	result := getPodOwner(ctx, pid)

	if result.Id() == pid {
		return result
	}

	return getPodRootOwner(ctx, getPodOwner(ctx, result.Id()).Id())
}
