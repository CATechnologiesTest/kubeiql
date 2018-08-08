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

// The base Kubernetes component
type pod struct {
	Metadata  metadata
	Spec      podSpec
	Owner     resource
	RootOwner resource
}

type podSpec struct {
	ActiveDeadlineSeconds         *int32
	Affinity                      *affinity
	AutomountServiceAccountToken  bool
	Containers                    []container
	DnsConfig                     *podDNSConfig
	DnsPolicy                     *string
	HostAliases                   *[]hostAlias
	HostIPC                       bool
	HostNetwork                   bool
	HostPID                       bool
	Hostname                      *string
	ImagePullSecrets              *[]localObjectReference
	InitContainers                *[]container
	NodeName                      *string
	NodeSelector                  *nodeSelector
	Priority                      *int32
	PriorityClassName             *string
	ReadinessGates                *[]podReadinessGate
	RestartPolicy                 string
	SchedulerName                 *string
	SecurityContext               *podSecurityContext
	ServiceAccountName            *string // XXX deprecated serviceAccount
	ShareProcessNamespace         bool
	Subdomain                     *string
	TerminationGracePeriodSeconds int32
	Tolerations                   *[]toleration
	Volumes                       *[]volume
}

// TBFIL
type affinity struct {
	NodeAffinity    *nodeAffinity
	PodAffinity     *podAffinity
	PodAntiAffinity *podAntiAffinity
}

type podDNSConfig struct {
	Nameservers *[]string
	Options     *[]podDNSConfigOption
	Searches    *[]string
}

type hostAlias struct {
	Hostnames []string
	Ip        string
}

type localObjectReference struct {
	Name string
}

type nodeSelector struct {
	NodeSelectorTerms []nodeSelectorTerm
}

type nodeSelectorTerm struct {
	MatchExpressions *[]nodeSelectorRequirement
	MatchFields      *[]nodeSelectorRequirement
}

type podReadinessGate struct {
	ConditionType string
}

type podSecurityContext struct {
	FsGroup            *int32
	RunAsGroup         *int32
	RunAsNonRoot       bool
	RunAsUser          *int32
	SeLinuxOptions     *sELinuxOptions
	SupplementalGroups *[]int32
	Sysctls            *[]sysctl
}

type sELinuxOptions struct {
	Level *string
	Role  *string
	Type  *string
	User  *string
}

type toleration struct {
	Effect            *string
	Key               *string
	Operator          *string
	TolerationSeconds *int32
	Value             *string
}

type nodeAffinity struct {
	PreferredDuringSchedulingIgnoredDuringExecution *[]preferredSchedulingTerm
	RequiredDuringSchedulingIgnoredDuringExecution  *nodeSelector
}

type podAffinity struct {
	PreferredDuringSchedulingIgnoredDuringExecution *[]weightedPodAffinityTerm
	RequiredDuringSchedulingIgnoredDuringExecution  *podAffinityTerm
}

type podAntiAffinity struct {
	PreferredDuringSchedulingIgnoredDuringExecution *[]weightedPodAffinityTerm
	RequiredDuringSchedulingIgnoredDuringExecution  *podAffinityTerm
}

type preferredSchedulingTerm struct {
	Preference nodeSelectorTerm
	Weight     int32
}

type nodeSelectorRequirement struct {
	Key      string
	Operator string
	Values   *[]string
}

type weightedPodAffinityTerm struct {
	PodAffinityTerm *podAffinityTerm
	Weight          int32
}

type podAffinityTerm struct {
	LabelSelector *labelSelector
	Namespaces    *[]string
	TopologyKey   string
}

type podDNSConfigOption struct {
	Name  string
	Value string
}

type sysctl struct {
	Name  string
	Value string
}

// Resolvers

type podResolver struct {
	ctx context.Context
	p   pod
}

type podSpecResolver struct {
	ctx context.Context
	p   podSpec
}

type affinityResolver struct {
	ctx context.Context
	a   affinity
}

type podDNSConfigResolver struct {
	ctx context.Context
	p   podDNSConfig
}

type hostAliasResolver struct {
	ctx context.Context
	h   hostAlias
}

type localObjectReferenceResolver struct {
	ctx context.Context
	l   localObjectReference
}

type nodeSelectorResolver struct {
	ctx context.Context
	n   nodeSelector
}

type podReadinessGateResolver struct {
	ctx context.Context
	p   podReadinessGate
}

type podSecurityContextResolver struct {
	ctx context.Context
	p   podSecurityContext
}

type tolerationResolver struct {
	ctx context.Context
	t   toleration
}

type nodeAffinityResolver struct {
	ctx context.Context
	n   nodeAffinity
}

type podAffinityResolver struct {
	ctx context.Context
	p   podAffinity
}

