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

// N.B.: this is the cache implementation whose functions are
// not intended to be called directly.
// The intended cache interface
// is in cache_public.go: namely "Lookup", "Add", "Remove",
// and the key-building functions.
// All access to the cache is intended to be via the server mailbox.
// That is our serialization mechanism (akin to an erlang gen_server).
//
// (Deliberately resisting the "internal" package goo...)
var cache map[string]interface{}

func runServer(mbox <-chan *CacheRequest) {
	for {
		req, ok := <-mbox
		if ok {
			req.replyChan <- req.operation()
		} else {
			break
		}
	}
}

func initCache() {
	cache = make(map[string]interface{})
	serverMbox := make(chan *CacheRequest)
	go runServer(serverMbox)
	initCacheClient(serverMbox)
}

func findInList(clist []*JsonObject, target *JsonObject) int {
	tname := getName(*target)
	tns := getNamespace(*target)
	tkind := getKind(*target)
	for idx, obj := range clist {
		name := getName(*obj)
		ns := getNamespace(*obj)
		kind := getKind(*obj)
		if tname == name && tns == ns && tkind == kind {
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
			cache[key] = append(clist[:idx], clist[idx+1:]...)
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
	return cacheKey(getKind(*obj), getNamespace(*obj), getName(*obj))
}

func addToCache(obj *JsonObject) {
	kind := getKind(*obj)
	ns := getNamespace(*obj)
	name := getName(*obj)
	cacheKey := cacheKey(kind, ns, name)
	nsCacheKey := nsCacheKey(kind, ns)
	cache[cacheKey] = obj
	addToCacheList(nsCacheKey, obj)
	addToCacheList(kind, obj)
}

func removeFromCache(obj *JsonObject) {
	kind := getKind(*obj)
	ns := getNamespace(*obj)
	name := getName(*obj)
	cacheKey := cacheKey(kind, ns, name)
	nsCacheKey := nsCacheKey(kind, ns)
	delete(cache, cacheKey)
	deleteFromCacheList(nsCacheKey, obj)
	deleteFromCacheList(kind, obj)
}

func cacheLookup(key string) interface{} {
	if val, ok := cache[key]; ok {
		if ref, ok := val.(*JsonObject); ok {
			return *ref
		} else if list, ok := val.([]*JsonObject); ok {
			// clone returned slice so that its "shape" can't be changed
			// while caller is holding it...
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
