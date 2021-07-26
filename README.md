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

-------------

# Generate fully working Go CRUD HTTP API with Ent

## Introduction

One of the major time consumers when setting up a new API is setting up the basic CRUD (Create, Read, Update, Delete)
operations that repeat itself for every new entity you add to your graph. Luckily there is an extension to the
[`ent`](entgo.io) framework aiming to provide such handlers, including level logging, validation of the request body,
eager loading relations and serializing, all while leaving reflection out of sight and maintaining
type-safety: [elk](github.com/masseelch/elk). Let’s dig in!

## Setting up elk

First make sure you have the latest release of `elk` installed in your project:

```shell
go get github.com/masseelch/elk
```

The next step is to enable
the `elk` [extension](https://github.com/ent/ent/blob/a19a89a141cf1a5e1b38c93d7898f218a1f86c94/entc/entc.go#L197). This
requires you to use `entc` (enc codegen) package as
described [here](https://entgo.io/docs/code-gen#use-entc-as-a-package). Follow the next 3 steps to enable it and tell
the generator to execute the `elk` templates:

1. Create a new Go file named `ent/entc.go`, and paste the following content:

```go
// +build ignore

package main

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/masseelch/elk"
	"log"
)

func main() {
	ex, err := elk.NewExtension()
	if err != nil {
		log.Fatalf("creating elk extension: %v", err)
	}
	err = entc.Generate("./schema", &gen.Config{}, entc.Extensions(ex))
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}

```

2. Edit the `ent/generate.go` file to execute the `ent/entc.go` file:

```go
package ent

//go:generate go run -mod=mod entc.go

```

3. Run codegen for your ent project:

```shell
go generate./...
```

Since now all is set up create a schema, add some data and make use of elk-empowered ent!

## Setting up a simple server

To show you what elk can do for you, we use the schema and data [`ent`](entgo.io) described
[in its docs](https://entgo.io/docs/traversals). Head over there and create the schema as mentioned. You should end up
with a graph like that below:

![](https://entgo.io/images/assets/er_traversal_graph.png)

The generated handlers use [go-chi](https://github.com/go-chi/chi) to parse path and query parameters. However the
handlers implement `net/http`s `HandleFunc` interface and therefore seamlessly integrate in most existing apis.  
Furthermore `elk` uses [zap](https://github.com/uber-go/zap) for logging and
[go-playgrounds validator](https://github.com/go-playground/validator) to validate create / update request bodies.
Rendering is done by [sheriff](github.com/liip/sheriff) and [render](github.com/masseelch/render). To hook up our api
with the generated handlers add the following file:

```go
// main.go
package main

import (
	"<project>/ent"
	elk "<project>/ent/http"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	// Create the ent client.
	c, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer c.Close()
	// Run the auto migration tool.
	if err := c.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	// Create a zap logger to use.
	l, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("failed creating logger: %v", err)
	}
	// Validator used by elks handlers.
	v := validator.New()
	// Create a router.
	r := chi.NewRouter()
	// Hook up our generated handlers.
	r.Route("/pets", func(r chi.Router) {
		elk.NewPetHandler(c, l, v).Mount(r, elk.PetRoutes)
	})
	r.Route("/users", func(r chi.Router) {
		// We dont allow user deletion.
		elk.NewUserHandler(c, l, v).Mount(r, elk.PetRoutes &^ elk.UserDelete)
	})
	r.Route("/groups", func(r chi.Router) {
		// Dont include sub-resource routes.
		elk.NewGroupHandler(c, l, v).Mount(r, elk.GroupCreate | elk.GroupRead | elk.GroupUpdate | elk.GroupDelete | elk.GroupList)
	})
	// Start listen to incoming requests.
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
```

You can find a ready to be copied example [here](https://github.com/masseelch/elk/tree/master/examples/pets).

## Examples

You find an extensive list of examples of `elk`s capabilities below.

<details>
<summary>List a resource</summary>

`elk` provides endpoints to list a resource. Pagination is already set up.

```shell
curl 'localhost:8080/pets?itemsPerPage=2&page=2'
```

```json
[
  {
    "id": 3,
    "name": "Coco",
    "edges": {}
  }
]
```

</details>

<details>
<summary>Read a resource</summary>

To get detailed information about a resource set a path parameter.

```shell
curl 'localhost:8080/pets/3
```

```json
  {
  "id": 3,
  "name": "Coco",
  "edges": {}
}
```

</details>

<details>
<summary>Create a resource</summary>

To create a new resource send an POST request with `application/json` encoded body.

```shell
curl -X 'POST' -H 'Content-Type: application/json' -d '{"name":"Bob","owner":2}' 'localhost:8080/pets'
```

```json
{
  "id": 4,
  "name": "Bob",
  "edges": {}
}
```

</details>

<details>
<summary>Update a resource</summary>

To update a resources property send an PATCH request with `application/json` encoded body.

```shell
curl -X 'PATCH' -H 'Content-Type: application/json' -d '{"name":"Bobs Changed Name"}' 'localhost:8080/pets/4'
```

```json
{
  "id": 4,
  "name": "Bobs Changed Name",
  "edges": {}
}
```

</details>

<details>
<summary>Delete a resource</summary>

The handlers return a 204 response.

```shell
curl -X 'DELETE' 'localhost:8080/pets/4'
```

</details>

<details>
<summary>Request validation</summary>

`elk` can validate data sent in POST or PATCH requests. Use the `elk.Annotation` to set validation rules on fields.
Head over to [go-playgrounds validator](https://github.com/go-playground/validator) to see what validation rules exist.

```
// ent/schema/user.go

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("age"),
		field.String("name").
			Annotations(elk.Annotation{
				// No numbers allowed in name and it has to be at least 3 chars long.
				CreateValidation: "alpha,min=3",
				UpdateValidation: "alpha,min=3",
		}),
	}
}
```

```shell
curl -X 'POST' -H 'Content-Type: application/json' -d '{"name":"A"}' 'localhost:8080/users'
```

```json
{
  "code": 400,
  "status": "Bad Request",
  "errors": {
    "name": "This value failed validation on 'min:3'."
  }
}
```

</details>

<details>
<summary>Error responses</summary>

You get meaningful error responses.

```shell
curl -X 'POST' -H 'Content-Type: application/json' -d 'foo bar wtf' 'localhost:8080/pets'
```

```json
{
  "code": 400,
  "status": "Bad Request",
  "errors": "invalid json string"
}
```

</details>

<details>
<summary>Subresource routes</summary>

`elk` provides endpoints to fetch a resources edges.

```shell
curl 'localhost:8080/users/2/pets'
```

```json
[
  {
    "id": 1,
    "name": "Pedro",
    "edges": {}
  },
  {
    "id": 2,
    "name": "Xabi",
    "edges": {}
  },
  {
    "id": 4,
    "name": "Bob",
    "edges": {}
  }
]
```

</details>

<details>
<summary>Eager load edges</summary>

You can tell `elk` to eager load edges on specific routes by the use of serialization groups. Use `elk.SchemaAnnotation`
to define what groups to load on what endpoint and `elk.Annotation` on fields and edges to tell the serializer what
fields and edges are included in which group. `elk` takes care of eager loading the correct nodes.

```go
// ent/schema/pet.go

package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/masseelch/elk"
)

// Pet holds the schema definition for the Pet entity.
type Pet struct {
	ent.Schema
}

// Fields of the Pet.
func (Pet) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Annotations(elk.Annotation{
				// Include the name on the "pet:list" group.
				Groups: []string{"pet:list"},
			}),
	}
}

// Edges of the Pet.
func (Pet) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("friends", Pet.Type),
		edge.From("owner", User.Type).
			Ref("pets").
			Unique().
			Annotations(elk.Annotation{
				// Include the owner on the "pet:list" group.
				Groups: []string{"pet:list"},
			}),
	}
}

// Annotations of the Pet.
func (Pet) Annotations() []schema.Annotation {
	return []schema.Annotation{
		elk.SchemaAnnotation{
			// Tell elk to use the "pet:list" group on list routes.
			ListGroups: []string{"pet:list"},
		},
	}
}
```

```shell
curl 'localhost:8080/pets'
```

```json
[
  {
    "id": 1,
    "name": "Pedro",
    "edges": {
      "owner": {
        "id": 2,
        "age": 30,
        "name": "Ariel",
        "edges": {}
      }
    }
  },
  {
    "id": 2,
    "name": "Xabi",
    "edges": {
      "owner": {
        "id": 2,
        "age": 30,
        "name": "Ariel",
        "edges": {}
      }
    }
  },
  {
    "id": 3,
    "name": "Coco",
    "edges": {
      "owner": {
        "id": 3,
        "age": 37,
        "name": "Alex",
        "edges": {}
      }
    }
  }
]
```

</details>

<details>
<summary>Skip handlers</summary>

`elk` does always generate all handlers. You can declare what routes to mount.

```go
elk.NewPetHandler(c, l, v).Mount(r, elk.PetCreate | elk.PetList | elk.PetRead)
```

The compiler will not include the unused handlers since they are never called.

</details>

<details>
<summary>Logging</summary>

`elk` does leveled logging with [zap](https://github.com/uber-go/zap). See the example output below.

```shell
2021-07-22T07:22:25.436+0200    INFO    http/create.go:167      pet rendered    {"handler": "PetHandler", "method": "Create", "id": 4}
2021-07-22T07:22:25.450+0200    INFO    http/create.go:198      validation failed       {"handler": "UserHandler", "method": "Create", "error": "Key: 'UserCreateRequest.name' Error:Field validation for 'name' failed on the 'min' tag"}
2021-07-22T07:22:25.463+0200    INFO    http/update.go:239      validation failed       {"handler": "UserHandler", "method": "Update", "error": "Key: 'UserUpdateRequest.name' Error:Field validation for 'name' failed on the 'min' tag"}
2021-07-22T07:22:25.489+0200    INFO    http/create.go:254      user rendered   {"handler": "UserHandler", "method": "Create", "id": 4}
2021-07-22T07:22:25.508+0200    INFO    http/read.go:150        user rendered   {"handler": "UserHandler", "method": "Read", "id": 2}
```

</details>

## Future and Known Issues

`elk` has many cool features already, but there are some issues to address in the near future.

Currently, `elk` does use this [render package](github.com/masseelch/render) on combination with
[sheriff](github.com/liip/sheriff) to render its output to the client.
`render` does use reflection under the hood since it calls `json.Marshal` / `xml.Marshal`, as well does `sheriff`.
The mapping of request values does currently only work for `application/json` bodies and uses `json.Unmarshal`.
The goal is to have elk provide interfaces
`Renderer` and `Binder` which will be implemented by the generated nodes / request structs. This allows type safe and
reflection-free transformation between json / xml / protobuf and go structs.

`ent` already has some builtin validation which is not yet reflected by the request validation `elk` generates.
Validation is only executed if there are tags given.

`elk`s generated bitmask to choose what handlers to mount are not typesafe yet.

Another not yet implemented feature is to give the developer the possibility to customize the generated code by
providing custom templates to `elk` like `ent` already does with
[External Templates](https://entgo.io/docs/templates).

In the future `elk` is meant to provide fully working and easily extendable integration tests for the generated code.

Initial work to generate a fully working flutter frontend has been done and will hopefully lead to a release soon.
