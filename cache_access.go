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

var client CacheClient

type CacheOp func() interface{}

type CacheRequest struct {
	operation CacheOp
	replyChan chan interface{}
}

type clientif interface {
	Lookup(key string) interface{}
	Remove(obj *JsonObject)
	Add(obj *JsonObject)
}

type CacheClient struct {
	serverMbox chan<- *CacheRequest
}

func initCacheClient(serverMbox chan *CacheRequest) {
	client = CacheClient{serverMbox}
}

func GetCache() *CacheClient {
	return &client
}

func (client *CacheClient) Lookup(key string) interface{} {
	replyChan := make(chan interface{})
	req := &CacheRequest{
		func() interface{} {
			return cacheLookup(key)
		}, replyChan,
	}
	client.serverMbox <- req
	retval := <-replyChan
	return retval
}

func (client *CacheClient) Remove(obj *JsonObject) {
	replyChan := make(chan interface{})
	req := &CacheRequest{
		func() interface{} {
			removeFromCache(obj)
			return true
		}, replyChan,
	}
	client.serverMbox <- req
	if retval := <-replyChan; retval != true {
		panic("bad return from cache Remove")
	}
}

func (client *CacheClient) Add(obj *JsonObject) {
	replyChan := make(chan interface{})
	req := &CacheRequest{
		func() interface{} {
			addToCache(obj)
			return true
		}, replyChan,
	}
	client.serverMbox <- req
	if retval := <-replyChan; retval != true {
		panic("bad return from cache Add")
	}
}

// compilation error if we don't implement the i/f properly
var _ clientif = (*CacheClient)(nil)

func cacheKey(kind, namespace, name string) string {
	return kind + "#" + namespace + "#" + name
}

func nsCacheKey(kind, namespace string) string {
	return kind + "#" + namespace
}
