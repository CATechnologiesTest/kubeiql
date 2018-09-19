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
	"fmt"
)

var cache map[string]interface{}

func initCache() {
	cache = make(map[string]interface{})
}

func findInList(clist []*JsonObject, target *JsonObject) int {
	targids := getObjIds(target)
	for idx, obj := range clist {
		ids := getObjIds(obj)
		if targids.name == ids.name && targids.namespace == ids.namespace {
			return idx
		}
	}
	return -1
}

func deleteFromCacheList(key string, obj *JsonObject) {
	if val, ok := cache[key]; ok {
		clist := val.([]*JsonObject)
		idx := findInList(clist, obj)
		if idx > -1 {
			clist = append(clist[:idx], clist[idx+1:]...)
		}
	}
}

func addToCacheList(key string, obj *JsonObject) {
	if val, ok := cache[key]; ok {
		clist := val.([]*JsonObject)
		idx := findInList(clist, obj)
		if idx == -1 {
			cache[key] = append(clist, obj)
		} else {
			clist[idx] = obj
		}
	} else {
		cache[key] = []*JsonObject{obj}
	}
}

func formattedName(obj *JsonObject) string {
	ids := getObjIds(obj)
	return cacheKey(ids.kind, ids.namespace, ids.name)
}

func addToCache(obj *JsonObject) {
	ids := getObjIds(obj)
	cacheKey := cacheKey(ids.kind, ids.namespace, ids.name)
	nsCacheKey := nsCacheKey(ids.kind, ids.namespace)
	cache[cacheKey] = obj
	addToCacheList(nsCacheKey, obj)
	addToCacheList(ids.kind, obj)
}

func removeFromCache(obj *JsonObject) {
	ids := getObjIds(obj)
	cacheKey := cacheKey(ids.kind, ids.namespace, ids.name)
	nsCacheKey := nsCacheKey(ids.kind, ids.namespace)
	delete(cache, cacheKey)
	deleteFromCacheList(nsCacheKey, obj)
	deleteFromCacheList(ids.kind, obj)
}

func cacheLookup(key string) interface{} {
	if val, ok := cache[key]; ok {
		if ref, ok := val.(*JsonObject); ok {
			return *ref
		} else if list, ok := val.([]*JsonObject); ok {
			// XXX: clone slices??
			retlist := make([]JsonObject, len(list))
			for idx, item := range list {
				retlist[idx] = *item
			}
			return retlist
		} else {
			panic(fmt.Sprintf("invalid type in cache: %T", val))
		}
	} else {
		return nil
	}
}

func cacheKey(kind, namespace, name string) string {
	return kind + "#" + namespace + "#" + name
}

func nsCacheKey(kind, namespace string) string {
	return kind + "#" + namespace
}
