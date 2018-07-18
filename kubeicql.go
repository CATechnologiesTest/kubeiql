package main

import (
	// "bytes"
	"context"
	// "crypto/tls"
	// "crypto/x509"
	// "database/sql"
	// "encoding/json"
	// "errors"
	// "fmt"
	// log "github.com/Sirupsen/logrus"
	// "github.com/jackc/pgx"
	// "github.com/jackc/pgx/stdlib"
	// "io/ioutil"
	// "math"
	// "os"
	// "runtime"
	// "strconv"
	// "strings"
	// "sync"
	// "time"
)

// The schema below defines the objects and relationships that can be
// queried and modified. Each "type" below can be returned from
// queries so a "resolver" must be implemented. The resolver has a
// method for each field of the object and the graphql server calls a
// resolver method as needed based on what is requested by the
// user. Each type has a struct that holds its scalar values while set
// or list values are constructed on demand.

var Schema = `
    schema {
       query: Query
       mutation: Mutation
    }
    # The query type, represents all of the entry points into our object graph
    type Query {
      # look up pods
      allPods(): [Pod]
      podByName(namespace: String!, name: String!): Pod
      # look up deployments
      allDeployments(): [Deployment]
      deploymentByName(namespace: String!, name: String!): Deployment
      # look up replica sets
      allReplicaSets(): [ReplicaSet]
      replicaSetByName(namespace: String!, name: String!): ReplicaSet
    }

    # The mutation type, represents all updates we can make to our data
    type Mutation {
    }

    # Available logging levels
    enum LogLevel {
      debug
      info
      warning
      error
      fatal
      panic
    }

    # A pod
    type Pod implements Resource {
      # The metadata for the pod (name, labels, namespace, etc.)
      metadata: Metadata!
      # The direct owner of the pod
      owner: Resource
      # The root owner of the pod
      rootOwner: Resource
    }

    # A replicaSet
    type ReplicaSet implements Resource {
      # The metadata for the replicaSet (name, labels, namespace, etc.)
      metadata: Metadata!
      # The direct owner of the replicaSet
      owner: Resource
      # The root owner of the replicaSet
      rootOwner: Resource
    }

    # A deployment
    type Deployment implements Resource {
      # The metadata for the deployment (name, labels, namespace, etc.)
      metadata: Metadata!
      # The direct owner of the deployment
      owner: Resource
      # The root owner of the deployment
      rootOwner: Resource
      # The replica sets that are children of this deployment
      replicaSets: [ReplicaSet!]!
    }

    # metadata
    type Metadata {
      # When was the decorated object created
      creationTimestamp: String
      # Prefix for generated names
      generateName: String
      # Top level labels
      labels: [Label!]!
      # Generated name
      name: String!
      # Namespace containing the object
      namespace: String!
      # All owners
      ownerReferences: [Resource!]!
      # Version
      resourceVersion: String!
      # How to find this object
      selfLink: String!
      # UUID
      uid: String!
    }

    # A label
    type Label {
      # label name
      name: String!
      # label value
      value: String!
    }

    # Any Kubernetes resource
    interface Resource {
      # type of resource
      kind: String!
      # resource metadata
      metadata: Metadata!
      # resource direct owner
      owner: Resource
      # resource root owner
      rootOwner: Resource
    }
`

const PodKind = "Pod"
const ReplicaSetKind = "ReplicaSet"
const DeploymentKind = "Deployment"

// The root of all queries and mutations. All defined queries and mutations
// start as methods on Resolver
type Resolver struct {
}

func (r *Resolver) AllPods(ctx context.Context) *[]*podResolver {
	var podResolvers []*podResolver
	return &podResolvers
}

func (r *Resolver) PodByName(
	ctx context.Context,
	args *struct {
		Namespace string
		Name      string
	}) *podResolver {
	if pmap := getK8sResource(PodKind, args.Namespace, args.Name); pmap != nil {
		return &podResolver{ctx, mapToPod(ctx, pmap)}
	}

	return nil
}

func (r *Resolver) AllDeployments(ctx context.Context) *[]*deploymentResolver {
	var deploymentResolvers []*deploymentResolver
	return &deploymentResolvers
}

func (r *Resolver) DeploymentByName(
	ctx context.Context,
	args *struct {
		Namespace string
		Name      string
	}) *deploymentResolver {
	if dmap := getK8sResource(
		DeploymentKind, args.Namespace, args.Name); dmap != nil {
		return &deploymentResolver{ctx, mapToDeployment(ctx, dmap)}
	}

	return nil
}

func (r *Resolver) AllReplicaSets(ctx context.Context) *[]*replicaSetResolver {
	var replicaSetResolvers []*replicaSetResolver
	return &replicaSetResolvers
}

func (r *Resolver) ReplicaSetByName(
	ctx context.Context,
	args *struct {
		Namespace string
		Name      string
	}) *replicaSetResolver {
	if rmap := getK8sResource(
		ReplicaSetKind, args.Namespace, args.Name); rmap != nil {
		return &replicaSetResolver{ctx, mapToReplicaSet(ctx, rmap)}
	}

	return nil
}
