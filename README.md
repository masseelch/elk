# elk

This package provides an extension to the [awesome entgo.io](https://github.com/ent/ent) code generator.

`elk` can do two things for you:

1. Generate a fully compliant, extendable [OpenAPI](https://spec.openapis.org/oas/v3.0.3.html)
   specification file to enable you to make use of the [Swagger Tooling](https://swagger.io/tools/) to generate RESTful
   server stubs and clients.
2. Generate a ready-to-use and extendable server implementation of the OpenAPI specification. The code generated
   by `elk` uses the Ent ORM while maintaining complete type-safety and leaving reflection out of sight.

> :warning: **This is work in progress**: The API may change without further notice!

This package depends on [Ent](https://entgo.io), an ORM project for Go. To learn more about Ent, how to connect to
different types of databases, run migrations or work with entities head over to
their [documentation](https://entgo.io/docs/getting-started).

## Getting Started

The first step is to add the `elk` package to your project:

```shell
go get github.com/masseelch/elk
```

`elk` uses the Ent [Extension API](https://entgo.io/docs/extensions) to integrate with Ent’s code-generation. This
requires that we use the `entc` (ent codegen) package as
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
	ex, err := elk.NewExtension(
		elk.GenerateSpec("openapi.json"),
		elk.GenerateHandlers(),
	)
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

3. _(Only required if server generation is enabled)_ `elk` uses some external packages in its generated code. Currently,
   you have to get those packages manually once when setting up `elk`:

```shell
go get github.com/mailru/easyjson github.com/go-chi/chi/v5 go.uber.org/zap
```

4. Run the code generator:

```shell
go generate ./...
```

In addition to the files Ent would normally generate, another directory named `ent/http` and a file named `openapi.json`
was created. The `ent/http` directory contains the code for the `elk`-generated HTTP CRUD handlers while `openapi.json`
contains the OpenAPI Specification. Feel free to have a look
at [this example spec file](https://github.com/masseelch/elk/tree/master/internal/simple/ent/openapo.json)
and [the implementing server code](https://github.com/masseelch/elk/tree/master/internal/simple/ent/http).

### Setting up a server

This section guides you to a very simple setup for an `elk`-powered Ent. The following two files define the two schemas
Pet and User with a Many-To-One relation: A Pet belongs to a User, and a User can have multiple Pets.

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

After regenerating the code you can spin up a runnable server with the below `main` function:

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
	// Start listen to incoming requests.
	if err := http.ListenAndServe(":8080", elk.NewHandler(c, zap.NewExample())); err != nil {
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

The response data on the creation operation does not include the User the new Pet belongs to. `elk` does not include
edges in its output by default. You can configure `elk` to render edges using a feature called **serialization groups**.

## Serialization Groups

`elk` by default includes every field of a schema in an endpoints output and excludes fields. This behaviour can be
changed by using **serialization groups**. You can configure `elk` what **serialization groups** to request on what
endpoint using a `elk.SchemaAnnotation`. With a `elk.Annotation` you configure what fields and edges to include. `elk`
follows the following rules to determine if a field or edge is included or not:

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
					return errors.New("name must begin with uppercase")
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

`elk` paginates all list endpoints. This is valid for both resource and sub-resources routes.

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

## Configuration

`elk` lets you decide what endpoints you want it to generate by the use of **generation policies**. You can either
expose all routes by default and hide some you are not interested in or exclude all routes by default and only expose
those you want generated:

***ent/entc.go***

```go
package main

import (
	"log"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/masseelch/elk"
	"github.com/masseelch/elk/policy"
	"github.com/masseelch/elk/spec"
)

func main() {
	ex, err := elk.NewExtension(
		elk.GenerateSpec("openapi.json"),
		elk.GenerateHandlers(),
		// Exclude all routes by default.
		elk.DefaultHandlerPolicy(elk.Exclude),
	)
	if err != nil {
		log.Fatalf("creating elk extension: %v", err)
	}
	err = entc.Generate("./schema", &gen.Config{}, entc.Extensions(ex))
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}

```

***ent/schema/user.go***

```go
package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/masseelch/elk"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Annotations of the User.
func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		// Generate creation and read endpoints.
		elk.Expose(elk.Create, elk.Read),
    }
}
```

For more information about how to configure `elk` and what it can do have a look at
the [docs](https://pkg.go.dev/github.com/masseelch/elk) [integration test setup](https://github.com/masseelch/elk/tree/master/internal)
.

## Known Issues and Outlook

- `elk` does currently only work with JSON. It is relatively easy to support XML as well and there are plans to provide
  conditional XML / JSON parsing and rendering based on the `Content-Type` and `Accept` headers.

- The generated code does not have very good automated tests yet.

## Contribution

`elk` has not reach its first release yet but the API can be considered somewhat stable. I welcome any suggestion or
feedback and if you are willing to help I'd be very glad. The [issues tab](https://github.com/masseelch/elk/issues) is a
wonderful place for you to reach out for help, feedback, suggestions and contribution.
