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
	"errors"
	"fmt"
)

// Utility methods for getting data out of nested maps

func mapItem(obj JsonObject, item string) JsonObject {
	if itemobj, ok := obj[item].(JsonObject); ok {
		return itemobj
	}
	return nil
}

func mapItemRef(obj JsonObject, item string) *JsonObject {
	if mitem, ok := obj[item].(JsonObject); ok {
		return &mitem
	}

	return nil
}

func getKind(resourceMap JsonObject) string {
	kind := resourceMap["kind"]
	if kindstr, ok := kind.(string); ok {
		return kindstr
	}

	panic(errors.New(
		fmt.Sprintf("Invalid Kubernetes resource: %v", resourceMap)))
}

func getNamespace(resourceMap JsonObject) string {
	namespace := getMetadataField(resourceMap, "namespace")
	if nsstr, ok := namespace.(string); ok {
		return nsstr
	}

	panic(errors.New(
		fmt.Sprintf("Invalid Kubernetes resource: %v", resourceMap)))
}

func getName(resourceMap JsonObject) string {
	name := getMetadataField(resourceMap, "name")
	if nsstr, ok := name.(string); ok {
		return nsstr
	}

	panic(errors.New(
		fmt.Sprintf("Invalid Kubernetes resource: %v", resourceMap)))
}

func getUid(resourceMap JsonObject) string {
	uid := getMetadataField(resourceMap, "uid")
	if uidstr, ok := uid.(string); ok {
		return uidstr
	}

	panic(errors.New(
		fmt.Sprintf("Invalid Kubernetes resource: %v", resourceMap)))
}

func getMetadataField(
	resourceMap JsonObject,
	field string) interface{} {
	if meta, ok := resourceMap["metadata"]; ok {
		if mmap, ok := meta.(JsonObject); ok {
			if val, ok := mmap[field]; ok {
				return val
			}
		}
	}
	return nil
}

func toStringArray(sa JsonArray) []string {
	strs := make([]string, len(sa))
	for idx, val := range sa {
		strs[idx] = val.(string)
	}
	return strs
}

func toStringArrayRef(sa *JsonArray) *[]string {
	if sa == nil {
		return nil
	}
	strs := make([]string, len(*sa))
	for idx, val := range *sa {
		strs[idx] = val.(string)
	}
	return &strs
}

func toIntArrayRef(ia *JsonArray) *[]int32 {
	if ia == nil {
		return nil
	}
	ints := make([]int32, len(*ia))
	for idx, val := range *ia {
		ints[idx] = toGQLInt(val)
	}
	return &ints
}

func toGQLInt(val interface{}) int32 {
	if num, ok := val.(float64); ok {
		return int32(num)
	} else {
		return val.(int32)
	}
}

type jsonGetter struct {
	jsonObj JsonObject
}

func (jg jsonGetter) boolItem(field string) bool {
	return jg.jsonObj[field].(bool)
}

func (jg jsonGetter) intItem(field string) int32 {
	return toGQLInt(jg.jsonObj[field])
}

func (jg jsonGetter) stringItem(field string) string {
	return jg.jsonObj[field].(string)
}

func (jg jsonGetter) arrayItem(field string) JsonArray {
	return jg.jsonObj[field].(JsonArray)
}

func (jg jsonGetter) objItem(field string) JsonObject {
	return jg.jsonObj[field].(JsonObject)
}

func (jg jsonGetter) boolItemOr(field string, defVal bool) bool {
	if val := jg.jsonObj[field]; val != nil {
		return val.(bool)
	}
	return defVal
}

func (jg jsonGetter) intRefItemOr(field string, defVal *int32) *int32 {
	if val := jg.jsonObj[field]; val != nil {
		intval := toGQLInt(val)
		return &intval
	}
	return defVal
}

func (jg jsonGetter) intItemOr(field string, defVal int32) int32 {
	if val := jg.jsonObj[field]; val != nil {
		if num, ok := val.(float64); ok {
			return int32(num)
		} else {
			return val.(int32)
		}
	}
	return defVal
}

func (jg jsonGetter) stringRefItemOr(field string, defVal *string) *string {
	if val := jg.jsonObj[field]; val != nil {
		strVal := val.(string)
		return &strVal
	}
	return defVal
}

func (jg jsonGetter) stringItemOr(field string, defVal string) string {
	if val := jg.jsonObj[field]; val != nil {
		return val.(string)
	}
	return defVal
}

func (jg jsonGetter) arrayItemOr(field string, defVal *JsonArray) *JsonArray {
	if val, ok := jg.jsonObj[field].(JsonArray); ok {
		arrVal := JsonArray(val)
		return &arrVal
	}
	return defVal
}

func (jg jsonGetter) objItemOr(field string, defVal *JsonObject) *JsonObject {
	if val, ok := jg.jsonObj[field].(JsonObject); ok {
		objVal := JsonObject(val)
		return &objVal
	}
	return defVal
}

func jgetter(jsonObj interface{}) jsonGetter {
	if mval, ok := jsonObj.(JsonObject); ok {
		return jsonGetter{mval}
	}
	panic(errors.New(
		fmt.Sprintf("Invalid Kubernetes object: %v", jsonObj)))
}
