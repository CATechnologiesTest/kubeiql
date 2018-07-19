// Copyright 2018 Yipee.io
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
	"fmt"
	"github.com/gorilla/mux"
	graphql "github.com/neelance/graphql-go"
	"github.com/neelance/graphql-go/relay"
	"log"
	"net/http"
)

var schema *graphql.Schema

type UserData string

func init() {
	var err error
	schema, err = graphql.ParseSchema(Schema, &Resolver{})
	fmt.Println(schema)
	if err != nil {
		panic(err)
	}
}

func main() {
	r := mux.NewRouter()
	handler := &relay.Handler{Schema: schema}
	r.HandleFunc(
		"/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write(page)
		}))
	r.Handle("/query",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cache := make(map[string]interface{})
			handler.ServeHTTP(w,
				r.WithContext(
					context.WithValue(
						r.Context(),
						"queryCache",
						&cache)))
		})).Methods("POST")
	r.HandleFunc("/query",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		})).Methods("OPTIONS")
	log.Fatal(http.ListenAndServe(":8128", r))
}

var page = []byte(`
<!DOCTYPE html>
<html>
    <head>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.7.8/graphiql.css" />
    <script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/1.0.0/fetch.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.3.2/react.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.3.2/react-dom.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.7.8/graphiql.js"></script>
    </head>
    <body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
    <div id="graphiql" style="height: 100vh;">Loading...</div>
    <script>
        function graphQLFetcher(graphQLParams) {
        graphQLParams.variables = graphQLParams.variables ? JSON.parse(graphQLParams.variables) : null;
        return fetch("/query", {
            method: "post",
            body: JSON.stringify(graphQLParams),
            credentials: "include",
        }).then(function (response) {
            return response.text();
        }).then(function (responseBody) {
            try {
            return JSON.parse(responseBody);
            } catch (error) {
            return responseBody;
            }
        });
        }

        ReactDOM.render(
        React.createElement(GraphiQL, {fetcher: graphQLFetcher}),
        document.getElementById("graphiql")
        );
    </script>
    </body>
</html>
`)
