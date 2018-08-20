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
)

// Services expose functionality to the outside world
type service struct {
	Metadata metadata
	Spec     serviceSpec
	Selected []*resource
}

type serviceSpec struct {
	ClusterIP                *string
	ExternalIPs              *[]string
	ExternalName             *string
	ExternalTrafficPolicy    *string
	HealthCheckNodePort      *int32
	LoadBalancerIP           *string
	LoadBalancerSourceRanges *[]string
	Ports                    *[]servicePort
	PublishNotReadyAddresses *bool
	Selector                 []label
	SessionAffinity          *string
	SessionAffinityConfig    *sessionAffinityConfig
	Type                     string
}

type sessionAffinityConfig struct {
	ClientIP clientIPConfig
}

type clientIPConfig struct {
	TimeoutSeconds int32
}

type servicePort struct {
	Name             *string
	NodePort         *int32
	Port             int32
	Protocol         *string
	TargetPortString *string
	TargetPortInt    *int32
}

type serviceResolver struct {
	ctx context.Context
	s   service
}

type serviceSpecResolver struct {
	ctx context.Context
	s   serviceSpec
}

type sessionAffinityConfigResolver struct {
	ctx context.Context
	s   sessionAffinityConfig
}

type clientIPConfigResolver struct {
	ctx context.Context
	s   clientIPConfig
}

type servicePortResolver struct {
	ctx context.Context
	s   servicePort
}

// Translate unmarshalled json into a deployment object
func mapToService(
	ctx context.Context,
	jsonObj JsonObject) service {
	ns := getNamespace(jsonObj)
	meta := mapToMetadata(ctx, ns, mapItem(jsonObj, "metadata"))
	return service{
		meta,
		extractServiceSpec(ctx, jsonObj),
		extractSelected(ctx, ns, jsonObj)}
}

func extractSelected(ctx context.Context, ns string, jsonObj JsonObject) []*resource {
	return []*resource{}
}

func extractServiceSpec(ctx context.Context, jsonObj JsonObject) serviceSpec {
	return serviceSpec{
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		[]label{},
		nil,
		nil,
		"ClientIP"}
}

// Service implementations
func (r *serviceResolver) Kind() string {
	return ServiceKind
}

func (r *serviceResolver) Metadata() metadataResolver {
	return metadataResolver{r.ctx, r.s.Metadata}
}

func (r *serviceResolver) Spec() serviceSpecResolver {
	return serviceSpecResolver{r.ctx, r.s.Spec}
}

func (r *serviceResolver) Owner() *resourceResolver {
	return &resourceResolver{r.ctx, &serviceResolver{r.ctx, r.s}}
}

func (r *serviceResolver) RootOwner() *resourceResolver {
	return &resourceResolver{r.ctx, &serviceResolver{r.ctx, r.s}}
}

func (r *serviceResolver) Selected() []*resourceResolver {
	return getSelectedResources(r)
}

func getSelectedResources(r *serviceResolver) []*resourceResolver {
	ns := (*r).s.Metadata.Namespace
	ls := (*r).s.Spec.Selector
	objs := getAllK8sObjsOfKindInNamespace(r.ctx, PodKind, *ns,
		func(jo JsonObject) bool {
			labels := getMatchLabels(jo)
			for _, label := range ls {
				seenMatch := false
				for k, v := range labels {
					if k == label.Name {
						seenMatch = true
						if v != label.Value {
							return false
						}
					}
				}
				if !seenMatch {
					return false
				}
			}
			return true
		})
	results := make([]*resourceResolver, len(objs))
	for idx, val := range objs {
		results[idx] = &resourceResolver{r.ctx, val}
	}
	return results
}

func getMatchLabels(jo JsonObject) JsonObject {
	kind := getKind(jo)

	switch kind {
	case PodKind:
		return getMetadataField(jo, "labels").(JsonObject)
	}

	return JsonObject{}
}

// Service Spec implementations
func (r serviceSpecResolver) ClusterIP() *string {
	return r.s.ClusterIP
}

func (r serviceSpecResolver) ExternalIPs() *[]string {
	return r.s.ExternalIPs
}

func (r serviceSpecResolver) ExternalName() *string {
	return r.s.ExternalName
}

func (r serviceSpecResolver) ExternalTrafficPolicy() *string {
	return r.s.ExternalTrafficPolicy
}

func (r serviceSpecResolver) HealthCheckNodePort() *int32 {
	return r.s.HealthCheckNodePort
}

func (r serviceSpecResolver) LoadBalancerIP() *string {
	return r.s.LoadBalancerIP
}

func (r serviceSpecResolver) LoadBalancerSourceRanges() *[]string {
	return r.s.LoadBalancerSourceRanges
}

func (r serviceSpecResolver) Ports() *[]servicePortResolver {
	s := r.s.Ports
	if s == nil || len(*s) == 0 {
		res := make([]servicePortResolver, 0)
		return &res
	}
	resolvers := make([]servicePortResolver, len(*s))
	for idx, val := range *s {
		resolvers[idx] = servicePortResolver{r.ctx, val}
	}
	return &resolvers
}

func (r serviceSpecResolver) PublishNotReadyAddresses() *bool {
	return r.s.PublishNotReadyAddresses
}

func (r serviceSpecResolver) Selector() []labelResolver {
	s := r.s.Selector
	if len(s) == 0 {
		res := make([]labelResolver, 0)
		return res
	}
	resolvers := make([]labelResolver, len(s))
	for idx, val := range s {
		l := val
		resolvers[idx] = labelResolver{r.ctx, &l}
	}
	return resolvers
}

func (r serviceSpecResolver) SessionAffinity() *string {
	return r.s.SessionAffinity
}

func (r serviceSpecResolver) SessionAffinityConfig() *sessionAffinityConfigResolver {
	if r.s.SessionAffinityConfig == nil {
		return nil
	}
	return &sessionAffinityConfigResolver{r.ctx, *r.s.SessionAffinityConfig}
}

func (r serviceSpecResolver) Type() string {
	return r.s.Type
}

// Service port implementations
func (r servicePortResolver) Name() *string {
	return r.s.Name
}

func (r servicePortResolver) NodePort() *int32 {
	return r.s.NodePort
}

func (r servicePortResolver) Port() int32 {
	return r.s.Port
}

func (r servicePortResolver) Protocol() *string {
	return r.s.Protocol
}

func (r servicePortResolver) TargetPortString() *string {
	return r.s.TargetPortString
}

func (r servicePortResolver) TargetPortInt() *int32 {
	return r.s.TargetPortInt
}

// SessionAffinityConfig implementations
func (r *sessionAffinityConfigResolver) ClientIP() clientIPConfigResolver {
	return clientIPConfigResolver{r.ctx, r.s.ClientIP}
}

// ClientIPConfig implementations
func (r clientIPConfigResolver) TimeoutSeconds() int32 {
	return r.s.TimeoutSeconds
}
