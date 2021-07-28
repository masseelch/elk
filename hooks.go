package elk

import (
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/edge"
	"fmt"
	"reflect"
	"strings"
)

// AddGroupsTag adds the serialization groups defined by the annotation of each field to the generated entity struct.
func AddGroupsTag(next gen.Generator) gen.Generator {
	return gen.GenerateFunc(func(g *gen.Graph) error {
		for _, n := range g.Nodes {
			// Groups for fields.
			for _, f := range n.Fields {
				tag := reflect.StructTag(f.StructTag)

				// If the field does not yet have a groups tag and there are groups defined on the annotation add those
				// groups to the fields struct-tag.
				if _, ok := tag.Lookup("groups"); !ok {
					a := Annotation{}
					if f.Annotations != nil && f.Annotations[a.Name()] != nil {
						if err := a.Decode(f.Annotations[a.Name()]); err != nil {
							return err
						}

						f.StructTag = fmt.Sprintf("%s %s", tag, a.Groups.StructTag())
					}
				}
			}

			// Groups for edges.
			gs := groups{}
			for _, e := range n.Edges {
				tag := reflect.StructTag(e.StructTag)

				// If the edge does not yet have a groups tag and there are groups defined on the annotation add those
				// groups to the edges struct-tag.
				if _, ok := tag.Lookup("groups"); !ok {
					a := Annotation{}
					if e.Annotations != nil && e.Annotations[a.Name()] != nil {
						if err := a.Decode(e.Annotations[a.Name()]); err != nil {
							return err
						}

						gs.Add(a.Groups...)
						e.StructTag = fmt.Sprintf("%s %s", tag, a.Groups.StructTag())
					}
				}
			}

			// Make sure to add all groups used on the edges to the Edges field of the generated node.
			var a edge.Annotation
			if n.Annotations[a.Name()] != nil {
				a = n.Annotations[a.Name()].(edge.Annotation)
			}
			if len(gs) > 0 {
				a.StructTag = fmt.Sprintf(`%s groups:"%s"`, a.StructTag, strings.Join(gs, ","))
			}
			if n.Annotations == nil {
				n.Annotations = make(gen.Annotations)
			}
			n.Annotations[a.Name()] = a
		}

		return next.Generate(g)
	})
}
