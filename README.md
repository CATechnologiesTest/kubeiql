# Kubeicql
A GraphQL interface for Kubernetes.

The goal of this project is to provide an alternative GraphQL
interface to a Kubernetes cluster. It is not intended to entirely replace the
ReST APIs as some of them (particularly the watch APIs) don't map well
onto GraphQL.

## Current Status

pre-alpha

* Queries are currently supported against Pods, Deployments,
ReplicaSets, StatefulSets, and DaemonSets.
* No mutations are yet implemented.
* The retrieval of data from the cluster is currently a hack on top of
kubectl to enable the development of the graph traversal code. The
plan is to use API Aggregation to put the cluster lookups into the API
Service.
* Tests are lacking
* Not yet built into a container (which we will need to deploy within
  the cluster's API service)

## Getting Started
To experiment with the API:

1. Download the code
2. Type <code>sh gobuild.sh
3. If your kubectl is located somewhere other than /usr/local/bin, set
the environment variable: KUBECTL_PATH to the location of your
executable (e.g. KUBECTL_PATH=/usr/share/local/bin/kubectl)
4. Run ./kubeicql

The server runs at port 8128. You can use curl to play with it as
shown in the examples below via the /query endpoint, or point your
browser at 'localhost:8128/' and experiment with the GraphiQL tool
(much more user-friendly).

## Examples

The query:

``` json
{
  daemonSetByName(namespace: "kube-system", name: "kube-proxy") {
    metadata {name namespace labels {name value}}
  }
}
```

<code>
curl -X POST -H"Content-Type: application/json"
http://localhost:8128/query -d
'{ "query": "{daemonSetByName(namespace: \"kube-system\", name: \"kube-proxy\") {    metadata {name namespace labels {name value}} pods {metadata {name}}}}"}'
</code>


returns:

```json
{
  "data": {
    "daemonSetByName": {
      "metadata": {
        "name": "kube-proxy",
        "namespace": "kube-system",
        "labels": [
          {
            "name": "k8s-app",
            "value": "kube-proxy"
          }
        ]
      },
      "pods": [
        {
          "metadata": {
            "name": "kube-proxy-7vhx5"
          }
        }
      ]
    }
  }
}

```

and the query:

``` json
{
  allPods() {
    owner {kind metadata {name}}
    rootOwner { kind metadata { name namespace }
      ... on StatefulSet {
        metadata { name }
      }
      ... on Deployment {
        replicaSets {
            metadata { name }
          pods { metadata { name } } }
      }
    }
  }
}
```

<code>
curl -X POST -H"Content-Type: application/json"
http://localhost:8128/query -d '{"query": "{allPods() {owner {kind
metadata {name}} rootOwner { kind metadata { name namespace } ... on
StatefulSet { metadata { name } } ... on Deployment { replicaSets {
metadata { name } pods { metadata { name } } } } } } }" }'
</code>


returns:

