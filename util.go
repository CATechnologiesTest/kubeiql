package main

import (
//	"context"
)

func getOwners(resourceMap map[string]interface{}) (*resource, *resource) {
	return nil, nil
}

func mapItem(obj map[string]interface{}, item string) map[string]interface{} {
	return obj[item].(map[string]interface{})
}
