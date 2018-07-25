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
)

// The schema below defines the objects and relationships for Kubernetes.
// It is not yet complete.

var Schema = `
    schema {
       query: Query
       mutation: Mutation
    }
    # The query type, represents all of the entry points into our object graph
    type Query {
      # look up pods
      allPods(): [Pod]
      allPodsInNamespace(namespace: String!): [Pod]
      podByName(namespace: String!, name: String!): Pod
      # look up deployments
      allDeployments(): [Deployment]
      allDeploymentsInNamespace(namespace: String!): [Deployment]
      deploymentByName(namespace: String!, name: String!): Deployment
      # look up replica sets
      allReplicaSets(): [ReplicaSet]
      allReplicaSetsInNamespace(namespace: String!): [ReplicaSet]
      replicaSetByName(namespace: String!, name: String!): ReplicaSet
      # look up daemon sets
      allDaemonSets(): [DaemonSet]
      allDaemonSetsInNamespace(namespace: String!): [DaemonSet]
      daemonSetByName(namespace: String!, name: String!): DaemonSet
      # look up stateful sets
      allStatefulSets(): [StatefulSet]
      allStatefulSetsInNamespace(namespace: String!): [StatefulSet]
      statefulSetByName(namespace: String!, name: String!): StatefulSet
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
      # The pods controlled by this replicaSet
      pods: [Pod!]!
    }

    # A statefulSet
    type StatefulSet implements Resource {
      # The metadata for the statefulSet (name, labels, namespace, etc.)
      metadata: Metadata!
      # The direct owner of the statefulSet
      owner: Resource
      # The root owner of the statefulSet
      rootOwner: Resource
      # The pods controlled by this statefulSet
      pods: [Pod!]!
    }

    # A daemonSet
    type DaemonSet implements Resource {
      # The metadata for the daemonSet (name, labels, namespace, etc.)
      metadata: Metadata!
      # The direct owner of the daemonSet
      owner: Resource
      # The root owner of the daemonSet
      rootOwner: Resource
      # The pods controlled by this daemonSet
      pods: [Pod!]!
    }

    # A deployment
    type Deployment implements Resource {
      # The metadata for the deployment (name, labels, namespace, etc.)
      metadata: Metadata!
      # Minimum number of seconds for which a newly created pod should be
      # ready without any of its container crashing, for it to be considered
      # available (defaults to 0)
      minReadySeconds: Int!
      # Whether or not the deployment is paused
      paused: Boolean!
      # The maximum time in seconds for a deployment to make progress before
      # it is considered to be failed.
      progressDeadlineSeconds: Int!
      # Number of desired pods (default: 1).
      replicas: Int!
      # The number of old ReplicaSets to retain to allow rollback (default: 10).
      revisionHistoryLimit: Int!
      # Label selector for pods.
      selector: LabelSelector
      # The deployment strategy to use to replace existing pods with new ones.
      strategy: DeploymentStrategy!
      # Template describing the pods that will be created.
      # template: PodTemplateSpec! XXXXX - not yet
      # The direct owner of the deployment
      owner: Resource
      # The root owner of the deployment
      rootOwner: Resource
      # The replicaSets that are children of this deployment
      replicaSets: [ReplicaSet!]!
    }

    # metadata
    type Metadata {
      # When was the decorated object created
      creationTimestamp: String
      # Prefix for generated names
      generateName: String
      # Sequence number for state transitions
      generation: Int
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

    # LabelSelector for matching pods
    type LabelSelector {
      # constraint expressions for labels
      matchExpressions: [LabelSelectorRequirement!]
      # key/value matches
      matchLabels: [Label!]
    }

    # Constraint expression for labels
    type LabelSelectorRequirement {
      # The label key that the selector applies to
      key: String!
      # The expression operator
      operator: Operator!
      # The values to match against
      values: [String!]!
    }

    # Constraint operators for labels
    enum Operator {
      In NotIn Exists DoesNotExist
    }

    # deployment strategy
    type DeploymentStrategy {
      # Rolling update config parameters
      rollingUpdate: RollingUpdateDeployment
      # Type of deployment
      type: DeploymentStrategyType
    }

    # Types of deployment strategy
    enum DeploymentStrategyType {
      Recreate RollingUpdate
    }

    # The following section is a mess due to the questionable decision by
    # the Kubernetes team to make certain fields contain either ints or
    # strings (WHY?????)

    # rolling update parameters
    type  RollingUpdateDeployment {
      # The maximum number of pods that can be scheduled above the desired
      # number of pods.
      maxSurgeInt: Int
      maxSurgeString: String
      # The maximum number of pods that can be unavailable during the update.
      maxUnavailableInt: Int
      maxUnavailableString: String
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
const StatefulSetKind = "StatefulSet"
const DaemonSetKind = "DaemonSet"
const DeploymentKind = "Deployment"

// The root of all queries and mutations. All defined queries and mutations
// start as methods on Resolver
type Resolver struct {
}

// Objects in json are unmarshalled into map[string]interface{}
type JsonObject map[string]interface{}
type JsonArray []interface{}

// Pod lookups
func (r *Resolver) AllPods(ctx context.Context) *[]*podResolver {
	pset := getAllK8sObjsOfKind(
		ctx,
		PodKind,
		func(jobj map[string]interface{}) bool { return true })

	results := make([]*podResolver, len(pset))

	for idx, p := range pset {
		results[idx] = p.(*podResolver)
	}

	return &results
}

func (r *Resolver) AllPodsInNamespace(
	ctx context.Context,
	args *struct {
		Namespace string
	}) *[]*podResolver {
	pset := getAllK8sObjsOfKindInNamespace(
		ctx,
		PodKind,
		args.Namespace,
		func(jobj map[string]interface{}) bool { return true })

	results := make([]*podResolver, len(pset))

	for idx, p := range pset {
		results[idx] = p.(*podResolver)
	}

	return &results
}

func (r *Resolver) PodByName(
	ctx context.Context,
	args *struct {
		Namespace string
		Name      string
	}) *podResolver {
	return getK8sResource(ctx, PodKind, args.Namespace, args.Name).(*podResolver)
}

// Deployment lookups
func (r *Resolver) AllDeployments(ctx context.Context) *[]*deploymentResolver {
	dset := getAllK8sObjsOfKind(
		ctx,
		DeploymentKind,
		func(jobj map[string]interface{}) bool { return true })

	results := make([]*deploymentResolver, len(dset))

	for idx, d := range dset {
		results[idx] = d.(*deploymentResolver)
	}

	return &results
}

func (r *Resolver) AllDeploymentsInNamespace(
	ctx context.Context,
	args *struct {
		Namespace string
	}) *[]*deploymentResolver {
	dset := getAllK8sObjsOfKindInNamespace(
		ctx,
		DeploymentKind,
		args.Namespace,
		func(jobj map[string]interface{}) bool { return true })

	results := make([]*deploymentResolver, len(dset))

	for idx, p := range dset {
		results[idx] = p.(*deploymentResolver)
	}

	return &results
}

func (r *Resolver) DeploymentByName(
	ctx context.Context,
	args *struct {
		Namespace string
		Name      string
	}) *deploymentResolver {
	return getK8sResource(
		ctx, DeploymentKind, args.Namespace, args.Name).(*deploymentResolver)
}

// ReplicaSet lookups
func (r *Resolver) AllReplicaSets(ctx context.Context) *[]*replicaSetResolver {
	rset := getAllK8sObjsOfKind(
		ctx,
		ReplicaSetKind,
		func(jobj map[string]interface{}) bool { return true })

	results := make([]*replicaSetResolver, len(rset))

	for idx, r := range rset {
		results[idx] = r.(*replicaSetResolver)
	}

	return &results
}

func (r *Resolver) AllReplicaSetsInNamespace(
	ctx context.Context,
	args *struct {
		Namespace string
	}) *[]*replicaSetResolver {
	rset := getAllK8sObjsOfKindInNamespace(
		ctx,
		ReplicaSetKind,
		args.Namespace,
		func(jobj map[string]interface{}) bool { return true })

	results := make([]*replicaSetResolver, len(rset))

	for idx, p := range rset {
		results[idx] = p.(*replicaSetResolver)
	}

	return &results
}

func (r *Resolver) ReplicaSetByName(
	ctx context.Context,
	args *struct {
		Namespace string
		Name      string
	}) *replicaSetResolver {
	return getK8sResource(
		ctx, ReplicaSetKind, args.Namespace, args.Name).(*replicaSetResolver)
}

// StatefulSet lookups
func (r *Resolver) AllStatefulSets(ctx context.Context) *[]*statefulSetResolver {
	sset := getAllK8sObjsOfKind(
		ctx,
		StatefulSetKind,
		func(jobj map[string]interface{}) bool { return true })

	results := make([]*statefulSetResolver, len(sset))

	for idx, s := range sset {
		results[idx] = s.(*statefulSetResolver)
	}

	return &results
}

func (r *Resolver) AllStatefulSetsInNamespace(
	ctx context.Context,
	args *struct {
		Namespace string
	}) *[]*statefulSetResolver {
	sset := getAllK8sObjsOfKindInNamespace(
		ctx,
		StatefulSetKind,
		args.Namespace,
		func(jobj map[string]interface{}) bool { return true })

	results := make([]*statefulSetResolver, len(sset))

	for idx, p := range sset {
		results[idx] = p.(*statefulSetResolver)
	}

	return &results
}

func (r *Resolver) StatefulSetByName(
	ctx context.Context,
	args *struct {
		Namespace string
		Name      string
	}) *statefulSetResolver {
	return getK8sResource(
		ctx, StatefulSetKind, args.Namespace, args.Name).(*statefulSetResolver)
}

// DaemonSet lookups
func (r *Resolver) AllDaemonSets(ctx context.Context) *[]*daemonSetResolver {
	dset := getAllK8sObjsOfKind(
		ctx,
		DaemonSetKind,
		func(jobj map[string]interface{}) bool { return true })

	results := make([]*daemonSetResolver, len(dset))

	for idx, d := range dset {
		results[idx] = d.(*daemonSetResolver)
	}

	return &results
}

func (r *Resolver) AllDaemonSetsInNamespace(
	ctx context.Context,
	args *struct {
		Namespace string
	}) *[]*daemonSetResolver {
	dset := getAllK8sObjsOfKindInNamespace(
		ctx,
		DaemonSetKind,
		args.Namespace,
		func(jobj map[string]interface{}) bool { return true })

	results := make([]*daemonSetResolver, len(dset))

	for idx, p := range dset {
		results[idx] = p.(*daemonSetResolver)
	}

	return &results
}

func (r *Resolver) DaemonSetByName(
	ctx context.Context,
	args *struct {
		Namespace string
		Name      string
	}) *daemonSetResolver {
	return getK8sResource(
		ctx, DaemonSetKind, args.Namespace, args.Name).(*daemonSetResolver)
}
