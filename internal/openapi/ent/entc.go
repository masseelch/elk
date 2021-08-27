// +build ignore

package main

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/masseelch/elk"
	"github.com/masseelch/elk/openapi/spec"
	"log"
)

func main() {
	spec, err := spec.New()
	if err != nil {
		log.Fatalf("creating openapi spec: %v", err)
	}
	ex, err := elk.NewExtension(elk.WithOpenAPISpec(spec))
	if err != nil {
		log.Fatalf("creating elk extension: %v", err)
	}
	err = entc.Generate("./schema", &gen.Config{}, entc.Extensions(ex))
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
