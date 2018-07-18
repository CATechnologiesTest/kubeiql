package main

import (
	"context"
	"encoding/json"
	//	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

func getK8sResource(kind, namespace, name string) map[string]interface{} {
	return fromJson(
		lookUpResource(kind, namespace, name)).(map[string]interface{})
}

func fromJson(val []byte) interface{} {
	var result interface{}

	if err := json.Unmarshal(val, &result); err != nil {
		panic(err)
	}

	return result
}

func lookUpResource(kind, namespace, name string) []byte {
	cmd := exec.Command("/usr/local/bin/kubectl", "get",
		"-o", "json", "--namespace", namespace, kind, name)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	return bytes
}

func getAllK8sObjsOfKindInNamespace(
	ctx context.Context,
	kind, ns string,
	test func(map[string]interface{}) bool) []resource {
	cmd := exec.Command("/usr/local/bin/kubectl", "get",
		"-o", "json", "--namespace", ns, kind)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	var results []resource
	arr := (fromJson(bytes).(map[string]interface{}))["items"].([]interface{})
	for _, res := range arr {
		if test(res.(map[string]interface{})) {
			results =
				append(results, mapToResource(ctx, res.(map[string]interface{})))
		}
	}
	if results == nil {
		results = make([]resource, 0)
	}
	return results
}
