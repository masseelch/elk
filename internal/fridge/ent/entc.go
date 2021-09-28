//go:build ignore
// +build ignore

package main

import (
	"log"
	"net/http"
	"strconv"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/masseelch/elk"
	"github.com/masseelch/elk/spec"
)

func main() {
	ex, err := elk.NewExtension(
		elk.GenerateHandlers(),
		elk.GenerateSpec(
			"openapi.json",
			elk.SpecTitle("Fridge CMS"), // It is a Content-Management-System ...
			elk.SpecDescription("API to manage fridges and their cooled contents. **ICY!**"), // You can use CommonMark syntax.
			elk.SpecVersion("0.0.1"),
			func(next elk.Generator) elk.Generator {
				return elk.GenerateFunc(func(s *spec.Spec) error {
					// Run the generator so the spec is filled.
					if err := next.Generate(s); err != nil {
						return err
					}
					// Add our custom path.
					s.Paths["/fridges/{id}/contents"] = &spec.Path{
						Get: &spec.Operation{
							Summary:     "Return everything stored in this fridge",
							Description: "List every item stored in every compartment belonging to this fridge.",
							Tags:        []string{"Fridge"},
							OperationID: "fridgeContents",
							Responses: map[string]*spec.OperationResponse{
								strconv.Itoa(http.StatusOK): {
									Response: spec.Response{
										Description: "All the contents",
										Content: &spec.Content{
											spec.JSON: spec.MediaTypeObject{
												Unique: false,
												Schema: spec.Schema{
													Name: "FridgeContents",
													Fields: spec.Fields{
														"id":   {Type: spec.Type{Type: "integer"}},
														"name": {Type: spec.Type{Type: "string"}},
													},
												},
											},
										},
									},
								},
							},
						},
						Parameters: []spec.Parameter{{
							Name:        "id",
							In:          spec.InPath,
							Description: "ID of the fridge",
							Required:    true,
							Schema: spec.Type{
								Type:   "integer",
								Format: "int32",
							},
						}},
					}
					return nil
				})
			},
		),
	)
	if err != nil {
		log.Fatalf("creating elk extension: %v", err)
	}
	err = entc.Generate("./schema", &gen.Config{}, entc.Extensions(ex))
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
