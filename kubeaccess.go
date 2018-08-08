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
	//	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
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
	cache := getCache(ctx)
	key := rawCacheKey(kind, namespace, name)
	cachedVal := (*cache)[key]
	var result JsonObject
	if cachedVal == nil {
		if isTest() {
			return map[string]interface{}{}
		}
		cmd := exec.Command(KubectlPath, "get",
			"-o", "json", "--namespace", namespace, kind, name)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		bytes, err := ioutil.ReadAll(stdout)
		if err != nil {
			log.Fatal(err)
		}
		if err := cmd.Wait(); err != nil {
			log.Fatal(err)
		}
		result = fromJson(bytes).(JsonObject)
		(*cache)[key] = result
	} else {
		result = cachedVal.(JsonObject)
	}
	return result
}

func lookUpResource(ctx context.Context, kind, namespace, name string) resource {
	mapval := lookUpMap(ctx, kind, namespace, name)

	if mapval == nil {
		return nil
	}

	return mapToResource(ctx, mapval)
}

// Get all resource instances of a specific kind
func getAllK8sObjsOfKind(
	ctx context.Context,
	kind string,
	test func(JsonObject) bool) []resource {
	cache := getCache(ctx)
	results := (*cache)[kind]
	if results == nil {
		if isTest() {
			return make([]resource, 0)
		}
		cmd := exec.Command(KubectlPath, "get",
			"-o", "json", "--all-namespaces", kind)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		bytes, err := ioutil.ReadAll(stdout)
		if err != nil {
			log.Fatal(err)
		}
		if err := cmd.Wait(); err != nil {
			log.Fatal(err)
		}
		var resources []resource
		arr := (fromJson(bytes).(JsonObject))["items"].(JsonArray)
		for _, res := range arr {
			val := mapToResource(ctx, res.(JsonObject))
			(*cache)[cacheKey(kind,
				*val.Metadata().Namespace(), *val.Metadata().Name())] = val
			if test(res.(JsonObject)) {
				resources = append(resources, val)
			}
		}
		results = resources
	}
	if results == nil {
		results = make([]resource, 0)
	}
	if (*cache)[kind] == nil {
		(*cache)[kind] = results
	}
	return results.([]resource)
}

// Get all resource instances of a specific kind in a specific namespace
func getAllK8sObjsOfKindInNamespace(
	ctx context.Context,
	kind, ns string,
	test func(JsonObject) bool) []resource {
	cache := getCache(ctx)
	results := (*cache)[kind]
	if results == nil {
		if isTest() {
			return make([]resource, 0)
		}
		cmd := exec.Command(KubectlPath, "get",
			"-o", "json", "--namespace", ns, kind)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		bytes, err := ioutil.ReadAll(stdout)
		if err != nil {
			log.Fatal(err)
		}
		if err := cmd.Wait(); err != nil {
			log.Fatal(err)
		}
		var resources []resource
		arr := (fromJson(bytes).(JsonObject))["items"].(JsonArray)
		for _, res := range arr {
			val := mapToResource(ctx, res.(JsonObject))
			if test(res.(JsonObject)) {
				resources = append(resources, val)
			}
		}
		results = resources
	}
	if results == nil {
		results = make([]resource, 0)
	}
	if (*cache)[kind] == nil {
		(*cache)[kind] = results
	}
	return results.([]resource)
}

func cacheKey(kind, namespace, name string) string {
	return kind + "#" + namespace + "#" + name
}

func rawCacheKey(kind, namespace, name string) string {
	return "raw#" + kind + "#" + namespace + "#" + name
}
