# elk

This package aims to extend the [awesome entgo.io](https://github.com/ent/ent) code generator to generate a
fully-functional HTTP API for your schema. `elk` strives to automate all the tedious work of setting up the basic CRUD
endpoints for every entity you add to your graph, including logging, validation of the request body, eager loading
relations and serializing, all while leaving reflection out of sight and maintaining type-safety.

> :warning: **This is work in progress**: The API may change without further notice!

This package depends on [Ent](https://entgo.io), an ORM project for Go. To learn more about Ent, how to connect to
different types of databases, run migrations or work with entities head over to
their [documentation](https://entgo.io/docs/getting-started).

## Getting Started

The first step is to add the `elk` package to your project:

```shell
go get github.com/masseelch/elk
```

`elk` uses the
Ent [extension API](https://github.com/ent/ent/blob/a19a89a141cf1a5e1b38c93d7898f218a1f86c94/entc/entc.go#L197) to
integrate with Ent’s code-generation. This requires that we use the `entc` (ent codegen) package as
described [here](https://entgo.io/docs/code-gen#use-entc-as-a-package). Follow the next four steps to enable it and to
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

3. `elk` uses some external packages in its generated code. Currently, you have to get those packages manually once when
   setting up `elk`:

```shell
go get github.com/mailru/easyjson github.com/masseelch/render github.com/go-chi/chi/v5 go.uber.org/zap
```

4. Run the code generator:

```shell
go generate ./...
```

In addition to the files Ent would normally generate, another directory names `ent/http` was created. There files
contain the code for the `elk`-generated HTTP CRUD handlers. An example of the generated code can be
found [here](https://github.com/masseelch/elk/tree/master/internal/integration/pets/ent/http).

### Setting up a server

This section guides you to a very simple setup for an `elk`-powered Ent. The following two files define two schemas Pet
and User with a Many-To-One relation: A Pet belongs to a User, and a User can have multiple Pets.

***ent/schema/pet.go***

```go
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Pet holds the schema definition for the Pet entity.
type Pet struct {
	ent.Schema
}

// Fields of the Pet.
func (Pet) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
	}
}

// Edges of the Pet.
func (Pet) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("pets").
			Unique(),
	}
}
```

***ent/schema/user.go***

```go
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("pets", Pet.Type),
	}
}
```

To spin up a runnable server you can use the below `main` function:

```go
package main

import (
	"context"
	"log"
	"net/http"

	"<your-project>/ent"
	elk "<your-project>/ent/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

func main() {
	// Create the ent client. This opens up a sqlite file named elk.db.
	c, err := ent.Open("sqlite3", "./elk.db?_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer c.Close()
	// Run the auto migration tool.
	if err := c.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	// Router and Logger.
	r, l := chi.NewRouter(), zap.NewExample()
	// Mount the generates handlers.
	r.Route("/pets", func(r chi.Router) {
		elk.NewPetHandler(c, l).Mount(r, elk.PetRoutes)
	})
	r.Route("/users", func(r chi.Router) {
		// Only register the create and read endpoints.
		elk.NewUserHandler(c, l).Mount(r, elk.UserCreate|elk.UserRead)
	})
	// Start listen to incoming requests.
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
```

Start the server:

```shell
go run -mod=mod main.go
```

Congratulations! You now have a running server serving the Pets API. The database is still empty though. the following
two curl requests create a new user and adds a pet, that belongs to that user.

```shell
curl -X 'POST' -H 'Content-Type: application/json' -d '{"name":"Elk"}' 'localhost:8080/users'
```

```json
{
  "id": 1,
  "name": "Elk"
}
```

```shell
curl -X 'POST' -H 'Content-Type: application/json' -d '{"name":"Kuro","owner":1}' 'localhost:8080/pets'
```

```json
{
  "id": 1,
  "name": "Kuro"
}
```

The response data on the create action does not include the User the new Pet belongs to. `elk` does not include edges in
its output by default. You can configure `elk` to render edges using a feature called **serialization groups**.

## Serialization Groups

`elk` by default includes every field of a schema in an endpoints output and excludes fields. This behaviour can be
changed by using **serialization groups**. You can configure `elk` what **serialization groups** to request on what
endpoint using a `elk.SchemaAnnotation`. With a `elk.Annotation` you configure what fields and edges to include. `elk`
follows the following rules to determine if a field or edge is inlcuded or not:

- If no groups are requested all fields are included and all edges are excluded
- If a group x is requested all fields with no groups and fields with group x are included. Edges with x are eager
  loaded and rendered.

Change the previously mentioned schemas and add **serialization groups**:

***ent/schema/pet.go***

```go
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
		field.String("name"),
	}
}

// Edges of the Pet.
func (Pet) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("pets").
			Unique().
			// render this edge if one of 'pet:read' or 'pet:list' is requested.
			Annotations(elk.Groups("pet:read", "pet:list")),
	}
}

// Annotations of the Pet.
func (Pet) Annotations() []schema.Annotation {
	return []schema.Annotation{
		// Request the 'pet:read' group when rendering the entity after creation.
		elk.CreateGroups("pet:read"),
		// You can request several groups per endpoint.
		elk.ReadGroups("pet:list", "pet:read"),
	}
}
```

***ent/schema/user.go***

```go
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/masseelch/elk"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			// render this field only if no groups or the 'owner:read' groups is requested.
			Annotations(elk.Groups("owner:read")),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("pets", Pet.Type),
	}
}
```

After regenerating the code and restarting the server `elk` renders the owner if you create a new pet.

```shell
curl -X 'POST' -H 'Content-Type: application/json' -d '{"name":"Martha","owner":1}' 'localhost:8080/pets'
```

```json
{
  "id": 2,
  "name": "Martha",
  "owner": {
    "id": 1,
    "name": "Elk"
  }
}
```

## Validation

`elk` supports the [validation](https://entgo.io/docs/schema-fields#validators) feature of Ent. For demonstration extend
the above Pet schema:

```go
package schema

import (
	"errors"
	"strings"

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
		field.Int("age").
			// Validator will only be called if the request body has a 
			// non nil value for the field 'age'.
			Optional().
			// Works for built-in validators.
			Positive(),
		field.String("name").
			// Works for built-in validators.
			MinLen(3).
			// Works for custom validators.
			Validate(func(s string) error {
				if strings.ToLower(s) == s {
					return errors.New("group name must begin with uppercase")
				}
				return nil
			}),
		// Enums are validated against the allowed values.
		field.Enum("color").
			Values("red", "blue", "green", "yellow"),
	}
}

// Edges of the Pet.
func (Pet) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("pets").
			Unique().
			// Works with edge validation.
			Required().
			// render this edge if one of 'pet:read' or 'pet:list' is requested.
			Annotations(elk.Groups("pet:read", "pet:list")),
	}
}

// Annotations of the Pet.
func (Pet) Annotations() []schema.Annotation {
	return []schema.Annotation{
		// Request the 'pet:read' group when rendering the entity after creation.
		elk.CreateGroups("pet:read"),
		// You can request several groups per endpoint.
		elk.ReadGroups("pet:list", "pet:read"),
	}
}
```

## Sub Resources

`elk` provides first level sub resource handlers for all your entities. With previously set up server, run the
following:

```shell
curl 'localhost:8080/pets/1/owner'
```

You'll get information about the Owner of the Pet with the id 1. `elk` uses `elk.SchemaAnnotation.ReadGroups` for a
unique edge and `elk.SchemaAnnotation.ListGroups` for a non-unique edge.

## Pagination

`elk` paginates list endpoints. This is valid for both resource and sub-resources routes.

```shell
curl 'localhost:8080/pets?page=2&itemsPerPage=1'
```

```json
[
  {
    "id": 2,
    "name": "Martha"
  }
]
```

## Known Issues and Outlook

- `elk` does currently only work with JSON. It is relatively easy to support XML as well and there are plans to provide
  conditional XML / JSON parsing and rendering based on the `Content-Type` and `Accept` headers.

- `elk`s generated bitmask to choose what handlers to mount is not typesafe yet. You can
  call `elk.FooHandler(c, l).Mount(r, elk.BarRoutes)` without any compile time errors.

- Customization of the generated handlers is not great yet. I'd like to provide something similar to Ent'
  s [External Templates](https://entgo.io/docs/templates) in the future.

- The generated code does not have very good automated tests yet.

- Initial work to generate a fully working flutter frontend has been done.

## Contribution

`elk` is in an early stage of development, we welcome any suggestion or feedback and if you are willing to help I'd be
very glad. The [issues tab](https://github.com/masseelch/elk/issues) is a wonderful place for you to reach out for help,
feedback, suggestions and contribution.