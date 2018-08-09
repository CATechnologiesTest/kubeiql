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
	"strings"
)

// Top level Kubernetes replicated controller. Deployments are built out
// of ReplicaSets.
type deployment struct {
	Metadata    metadata
	Spec        deploymentSpec
	Owner       resource
	RootOwner   resource
	ReplicaSets *[]replicaSet
}

type deploymentSpec struct {
	MinReadySeconds         int32
	Paused                  bool
	ProgressDeadlineSeconds int32
	Replicas                int32
	RevisionHistoryLimit    int32
	Selector                *labelSelector
	Strategy                *deploymentStrategy
	Template                podTemplateSpec
}

type podTemplateSpec struct {
	Metadata metadata
	Spec     podSpec
}

type labelSelector struct {
	MatchExpressions *[]labelSelectorRequirement
	MatchLabels      *[]label
}

type labelSelectorRequirement struct {
	Key      string
	Operator string
	Values   []string
}

type deploymentStrategy struct {
	RollingUpdate *rollingUpdateDeployment
	Type          *string
}

type rollingUpdateDeployment struct {
	MaxSurgeInt          *int32
	MaxSurgeString       *string
	MaxUnavailableInt    *int32
	MaxUnavailableString *string
}

type deploymentResolver struct {
	ctx context.Context
	d   deployment
}

type deploymentSpecResolver struct {
	ctx context.Context
	d   deploymentSpec
}

type labelSelectorResolver struct {
	ctx context.Context
	l   labelSelector
}

type deploymentStrategyResolver struct {
	ctx context.Context
	d   *deploymentStrategy
}

type labelSelectorRequirementResolver struct {
	ctx context.Context
	l   labelSelectorRequirement
}

type rollingUpdateDeploymentResolver struct {
	ctx context.Context
	r   rollingUpdateDeployment
}

type podTemplateSpecResolver struct {
	ctx context.Context
	p   podTemplateSpec
}

// Translate unmarshalled json into a deployment object
func mapToDeployment(
	ctx context.Context,
	jsonObj JsonObject) deployment {
	ns := getNamespace(jsonObj)
	return deployment{
		mapToMetadata(ctx, ns, mapItem(jsonObj, "metadata")),
		extractDeploymentSpec(ctx, ns, mapItem(jsonObj, "spec")),
		nil,
		nil,
		nil}
}

func extractDeploymentSpec(ctx context.Context, ns string, jsonObj JsonObject) deploymentSpec {
	jg := jgetter(jsonObj)
	template := mapItem(jsonObj, "template")
	return deploymentSpec{
		jg.intItemOr("minReadySeconds", 0),
		jg.boolItemOr("paused", false),
		jg.intItemOr("progressDeadlineSeconds", 600),
		jg.intItemOr("replicas", 1),
		jg.intItemOr("revisionHistoryLimit", 10),
		mapToSelector(jg.objItemOr("selector", nil)),
		mapToStrategy(jg.objItemOr("strategy", nil)),
		podTemplateSpec{
			mapToMetadata(ctx, ns, mapItem(template, "metadata")),
			mapToPodSpec(ctx, mapItem(template, "spec"))}}
}

func mapToSelector(sel *JsonObject) *labelSelector {
	if sel == nil {
		return nil
	}
	jg := jgetter(*sel)
	exprs := jg.arrayItemOr("matchExpressions", nil)
	labels := jg.objItemOr("matchLabels", nil)

	var exprsVal *[]labelSelectorRequirement

	if exprs != nil {
		eslice := make([]labelSelectorRequirement, len(*exprs))
		exprsVal = &eslice
		for idx, lsr := range *exprs {
			jg := jgetter(lsr.(JsonObject))
			vals := jg.arrayItem("values")
			strVals := make([]string, len(vals))
			for sidx, sval := range vals {
				strVals[sidx] = sval.(string)
			}
			(*exprsVal)[idx] = labelSelectorRequirement{
				jg.stringItem("key"),
				jg.stringItem("operator"),
				strVals}
		}
	}

	if labels != nil {
		lslice := make([]label, len(*labels))
		i := 0
		for k, v := range *labels {
			lslice[i] = label{k, v.(string)}
			i = i + 1
		}

		return &labelSelector{exprsVal, &lslice}
	}

	empty := make([]label, 0)
	return &labelSelector{exprsVal, &empty}
}

func mapToStrategy(strat *JsonObject) *deploymentStrategy {
	if strat == nil {
		return nil
	}
	jg := jgetter(*strat)
	sType := jg.stringRefItemOr("type", nil)
	var updateItem *JsonObject

	if *sType == "RollingUpdate" {
		updateItem = jg.objItemOr("rollingUpdate", nil)
	}

	if updateItem == nil {
		return &deploymentStrategy{nil, sType}
	}

	sval, spresent := (*updateItem)["maxSurge"]
	uval, upresent := (*updateItem)["maxUnavailable"]
	var ss, su string
	var is, iu int32
	var ssptr *string = nil
	var suptr *string = nil
	var isptr *int32 = nil
	var iuptr *int32 = nil
	defval := "25%"

	if !spresent {
		ss = defval
		ssptr = &ss
	} else if ssval, ok := sval.(string); ok {
		ss = ssval
		ssptr = &ss
	} else {
		is = toGQLInt(sval)
		isptr = &is
	}

	if !upresent {
		su = defval
		suptr = &su
	} else if suval, ok := uval.(string); ok {
		su = suval
		suptr = &su
	} else {
		iu = toGQLInt(uval)
		iuptr = &iu
	}

	return &deploymentStrategy{
		&rollingUpdateDeployment{isptr, ssptr, iuptr, suptr},
		sType}
}

