# elk

This package aims to extend the [awesome entgo.io](https://github.com/ent/ent) code generator to generate fully functional code on a defined set of entities.

> :warning: **This is work in progress**: The API may change without further notice!
> 
### Features
- Generate http crud handlers 
- Generate flutter models and http client to consume the generated http api

### How to use

#### 1. Create a new Go file named `ent/elk.go`, and paste the following content:
```go
// +build ignore

package main

import (
	"log"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/masseelch/elk"
)

func main() {
	// ent plus http
	err := entc.Generate("./schema", &gen.Config{
		Templates: elk.HTTPTemplates,
		Hooks: []gen.Hook{
			elk.AddGroupsTag,
		},
	})
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
	// flutter
	if err := elk.Flutter("./schema", ""); err != nil {
		log.Fatalf("running flutter codegen: %v", err)
	}
}
```

#### 2. Edit the `ent/generate.go` file to execute the `ent/elk.go` file:
```go
package ent

//go:generate go run -mod=mod elk.go
```

#### 3. Run codegen for your ent project:
```shell
go generate ./...
```