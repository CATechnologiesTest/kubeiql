// Copyright (c) 2018 CA. All rights reserved.
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
	"encoding/json"
	"fmt"
)

// Functions for retrieving Kubernetes information from a cluster

// Get a single resource instance from a namespace
func getK8sResource(ctx context.Context, kind, namespace, name string) resource {
	return lookUpResource(ctx, kind, namespace, name)
}

func getRawK8sResource(
	ctx context.Context, kind, namespace, name string) JsonObject {
	return lookUpMap(ctx, kind, namespace, name)
}

func fromJson(val []byte) interface{} {
	var result interface{}

	if err := json.Unmarshal(val, &result); err != nil {
		panic(err)
	}

	return result
}

var testContext *context.Context = nil

func setTestContext(ctx *context.Context) {
	testContext = ctx
}

func getTestContext() *context.Context {
	return testContext
}

func getCache(inctx context.Context) *JsonObject {
	ctx := &inctx
	if isTest() {
		ctx = getTestContext()
	}
	return (*ctx).Value("queryCache").(*JsonObject)
}

func isTest() bool {
	return testContext != nil
}

func lookUpMap(
	ctx context.Context,
	kind, namespace, name string) JsonObject {
	key := cacheKey(kind, namespace, name)
	var cachedVal interface{}
	if !isWatchedKind(kind) {
		panic(fmt.Sprintf("Add watcher for kind '%s'", kind))
	}
	cachedVal = GetCache().Lookup(key)
	return cachedVal.(JsonObject)
}

func lookUpResource(ctx context.Context, kind, namespace, name string) resource {
	mapval := lookUpMap(ctx, kind, namespace, name)

	if mapval == nil {
		return nil
	}

	return mapToResource(ctx, mapval)
}

func getCachedResourceList(
	ctx context.Context,
	cacheKey string,
	test func(JsonObject) bool) []resource {

	var cachedJsonObjs []JsonObject
	var results []resource

	if objs := GetCache().Lookup(cacheKey); objs != nil {
		cachedJsonObjs = objs.([]JsonObject)
	}
	for _, res := range cachedJsonObjs {
		val := mapToResource(ctx, res)
		if test(res) {
			results = append(results, val)
		}
	}
	if results == nil {
		results = make([]resource, 0)
	}
	return results
}

// Get all resource instances of a specific kind
func getAllK8sObjsOfKind(
	ctx context.Context,
	kind string,
	test func(JsonObject) bool) []resource {

	if !isWatchedKind(kind) {
		panic(fmt.Sprintf("Add watcher for kind '%s'", kind))
	}
	return getCachedResourceList(ctx, kind, test)
}

// Get all resource instances of a specific kind in a specific namespace
func getAllK8sObjsOfKindInNamespace(
	ctx context.Context,
	kind, ns string,
	test func(JsonObject) bool) []resource {
	key := nsCacheKey(kind, ns)
	if !isWatchedKind(kind) {
		panic(fmt.Sprintf("Add watcher for kind '%s'", kind))
	}
	return getCachedResourceList(ctx, key, test)
}
