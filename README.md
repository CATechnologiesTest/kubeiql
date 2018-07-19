# kubeicql
A GraphQL interface for Kubernetes.

The goal of this project is to provide an alternative GraphQL
interface to a Kubernetes cluster. It is not intended to entirely replace the
ReST APIs as some of them (particularly the watch APIs) don't map well
onto GraphQL.

## Current Status:

pre-alpha

* Queries are currently supported against Pods, Deployments,
ReplicaSets, StatefulSets, and DaemonSets.
* No mutations are yet implemented.
* The retrieval of data from the cluster is currently a hack on top of
kubectl to enable the development of the graph traversal code. The
plan is to use API Aggregation to put the cluster lookups into the API
Service.

Examples:

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
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-94jlc"
                  }
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-vcst6"
                  }
                },
                {
                  "metadata": {
                    "name": "etcd-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-addon-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-apiserver-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-controller-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-dns-86f4d74b45-54w7m"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-proxy-7vhx5"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-scheduler-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kubernetes-dashboard-5498ccf677-mss7v"
                  }
                },
                {
                  "metadata": {
                    "name": "storage-provisioner"
                  }
                },
                {
                  "metadata": {
                    "name": "tiller-deploy-7f5f67578d-6ltwl"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-6cd457d86c-x6cdm"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-db-784446fdd6-tgqkd"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-779cd58f9b-m4gkp"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-db-6794f65f5d-btflx"
                  }
                },
                {
                  "metadata": {
                    "name": "front-end-679d7bcb77-cbkrs"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-755bd9f786-s5n4r"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-db-84bb8f48d6-drgts"
                  }
                },
                {
                  "metadata": {
                    "name": "payment-674658f686-rx28n"
                  }
                },
                {
                  "metadata": {
                    "name": "queue-master-5f98bbd67-bj8qx"
                  }
                },
                {
                  "metadata": {
                    "name": "rabbitmq-86d44dd846-f6vjx"
                  }
                },
                {
                  "metadata": {
                    "name": "shipping-79786fb956-9p859"
                  }
                },
                {
                  "metadata": {
                    "name": "user-6995984547-t6jrb"
                  }
                },
                {
                  "metadata": {
                    "name": "user-db-fc7b47fb9-44k4s"
                  }
                }
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
                },
                {
                  "metadata": {
                    "name": "clunky-sabertooth-mariadb-0"
                  }
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-94jlc"
                  }
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-vcst6"
                  }
                },
                {
                  "metadata": {
                    "name": "etcd-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-addon-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-apiserver-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-controller-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-dns-86f4d74b45-54w7m"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-proxy-7vhx5"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-scheduler-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kubernetes-dashboard-5498ccf677-mss7v"
                  }
                },
                {
                  "metadata": {
                    "name": "storage-provisioner"
                  }
                },
                {
                  "metadata": {
                    "name": "tiller-deploy-7f5f67578d-6ltwl"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-6cd457d86c-x6cdm"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-db-784446fdd6-tgqkd"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-779cd58f9b-m4gkp"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-db-6794f65f5d-btflx"
                  }
                },
                {
                  "metadata": {
                    "name": "front-end-679d7bcb77-cbkrs"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-755bd9f786-s5n4r"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-db-84bb8f48d6-drgts"
                  }
                },
                {
                  "metadata": {
                    "name": "payment-674658f686-rx28n"
                  }
                },
                {
                  "metadata": {
                    "name": "queue-master-5f98bbd67-bj8qx"
                  }
                },
                {
                  "metadata": {
                    "name": "rabbitmq-86d44dd846-f6vjx"
                  }
                },
                {
                  "metadata": {
                    "name": "shipping-79786fb956-9p859"
                  }
                },
                {
                  "metadata": {
                    "name": "user-6995984547-t6jrb"
                  }
                },
                {
                  "metadata": {
                    "name": "user-db-fc7b47fb9-44k4s"
                  }
                }
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
                },
                {
                  "metadata": {
                    "name": "clunky-sabertooth-mariadb-0"
                  }
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-94jlc"
                  }
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-vcst6"
                  }
                },
                {
                  "metadata": {
                    "name": "etcd-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-addon-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-apiserver-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-controller-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-dns-86f4d74b45-54w7m"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-proxy-7vhx5"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-scheduler-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kubernetes-dashboard-5498ccf677-mss7v"
                  }
                },
                {
                  "metadata": {
                    "name": "storage-provisioner"
                  }
                },
                {
                  "metadata": {
                    "name": "tiller-deploy-7f5f67578d-6ltwl"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-6cd457d86c-x6cdm"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-db-784446fdd6-tgqkd"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-779cd58f9b-m4gkp"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-db-6794f65f5d-btflx"
                  }
                },
                {
                  "metadata": {
                    "name": "front-end-679d7bcb77-cbkrs"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-755bd9f786-s5n4r"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-db-84bb8f48d6-drgts"
                  }
                },
                {
                  "metadata": {
                    "name": "payment-674658f686-rx28n"
                  }
                },
                {
                  "metadata": {
                    "name": "queue-master-5f98bbd67-bj8qx"
                  }
                },
                {
                  "metadata": {
                    "name": "rabbitmq-86d44dd846-f6vjx"
                  }
                },
                {
                  "metadata": {
                    "name": "shipping-79786fb956-9p859"
                  }
                },
                {
                  "metadata": {
                    "name": "user-6995984547-t6jrb"
                  }
                },
                {
                  "metadata": {
                    "name": "user-db-fc7b47fb9-44k4s"
                  }
                }
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
                },
                {
                  "metadata": {
                    "name": "clunky-sabertooth-mariadb-0"
                  }
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-94jlc"
                  }
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-vcst6"
                  }
                },
                {
                  "metadata": {
                    "name": "etcd-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-addon-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-apiserver-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-controller-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-dns-86f4d74b45-54w7m"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-proxy-7vhx5"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-scheduler-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kubernetes-dashboard-5498ccf677-mss7v"
                  }
                },
                {
                  "metadata": {
                    "name": "storage-provisioner"
                  }
                },
                {
                  "metadata": {
                    "name": "tiller-deploy-7f5f67578d-6ltwl"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-6cd457d86c-x6cdm"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-db-784446fdd6-tgqkd"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-779cd58f9b-m4gkp"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-db-6794f65f5d-btflx"
                  }
                },
                {
                  "metadata": {
                    "name": "front-end-679d7bcb77-cbkrs"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-755bd9f786-s5n4r"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-db-84bb8f48d6-drgts"
                  }
                },
                {
                  "metadata": {
                    "name": "payment-674658f686-rx28n"
                  }
                },
                {
                  "metadata": {
                    "name": "queue-master-5f98bbd67-bj8qx"
                  }
                },
                {
                  "metadata": {
                    "name": "rabbitmq-86d44dd846-f6vjx"
                  }
                },
                {
                  "metadata": {
                    "name": "shipping-79786fb956-9p859"
                  }
                },
                {
                  "metadata": {
                    "name": "user-6995984547-t6jrb"
                  }
                },
                {
                  "metadata": {
                    "name": "user-db-fc7b47fb9-44k4s"
                  }
                }
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
                },
                {
                  "metadata": {
                    "name": "clunky-sabertooth-mariadb-0"
                  }
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-94jlc"
                  }
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-vcst6"
                  }
                },
                {
                  "metadata": {
                    "name": "etcd-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-addon-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-apiserver-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-controller-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-dns-86f4d74b45-54w7m"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-proxy-7vhx5"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-scheduler-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kubernetes-dashboard-5498ccf677-mss7v"
                  }
                },
                {
                  "metadata": {
                    "name": "storage-provisioner"
                  }
                },
                {
                  "metadata": {
                    "name": "tiller-deploy-7f5f67578d-6ltwl"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-6cd457d86c-x6cdm"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-db-784446fdd6-tgqkd"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-779cd58f9b-m4gkp"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-db-6794f65f5d-btflx"
                  }
                },
                {
                  "metadata": {
                    "name": "front-end-679d7bcb77-cbkrs"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-755bd9f786-s5n4r"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-db-84bb8f48d6-drgts"
                  }
                },
                {
                  "metadata": {
                    "name": "payment-674658f686-rx28n"
                  }
                },
                {
                  "metadata": {
                    "name": "queue-master-5f98bbd67-bj8qx"
                  }
                },
                {
                  "metadata": {
                    "name": "rabbitmq-86d44dd846-f6vjx"
                  }
                },
                {
                  "metadata": {
                    "name": "shipping-79786fb956-9p859"
                  }
                },
                {
                  "metadata": {
                    "name": "user-6995984547-t6jrb"
                  }
                },
                {
                  "metadata": {
                    "name": "user-db-fc7b47fb9-44k4s"
                  }
                }
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
      },
      {
        "owner": {
          "kind": "ReplicaSet",
          "metadata": {
            "name": "kubernetes-dashboard-5498ccf677"
          }
        },
        "rootOwner": {
          "kind": "Deployment",
          "metadata": {
            "name": "kubernetes-dashboard",
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
                },
                {
                  "metadata": {
                    "name": "clunky-sabertooth-mariadb-0"
                  }
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-94jlc"
                  }
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-vcst6"
                  }
                },
                {
                  "metadata": {
                    "name": "etcd-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-addon-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-apiserver-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-controller-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-dns-86f4d74b45-54w7m"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-proxy-7vhx5"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-scheduler-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kubernetes-dashboard-5498ccf677-mss7v"
                  }
                },
                {
                  "metadata": {
                    "name": "storage-provisioner"
                  }
                },
                {
                  "metadata": {
                    "name": "tiller-deploy-7f5f67578d-6ltwl"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-6cd457d86c-x6cdm"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-db-784446fdd6-tgqkd"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-779cd58f9b-m4gkp"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-db-6794f65f5d-btflx"
                  }
                },
                {
                  "metadata": {
                    "name": "front-end-679d7bcb77-cbkrs"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-755bd9f786-s5n4r"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-db-84bb8f48d6-drgts"
                  }
                },
                {
                  "metadata": {
                    "name": "payment-674658f686-rx28n"
                  }
                },
                {
                  "metadata": {
                    "name": "queue-master-5f98bbd67-bj8qx"
                  }
                },
                {
                  "metadata": {
                    "name": "rabbitmq-86d44dd846-f6vjx"
                  }
                },
                {
                  "metadata": {
                    "name": "shipping-79786fb956-9p859"
                  }
                },
                {
                  "metadata": {
                    "name": "user-6995984547-t6jrb"
                  }
                },
                {
                  "metadata": {
                    "name": "user-db-fc7b47fb9-44k4s"
                  }
                }
              ]
            }
          ]
        }
      },
      {
        "owner": {
          "kind": "Pod",
          "metadata": {
            "name": "storage-provisioner"
          }
        },
        "rootOwner": {
          "kind": "Pod",
          "metadata": {
            "name": "storage-provisioner",
            "namespace": "kube-system"
          }
        }
      },
      {
        "owner": {
          "kind": "ReplicaSet",
          "metadata": {
            "name": "tiller-deploy-7f5f67578d"
          }
        },
        "rootOwner": {
          "kind": "Deployment",
          "metadata": {
            "name": "tiller-deploy",
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
                },
                {
                  "metadata": {
                    "name": "clunky-sabertooth-mariadb-0"
                  }
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-94jlc"
                  }
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-vcst6"
                  }
                },
                {
                  "metadata": {
                    "name": "etcd-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-addon-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-apiserver-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-controller-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-dns-86f4d74b45-54w7m"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-proxy-7vhx5"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-scheduler-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kubernetes-dashboard-5498ccf677-mss7v"
                  }
                },
                {
                  "metadata": {
                    "name": "storage-provisioner"
                  }
                },
                {
                  "metadata": {
                    "name": "tiller-deploy-7f5f67578d-6ltwl"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-6cd457d86c-x6cdm"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-db-784446fdd6-tgqkd"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-779cd58f9b-m4gkp"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-db-6794f65f5d-btflx"
                  }
                },
                {
                  "metadata": {
                    "name": "front-end-679d7bcb77-cbkrs"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-755bd9f786-s5n4r"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-db-84bb8f48d6-drgts"
                  }
                },
                {
                  "metadata": {
                    "name": "payment-674658f686-rx28n"
                  }
                },
                {
                  "metadata": {
                    "name": "queue-master-5f98bbd67-bj8qx"
                  }
                },
                {
                  "metadata": {
                    "name": "rabbitmq-86d44dd846-f6vjx"
                  }
                },
                {
                  "metadata": {
                    "name": "shipping-79786fb956-9p859"
                  }
                },
                {
                  "metadata": {
                    "name": "user-6995984547-t6jrb"
                  }
                },
                {
                  "metadata": {
                    "name": "user-db-fc7b47fb9-44k4s"
                  }
                }
              ]
            }
          ]
        }
      },
      {
        "owner": {
          "kind": "ReplicaSet",
          "metadata": {
            "name": "carts-6cd457d86c"
          }
        },
        "rootOwner": {
          "kind": "Deployment",
          "metadata": {
            "name": "carts",
            "namespace": "sock-shop"
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
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-94jlc"
                  }
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-vcst6"
                  }
                },
                {
                  "metadata": {
                    "name": "etcd-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-addon-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-apiserver-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-controller-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-dns-86f4d74b45-54w7m"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-proxy-7vhx5"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-scheduler-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kubernetes-dashboard-5498ccf677-mss7v"
                  }
                },
                {
                  "metadata": {
                    "name": "storage-provisioner"
                  }
                },
                {
                  "metadata": {
                    "name": "tiller-deploy-7f5f67578d-6ltwl"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-6cd457d86c-x6cdm"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-db-784446fdd6-tgqkd"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-779cd58f9b-m4gkp"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-db-6794f65f5d-btflx"
                  }
                },
                {
                  "metadata": {
                    "name": "front-end-679d7bcb77-cbkrs"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-755bd9f786-s5n4r"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-db-84bb8f48d6-drgts"
                  }
                },
                {
                  "metadata": {
                    "name": "payment-674658f686-rx28n"
                  }
                },
                {
                  "metadata": {
                    "name": "queue-master-5f98bbd67-bj8qx"
                  }
                },
                {
                  "metadata": {
                    "name": "rabbitmq-86d44dd846-f6vjx"
                  }
                },
                {
                  "metadata": {
                    "name": "shipping-79786fb956-9p859"
                  }
                },
                {
                  "metadata": {
                    "name": "user-6995984547-t6jrb"
                  }
                },
                {
                  "metadata": {
                    "name": "user-db-fc7b47fb9-44k4s"
                  }
                }
              ]
            }
          ]
        }
      },
      {
        "owner": {
          "kind": "ReplicaSet",
          "metadata": {
            "name": "carts-db-784446fdd6"
          }
        },
        "rootOwner": {
          "kind": "Deployment",
          "metadata": {
            "name": "carts-db",
            "namespace": "sock-shop"
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
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-94jlc"
                  }
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-vcst6"
                  }
                },
                {
                  "metadata": {
                    "name": "etcd-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-addon-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-apiserver-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-controller-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-dns-86f4d74b45-54w7m"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-proxy-7vhx5"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-scheduler-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kubernetes-dashboard-5498ccf677-mss7v"
                  }
                },
                {
                  "metadata": {
                    "name": "storage-provisioner"
                  }
                },
                {
                  "metadata": {
                    "name": "tiller-deploy-7f5f67578d-6ltwl"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-6cd457d86c-x6cdm"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-db-784446fdd6-tgqkd"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-779cd58f9b-m4gkp"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-db-6794f65f5d-btflx"
                  }
                },
                {
                  "metadata": {
                    "name": "front-end-679d7bcb77-cbkrs"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-755bd9f786-s5n4r"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-db-84bb8f48d6-drgts"
                  }
                },
                {
                  "metadata": {
                    "name": "payment-674658f686-rx28n"
                  }
                },
                {
                  "metadata": {
                    "name": "queue-master-5f98bbd67-bj8qx"
                  }
                },
                {
                  "metadata": {
                    "name": "rabbitmq-86d44dd846-f6vjx"
                  }
                },
                {
                  "metadata": {
                    "name": "shipping-79786fb956-9p859"
                  }
                },
                {
                  "metadata": {
                    "name": "user-6995984547-t6jrb"
                  }
                },
                {
                  "metadata": {
                    "name": "user-db-fc7b47fb9-44k4s"
                  }
                }
              ]
            }
          ]
        }
      },
      {
        "owner": {
          "kind": "ReplicaSet",
          "metadata": {
            "name": "catalogue-779cd58f9b"
          }
        },
        "rootOwner": {
          "kind": "Deployment",
          "metadata": {
            "name": "catalogue",
            "namespace": "sock-shop"
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
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-94jlc"
                  }
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-vcst6"
                  }
                },
                {
                  "metadata": {
                    "name": "etcd-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-addon-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-apiserver-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-controller-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-dns-86f4d74b45-54w7m"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-proxy-7vhx5"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-scheduler-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kubernetes-dashboard-5498ccf677-mss7v"
                  }
                },
                {
                  "metadata": {
                    "name": "storage-provisioner"
                  }
                },
                {
                  "metadata": {
                    "name": "tiller-deploy-7f5f67578d-6ltwl"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-6cd457d86c-x6cdm"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-db-784446fdd6-tgqkd"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-779cd58f9b-m4gkp"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-db-6794f65f5d-btflx"
                  }
                },
                {
                  "metadata": {
                    "name": "front-end-679d7bcb77-cbkrs"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-755bd9f786-s5n4r"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-db-84bb8f48d6-drgts"
                  }
                },
                {
                  "metadata": {
                    "name": "payment-674658f686-rx28n"
                  }
                },
                {
                  "metadata": {
                    "name": "queue-master-5f98bbd67-bj8qx"
                  }
                },
                {
                  "metadata": {
                    "name": "rabbitmq-86d44dd846-f6vjx"
                  }
                },
                {
                  "metadata": {
                    "name": "shipping-79786fb956-9p859"
                  }
                },
                {
                  "metadata": {
                    "name": "user-6995984547-t6jrb"
                  }
                },
                {
                  "metadata": {
                    "name": "user-db-fc7b47fb9-44k4s"
                  }
                }
              ]
            }
          ]
        }
      },
      {
        "owner": {
          "kind": "ReplicaSet",
          "metadata": {
            "name": "catalogue-db-6794f65f5d"
          }
        },
        "rootOwner": {
          "kind": "Deployment",
          "metadata": {
            "name": "catalogue-db",
            "namespace": "sock-shop"
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
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-94jlc"
                  }
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-vcst6"
                  }
                },
                {
                  "metadata": {
                    "name": "etcd-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-addon-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-apiserver-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-controller-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-dns-86f4d74b45-54w7m"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-proxy-7vhx5"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-scheduler-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kubernetes-dashboard-5498ccf677-mss7v"
                  }
                },
                {
                  "metadata": {
                    "name": "storage-provisioner"
                  }
                },
                {
                  "metadata": {
                    "name": "tiller-deploy-7f5f67578d-6ltwl"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-6cd457d86c-x6cdm"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-db-784446fdd6-tgqkd"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-779cd58f9b-m4gkp"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-db-6794f65f5d-btflx"
                  }
                },
                {
                  "metadata": {
                    "name": "front-end-679d7bcb77-cbkrs"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-755bd9f786-s5n4r"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-db-84bb8f48d6-drgts"
                  }
                },
                {
                  "metadata": {
                    "name": "payment-674658f686-rx28n"
                  }
                },
                {
                  "metadata": {
                    "name": "queue-master-5f98bbd67-bj8qx"
                  }
                },
                {
                  "metadata": {
                    "name": "rabbitmq-86d44dd846-f6vjx"
                  }
                },
                {
                  "metadata": {
                    "name": "shipping-79786fb956-9p859"
                  }
                },
                {
                  "metadata": {
                    "name": "user-6995984547-t6jrb"
                  }
                },
                {
                  "metadata": {
                    "name": "user-db-fc7b47fb9-44k4s"
                  }
                }
              ]
            }
          ]
        }
      },
      {
        "owner": {
          "kind": "ReplicaSet",
          "metadata": {
            "name": "front-end-679d7bcb77"
          }
        },
        "rootOwner": {
          "kind": "Deployment",
          "metadata": {
            "name": "front-end",
            "namespace": "sock-shop"
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
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-94jlc"
                  }
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-vcst6"
                  }
                },
                {
                  "metadata": {
                    "name": "etcd-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-addon-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-apiserver-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-controller-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-dns-86f4d74b45-54w7m"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-proxy-7vhx5"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-scheduler-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kubernetes-dashboard-5498ccf677-mss7v"
                  }
                },
                {
                  "metadata": {
                    "name": "storage-provisioner"
                  }
                },
                {
                  "metadata": {
                    "name": "tiller-deploy-7f5f67578d-6ltwl"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-6cd457d86c-x6cdm"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-db-784446fdd6-tgqkd"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-779cd58f9b-m4gkp"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-db-6794f65f5d-btflx"
                  }
                },
                {
                  "metadata": {
                    "name": "front-end-679d7bcb77-cbkrs"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-755bd9f786-s5n4r"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-db-84bb8f48d6-drgts"
                  }
                },
                {
                  "metadata": {
                    "name": "payment-674658f686-rx28n"
                  }
                },
                {
                  "metadata": {
                    "name": "queue-master-5f98bbd67-bj8qx"
                  }
                },
                {
                  "metadata": {
                    "name": "rabbitmq-86d44dd846-f6vjx"
                  }
                },
                {
                  "metadata": {
                    "name": "shipping-79786fb956-9p859"
                  }
                },
                {
                  "metadata": {
                    "name": "user-6995984547-t6jrb"
                  }
                },
                {
                  "metadata": {
                    "name": "user-db-fc7b47fb9-44k4s"
                  }
                }
              ]
            }
          ]
        }
      },
      {
        "owner": {
          "kind": "ReplicaSet",
          "metadata": {
            "name": "orders-755bd9f786"
          }
        },
        "rootOwner": {
          "kind": "Deployment",
          "metadata": {
            "name": "orders",
            "namespace": "sock-shop"
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
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-94jlc"
                  }
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-vcst6"
                  }
                },
                {
                  "metadata": {
                    "name": "etcd-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-addon-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-apiserver-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-controller-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-dns-86f4d74b45-54w7m"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-proxy-7vhx5"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-scheduler-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kubernetes-dashboard-5498ccf677-mss7v"
                  }
                },
                {
                  "metadata": {
                    "name": "storage-provisioner"
                  }
                },
                {
                  "metadata": {
                    "name": "tiller-deploy-7f5f67578d-6ltwl"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-6cd457d86c-x6cdm"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-db-784446fdd6-tgqkd"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-779cd58f9b-m4gkp"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-db-6794f65f5d-btflx"
                  }
                },
                {
                  "metadata": {
                    "name": "front-end-679d7bcb77-cbkrs"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-755bd9f786-s5n4r"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-db-84bb8f48d6-drgts"
                  }
                },
                {
                  "metadata": {
                    "name": "payment-674658f686-rx28n"
                  }
                },
                {
                  "metadata": {
                    "name": "queue-master-5f98bbd67-bj8qx"
                  }
                },
                {
                  "metadata": {
                    "name": "rabbitmq-86d44dd846-f6vjx"
                  }
                },
                {
                  "metadata": {
                    "name": "shipping-79786fb956-9p859"
                  }
                },
                {
                  "metadata": {
                    "name": "user-6995984547-t6jrb"
                  }
                },
                {
                  "metadata": {
                    "name": "user-db-fc7b47fb9-44k4s"
                  }
                }
              ]
            }
          ]
        }
      },
      {
        "owner": {
          "kind": "ReplicaSet",
          "metadata": {
            "name": "orders-db-84bb8f48d6"
          }
        },
        "rootOwner": {
          "kind": "Deployment",
          "metadata": {
            "name": "orders-db",
            "namespace": "sock-shop"
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
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-94jlc"
                  }
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-vcst6"
                  }
                },
                {
                  "metadata": {
                    "name": "etcd-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-addon-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-apiserver-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-controller-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-dns-86f4d74b45-54w7m"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-proxy-7vhx5"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-scheduler-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kubernetes-dashboard-5498ccf677-mss7v"
                  }
                },
                {
                  "metadata": {
                    "name": "storage-provisioner"
                  }
                },
                {
                  "metadata": {
                    "name": "tiller-deploy-7f5f67578d-6ltwl"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-6cd457d86c-x6cdm"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-db-784446fdd6-tgqkd"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-779cd58f9b-m4gkp"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-db-6794f65f5d-btflx"
                  }
                },
                {
                  "metadata": {
                    "name": "front-end-679d7bcb77-cbkrs"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-755bd9f786-s5n4r"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-db-84bb8f48d6-drgts"
                  }
                },
                {
                  "metadata": {
                    "name": "payment-674658f686-rx28n"
                  }
                },
                {
                  "metadata": {
                    "name": "queue-master-5f98bbd67-bj8qx"
                  }
                },
                {
                  "metadata": {
                    "name": "rabbitmq-86d44dd846-f6vjx"
                  }
                },
                {
                  "metadata": {
                    "name": "shipping-79786fb956-9p859"
                  }
                },
                {
                  "metadata": {
                    "name": "user-6995984547-t6jrb"
                  }
                },
                {
                  "metadata": {
                    "name": "user-db-fc7b47fb9-44k4s"
                  }
                }
              ]
            }
          ]
        }
      },
      {
        "owner": {
          "kind": "ReplicaSet",
          "metadata": {
            "name": "payment-674658f686"
          }
        },
        "rootOwner": {
          "kind": "Deployment",
          "metadata": {
            "name": "payment",
            "namespace": "sock-shop"
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
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-94jlc"
                  }
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-vcst6"
                  }
                },
                {
                  "metadata": {
                    "name": "etcd-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-addon-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-apiserver-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-controller-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-dns-86f4d74b45-54w7m"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-proxy-7vhx5"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-scheduler-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kubernetes-dashboard-5498ccf677-mss7v"
                  }
                },
                {
                  "metadata": {
                    "name": "storage-provisioner"
                  }
                },
                {
                  "metadata": {
                    "name": "tiller-deploy-7f5f67578d-6ltwl"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-6cd457d86c-x6cdm"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-db-784446fdd6-tgqkd"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-779cd58f9b-m4gkp"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-db-6794f65f5d-btflx"
                  }
                },
                {
                  "metadata": {
                    "name": "front-end-679d7bcb77-cbkrs"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-755bd9f786-s5n4r"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-db-84bb8f48d6-drgts"
                  }
                },
                {
                  "metadata": {
                    "name": "payment-674658f686-rx28n"
                  }
                },
                {
                  "metadata": {
                    "name": "queue-master-5f98bbd67-bj8qx"
                  }
                },
                {
                  "metadata": {
                    "name": "rabbitmq-86d44dd846-f6vjx"
                  }
                },
                {
                  "metadata": {
                    "name": "shipping-79786fb956-9p859"
                  }
                },
                {
                  "metadata": {
                    "name": "user-6995984547-t6jrb"
                  }
                },
                {
                  "metadata": {
                    "name": "user-db-fc7b47fb9-44k4s"
                  }
                }
              ]
            }
          ]
        }
      },
      {
        "owner": {
          "kind": "ReplicaSet",
          "metadata": {
            "name": "queue-master-5f98bbd67"
          }
        },
        "rootOwner": {
          "kind": "Deployment",
          "metadata": {
            "name": "queue-master",
            "namespace": "sock-shop"
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
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-94jlc"
                  }
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-vcst6"
                  }
                },
                {
                  "metadata": {
                    "name": "etcd-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-addon-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-apiserver-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-controller-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-dns-86f4d74b45-54w7m"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-proxy-7vhx5"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-scheduler-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kubernetes-dashboard-5498ccf677-mss7v"
                  }
                },
                {
                  "metadata": {
                    "name": "storage-provisioner"
                  }
                },
                {
                  "metadata": {
                    "name": "tiller-deploy-7f5f67578d-6ltwl"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-6cd457d86c-x6cdm"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-db-784446fdd6-tgqkd"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-779cd58f9b-m4gkp"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-db-6794f65f5d-btflx"
                  }
                },
                {
                  "metadata": {
                    "name": "front-end-679d7bcb77-cbkrs"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-755bd9f786-s5n4r"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-db-84bb8f48d6-drgts"
                  }
                },
                {
                  "metadata": {
                    "name": "payment-674658f686-rx28n"
                  }
                },
                {
                  "metadata": {
                    "name": "queue-master-5f98bbd67-bj8qx"
                  }
                },
                {
                  "metadata": {
                    "name": "rabbitmq-86d44dd846-f6vjx"
                  }
                },
                {
                  "metadata": {
                    "name": "shipping-79786fb956-9p859"
                  }
                },
                {
                  "metadata": {
                    "name": "user-6995984547-t6jrb"
                  }
                },
                {
                  "metadata": {
                    "name": "user-db-fc7b47fb9-44k4s"
                  }
                }
              ]
            }
          ]
        }
      },
      {
        "owner": {
          "kind": "ReplicaSet",
          "metadata": {
            "name": "rabbitmq-86d44dd846"
          }
        },
        "rootOwner": {
          "kind": "Deployment",
          "metadata": {
            "name": "rabbitmq",
            "namespace": "sock-shop"
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
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-94jlc"
                  }
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-vcst6"
                  }
                },
                {
                  "metadata": {
                    "name": "etcd-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-addon-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-apiserver-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-controller-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-dns-86f4d74b45-54w7m"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-proxy-7vhx5"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-scheduler-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kubernetes-dashboard-5498ccf677-mss7v"
                  }
                },
                {
                  "metadata": {
                    "name": "storage-provisioner"
                  }
                },
                {
                  "metadata": {
                    "name": "tiller-deploy-7f5f67578d-6ltwl"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-6cd457d86c-x6cdm"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-db-784446fdd6-tgqkd"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-779cd58f9b-m4gkp"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-db-6794f65f5d-btflx"
                  }
                },
                {
                  "metadata": {
                    "name": "front-end-679d7bcb77-cbkrs"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-755bd9f786-s5n4r"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-db-84bb8f48d6-drgts"
                  }
                },
                {
                  "metadata": {
                    "name": "payment-674658f686-rx28n"
                  }
                },
                {
                  "metadata": {
                    "name": "queue-master-5f98bbd67-bj8qx"
                  }
                },
                {
                  "metadata": {
                    "name": "rabbitmq-86d44dd846-f6vjx"
                  }
                },
                {
                  "metadata": {
                    "name": "shipping-79786fb956-9p859"
                  }
                },
                {
                  "metadata": {
                    "name": "user-6995984547-t6jrb"
                  }
                },
                {
                  "metadata": {
                    "name": "user-db-fc7b47fb9-44k4s"
                  }
                }
              ]
            }
          ]
        }
      },
      {
        "owner": {
          "kind": "ReplicaSet",
          "metadata": {
            "name": "shipping-79786fb956"
          }
        },
        "rootOwner": {
          "kind": "Deployment",
          "metadata": {
            "name": "shipping",
            "namespace": "sock-shop"
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
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-94jlc"
                  }
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-vcst6"
                  }
                },
                {
                  "metadata": {
                    "name": "etcd-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-addon-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-apiserver-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-controller-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-dns-86f4d74b45-54w7m"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-proxy-7vhx5"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-scheduler-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kubernetes-dashboard-5498ccf677-mss7v"
                  }
                },
                {
                  "metadata": {
                    "name": "storage-provisioner"
                  }
                },
                {
                  "metadata": {
                    "name": "tiller-deploy-7f5f67578d-6ltwl"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-6cd457d86c-x6cdm"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-db-784446fdd6-tgqkd"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-779cd58f9b-m4gkp"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-db-6794f65f5d-btflx"
                  }
                },
                {
                  "metadata": {
                    "name": "front-end-679d7bcb77-cbkrs"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-755bd9f786-s5n4r"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-db-84bb8f48d6-drgts"
                  }
                },
                {
                  "metadata": {
                    "name": "payment-674658f686-rx28n"
                  }
                },
                {
                  "metadata": {
                    "name": "queue-master-5f98bbd67-bj8qx"
                  }
                },
                {
                  "metadata": {
                    "name": "rabbitmq-86d44dd846-f6vjx"
                  }
                },
                {
                  "metadata": {
                    "name": "shipping-79786fb956-9p859"
                  }
                },
                {
                  "metadata": {
                    "name": "user-6995984547-t6jrb"
                  }
                },
                {
                  "metadata": {
                    "name": "user-db-fc7b47fb9-44k4s"
                  }
                }
              ]
            }
          ]
        }
      },
      {
        "owner": {
          "kind": "ReplicaSet",
          "metadata": {
            "name": "user-6995984547"
          }
        },
        "rootOwner": {
          "kind": "Deployment",
          "metadata": {
            "name": "user",
            "namespace": "sock-shop"
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
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-94jlc"
                  }
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-vcst6"
                  }
                },
                {
                  "metadata": {
                    "name": "etcd-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-addon-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-apiserver-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-controller-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-dns-86f4d74b45-54w7m"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-proxy-7vhx5"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-scheduler-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kubernetes-dashboard-5498ccf677-mss7v"
                  }
                },
                {
                  "metadata": {
                    "name": "storage-provisioner"
                  }
                },
                {
                  "metadata": {
                    "name": "tiller-deploy-7f5f67578d-6ltwl"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-6cd457d86c-x6cdm"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-db-784446fdd6-tgqkd"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-779cd58f9b-m4gkp"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-db-6794f65f5d-btflx"
                  }
                },
                {
                  "metadata": {
                    "name": "front-end-679d7bcb77-cbkrs"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-755bd9f786-s5n4r"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-db-84bb8f48d6-drgts"
                  }
                },
                {
                  "metadata": {
                    "name": "payment-674658f686-rx28n"
                  }
                },
                {
                  "metadata": {
                    "name": "queue-master-5f98bbd67-bj8qx"
                  }
                },
                {
                  "metadata": {
                    "name": "rabbitmq-86d44dd846-f6vjx"
                  }
                },
                {
                  "metadata": {
                    "name": "shipping-79786fb956-9p859"
                  }
                },
                {
                  "metadata": {
                    "name": "user-6995984547-t6jrb"
                  }
                },
                {
                  "metadata": {
                    "name": "user-db-fc7b47fb9-44k4s"
                  }
                }
              ]
            }
          ]
        }
      },
      {
        "owner": {
          "kind": "ReplicaSet",
          "metadata": {
            "name": "user-db-fc7b47fb9"
          }
        },
        "rootOwner": {
          "kind": "Deployment",
          "metadata": {
            "name": "user-db",
            "namespace": "sock-shop"
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
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-94jlc"
                  }
                },
                {
                  "metadata": {
                    "name": "ui-9c6c8d79-vcst6"
                  }
                },
                {
                  "metadata": {
                    "name": "etcd-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-addon-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-apiserver-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-controller-manager-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-dns-86f4d74b45-54w7m"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-proxy-7vhx5"
                  }
                },
                {
                  "metadata": {
                    "name": "kube-scheduler-minikube"
                  }
                },
                {
                  "metadata": {
                    "name": "kubernetes-dashboard-5498ccf677-mss7v"
                  }
                },
                {
                  "metadata": {
                    "name": "storage-provisioner"
                  }
                },
                {
                  "metadata": {
                    "name": "tiller-deploy-7f5f67578d-6ltwl"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-6cd457d86c-x6cdm"
                  }
                },
                {
                  "metadata": {
                    "name": "carts-db-784446fdd6-tgqkd"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-779cd58f9b-m4gkp"
                  }
                },
                {
                  "metadata": {
                    "name": "catalogue-db-6794f65f5d-btflx"
                  }
                },
                {
                  "metadata": {
                    "name": "front-end-679d7bcb77-cbkrs"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-755bd9f786-s5n4r"
                  }
                },
                {
                  "metadata": {
                    "name": "orders-db-84bb8f48d6-drgts"
                  }
                },
                {
                  "metadata": {
                    "name": "payment-674658f686-rx28n"
                  }
                },
                {
                  "metadata": {
                    "name": "queue-master-5f98bbd67-bj8qx"
                  }
                },
                {
                  "metadata": {
                    "name": "rabbitmq-86d44dd846-f6vjx"
                  }
                },
                {
                  "metadata": {
                    "name": "shipping-79786fb956-9p859"
                  }
                },
                {
                  "metadata": {
                    "name": "user-6995984547-t6jrb"
                  }
                },
                {
                  "metadata": {
                    "name": "user-db-fc7b47fb9-44k4s"
                  }
                }
              ]
            }
          ]
        }
      }
    ]
  }
}

```
