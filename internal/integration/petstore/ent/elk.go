// +build ignore

package main

import (
	"log"

	"github.com/masseelch/elk"
)

func main() {
	// ent plus http
	// err := entc.Generate("./schema", &gen.Config{
	// 	Templates: elk.HTTPTemplates,
	// 	Hooks: []gen.Hook{
	// 		elk.AddGroupsTag,
	// 	},
	// })
	// if err != nil {
	// 	log.Fatalf("running ent codegen: %v", err)
	// }
	// flutter
	if err := elk.Flutter("./schema", "../../client/generated"); err != nil {
		log.Fatalf("running flutter codegen: %v", err)
	}
}