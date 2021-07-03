// +build ignore

package main

import (
	"log"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/masseelch/elk"
)

func main() {
	err := entc.Generate("./schema", &gen.Config{
		Templates: elk.HTTPTemplates,
		Hooks: []gen.Hook{
			elk.AddGroupsTag,
		},
	})
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}