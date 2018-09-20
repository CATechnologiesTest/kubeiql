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
	//	"fmt"
	graphql "github.com/neelance/graphql-go"
	"io/ioutil"
	"log"
	"testing"
)

var testschema *graphql.Schema = graphql.MustParseSchema(Schema, &Resolver{})

func simpletest(t *testing.T, query string, result string) {
	var stest Test
	stest.Schema = testschema
	stest.Query = query
	stest.ExpectedResult = result
	RunTest(t, &stest)
}

func init() {
	cache := make(map[string]interface{})
	ctx := context.WithValue(context.Background(), "queryCache", &cache)
	setTestContext(&ctx)
	for _, fname := range []string{
		"deployment.json",
		"replicaset.json",
		"daemonset.json",
		"statefulset.json",
		"service.json",
		"pod1.json",
		"pod2.json",
		"pod3.json"} {
		addToTestCache(&cache, "testdata/"+fname)
	}
}

func addToTestCache(cacheref *map[string]interface{}, fname string) {
	bytes, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}
	data := fromJson(bytes).(JsonObject)
	cache := GetCache()
	cache.Add(&data)
}

func TestPods(t *testing.T) {
	simpletest(
		t,
		`{
           allDeployments() {
             metadata {
               creationTimestamp
               generation
               labels { name value }
             }
             spec {
               minReadySeconds
               paused
               progressDeadlineSeconds
               replicas
               revisionHistoryLimit
               selector {
                 matchLabels { name value }
                 matchExpressions {
                   key
                   operator
                   values
                 }
               }
               strategy {
                 type
                 rollingUpdate {
                   maxSurgeInt
                   maxSurgeString
                   maxUnavailableInt
                   maxUnavailableString
                 }
               }
               template {
                 metadata {
                   creationTimestamp
                   labels { name value }
                 },
                 spec {
                   dnsPolicy
                   restartPolicy
                   schedulerName
                   terminationGracePeriodSeconds
                   volumes {
                     name
                     persistentVolumeClaim { claimName readOnly }
                   }
                 }
               }
             }
             replicaSets {
               metadata {
                 name
               }
               pods {
                 metadata {
                   name
                   namespace
                   labels {
                     name
                     value
                   }
                 }
                 spec {
                   dnsPolicy
                   nodeName
                   restartPolicy
                   schedulerName
                   serviceAccountName
                   terminationGracePeriodSeconds
                   tolerations {
                     effect
                     key
                     operator
                     tolerationSeconds
                   }
                   volumes {
                     name
                     persistentVolumeClaim { claimName readOnly }
                     secret { defaultMode secretName }
                   }
                 }
               }
             }
          }}`,
		`{
           "allDeployments": [
            {
               "metadata": {
                 "creationTimestamp": "2018-07-02T14:53:53Z",
                 "generation": 1,
                 "labels": [
                   {"name": "app", "value": "clunky-sabertooth-joomla"},
                   {"name": "chart", "value": "joomla-2.0.2"},
                   {"name": "heritage", "value": "Tiller"},
                   {"name": "release", "value": "clunky-sabertooth"}
                  ]
               },
               "spec": {
                 "minReadySeconds": 0,
                 "paused": false,
                 "progressDeadlineSeconds": 600,
                 "replicas": 1,
                 "revisionHistoryLimit": 10,
                 "selector": {
                   "matchExpressions": [],
                   "matchLabels": [
                     {"name": "app", "value": "clunky-sabertooth-joomla"}
                   ]
                 },
                 "strategy": {
                   "rollingUpdate": {
                     "maxSurgeInt": 1,
                     "maxSurgeString": null,
                     "maxUnavailableInt": 1,
                     "maxUnavailableString": null
                   },
                   "type": "RollingUpdate"
                 },
                 "template": {
                   "metadata": {
                     "creationTimestamp": null,
                     "labels": [
                       {"name": "app", "value": "clunky-sabertooth-joomla"}
                     ]
                   },
                   "spec": {
                     "dnsPolicy": "ClusterFirst",
                     "restartPolicy": "Always",
                     "schedulerName": "default-scheduler",
                     "terminationGracePeriodSeconds": 30,
                     "volumes": [
                       {
                         "name": "joomla-data",
                         "persistentVolumeClaim": {
                           "claimName": "clunky-sabertooth-joomla-joomla",
                           "readOnly": false
                         }
                       },
                       {
                         "name": "apache-data",
                         "persistentVolumeClaim": {
                           "claimName": "clunky-sabertooth-joomla-apache",
                           "readOnly": false
                         }
                       }
                     ]
                   }
                 }
               },
               "replicaSets": [
                 {
                   "metadata": {
                     "name": "clunky-sabertooth-joomla-5d4ddc985d"
                   },
                   "pods": [
                     {
                       "metadata": {
                         "name": "clunky-sabertooth-joomla-5d4ddc985d-fpddz",
                         "namespace": "default",
                         "labels": [
                           {"name": "app",  "value": "clunky-sabertooth-joomla"},
                           {"name": "pod-template-hash", "value": "1808875418"}
                         ]
                       },
                       "spec": {
                         "dnsPolicy": "ClusterFirst",
                         "nodeName": "minikube",
                         "restartPolicy": "Always",
                         "schedulerName": "default-scheduler",
                         "serviceAccountName": "default",
                         "terminationGracePeriodSeconds": 30,
                         "tolerations": [
                           {
                             "effect": "NoExecute",
                             "key": "node.kubernetes.io/not-ready",
                             "operator": "Exists",
                             "tolerationSeconds": 300
                           },
                           {
                             "effect": "NoExecute",
                             "key": "node.kubernetes.io/unreachable",
                             "operator": "Exists",
                             "tolerationSeconds": 300
                           }
                         ],
                         "volumes": [
                           {
                             "name": "joomla-data",
                             "persistentVolumeClaim": {
                               "claimName": "clunky-sabertooth-joomla-joomla",
                               "readOnly": false
                             },
                             "secret": null
                           },
                           {
                             "name": "apache-data",
                             "persistentVolumeClaim": {
                               "claimName": "clunky-sabertooth-joomla-apache",
                               "readOnly": false
                             },
                             "secret": null
                           },
                           {
                             "name": "default-token-l6lb2",
                             "persistentVolumeClaim": null,
                             "secret": {
                               "defaultMode": 420,
                               "secretName": "default-token-l6lb2"
                             }
                           }
                         ]
                       }
                     }
                   ]
                 }
               ]
             }
           ]
        }`)
	simpletest(
		t,
		`{
           podByName(namespace: "default",
                     name: "clunky-sabertooth-joomla-5d4ddc985d-fpddz") {
             owner { metadata { name } }
             rootOwner { metadata { name } }
           }
         }`,
		`{
           "podByName": {
             "owner": {
               "metadata": { "name": "clunky-sabertooth-joomla-5d4ddc985d" }
             },
             "rootOwner": {
               "metadata": { "name": "clunky-sabertooth-joomla" }
             }
           }
         }`)
	simpletest(
		t,
		`{
           allPodsInNamespace(namespace: "default") {
             owner { metadata { name } }
             rootOwner { metadata { name } }
           }
         }`,
		`{
           "allPodsInNamespace": [
             {
               "owner": {
                 "metadata": { "name": "clunky-sabertooth-joomla-5d4ddc985d" }
               },
               "rootOwner": {
                 "metadata": { "name": "clunky-sabertooth-joomla" }
               }
             }
           ]
         }`)
	simpletest(
		t,
		`{
           allReplicaSets() {
             owner { metadata { name } }
             rootOwner { metadata { name } }
           }
         }`,
		`{
           "allReplicaSets": [
             {
               "owner": {
                 "metadata": { "name": "clunky-sabertooth-joomla" }
               },
               "rootOwner": {
                 "metadata": { "name": "clunky-sabertooth-joomla" }
               }
             }
           ]
         }`)
	simpletest(
		t,
		`{
           allDaemonSets() {
             owner { metadata { name namespace } }
             rootOwner { metadata { name namespace } }
             pods { metadata { name labels { name value } } }
           }
         }`,
		`{
           "allDaemonSets": [
             {
               "owner": {
                 "metadata": {
                   "name": "calico-node",
                   "namespace": "kube-system"
                 }
               },
               "rootOwner": {
                 "metadata": {
                   "name": "calico-node",
                   "namespace": "kube-system"
                 }
               },
               "pods": [
                 {
                   "metadata": {
                     "name": "calico-node-ddxfj",
                     "labels": [
                       {"name": "controller-revision-hash",
                        "value": "3909226423"},
                       {"name": "k8s-app", "value": "calico-node"},
                       {"name": "pod-template-generation", "value": "1"}
                     ]
                   }
                 }
               ]
             }
           ]
         }`)
	simpletest(
		t,
		`{
           allServices() {
             owner { metadata { name namespace } }
             rootOwner { metadata { name namespace } }
             selected { metadata { name labels { name value } } }
           }
         }`,
		`{
           "allServices": [
             {
               "owner": {
                 "metadata": {
                   "name": "mongo",
                   "namespace": "flonjella"
                 }
               },
               "rootOwner": {
                 "metadata": {
                   "name": "mongo",
                   "namespace": "flonjella"
                 }
               },
               "selected": [
                 {
                   "metadata": {
                     "name": "mongo-0",
                     "labels": [
                       {"name": "app", "value": "mongo"},
                       {"name": "controller-revision-hash",
                        "value": "mongo-fdd786d"},
                       {"name": "name", "value": "mongo"},
                       {"name": "statefulset.kubernetes.io/pod-name",
                        "value": "mongo-0"}
                     ]
                   }
                 }
               ]
             }
           ]
         }`)
	simpletest(
		t,
		`{
           allServicesInNamespace(namespace: "flonjella") {
             owner { metadata { name namespace } }
             rootOwner { metadata { name namespace } }
             selected { metadata { name labels { name value } } }
           }
         }`,
		`{
           "allServicesInNamespace": [
             {
               "owner": {
                 "metadata": {
                   "name": "mongo",
                   "namespace": "flonjella"
                 }
               },
               "rootOwner": {
                 "metadata": {
                   "name": "mongo",
                   "namespace": "flonjella"
                 }
               },
               "selected": [
                 {
                   "metadata": {
                     "name": "mongo-0",
                     "labels": [
                       {"name": "app", "value": "mongo"},
                       {"name": "controller-revision-hash",
                        "value": "mongo-fdd786d"},
                       {"name": "name", "value": "mongo"},
                       {"name": "statefulset.kubernetes.io/pod-name",
                        "value": "mongo-0"}
                     ]
                   }
                 }
               ]
             }
           ]
         }`)
	simpletest(
		t,
		`{
           serviceByName(namespace: "flonjella", name: "mongo") {
             owner { metadata { name namespace } }
             rootOwner { metadata { name namespace } }
             selected { metadata { name labels { name value } } }
           }
         }`,
		`{
           "serviceByName": {
             "owner": {
               "metadata": {
                 "name": "mongo",
                 "namespace": "flonjella"
               }
             },
             "rootOwner": {
               "metadata": {
                 "name": "mongo",
                 "namespace": "flonjella"
               }
             },
             "selected": [
               {
                 "metadata": {
                   "name": "mongo-0",
                   "labels": [
                     {"name": "app", "value": "mongo"},
                     {"name": "controller-revision-hash",
                      "value": "mongo-fdd786d"},
                     {"name": "name", "value": "mongo"},
                     {"name": "statefulset.kubernetes.io/pod-name",
                      "value": "mongo-0"}
                   ]
                 }
               }
             ]
           }
         }`)
}
