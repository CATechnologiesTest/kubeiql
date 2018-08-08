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

// Volumes hold the actual data in Kubernetes

type volume struct {
	// XXX incomplete
	ConfigMap             *configMapVolumeSource
	Name                  string
	HostPath              *hostPathVolumeSource
	PersistentVolumeClaim *persistentVolumeClaimVolumeSource
	Secret                *secretVolumeSource
}

type configMapVolumeSource struct {
	DefaultMode *int32
	Items       *[]keyToPath
	Name        string
	Optional    bool
}

type hostPathVolumeSource struct {
	Path string
	Type *string
}

type persistentVolumeClaimVolumeSource struct {
	ClaimName string
	ReadOnly  bool
}

type secretVolumeSource struct {
	DefaultMode *int32
	Items       *[]keyToPath
	Optional    bool
	SecretName  string
}

type keyToPath struct {
	Key  string
	Mode *int32
	Path string
}

type volumeResolver struct {
	ctx context.Context
	v   volume
}

type configMapVolumeSourceResolver struct {
	ctx context.Context
	c   configMapVolumeSource
}

type hostPathVolumeSourceResolver struct {
	ctx context.Context
	h   hostPathVolumeSource
}

type persistentVolumeClaimVolumeSourceResolver struct {
	ctx context.Context
	p   persistentVolumeClaimVolumeSource
}

type secretVolumeSourceResolver struct {
	ctx context.Context
	s   secretVolumeSource
}

type keyToPathResolver struct {
	ctx context.Context
	k   keyToPath
}

func (r volumeResolver) ConfigMap() *configMapVolumeSourceResolver {
	if r.v.ConfigMap == nil {
		return nil
	}
	return &configMapVolumeSourceResolver{r.ctx, *r.v.ConfigMap}
}

func (r volumeResolver) Name() string {
	return r.v.Name
}

func (r volumeResolver) HostPath() *hostPathVolumeSourceResolver {
	if r.v.HostPath == nil {
		return nil
	}
	return &hostPathVolumeSourceResolver{r.ctx, *r.v.HostPath}
}

func (r volumeResolver) PersistentVolumeClaim() *persistentVolumeClaimVolumeSourceResolver {
	if r.v.PersistentVolumeClaim == nil {
		return nil
	}
	return &persistentVolumeClaimVolumeSourceResolver{r.ctx, *r.v.PersistentVolumeClaim}
}

func (r volumeResolver) Secret() *secretVolumeSourceResolver {
	if r.v.Secret == nil {
		return nil
	}
	return &secretVolumeSourceResolver{r.ctx, *r.v.Secret}
}

// configMapVolumeSourceResolver implementations
func (r *configMapVolumeSourceResolver) DefaultMode() *int32 {
	return r.c.DefaultMode
}

func (r *configMapVolumeSourceResolver) Items() *[]keyToPathResolver {
	c := r.c.Items
	if c == nil || len(*c) == 0 {
		res := make([]keyToPathResolver, 0)
		return &res
	}
	resolvers := make([]keyToPathResolver, len(*c))
	for idx, val := range *c {
		resolvers[idx] = keyToPathResolver{r.ctx, val}
	}
	return &resolvers
}

func (r *configMapVolumeSourceResolver) Name() string {
	return r.c.Name
}

func (r *configMapVolumeSourceResolver) Optional() bool {
	return r.c.Optional
}

// keyToPath implementations
func (r keyToPathResolver) Key() string {
	return r.k.Key
}

func (r keyToPathResolver) Mode() *int32 {
	return r.k.Mode
}

func (r keyToPathResolver) Path() string {
	return r.k.Path
}

// hostPathVolumeSource implementations
func (r *hostPathVolumeSourceResolver) Path() string {
	return r.h.Path
}

func (r *hostPathVolumeSourceResolver) Type() *string {
	return r.h.Type
}

// persistentVolumeClaimVolumeSource implementations
func (r *persistentVolumeClaimVolumeSourceResolver) ClaimName() string {
	return r.p.ClaimName
}

func (r *persistentVolumeClaimVolumeSourceResolver) ReadOnly() bool {
	return r.p.ReadOnly
}

// secretVolumeSource implementations
func (r *secretVolumeSourceResolver) DefaultMode() *int32 {
	return r.s.DefaultMode
}

func (r *secretVolumeSourceResolver) Items() *[]keyToPathResolver {
	s := r.s.Items
	if s == nil || len(*s) == 0 {
		res := make([]keyToPathResolver, 0)
		return &res
	}
	resolvers := make([]keyToPathResolver, len(*s))
	for idx, val := range *s {
		resolvers[idx] = keyToPathResolver{r.ctx, val}
	}
	return &resolvers
}

func (r *secretVolumeSourceResolver) Optional() bool {
	return r.s.Optional
}

func (r *secretVolumeSourceResolver) SecretName() string {
	return r.s.SecretName
}

// Translate unmarshalled json into a deployment object
func mapToVolume(ctx context.Context, _ JsonObject) volume {
	return volume{}
}