type podAntiAffinityResolver struct {
	ctx context.Context
	p   podAntiAffinity
}

type preferredSchedulingTermResolver struct {
	ctx context.Context
	p   preferredSchedulingTerm
}

type nodeSelectorTermResolver struct {
	ctx context.Context
	n   nodeSelectorTerm
}

type nodeSelectorRequirementResolver struct {
	ctx context.Context
	n   nodeSelectorRequirement
}

type weightedPodAffinityTermResolver struct {
	ctx context.Context
	w   weightedPodAffinityTerm
}

type podAffinityTermResolver struct {
	ctx context.Context
	p   podAffinityTerm
}

type podDNSConfigOptionResolver struct {
	ctx context.Context
	p   podDNSConfigOption
}

type sELinuxOptionsResolver struct {
	ctx context.Context
	s   sELinuxOptions
}

type sysctlResolver struct {
	ctx context.Context
	s   sysctl
}

// Translate unmarshalled json into a metadata object
func mapToPod(ctx context.Context, jsonObj JsonObject) pod {
	placeholder := &ownerRef{ctx, jsonObj, nil}
	owner := placeholder
	rootOwner := placeholder
	meta :=
		mapToMetadata(ctx, getNamespace(jsonObj), mapItem(jsonObj, "metadata"))
	spec := mapToPodSpec(ctx, mapItem(jsonObj, "spec"))
	return pod{Metadata: meta, Spec: spec, Owner: owner, RootOwner: rootOwner}
}

func mapToPodSpec(ctx context.Context, jsonObj JsonObject) podSpec {
	jg := jgetter(jsonObj)
	adsp := jg.intRefItemOr("activeDeadlineSeconds", nil)
	aff := extractAffinity(jg)
	asat := jg.boolItemOr("automountServiceAccountToken", false)
	cont := arrayToContainers(jg)
	dc := extractPodDNSConfig(jg)
	dp := jg.stringRefItemOr("dnsPolicy", nil)
	ha := arrayToHostAliases(jg)
	hipc := jg.boolItemOr("hostIPC", false)
	hnet := jg.boolItemOr("hostNetwork", false)
	hpid := jg.boolItemOr("hostPID", false)
	hname := jg.stringRefItemOr("hostname", nil)
	ips := arrayToLOR(jg)
	ic := arrayToInitContainers(jg)
	nn := jg.stringRefItemOr("nodeName", nil)
	ns := extractNodeSelector(jg)
	p := jg.intRefItemOr("priority", nil)
	pcn := jg.stringRefItemOr("priorityClassName", nil)
	rg := arrayToReadinessGates(jg)
	rp := jg.stringItemOr("restartPolicy", "Always")
	sn := jg.stringRefItemOr("schedulerName", nil)
	sc := extractPodSecurityContext(jg)
	san := jg.stringRefItemOr("serviceAccountName", nil)
	spn := jg.boolItemOr("shareProcessNamespace", false)
	s := jg.stringRefItemOr("subdomain", nil)
	tgps := jg.intItemOr("terminationGracePeriodSeconds", 30)
	tols := arrayToTolerations(jg)
	vols := arrayToVolumes(jg)
	return podSpec{
		adsp,
		aff,
		asat,
		cont,
		dc,
		dp,
		ha,
		hipc,
		hnet,
		hpid,
		hname,
		ips,
		ic,
		nn,
		ns,
		p,
		pcn,
		rg,
		rp,
		sn,
		sc,
		san,
		spn,
		s,
		tgps,
		tols,
		vols}
}

// extraction methods for pod specs
func extractAffinity(jg jsonGetter) *affinity {
	if ref := jg.objItemOr("affinity", nil); ref != nil {
		lg := jgetter(*ref)
		return &affinity{
			extractNodeAffinity(lg),
			extractPodAffinity(lg),
			extractPodAntiAffinity(lg)}
	}
	return nil
}

func extractNodeAffinity(jg jsonGetter) *nodeAffinity {
	if ref := jg.objItemOr("nodeAffinity", nil); ref != nil {
		lg := jgetter(*ref)
		return &nodeAffinity{
			arrayToPreferredDuringScheduling(lg),
			extractRequiredDuringScheduling(lg)}
	}
	return nil
}

func extractPodAffinity(jg jsonGetter) *podAffinity {
	if ref := jg.objItemOr("podAffinity", nil); ref != nil {
		lg := jgetter(*ref)
		return &podAffinity{
			arrayToPodPreferredDuringScheduling(lg),
			extractPodRequiredDuringScheduling(lg)}
	}
	return nil
}