```json
{
  "data": {
    "allPods": [
      {
        "owner": {
          "kind": "ReplicaSet",
          "metadata": {
            "name": "backend-549447ccf"
          }
        },
        "rootOwner": {
          "kind": "Deployment",
          "metadata": {
            "name": "backend",
            "namespace": "default"
          },
          "replicaSets": [
            {
              "metadata": {
                "name": "backend-549447ccf"
              },
              "pods": [
                {
                  "metadata": {
                    "name": "backend-549447ccf-4zphf"
                  }
                },
                {
                  "metadata": {
                    "name": "clunky-sabertooth-joomla-5d4ddc985d-fpddz"
                  }
                },
                {
                  "metadata": {
                    "name": "clunky-sabertooth-mariadb-0"
                  }
                }
                // ...
              ]
            }
          ]
        }
      },
      {
        "owner": {
          "kind": "ReplicaSet",
          "metadata": {
            "name": "clunky-sabertooth-joomla-5d4ddc985d"
          }
        },
        "rootOwner": {
          "kind": "Deployment",
          "metadata": {
            "name": "clunky-sabertooth-joomla",
            "namespace": "default"
          },
          "replicaSets": [
            {
              "metadata": {
                "name": "backend-549447ccf"
              },
              "pods": [
                {
                  "metadata": {
                    "name": "backend-549447ccf-4zphf"
                  }
                },
                {
                  "metadata": {
                    "name": "clunky-sabertooth-joomla-5d4ddc985d-fpddz"
                  }
                }
                //...
              ]
            }
          ]
        }
      },
      {
        "owner": {
          "kind": "StatefulSet",
          "metadata": {
            "name": "clunky-sabertooth-mariadb"
          }
        },
        "rootOwner": {
          "kind": "StatefulSet",
          "metadata": {
            "name": "clunky-sabertooth-mariadb",
            "namespace": "default"
          }
        }
      },
      {
        "owner": {
          "kind": "ReplicaSet",
          "metadata": {
            "name": "ui-9c6c8d79"
          }
        },
        "rootOwner": {
          "kind": "Deployment",
          "metadata": {
            "name": "ui",
            "namespace": "default"
          },
          "replicaSets": [
            {
              "metadata": {
                "name": "backend-549447ccf"
              },
              "pods": [
                {
                  "metadata": {
                    "name": "backend-549447ccf-4zphf"
                  }
                },
                {
                  "metadata": {
                    "name": "clunky-sabertooth-joomla-5d4ddc985d-fpddz"
                  }
                }
                // ...
              ]
            }
          ]
        }
      },
      {
        "owner": {
          "kind": "ReplicaSet",
          "metadata": {
            "name": "ui-9c6c8d79"
          }
        },
        "rootOwner": {
          "kind": "Deployment",
          "metadata": {
            "name": "ui",
            "namespace": "default"
          },
          "replicaSets": [
            {
              "metadata": {
                "name": "backend-549447ccf"
              },
              "pods": [
                {
                  "metadata": {
                    "name": "backend-549447ccf-4zphf"
                  }
                },
                {
                  "metadata": {
                    "name": "clunky-sabertooth-joomla-5d4ddc985d-fpddz"
                  }
                }
                // ...
              ]
            }
          ]
        }
      },
      {
        "owner": {
          "kind": "Pod",
          "metadata": {
            "name": "etcd-minikube"
          }
        },
        "rootOwner": {
          "kind": "Pod",
          "metadata": {
            "name": "etcd-minikube",
            "namespace": "kube-system"
          }
        }
      },
      {
        "owner": {
          "kind": "Pod",
          "metadata": {
            "name": "kube-addon-manager-minikube"
          }
        },
        "rootOwner": {
          "kind": "Pod",
          "metadata": {
            "name": "kube-addon-manager-minikube",
            "namespace": "kube-system"
          }
        }
      },
      {
        "owner": {
          "kind": "Pod",
          "metadata": {
            "name": "kube-apiserver-minikube"
          }
        },
        "rootOwner": {
          "kind": "Pod",
          "metadata": {
            "name": "kube-apiserver-minikube",
            "namespace": "kube-system"
          }
        }
      },
      {
        "owner": {
          "kind": "Pod",
          "metadata": {
            "name": "kube-controller-manager-minikube"
          }
        },
        "rootOwner": {
          "kind": "Pod",
          "metadata": {
            "name": "kube-controller-manager-minikube",
            "namespace": "kube-system"
          }
        }
      },
      {
        "owner": {
          "kind": "ReplicaSet",
          "metadata": {
            "name": "kube-dns-86f4d74b45"
          }
        },
        "rootOwner": {
          "kind": "Deployment",
          "metadata": {
            "name": "kube-dns",
            "namespace": "kube-system"
          },
          "replicaSets": [
            {
              "metadata": {
                "name": "backend-549447ccf"
              },
              "pods": [
                {
                  "metadata": {
                    "name": "backend-549447ccf-4zphf"
                  }
                },
                {
                  "metadata": {
                    "name": "clunky-sabertooth-joomla-5d4ddc985d-fpddz"
                  }
                }
                // ...
              ]
            }
          ]
        }
      },
      {
        "owner": {
          "kind": "DaemonSet",
          "metadata": {
            "name": "kube-proxy"
          }
        },
        "rootOwner": {
          "kind": "DaemonSet",
          "metadata": {
            "name": "kube-proxy",
            "namespace": "kube-system"
          }
        }
      },
      {
        "owner": {
          "kind": "Pod",
          "metadata": {
            "name": "kube-scheduler-minikube"
          }
        },
        "rootOwner": {
          "kind": "Pod",
          "metadata": {
            "name": "kube-scheduler-minikube",
            "namespace": "kube-system"
          }
        }
      }
      // ...
    ]
  }
}

```
## License

The work done has been licensed under Apache License 2.0. The license file can be found [here](LICENSE). You can find
out more about the license at [www.apache.org/licenses/LICENSE-2.0](//www.apache.org/licenses/LICENSE-2.0).

## Questions?

Feel free to [contact us](mailto:support@yipee.io).
