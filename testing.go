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
	"bytes"
	"context"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	graphql "github.com/neelance/graphql-go"
	"strings"
	"testing"
)

type Test struct {
	Schema         *graphql.Schema
	Query          string
	OperationName  string
	Variables      map[string]interface{}
	ExpectedResult string
}

func RunTest(t *testing.T, test *Test) {
	result := test.Schema.Exec(context.Background(), test.Query, test.OperationName, test.Variables)
	if len(result.Errors) != 0 {
		t.Fatal(result.Errors[0])
	}

	var data interface{}
	dd := json.NewDecoder(bytes.NewReader([]byte(result.Data)))
	dd.UseNumber()
	if err := dd.Decode(&data); err != nil {
		t.Fatalf("invalid JSON return value: %s", err)
	}
	got, _ := json.Marshal(data)
	var v interface{}
	d := json.NewDecoder(strings.NewReader(test.ExpectedResult))
	d.UseNumber()
	if err := d.Decode(&v); err != nil {
		t.Fatalf("invalid JSON for ExpectedResult: %s", err)
	}
	want, _ := json.Marshal(v)

	if !bytes.Equal(want, got) {
		spew.Dump(data, v)
		t.Logf("want: %s", want)
		t.Logf("got:  %s", got)
		t.Fail()
	}
}
