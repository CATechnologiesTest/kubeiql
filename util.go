package main

import (
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