func extractPodAntiAffinity(jg jsonGetter) *podAntiAffinity {
	if ref := jg.objItemOr("podAntiAffinity", nil); ref != nil {
		lg := jgetter(*ref)
		return &podAntiAffinity{
			arrayToPodPreferredDuringScheduling(lg),
			extractPodRequiredDuringScheduling(lg)}
	}
	return nil
}

func extractNodeSelector(jg jsonGetter) *nodeSelector {
	if ref := jg.objItemOr("nodeSelector", nil); ref != nil {
		lg := jgetter(*ref)
		return &nodeSelector{arrayToNodeSelectorTerms(lg)}
	}
	return nil
}

func extractSELinuxOptions(jg jsonGetter) *sELinuxOptions {
	if ref := jg.objItemOr("seLinuxOptions", nil); ref != nil {
		lg := jgetter(*ref)
		return &sELinuxOptions{
			lg.stringRefItemOr("level", nil),
			lg.stringRefItemOr("role", nil),
			lg.stringRefItemOr("type", nil),
			lg.stringRefItemOr("user", nil)}
	}
	return nil
}

func extractPodDNSConfig(jg jsonGetter) *podDNSConfig {
	if ref := jg.objItemOr("podDNSConfig", nil); ref != nil {
		lg := jgetter(*ref)
		return &podDNSConfig{
			toStringArrayRef(lg.arrayItemOr("nameservers", nil)),
			arrayToPodDNSConfigOptions(lg),
			toStringArrayRef(lg.arrayItemOr("searches", nil))}
	}
	return nil
}

func extractPodSecurityContext(jg jsonGetter) *podSecurityContext {
	if ref := jg.objItemOr("podSecurityContext", nil); ref != nil {
		lg := jgetter(*ref)
		return &podSecurityContext{
			lg.intRefItemOr("fsGroup", nil),
			lg.intRefItemOr("runAsGroup", nil),
			lg.boolItemOr("runAsNonRoot", false),
			lg.intRefItemOr("runAsUser", nil),
			extractSELinuxOptions(lg),
			toIntArrayRef(lg.arrayItemOr("supplementalGroups", nil)),
			arrayToSysctls(lg)}
	}
	return nil
}

func arrayToNodeSelectorTerms(jg jsonGetter) []nodeSelectorTerm {
	jterms := jg.arrayItem("nodeSelectorTerms")
	terms := make([]nodeSelectorTerm, len(jterms))
	for idx, val := range jterms {
		lg := jgetter(val.(JsonObject))
		terms[idx] = nodeSelectorTerm{
			arrayToMatches(lg, "matchExpressions"),
			arrayToMatches(lg, "matchFields")}
	}
	return terms
}

func arrayToMatches(jg jsonGetter, field string) *[]nodeSelectorRequirement {
	if jmatches := jg.arrayItemOr(field, nil); jmatches != nil {
		matches := make([]nodeSelectorRequirement, len(*jmatches))
		for idx, val := range *jmatches {
			lg := jgetter(val.(JsonObject))
			matches[idx] = nodeSelectorRequirement{
				lg.stringItem("key"),
				lg.stringItem("operator"),
				toStringArrayRef(lg.arrayItemOr("values", nil))}
		}
		return &matches
	}
	return nil
}

func arrayToPodDNSConfigOptions(jg jsonGetter) *[]podDNSConfigOption {
	if jopts := jg.arrayItemOr("options", nil); jopts != nil {
		opts := make([]podDNSConfigOption, len(*jopts))
		for idx, val := range *jopts {
			lg := jgetter(val.(JsonObject))
			opts[idx] = podDNSConfigOption{
				lg.stringItem("name"),
				lg.stringItem("value")}
		}
		return &opts
	}
	return nil
}

func arrayToPreferredDuringScheduling(jg jsonGetter) *[]preferredSchedulingTerm {
	if jprefs := jg.arrayItemOr(
		"preferredDuringSchedulingIgnoredDuringExecution", nil); jprefs != nil {
		prefs := make([]preferredSchedulingTerm, len(*jprefs))
		for idx, val := range *jprefs {
			lg := jgetter(val.(JsonObject))
			prefs[idx] = preferredSchedulingTerm{
				extractPreference(lg),
				lg.intItem("weight")}
		}
		return &prefs
	}
	return nil
}

func extractPreference(jg jsonGetter) nodeSelectorTerm {
	ns := jg.objItem("preference")
	lg := jgetter(ns)
	return nodeSelectorTerm{
		arrayToMatches(lg, "matchExpressions"),
		arrayToMatches(lg, "matchFields")}
}

