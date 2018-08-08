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
      # behavior specification
      spec: PodSpec!
      # The direct owner of the pod
      owner: Resource
      # The root owner of the pod
      rootOwner: Resource
      # XXX PodStatus
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
      # Description of the deployment
      spec: DeploymentSpec!
      # The direct owner of the deployment
      owner: Resource
      # The root owner of the deployment
      rootOwner: Resource
      # The replicaSets that are children of this deployment
      replicaSets: [ReplicaSet!]!
    }

    # A deployment specification
    type DeploymentSpec {
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
      template: PodTemplateSpec!
    }

    # metadata
    type Metadata { # annotations??
      # When was the decorated object created
      creationTimestamp: String
      # Prefix for generated names
      generateName: String
      # Sequence number for state transitions
      generation: Int
      # Top level labels
      labels: [Label!]!
      # Generated name
      name: String
      # Namespace containing the object
      namespace: String
      # All owners
      ownerReferences: [Resource!]
      # Version
      resourceVersion: String
      # How to find this object
      selfLink: String
      # UUID
      uid: String
    }

    # PodTemplateSpec
    type PodTemplateSpec {
      # standard metadata
      metadata: Metadata!
      # Specification for generated pods
      spec: PodSpec!
    }

    # PodSpec
    type PodSpec {
      # Optional duration in seconds the pod may be active on the node
      # relative to StartTime before the system will actively try to mark it
      # failed and kill associated containers. Value must be a positive integer.
      activeDeadlineSeconds: Int
      # Scheduling constraints for the pod, if specified
      affinity: Affinity
      # AutomountServiceAccountToken indicates whether a service account token
      # should be automatically mounted.
      automountServiceAccountToken: Boolean!
      # List of containers belonging to the pod
      containers: [Container!]!
      # DNS parameters of a pod
      dnsConfig: PodDNSConfig
      # DNS policy for the pod -- defaults to "ClusterFirst"
      dnsPolicy: DNSPolicy
      # List of hosts and IPs to inject into pod host file
      hostAliases: [HostAlias!]
      # Whether or not to use the host ipc namespace -- default: false
      hostIPC: Boolean!
      # Whether or not to use host networking -- default: false
      hostNetwork: Boolean!
      # Whether or not to use the host pid namespace -- default: false
      hostPID: Boolean!
      # Hostname for the pod -- defaults to system-generated value
      hostname: String
      # List of references to secrets in the same namespace used to pull images
      imagePullSecrets: [LocalObjectReference!]
      # Initialization containers for the pod
      initContainers: [Container!]
      # Name of specific host on which to schedule the pod (if any)
      nodeName: String
      # A selector that must be true for the pod to fit on a node
      nodeSelector: NodeSelector
      # Priority of the pod
      priority: Int
      # Class name of priority class for the pod
      priorityClassName: String
      # Conditions to evaluate for pod readiness
      readinessGates: [PodReadinessGate!]
      # Restart policy for all containers within the pod -- default: Always
      restartPolicy: RestartPolicy!
      # Specific scheduler to use for pod
      schedulerName: String
      # Pod-level security attributes and common container settings
      securityContext: PodSecurityContext
      # Name of service account to use when running this pod
      serviceAccountName: String
      # Whether or not to share a single process namespace between all containers
      # default: false
      shareProcessNamespace: Boolean!
      # Desired pod subdomain, if any
      subdomain: String
      # Duration in seconds the pod needs to terminate gracefully -- default: 30s
      terminationGracePeriodSeconds: Int!
      # Tolerations for the pod
      tolerations: [Toleration!]
      # List of volumes that can be mounted by containers in the pod
      volumes: [Volume!]
    }

    # Persistent storage
    type Volume {
      #
      # INCOMPLETE
      #

      # config map used to populate the volume
      configMap: ConfigMapVolumeSource
      # volume name -- must be a DNS_LABEL and unique within the pod
      name: String!
      # a pre-existing file or directory on the host machine
      hostPath: HostPathVolumeSource
      # reference to a persistent volume claim
      persistentVolumeClaim: PersistentVolumeClaimVolumeSource
      # a secret that should populate the volume
      secret: SecretVolumeSource
    }

    # Volume sources
    type ConfigMapVolumeSource {
      # mode bits to use on created files by default -- default value: 0644
      defaultMode: Int
      # specific config map items to expose -- if unset, expose all
      items: [KeyToPath!]
      # name of map
      name: String!
      # Is it okay if the config map does not exist? -- default: false
      optional: Boolean!
    }

    type HostPathVolumeSource {
      # path of the directory on the host
      path: String!
      # type for host path volume -- defaults to: ""
      type: String
    }

    type PersistentVolumeClaimVolumeSource {
      # name of a PersistentVolumeClaim in the same namespace as the pod
      claimName: String!
      # Is the claim read only? -- default: false
      readOnly: Boolean!
    }

    type SecretVolumeSource {
      # mode bits to use on created files by default -- default value: 0644
      defaultMode: Int
      # specific secret items to expose -- if unset, expose all
      items: [KeyToPath!]
      # Is it okay if the secret does not exist? -- default: false
      optional: Boolean!
      # name of secret
      secretName: String!
    }

    # kernel parameters to set
    type Sysctl {
      # name of a parameter
      name: String!
      # parameter value
      value: String!
    }

    # mapping of string key to path within a volume
    type KeyToPath {
      # the key to project
      key: String!
      # mode bits to use on the file [0 .. 0777]; if empty, uses volume default
      mode: Int
      # relative path of the file to map
      path: String!
    }

    # Reference to a pod condition
    type PodReadinessGate {
      # a condition in the condition list for the pod
      conditionType: String!
    }

    # Pod-level security attributes and common container settings
    type PodSecurityContext {
      # supplemental group that applies to all containers in a pod
      fsGroup: Int
      # GID for container process entrypoint
      runAsGroup: Int
      # Must the container run as non-root?
      runAsNonRoot: Boolean!
      # UID for container process entrypoint
      runAsUser: Int
      # SELinux context to apply to all containers
      seLinuxOptions: SELinuxOptions
      # list of groups applied to the first process run in each container in
      # addition to the primary group
      supplementalGroups: [Int!]
      # list of namespaced sysctls user for the pod
      sysctls: [Sysctl!]
    }

    # Conditions determining whether a node can host a particular pod
    type Toleration {
      # the taint effect to match
      effect: TolerationEffect
      # the taint key the toleration references; empty means match all
      key: String
      # relationship between key and value
      operator: TolerationOperator
      # the period of time the toleration tolerates the taint; unset means
      # forever, zero or negative means immediately evict
      tolerationSeconds: Int
      # taint value matched by the toleration -- should be empty if operator
      # is "Exists"
      value: String
    }

    # Operators for tolerations
    enum TolerationOperator {
      Exists Equal
    }

    # Toleration effect possible values
    enum TolerationEffect {
      NoSchedule PreferNoSchedule NoExecute
    }

    # SELinux labels to apply to a container
    type SELinuxOptions {
      # level label
      level: String
      # role label
      role: String
      # type label
      type: String
      #user label
      user: String
    }

    # LocalObjectReference
    type LocalObjectReference {
      # name of referenced object
      name: String!
    }

    # DNS support
    type PodDNSConfig {
      # list of DNS name server IP addresses
      nameservers: [String!]
      # DNS resolver options
      options: [PodDNSConfigOption!]
      # DNS search domains for host-name lookup
      searches: [String!]
    }

    # Options for DNS config
    type PodDNSConfigOption {
      name: String!
      value: String!
    }

    # HostAlias
    type HostAlias {
      # Hostnames for the IP address
      hostnames: [String!]!
      # IP address of the host
      ip: String!
    }

    # Container
    type Container {
      # XXX not yet
    }

    # Affinity
    type Affinity {
      # affinity with nodes
      nodeAffinity: NodeAffinity
      # affinity with other pods
      podAffinity: PodAffinity
      # anti-affinity with other pods
      podAntiAffinity: PodAntiAffinity
    }

    # NodeAffinity
    type NodeAffinity {
      # nodes satisfying these conditions will be preferred
      preferredDuringSchedulingIgnoredDuringExecution: [PreferredSchedulingTerm!]
      # nodes satisfying these conditions will be required
      requiredDuringSchedulingIgnoredDuringExecution: NodeSelector
    }

    # NodeSelector
    type NodeSelector {
      # list of terms used to match nodes for deployment (ORed together)
      nodeSelectorTerms: [NodeSelectorTerm!]!
    }

    # NodeSelectorTerm
    type NodeSelectorTerm {
      # node requirements based on node labels
      matchExpressions: [NodeSelectorRequirement!]
      # node requirements based on node fields
      matchFields: [NodeSelectorRequirement!]
    }

    # Requirement expression matched against node labels
    type NodeSelectorRequirement {
      # The node key that the selector applies to
      key: String!
      # The expression operator
      operator: NodeOperator!
      # The values to match against
      values: [String!]
    }

    # PodAffinity
    type PodAffinity {
      # nodes satisfying these conditions will be preferred
      preferredDuringSchedulingIgnoredDuringExecution: [WeightedPodAffinityTerm!]
      # nodes satisfying these conditions will be required
      requiredDuringSchedulingIgnoredDuringExecution: PodAffinityTerm
    }

    # PodAntiAffinity
    type PodAntiAffinity {
      # nodes satisfying these conditions will be preferred
      preferredDuringSchedulingIgnoredDuringExecution: [WeightedPodAffinityTerm!]
      # nodes satisfying these conditions will be required
      requiredDuringSchedulingIgnoredDuringExecution: PodAffinityTerm
    }

    # WeightedPodAffinityTerm
    type WeightedPodAffinityTerm {
      # term associated with a weight
      podAffinityTerm: PodAffinityTerm!
      # weight for the term
      weight: Int!
    }

    # PodAffinityTerm
    type PodAffinityTerm {
      # selector to match other pod labels
      labelSelector: LabelSelector
      # which namespaces to match against -- defaults to the pod namespace
      namespaces: [String!]
      # whether the pod should be co-located or not co-located with pods
      # matching the selector
      topologyKey: String!
    }

    # PreferredSchedulingTerm
    type PreferredSchedulingTerm {
      # node preference
      preference: NodeSelectorTerm!
      # weight associated with the term
      weight: Int!
    }

    # Node operator values
    enum NodeOperator {
      In NotIn Exists DoesNotExist Gt Lt
    }

    # RestartPolicy values
    enum RestartPolicy {
      Always OnFailure Never
    }

    # DNSPolicy values
    enum DNSPolicy {
      ClusterFirstWithHostNet ClusterFirst Default None
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
      operator: LabelOperator!
      # The values to match against
      values: [String!]!
    }

    # Constraint operators for labels
    enum LabelOperator {
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
type JsonObject = map[string]interface{}
type JsonArray = []interface{}

// Pod lookups
func (r *Resolver) AllPods(ctx context.Context) *[]*podResolver {
	pset := getAllK8sObjsOfKind(
		ctx,
		PodKind,
		func(jobj JsonObject) bool { return true })

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
		func(jobj JsonObject) bool { return true })

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
		func(jobj JsonObject) bool { return true })

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
		func(jobj JsonObject) bool { return true })

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
		func(jobj JsonObject) bool { return true })

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
		func(jobj JsonObject) bool { return true })

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
		func(jobj JsonObject) bool { return true })

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
		func(jobj JsonObject) bool { return true })

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
		func(jobj JsonObject) bool { return true })

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
		func(jobj JsonObject) bool { return true })

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
