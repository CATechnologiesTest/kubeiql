package main

import (
	"context"
	"errors"
	"fmt"
)

func mapItem(obj map[string]interface{}, item string) map[string]interface{} {
	return obj[item].(map[string]interface{})
}

func getMetadata(resourceMap map[string]interface{}) metadata {
	if meta, ok := resourceMap["metadata"]; ok {
		if mmap, ok := meta.(map[string]interface{}); ok {
			return mapToMetadata(mmap)
		}
	}

	panic(errors.New(
		fmt.Sprintf("Invalid Kubernetes resource: %v", resourceMap)))
}

func getKind(resourceMap map[string]interface{}) string {
	if meta, ok := resourceMap["metadata"]; ok {
		if mmap, ok := meta.(map[string]interface{}); ok {
			if kind, ok := mmap["kind"]; ok {
				if kindstr, ok := kind.(string); ok {
					return kindstr
				}
			}
		}
	}

	panic(errors.New(
		fmt.Sprintf("Invalid Kubernetes resource: %v", resourceMap)))
}

func getUid(resourceMap map[string]interface{}) string {
	if meta, ok := resourceMap["metadata"]; ok {
		if mmap, ok := meta.(map[string]interface{}); ok {
			if uid, ok := mmap["uid"]; ok {
				if uidstr, ok := uid.(string); ok {
					return uidstr
				}
			}
		}
	}

	panic(errors.New(
		fmt.Sprintf("Invalid Kubernetes resource: %v", resourceMap)))
}

func getRawOwner(val map[string]interface{}) map[string]interface{} {
	if orefs := val["OwnerReferences"]; orefs != nil {
		oArray := orefs.([]map[string]interface{})
		if len(oArray) > 0 {
			if res := getK8sResource(
				getKind(val),
				oArray[0]["uid"].(string)); res != nil {
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