func arrayToPodPreferredDuringScheduling(jg jsonGetter) *[]weightedPodAffinityTerm {
	if jprefs := jg.arrayItemOr(
		"preferredDuringSchedulingIgnoredDuringExecution", nil); jprefs != nil {
		prefs := make([]weightedPodAffinityTerm, len(*jprefs))
		for idx, val := range *jprefs {
			lg := jgetter(val.(JsonObject))
			prefs[idx] = weightedPodAffinityTerm{
				extractPodAffinityTerm(lg, "podAffinity"),
				lg.intItem("weight")}
		}
		return &prefs
	}
	return nil
}

func extractRequiredDuringScheduling(jg jsonGetter) *nodeSelector {
	if ref := jg.objItemOr("requiredDuringSchedulingIgnoredDuringExecution",
		nil); ref != nil {
		lg := jgetter(*ref)
		return &nodeSelector{arrayToNodeSelectorTerms(lg)}
	}
	return nil
}

func extractPodAffinityTerm(jg jsonGetter, field string) *podAffinityTerm {
	if ref := jg.objItemOr(field, nil); ref != nil {
		lg := jgetter(ref)
		return &podAffinityTerm{
			mapToSelector(ref),
			toStringArrayRef(lg.arrayItemOr("namespaces", nil)),
			lg.stringItem("topologyKey")}
	}
	return nil
}

func extractPodRequiredDuringScheduling(jg jsonGetter) *podAffinityTerm {
	return extractPodAffinityTerm(
		jg,
		"requiredDuringSchedulingIgnoredDuringExecution")
}

func arrayToContainers(jg jsonGetter) []container {
	jconts := jg.arrayItem("containers")
	conts := make([]container, len(jconts))
	for idx, _ := range jconts {
		conts[idx] = container{} // XXX
	}
	return conts
}

func arrayToInitContainers(jg jsonGetter) *[]container {
	if jconts := jg.arrayItemOr("initContainers", nil); jconts != nil {
		conts := make([]container, len(*jconts))
		for idx, _ := range *jconts {
			conts[idx] = container{} // XXX
		}
		return &conts
	}
	return nil
}

func arrayToReadinessGates(jg jsonGetter) *[]podReadinessGate {
	if jgates := jg.arrayItemOr("readinessGates", nil); jgates != nil {
		gates := make([]podReadinessGate, len(*jgates))
		for idx, val := range *jgates {
			mval := val.(JsonObject)
			if gval, ok := mval["podReadinessGate"]; ok {
				mgval := gval.(JsonObject)
				gates[idx] = podReadinessGate{mgval["conditionType"].(string)}
			}
		}
		return &gates
	}
	return nil
}

func arrayToTolerations(jg jsonGetter) *[]toleration {
	if jtols := jg.arrayItemOr("tolerations", nil); jtols != nil {
		tols := make([]toleration, len(*jtols))
		for idx, val := range *jtols {
			mval := val.(JsonObject)
			lg := jgetter(mval)
			tols[idx] = toleration{
				lg.stringRefItemOr("effect", nil),
				lg.stringRefItemOr("key", nil),
				lg.stringRefItemOr("operator", nil),
				lg.intRefItemOr("tolerationSeconds", nil),
				lg.stringRefItemOr("value", nil)}
		}
		return &tols
	}
	return nil
}

func arrayToSysctls(jg jsonGetter) *[]sysctl {
	if jctls := jg.arrayItemOr("sysctls", nil); jctls != nil {
		ctls := make([]sysctl, len(*jctls))
		for idx, val := range *jctls {
			mval := val.(JsonObject)
			if gval, ok := mval["sysctl"]; ok {
				mgval := gval.(JsonObject)
				lg := jgetter(mgval)
				ctls[idx] = sysctl{lg.stringItem("name"), lg.stringItem("value")}
			}
		}
		return &ctls
	}
	return nil
}

func arrayToVolumes(jg jsonGetter) *[]volume {
	if jvols := jg.arrayItemOr("volumes", nil); jvols != nil {
		vols := make([]volume, len(*jvols))
		for idx, val := range *jvols {
			mval := val.(JsonObject)
			lg := jgetter(mval)
			vols[idx] = volume{
				extractConfigMapVolumeSource(lg),
				lg.stringItem("name"),
				extractHostPathVolumeSource(lg),
				extractPersistentVolumeClaimVolumeSource(lg),
				extractSecretVolumeSource(lg)}
		}
		return &vols
	}
	return nil
}

func extractConfigMapVolumeSource(jg jsonGetter) *configMapVolumeSource {
	if cm := jg.objItemOr("configMap", nil); cm != nil {
		lg := jgetter(*cm)
		return &configMapVolumeSource{
			lg.intRefItemOr("defaultMode", nil),
			arrayToItems(lg),
			lg.stringItem("name"),
			lg.boolItemOr("optional", false)}
	}
	return nil
}

