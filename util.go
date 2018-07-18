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
	"errors"
	"fmt"
)

// Utility methods for getting data out of nested maps

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
