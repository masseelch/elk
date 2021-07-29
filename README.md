# elk

This package aims to extend the [awesome entgo.io](https://github.com/ent/ent) code generator to generate a
fully-functional HTTP API on a defined set of entities.

> :warning: **This is work in progress**: The API may change without further notice!
>

## Getting Started

First make sure you have the latest version of `elk` installed in your project:

```shell
go get -u github.com/masseelch/elk
```

`elk` uses the
Ent [extension API](https://github.com/ent/ent/blob/a19a89a141cf1a5e1b38c93d7898f218a1f86c94/entc/entc.go#L197) to
integrate with Ent’s code-generation. This requires that we use the `entc` (ent codegen) package as
described [here](https://entgo.io/docs/code-gen#use-entc-as-a-package). Follow the next three steps to enable it and to
configure Ent to work with the `elk` extension:

1. Create a new Go file named `ent/entc.go` and paste the following content:

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

3. `elk` uses a not yet released version of Ent. To have the dependencies up to date run the following:

```shell
go mod tidy
```

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

`elk` can validate data sent in POST or PATCH requests. Use the `elk.Annotation` to set validation rules on fields. Head
over to [go-playgrounds validator](https://github.com/go-playground/validator) to see what validation rules exist.

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
`render` does use reflection under the hood since it calls `json.Marshal` / `xml.Marshal`, as well does `sheriff`. The
mapping of request values does currently only work for `application/json` bodies and uses `json.Unmarshal`. The goal is
to have elk provide interfaces
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