func extractHostPathVolumeSource(jg jsonGetter) *hostPathVolumeSource {
	if hp := jg.objItemOr("hostPath", nil); hp != nil {
		lg := jgetter(*hp)
		return &hostPathVolumeSource{
			lg.stringItem("path"),
			lg.stringRefItemOr("type", nil)}
	}
	return nil
}

func extractPersistentVolumeClaimVolumeSource(jg jsonGetter) *persistentVolumeClaimVolumeSource {
	if pv := jg.objItemOr("persistentVolumeClaim", nil); pv != nil {
		lg := jgetter(*pv)
		return &persistentVolumeClaimVolumeSource{
			lg.stringItem("claimName"),
			lg.boolItemOr("type", false)}
	}
	return nil
}

func extractSecretVolumeSource(jg jsonGetter) *secretVolumeSource {
	if s := jg.objItemOr("secret", nil); s != nil {
		lg := jgetter(*s)
		return &secretVolumeSource{
			lg.intRefItemOr("defaultMode", nil),
			arrayToItems(jg),
			lg.boolItemOr("optional", false),
			lg.stringItem("secretName")}
	}
	return nil
}

func arrayToItems(jg jsonGetter) *[]keyToPath {
	if jitems := jg.arrayItemOr("items", nil); jitems != nil {
		items := make([]keyToPath, len(*jitems))
		for idx, val := range *jitems {
			lg := jgetter(val.(JsonObject))
			items[idx] = keyToPath{
				lg.stringItem("key"),
				lg.intRefItemOr("mode", nil),
				lg.stringItem("path")}
		}
		return &items
	}
	return nil
}

func arrayToHostAliases(jg jsonGetter) *[]hostAlias {
	if jaliases := jg.arrayItemOr("hostAliases", nil); jaliases != nil {
		aliases := make([]hostAlias, len(*jaliases))
		for idx, val := range *jaliases {
			lg := jgetter(val.(JsonObject))
			aliases[idx] = hostAlias{
				toStringArray(lg.arrayItem("hostnames")),
				lg.stringItem("ip")}

		}
		return &aliases
	}
	return nil
}

func arrayToLOR(jg jsonGetter) *[]localObjectReference {
	if jips := jg.arrayItemOr("imagePullSecrets", nil); jips != nil {
		ips := make([]localObjectReference, len(*jips))
		for idx, val := range *jips {
			lg := jgetter(val.(JsonObject))
			ips[idx] = localObjectReference{lg.stringItem("name")}
		}
		return &ips
	}
	return nil
}

// Affinity method implementations
func (r *affinityResolver) NodeAffinity() *nodeAffinityResolver {
	if r.a.NodeAffinity == nil {
		return nil
	}
	return &nodeAffinityResolver{r.ctx, *r.a.NodeAffinity}
}

func (r *affinityResolver) PodAffinity() *podAffinityResolver {
	if r.a.PodAffinity == nil {
		return nil
	}
	return &podAffinityResolver{r.ctx, *r.a.PodAffinity}
}

func (r *affinityResolver) PodAntiAffinity() *podAntiAffinityResolver {
	if r.a.PodAntiAffinity == nil {
		return nil
	}
	return &podAntiAffinityResolver{r.ctx, *r.a.PodAntiAffinity}
}

// NodeAffinity method implementations
func (r *nodeAffinityResolver) PreferredDuringSchedulingIgnoredDuringExecution() *[]preferredSchedulingTermResolver {
	n := r.n.PreferredDuringSchedulingIgnoredDuringExecution
	if n == nil || len(*n) == 0 {
		res := make([]preferredSchedulingTermResolver, 0)
		return &res
	}
	resolvers := make([]preferredSchedulingTermResolver, len(*n))
	for idx, val := range *n {
		resolvers[idx] = preferredSchedulingTermResolver{r.ctx, val}
	}
	return &resolvers
}

func (r *nodeAffinityResolver) RequiredDuringSchedulingIgnoredDuringExecution() *nodeSelectorResolver {
	if r.n.RequiredDuringSchedulingIgnoredDuringExecution == nil {
		return nil
	}
	return &nodeSelectorResolver{
		r.ctx, *r.n.RequiredDuringSchedulingIgnoredDuringExecution}
}

// PreferredSchedulingTerm method implementations
func (r preferredSchedulingTermResolver) Preference() nodeSelectorTermResolver {
	return nodeSelectorTermResolver{r.ctx, r.p.Preference}
}

func (r preferredSchedulingTermResolver) Weight() int32 {
	return r.p.Weight
}

