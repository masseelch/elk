//go:build ignore
// +build ignore

package main

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/masseelch/elk"
	"log"
)

func main() {
	ex, err := elk.NewExtension(
		elk.EnableSpecGenerator(nil),
		elk.SpecTitle("My Pets API"),
		elk.SpecDescription("Awesome, Mega Cool API to manage Ariel's Pet Leopards!"),
		elk.SpecVersion("0.0.1"),
	)
	if err != nil {
		log.Fatalf("creating elk extension: %v", err)
	}
	err = entc.Generate("./schema", &gen.Config{}, entc.Extensions(ex))
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
