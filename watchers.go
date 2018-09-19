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
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Notification struct {
	Type   string
	Object JsonObject
}

type ObjId struct {
	kind      string
	namespace string
	name      string
}

// const secretDir = "/var/run/secrets/kubernetes.io/serviceaccount"
const secretDir = "/Users/engji01"

// const apiHost = "https://kubernetes.default.svc"
const apiHost = "https://192.168.99.100:8443"

var watchUrlByKind = map[string]string{
	"Pod":         "/api/v1/watch/pods?watch=true",
	"Deployment":  "/apis/apps/v1/watch/deployments?watch=true",
	"ReplicaSet":  "/apis/apps/v1/watch/replicasets?watch=true",
	"StatefulSet": "/apis/apps/v1/watch/statefulsets?watch=true",
	"DaemonSet":   "/apis/apps/v1/watch/daemonsets?watch=true",
	"Service":     "/api/v1/watch/services?watch=true",
}

func isWatchedKind(kind string) bool {
	_, ok := watchUrlByKind[kind]
	return ok
}

func buildWatchUrl(kind string) string {
	if watchurl, ok := watchUrlByKind[kind]; ok {
		return apiHost + watchurl
	}
	panic(fmt.Sprintf("no url for kind: '%s'", kind))
}

func readSecret(name string) []byte {
	b, err := ioutil.ReadFile(secretDir + "/" + name)
	if err != nil {
		panic(fmt.Sprintf("error reading secret %s: %s\n",
			name, err.Error()))
	}
	return b
}

func getTlsConfig() *tls.Config {
	roots := x509.NewCertPool()
	roots.AppendCertsFromPEM(readSecret("ca.crt"))
	tlsConfig := &tls.Config{}
	tlsConfig.RootCAs = roots

	return tlsConfig
}

func getObjIds(obj *JsonObject) ObjId {
	var kind, namespace, name string
	if k, ok := (*obj)["kind"]; ok {
		kind = k.(string)
	}
	if meta, ok := (*obj)["metadata"]; ok {
		if md, ok := meta.(map[string]interface{}); ok {
			if nm, ok := md["name"]; ok {
				name = nm.(string)
			}
			if ns, ok := md["namespace"]; ok {
				namespace = ns.(string)
			}
		}
	}
	return ObjId{kind: kind, namespace: namespace, name: name}
}

var k8sClient *http.Client
var token string

func initClient() {
	tr := &http.Transport{
		TLSClientConfig: getTlsConfig(),
	}
	k8sClient = &http.Client{Transport: tr}
	token = string(readSecret("token"))
}

func makeWatchRequest(kind string) *http.Request {
	req, err := http.NewRequest("GET", buildWatchUrl(kind), nil)
	if err != nil {
		panic(fmt.Sprintf("http.NewRequest error for kind %s: %s\n",
			kind, err.Error()))
	}
	req.Header.Add("Authorization", "Bearer "+token)
	return req
}

func runWatcher(kind string) {
	req := makeWatchRequest(kind)
	resp, err := k8sClient.Do(req)
	if err != nil {
		panic(fmt.Sprintf("http GET error for watch of kind %s: %s",
			kind, err.Error()))
	}
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			panic(err)
		}
		var notif Notification
		if err := json.Unmarshal(line, &notif); err != nil {
			panic(err)
		}

		if notif.Type == "ADDED" || notif.Type == "MODIFIED" {
			addToCache(&notif.Object)
		} else if notif.Type == "DELETED" {
			removeFromCache(&notif.Object)
		}
	}
	fmt.Printf("watcher terminates...\n")
}

func initWatchers() {
	initClient()
	for key, _ := range watchUrlByKind {
		go runWatcher(key)
	}
}
