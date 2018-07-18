package main

import (
	"context"
	"errors"
	"fmt"
)

func mapItem(obj map[string]interface{}, item string) map[string]interface{} {
	return obj[item].(map[string]interface{})
}

func getKind(resourceMap map[string]interface{}) string {
	kind := resourceMap["kind"]
	if kindstr, ok := kind.(string); ok {
		return kindstr
	}

	panic(errors.New(
		fmt.Sprintf("Invalid Kubernetes resource: %v", resourceMap)))
}

func getNamespace(resourceMap map[string]interface{}) string {
	namespace := getMetadataField(resourceMap, "namespace")
	if nsstr, ok := namespace.(string); ok {
		return nsstr
	}

	panic(errors.New(
		fmt.Sprintf("Invalid Kubernetes resource: %v", resourceMap)))
}

func getName(resourceMap map[string]interface{}) string {
	name := getMetadataField(resourceMap, "name")
	if nsstr, ok := name.(string); ok {
		return nsstr
	}

	panic(errors.New(
		fmt.Sprintf("Invalid Kubernetes resource: %v", resourceMap)))
}

func getUid(resourceMap map[string]interface{}) string {
	uid := getMetadataField(resourceMap, "uid")
	if uidstr, ok := uid.(string); ok {
		return uidstr
	}

	panic(errors.New(
		fmt.Sprintf("Invalid Kubernetes resource: %v", resourceMap)))
}

func getMetadataField(
	resourceMap map[string]interface{},
	field string) interface{} {
	if meta, ok := resourceMap["metadata"]; ok {
		if mmap, ok := meta.(map[string]interface{}); ok {
			if val, ok := mmap[field]; ok {
				return val
			}
		}
	}

	return nil
}

func getRawOwner(val map[string]interface{}) map[string]interface{} {
	if orefs := getMetadataField(val, "ownerReferences"); orefs != nil {
		oArray := orefs.([]interface{})
		if len(oArray) > 0 {
			owner := oArray[0].(map[string]interface{})
			if res := getK8sResource(owner["kind"].(string),
				getNamespace(val),
				owner["name"].(string)); res != nil {
				return res
			}
		}
	}

	return val
}

func getOwner(ctx context.Context, val map[string]interface{}) resource {
	return mapToResource(ctx, getRawOwner(val))
}

func getRootOwner(ctx context.Context, val map[string]interface{}) resource {
	result := getRawOwner(val)

	if getUid(result) == getUid(val) {
		return mapToResource(ctx, result)
	}

	return getRootOwner(ctx, getRawOwner(result))
}

func hasMatchingOwner(jsonObj map[string]interface{}, name, kind string) bool {
	if orefs := getMetadataField(jsonObj, "ownerReferences"); orefs != nil {
		oArray := orefs.([]interface{})
		for _, oref := range oArray {
			owner := oref.(map[string]interface{})
			if owner["name"] == name && owner["kind"] == kind {
				return true
			}
		}
	}

	return false
}
