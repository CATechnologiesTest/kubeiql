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

// Single label value within a Kubernetes object
type label struct {
	Name  string
	Value string
}

type labelResolver struct {
	ctx context.Context
	l   *label
}

// Translate unmarshalled json into a set of labels
func mapToLabels(lMap map[string]interface{}) *[]label {
	var labels []label

	for k, v := range lMap {
		labels = append(labels, label{k, v.(string)})
	}

	return &labels
}

func (r labelResolver) Name() string {
	return r.l.Name
}

func (r labelResolver) Value() string {
	return r.l.Value
}
