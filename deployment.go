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
	Metadata                metadata
	MinReadySeconds         int32
	Paused                  bool
	ProgressDeadlineSeconds int32
	Replicas                int32
	RevisionHistoryLimit    int32
	Selector                *labelSelector
	Strategy                *deploymentStrategy
	//  Template PodTemplateSpec
	Owner       resource
	RootOwner   resource
	ReplicaSets *[]replicaSet
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
	r   *rollingUpdateDeployment
}

// Translate unmarshalled json into a deployment object
func mapToDeployment(
	ctx context.Context,
	jsonObj map[string]interface{}) deployment {
	jg := jgetter(jsonObj["spec"].(map[string]interface{}))
	return deployment{
		mapToMetadata(ctx, getNamespace(jsonObj), mapItem(jsonObj, "metadata")),
		jg.intItemOr("minReadySeconds", 0),
		jg.boolItemOr("paused", false),
		jg.intItemOr("progressDeadlineSeconds", 600),
		jg.intItemOr("replicas", 1),
		jg.intItemOr("revisionHistoryLimit", 10),
		mapToSelector(jg.objItemOr("selector", nil)),
		mapToStrategy(jg.objItemOr("strategy", nil)),
		nil,
		nil,
		nil}
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

	rudg := jgetter(*updateItem)
	ms := rudg.stringItemOr("maxSurgeString", "25%")
	mu := rudg.stringItemOr("maxUnavailableString", "25%")
	return &deploymentStrategy{&rollingUpdateDeployment{nil, &ms, nil, &mu},
		sType}
}

// Retrieve the ReplicaSets comprising the deployment
func getReplicaSets(ctx context.Context, d deployment) *[]replicaSet {
	depName := d.Metadata.Name
	depNamePrefix := depName + "-"
	depNamespace := d.Metadata.Namespace

	rsets := getAllK8sObjsOfKindInNamespace(
		ctx,
		"ReplicaSet",
		depNamespace,
		func(jobj map[string]interface{}) bool {
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

// Resource method implementations
func (r *deploymentResolver) Kind() string {
	return DeploymentKind
}

func (r *deploymentResolver) Metadata() *metadataResolver {
	return &metadataResolver{r.ctx, r.d.Metadata}
}

func (r *deploymentResolver) MinReadySeconds() int32 {
	return r.d.MinReadySeconds
}

func (r *deploymentResolver) Paused() bool {
	return r.d.Paused
}

func (r *deploymentResolver) ProgressDeadlineSeconds() int32 {
	return r.d.ProgressDeadlineSeconds
}

func (r *deploymentResolver) Replicas() int32 {
	return r.d.Replicas
}

func (r *deploymentResolver) RevisionHistoryLimit() int32 {
	return r.d.RevisionHistoryLimit
}

func (r *deploymentResolver) Selector() *labelSelectorResolver {
	if r.d.Selector != nil {
		return &labelSelectorResolver{r.ctx, *r.d.Selector}
	}

	return nil
}

func (r *deploymentResolver) Strategy() *deploymentStrategyResolver {
	return &deploymentStrategyResolver{r.ctx, r.d.Strategy}
}

func (r *deploymentResolver) Owner() *resourceResolver {
	return &resourceResolver{r.ctx, &deploymentResolver{r.ctx, r.d}}
}

func (r *deploymentResolver) RootOwner() *resourceResolver {
	return &resourceResolver{r.ctx, &deploymentResolver{r.ctx, r.d}}
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
	return &rollingUpdateDeploymentResolver{r.ctx, r.d.RollingUpdate}
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
