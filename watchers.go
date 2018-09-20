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
		return ApiHost + watchurl
	}
	panic(fmt.Sprintf("no watcher url for kind: '%s'", kind))
}

func readSecret(name string) []byte {
	if ApiSecretPath == "" {
		return nil
	}
	b, err := ioutil.ReadFile(ApiSecretPath + "/" + name)
	if err != nil {
		panic(fmt.Sprintf("watcher error reading secret %s: %s\n",
			name, err.Error()))
	}
	return b
}

func getTlsConfig() *tls.Config {
	cert := readSecret("ca.crt")
	if cert != nil {
		roots := x509.NewCertPool()
		roots.AppendCertsFromPEM(cert)
		tlsConfig := &tls.Config{}
		tlsConfig.RootCAs = roots

		return tlsConfig
	}
	return nil
}

var k8sClient *http.Client
var token = ""

func initClient() {
	tr := &http.Transport{}
	if tlsConfig := getTlsConfig(); tlsConfig != nil {
		tr.TLSClientConfig = tlsConfig
	}
	k8sClient = &http.Client{Transport: tr}
	if tok := readSecret("token"); tok != nil {
		token = string(tok)
	}
}

func makeWatchRequest(kind string) *http.Request {
	req, err := http.NewRequest("GET", buildWatchUrl(kind), nil)
	if err != nil {
		panic(fmt.Sprintf("watcher http.NewRequest error for kind %s: %s\n",
			kind, err.Error()))
	}
	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}
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
			panic(fmt.Sprintf("read error on watcher stream: %s\n", err.Error()))
		}
		var notif Notification
		if err := json.Unmarshal(line, &notif); err != nil {
			panic(fmt.Sprintf("JSON unmarshal error on watcher input: %s\n",
				err.Error()))
		}

		if notif.Type == "ADDED" || notif.Type == "MODIFIED" {
			GetCache().Add(&notif.Object)
		} else if notif.Type == "DELETED" {
			GetCache().Remove(&notif.Object)
		}
	}
}

func initWatchers() {
	initClient()
	for key, _ := range watchUrlByKind {
		go runWatcher(key)
	}
}