// NodeSelectorTerm method implementations
func (r nodeSelectorTermResolver) MatchExpressions() *[]nodeSelectorRequirementResolver {
	n := r.n.MatchExpressions
	if n == nil || len(*n) == 0 {
		res := make([]nodeSelectorRequirementResolver, 0)
		return &res
	}
	resolvers := make([]nodeSelectorRequirementResolver, len(*n))
	for idx, val := range *n {
		resolvers[idx] = nodeSelectorRequirementResolver{r.ctx, val}
	}
	return &resolvers
}

func (r nodeSelectorTermResolver) MatchFields() *[]nodeSelectorRequirementResolver {
	n := r.n.MatchFields
	if n == nil || len(*n) == 0 {
		res := make([]nodeSelectorRequirementResolver, 0)
		return &res
	}
	resolvers := make([]nodeSelectorRequirementResolver, len(*n))
	for idx, val := range *n {
		resolvers[idx] = nodeSelectorRequirementResolver{r.ctx, val}
	}
	return &resolvers
}

// NodeSelectorRequirement method implementations
func (r nodeSelectorRequirementResolver) Key() string {
	return r.n.Key
}

func (r nodeSelectorRequirementResolver) Operator() string {
	return r.n.Operator
}

func (r nodeSelectorRequirementResolver) Values() *[]string {
	return r.n.Values
}

// NodeSelector method implementations
func (r *nodeSelectorResolver) NodeSelectorTerms() []nodeSelectorTermResolver {
	n := r.n.NodeSelectorTerms
	if len(n) == 0 {
		res := make([]nodeSelectorTermResolver, 0)
		return res
	}
	resolvers := make([]nodeSelectorTermResolver, len(n))
	for idx, val := range n {
		resolvers[idx] = nodeSelectorTermResolver{r.ctx, val}
	}
	return resolvers
}

// PodAffinity method implementations
func (r *podAffinityResolver) PreferredDuringSchedulingIgnoredDuringExecution() *[]weightedPodAffinityTermResolver {
	p := r.p.PreferredDuringSchedulingIgnoredDuringExecution
	if p == nil || len(*p) == 0 {
		res := make([]weightedPodAffinityTermResolver, 0)
		return &res
	}
	resolvers := make([]weightedPodAffinityTermResolver, len(*p))
	for idx, val := range *p {
		resolvers[idx] = weightedPodAffinityTermResolver{r.ctx, val}
	}
	return &resolvers
}

func (r *podAffinityResolver) RequiredDuringSchedulingIgnoredDuringExecution() *podAffinityTermResolver {
	p := r.p.RequiredDuringSchedulingIgnoredDuringExecution
	if p == nil {
		return nil
	}
	return &podAffinityTermResolver{r.ctx, *p}
}

// PodAntiAffinity method implementations
func (r *podAntiAffinityResolver) PreferredDuringSchedulingIgnoredDuringExecution() *[]weightedPodAffinityTermResolver {
	p := r.p.PreferredDuringSchedulingIgnoredDuringExecution
	if p == nil || len(*p) == 0 {
		res := make([]weightedPodAffinityTermResolver, 0)
		return &res
	}
	resolvers := make([]weightedPodAffinityTermResolver, len(*p))
	for idx, val := range *p {
		resolvers[idx] = weightedPodAffinityTermResolver{r.ctx, val}
	}
	return &resolvers
}

func (r *podAntiAffinityResolver) RequiredDuringSchedulingIgnoredDuringExecution() *podAffinityTermResolver {
	p := r.p.RequiredDuringSchedulingIgnoredDuringExecution
	if p == nil {
		return nil
	}
	return &podAffinityTermResolver{r.ctx, *p}
}

// WeightedPodAffinityTerm method implementations
func (r weightedPodAffinityTermResolver) PodAffinityTerm() *podAffinityTermResolver {
	if r.w.PodAffinityTerm == nil {
		return nil
	}
	return &podAffinityTermResolver{r.ctx, *r.w.PodAffinityTerm}
}

func (r weightedPodAffinityTermResolver) Weight() int32 {
	return r.w.Weight
}

// PodAffinityTerm method implementations
func (r *podAffinityTermResolver) LabelSelector() *labelSelectorResolver {
	if r.p.LabelSelector == nil {
		return nil
	}
	return &labelSelectorResolver{r.ctx, *r.p.LabelSelector}
}

func (r *podAffinityTermResolver) Namespaces() *[]string {
	return r.p.Namespaces
}

func (r *podAffinityTermResolver) TopologyKey() string {
	return r.p.TopologyKey
}

// PodDNSConfig method implementations
func (r *podDNSConfigResolver) Nameservers() *[]string {
	return r.p.Nameservers
}

func (r *podDNSConfigResolver) Options() *[]podDNSConfigOptionResolver {
	p := r.p.Options
	if p == nil || len(*p) == 0 {
		res := make([]podDNSConfigOptionResolver, 0)
		return &res
	}
	resolvers := make([]podDNSConfigOptionResolver, len(*p))
	for idx, val := range *p {
		resolvers[idx] = podDNSConfigOptionResolver{r.ctx, val}
	}
	return &resolvers
}

