package main

import (
	//	"context"
	"errors"
	"fmt"
)

func getOwners(resourceMap map[string]interface{}) (resource, resource) {
	return nil, nil
}

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
