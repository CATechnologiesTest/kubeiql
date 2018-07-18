package main

import (
	"context"
)

type label struct {
	Name  string
	Value string
}

type labelResolver struct {
	ctx context.Context
	l   *label
}

func mapToLabels(lMap map[string]interface{}) *[]label {
	var labels []label

	for k, v := range lMap {
		labels = append(labels, label{k, v.(string)})
	}

	return &labels
}

func (r *labelResolver) Name() string {
	return r.l.Name
}

func (r *labelResolver) Value() string {
	return r.l.Value
}