// Retrieve the ReplicaSets comprising the deployment
func getReplicaSets(ctx context.Context, d deployment) *[]replicaSet {
	depName := *d.Metadata.Name
	depNamePrefix := depName + "-"
	depNamespace := *d.Metadata.Namespace

	rsets := getAllK8sObjsOfKindInNamespace(
		ctx,
		"ReplicaSet",
		depNamespace,
		func(jobj JsonObject) bool {
			return (strings.HasPrefix(getName(jobj), depNamePrefix) &&
				hasMatchingOwner(jobj, depName, DeploymentKind))
		})

	results := make([]replicaSet, len(rsets))

	for idx, rs := range rsets {
		rsr := rs.(*replicaSetResolver)
		results[idx] = rsr.r
	}

	return &results
}

func (r *labelSelectorResolver) MatchExpressions() *[]labelSelectorRequirementResolver {
	if r.l.MatchExpressions == nil {
		empty := make([]labelSelectorRequirementResolver, 0)
		return &empty
	}
	resolvers := make([]labelSelectorRequirementResolver,
		len(*r.l.MatchExpressions))
	for idx, val := range *r.l.MatchExpressions {
		resolvers[idx] = labelSelectorRequirementResolver{r.ctx, val}
	}
	return &resolvers
}

func (r *labelSelectorResolver) MatchLabels() *[]labelResolver {
	if r.l.MatchLabels == nil {
		empty := make([]labelResolver, 0)
		return &empty
	}
	resolvers := make([]labelResolver, len(*r.l.MatchLabels))
	for idx, val := range *r.l.MatchLabels {
		labelVal := val
		resolvers[idx] = labelResolver{r.ctx, &labelVal}
	}
	return &resolvers
}

func (r labelSelectorRequirementResolver) Key() string {
	return r.l.Key
}

func (r labelSelectorRequirementResolver) Operator() string {
	return r.l.Operator
}

func (r labelSelectorRequirementResolver) Values() []string {
	return r.l.Values
}

// Pod template spec implementations
func (r podTemplateSpecResolver) Metadata() metadataResolver {
	return metadataResolver{r.ctx, r.p.Metadata}
}

func (r podTemplateSpecResolver) Spec() podSpecResolver {
	return podSpecResolver{r.ctx, r.p.Spec}
}

// Resource method implementations
func (r *deploymentResolver) Kind() string {
	return DeploymentKind
}

func (r *deploymentResolver) Metadata() metadataResolver {
	return metadataResolver{r.ctx, r.d.Metadata}
}

func (r *deploymentResolver) Spec() deploymentSpecResolver {
	return deploymentSpecResolver{r.ctx, r.d.Spec}
}

func (r *deploymentResolver) Owner() *resourceResolver {
	return &resourceResolver{r.ctx, &deploymentResolver{r.ctx, r.d}}
}

func (r *deploymentResolver) RootOwner() *resourceResolver {
	return &resourceResolver{r.ctx, &deploymentResolver{r.ctx, r.d}}
}

// Deployment spec implementations
func (r deploymentSpecResolver) MinReadySeconds() int32 {
	return r.d.MinReadySeconds
}

func (r deploymentSpecResolver) Paused() bool {
	return r.d.Paused
}

func (r deploymentSpecResolver) ProgressDeadlineSeconds() int32 {
	return r.d.ProgressDeadlineSeconds
}

func (r deploymentSpecResolver) Replicas() int32 {
	return r.d.Replicas
}

func (r deploymentSpecResolver) RevisionHistoryLimit() int32 {
	return r.d.RevisionHistoryLimit
}

func (r deploymentSpecResolver) Selector() *labelSelectorResolver {
	if r.d.Selector != nil {
		return &labelSelectorResolver{r.ctx, *r.d.Selector}
	}

	return nil
}

func (r deploymentSpecResolver) Template() podTemplateSpecResolver {
	return podTemplateSpecResolver{r.ctx, r.d.Template}
}

func (r deploymentSpecResolver) Strategy() *deploymentStrategyResolver {
	return &deploymentStrategyResolver{r.ctx, r.d.Strategy}
}

// Resolve child ReplicaSets
func (r *deploymentResolver) ReplicaSets() []*replicaSetResolver {
	if r.d.ReplicaSets == nil {
		r.d.ReplicaSets = getReplicaSets(r.ctx, r.d)
	}

	var res []*replicaSetResolver
	for _, rs := range *r.d.ReplicaSets {
		res = append(res, &replicaSetResolver{r.ctx, rs})
	}
	if res == nil {
		res = make([]*replicaSetResolver, 0)
	}
	return res
}

func (r deploymentStrategyResolver) RollingUpdate() *rollingUpdateDeploymentResolver {
	if r.d.RollingUpdate == nil {
		return nil
	}
	return &rollingUpdateDeploymentResolver{r.ctx, *r.d.RollingUpdate}
}

func (r deploymentStrategyResolver) Type() *string {
	val := "RollingUpdate"
	if r.d == nil || r.d.Type == nil {
		return &val
	}
	return r.d.Type
}

func (r rollingUpdateDeploymentResolver) MaxSurgeInt() *int32 {
	return r.r.MaxSurgeInt
}

func (r rollingUpdateDeploymentResolver) MaxUnavailableInt() *int32 {
	return r.r.MaxUnavailableInt
}

func (r rollingUpdateDeploymentResolver) MaxSurgeString() *string {
	return r.r.MaxSurgeString
}

func (r rollingUpdateDeploymentResolver) MaxUnavailableString() *string {
	return r.r.MaxUnavailableString
}