func (r *podDNSConfigResolver) Searches() *[]string {
	return r.p.Searches
}

// PodDNSConfigOption method implementations
func (r podDNSConfigOptionResolver) Name() string {
	return r.p.Name
}

func (r podDNSConfigOptionResolver) Value() string {
	return r.p.Value
}

// HostAlias method implementations
func (r hostAliasResolver) Hostnames() []string {
	return r.h.Hostnames
}

func (r hostAliasResolver) Ip() string {
	return r.h.Ip
}

// PodSpec method implementations
func (r podSpecResolver) ActiveDeadlineSeconds() *int32 {
	return r.p.ActiveDeadlineSeconds
}

func (r podSpecResolver) Affinity() *affinityResolver {
	if r.p.Affinity == nil {
		return nil
	}
	return &affinityResolver{r.ctx, *r.p.Affinity}
}

// LocalObjectReference method implementations
func (r localObjectReferenceResolver) Name() string {
	return r.l.Name
}

// PodReadinessGate method implementations
func (r podReadinessGateResolver) ConditionType() string {
	return r.p.ConditionType
}

// PodSecurityContext method implementations
func (r *podSecurityContextResolver) FsGroup() *int32 {
	return r.p.FsGroup
}

func (r *podSecurityContextResolver) RunAsGroup() *int32 {
	return r.p.RunAsGroup
}

func (r *podSecurityContextResolver) RunAsNonRoot() bool {
	return r.p.RunAsNonRoot
}

func (r *podSecurityContextResolver) RunAsUser() *int32 {
	return r.p.RunAsUser
}

func (r *podSecurityContextResolver) SeLinuxOptions() *sELinuxOptionsResolver {
	return &sELinuxOptionsResolver{r.ctx, *r.p.SeLinuxOptions}
}

func (r *podSecurityContextResolver) SupplementalGroups() *[]int32 {
	return r.p.SupplementalGroups
}

func (r *podSecurityContextResolver) Sysctls() *[]sysctlResolver {
	p := r.p.Sysctls
	if p == nil || len(*p) == 0 {
		res := make([]sysctlResolver, 0)
		return &res
	}
	resolvers := make([]sysctlResolver, len(*p))
	for idx, val := range *p {
		resolvers[idx] = sysctlResolver{r.ctx, val}
	}
	return &resolvers
}

// SELinuxOptions method implementations
func (r *sELinuxOptionsResolver) Level() *string {
	return r.s.Level
}

func (r *sELinuxOptionsResolver) Role() *string {
	return r.s.Role
}

func (r *sELinuxOptionsResolver) Type() *string {
	return r.s.Type
}

func (r *sELinuxOptionsResolver) User() *string {
	return r.s.User
}

// Sysctl method implementations
func (r sysctlResolver) Name() string {
	return r.s.Name
}

func (r sysctlResolver) Value() string {
	return r.s.Value
}

// Toleration method implementations
func (r tolerationResolver) Effect() *string {
	return r.t.Effect
}

func (r tolerationResolver) Key() *string {
	return r.t.Key
}

func (r tolerationResolver) Operator() *string {
	return r.t.Operator
}

func (r tolerationResolver) TolerationSeconds() *int32 {
	return r.t.TolerationSeconds
}

func (r tolerationResolver) Value() *string {
	return r.t.Value
}

// PodSpec method implementations
func (r podSpecResolver) AutomountServiceAccountToken() bool {
	return r.p.AutomountServiceAccountToken
}

func (r podSpecResolver) Containers() []containerResolver {
	cont := r.p.Containers
	resolvers := make([]containerResolver, len(cont))
	for idx, val := range cont {
		resolvers[idx] = containerResolver{r.ctx, val}
	}
	return resolvers
}

func (r podSpecResolver) DnsConfig() *podDNSConfigResolver {
	if r.p.DnsConfig == nil {
		return nil
	}
	return &podDNSConfigResolver{r.ctx, *r.p.DnsConfig}
}

func (r podSpecResolver) DnsPolicy() *string {
	return r.p.DnsPolicy
}

func (r podSpecResolver) HostAliases() *[]hostAliasResolver {
	ha := r.p.HostAliases
	if ha == nil || len(*ha) == 0 {
		res := make([]hostAliasResolver, 0)
		return &res
	}
	resolvers := make([]hostAliasResolver, len(*ha))
	for idx, val := range *ha {
		resolvers[idx] = hostAliasResolver{r.ctx, val}
	}
	return &resolvers
}

func (r podSpecResolver) HostIPC() bool {
	return r.p.HostIPC
}

