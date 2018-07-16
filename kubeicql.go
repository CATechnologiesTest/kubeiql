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
// or list values are constructed on demand. Each "input" below
// represents an argument to a query rather than a returnable
// Ids are represented by the String type (which is convertible
// to/from string). This is the type the graphql server expects.

var Schema = `
    schema {
       query: Query
       mutation: Mutation
    }
    # The query type, represents all of the entry points into our object graph
    type Query {
      # look up pods
      allPods(): [Pod]
      podById(id: String!): Pod
      # look up deployments
      allDeployments(): [Deployment]
      deploymentById(id: String!): Deployment
      # look up replica sets
      allReplicaSets(): [ReplicaSet]
      replicaSetById(id: String!): ReplicaSet
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
    type Pod {
      # The ID of the pod
      id: String!
      # The metadata for the pod (name, labels, namespace, etc.)
      metadata: Metadata!
      # The direct owner of the pod
      owner: Resource
      # The root owner of the pod
      rootOwner: Resource
    }

    # A replicaSet
    type ReplicaSet {
      # The ID of the replicaSet
      id: String!
      # The metadata for the replicaSet (name, labels, namespace, etc.)
      metadata: Metadata!
      # The direct owner of the replicaSet
      owner: Resource
      # The root owner of the replicaSet
      rootOwner: Resource
    }

    # A deployment
    type Deployment {
      # The ID of the deployment
      id: String!
      # The metadata for the deployment (name, labels, namespace, etc.)
      metadata: Metadata!
      # The direct owner of the deployment
      owner: Resource
      # The root owner of the deployment
      rootOwner: Resource
    }

    # metadata
    type Metadata {
      # When was the decorated object created
      creationTimestamp: String
      # Prefix for generated names
      generateName: String!
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
      # resource id
      id: String!
      # resource metadata
      metadata: Metadata!
      # resource direct owner
      owner: Resource
      # resource root owner
      rootOwner: Resource
    }
`

// The root of all queries and mutations. All defined queries and mutations
// start as methods on Resolver
type Resolver struct {
}

func (r *Resolver) AllPods(ctx context.Context) *[]*podResolver {
	var podResolvers []*podResolver
	return &podResolvers
}

func (r *Resolver) PodById(
	ctx context.Context,
	args *struct{ ID string }) *podResolver {
	return nil
}

func (r *Resolver) AllDeployments(ctx context.Context) *[]*deploymentResolver {
	var deploymentResolvers []*deploymentResolver
	return &deploymentResolvers
}

func (r *Resolver) DeploymentById(
	ctx context.Context,
	args *struct{ ID string }) *deploymentResolver {
	return nil
}

func (r *Resolver) AllReplicaSets(ctx context.Context) *[]*replicaSetResolver {
	var replicaSetResolvers []*replicaSetResolver
	return &replicaSetResolvers
}

func (r *Resolver) ReplicaSetById(
	ctx context.Context,
	args *struct{ ID string }) *replicaSetResolver {
	return nil
}