func (r podSpecResolver) HostNetwork() bool {
	return r.p.HostNetwork
}

func (r podSpecResolver) HostPID() bool {
	return r.p.HostPID
}

func (r podSpecResolver) Hostname() *string {
	return r.p.Hostname
}

func (r podSpecResolver) ImagePullSecrets() *[]localObjectReferenceResolver {
	ips := r.p.ImagePullSecrets
	if ips == nil || len(*ips) == 0 {
		res := make([]localObjectReferenceResolver, 0)
		return &res
	}
	resolvers := make([]localObjectReferenceResolver, len(*ips))
	for idx, val := range *ips {
		resolvers[idx] = localObjectReferenceResolver{r.ctx, val}
	}
	return &resolvers
}

func (r podSpecResolver) InitContainers() *[]containerResolver {
	cont := r.p.InitContainers
	if cont == nil || len(*cont) == 0 {
		res := make([]containerResolver, 0)
		return &res
	}
	resolvers := make([]containerResolver, len(*cont))
	for idx, val := range *cont {
		resolvers[idx] = containerResolver{r.ctx, val}
	}
	return &resolvers
}

func (r podSpecResolver) NodeName() *string {
	return r.p.NodeName
}

func (r podSpecResolver) NodeSelector() *nodeSelectorResolver {
	if r.p.NodeSelector == nil {
		return nil
	}
	return &nodeSelectorResolver{r.ctx, *r.p.NodeSelector}
}

func (r podSpecResolver) Priority() *int32 {
	return r.p.Priority
}

func (r podSpecResolver) PriorityClassName() *string {
	return r.p.PriorityClassName
}

func (r podSpecResolver) ReadinessGates() *[]podReadinessGateResolver {
	prg := r.p.ReadinessGates
	if prg == nil || len(*prg) == 0 {
		res := make([]podReadinessGateResolver, 0)
		return &res
	}
	resolvers := make([]podReadinessGateResolver, len(*prg))
	for idx, val := range *prg {
		resolvers[idx] = podReadinessGateResolver{r.ctx, val}
	}
	return &resolvers
}

func (r podSpecResolver) RestartPolicy() string {
	return r.p.RestartPolicy
}

func (r podSpecResolver) SchedulerName() *string {
	return r.p.SchedulerName
}

func (r podSpecResolver) SecurityContext() *podSecurityContextResolver {
	if r.p.SecurityContext == nil {
		return nil
	}

	return &podSecurityContextResolver{r.ctx, *r.p.SecurityContext}
}

func (r podSpecResolver) ServiceAccountName() *string {
	return r.p.ServiceAccountName
}

func (r podSpecResolver) ShareProcessNamespace() bool {
	return r.p.ShareProcessNamespace
}

func (r podSpecResolver) Subdomain() *string {
	return r.p.Subdomain
}

func (r podSpecResolver) TerminationGracePeriodSeconds() int32 {
	return r.p.TerminationGracePeriodSeconds
}

func (r podSpecResolver) Tolerations() *[]tolerationResolver {
	t := r.p.Tolerations
	if t == nil || len(*t) == 0 {
		res := make([]tolerationResolver, 0)
		return &res
	}
	resolvers := make([]tolerationResolver, len(*t))
	for idx, val := range *t {
		resolvers[idx] = tolerationResolver{r.ctx, val}
	}
	return &resolvers
}

func (r podSpecResolver) Volumes() *[]volumeResolver {
	v := r.p.Volumes
	if v == nil || len(*v) == 0 {
		res := make([]volumeResolver, 0)
		return &res
	}
	resolvers := make([]volumeResolver, len(*v))
	for idx, val := range *v {
		resolvers[idx] = volumeResolver{r.ctx, val}
	}
	return &resolvers
}

// Resource method implementations
func (r *podResolver) Kind() string {
	return PodKind
}

func (r *podResolver) Metadata() metadataResolver {
	meta := r.p.Metadata
	return metadataResolver{r.ctx, meta}
}

func (r *podResolver) Spec() *podSpecResolver {
	spec := r.p.Spec
	return &podSpecResolver{r.ctx, spec}
}

func (r *podResolver) Owner() *resourceResolver {
	if oref, ok := r.p.Owner.(*ownerRef); ok {
		r.p.Owner = getOwner(oref.ctx, oref.ref)
	}
	return &resourceResolver{r.ctx, r.p.Owner}
}

func (r *podResolver) RootOwner() *resourceResolver {
	if oref, ok := r.p.RootOwner.(*ownerRef); ok {
		r.p.RootOwner = getRootOwner(oref.ctx, oref.ref)
	}
	return &resourceResolver{r.ctx, r.p.RootOwner}
}
