package elk

import (
	"entgo.io/ent/entc/gen"
	"fmt"
	inflect "github.com/go-openapi/inflect"
	"github.com/masseelch/elk/serialization"
	"log"
	"strings"
)

type (
	// A view is a subset of a Node. It may hold fewer Fields and Edges than the Node it is derived from. It has no
	// knowledge of the serialization.Groups that led to its creation.
	view struct {
		Node   *gen.Type
		Fields []*gen.Field
		Edges  []*viewEdge
		Groups serialization.Groups
	}
	// A viewEdge wraps an ent Edge and holds the name of the view to use when this viewEdge is serialized. A viewEdge
	// is only valid in its views' context.
	viewEdge struct {
		*gen.Edge
		Name string
	}
	// A mergedView is essentially the same as a view but has information about what serialization.Groups where used to
	// create it.
	mergedView struct {
		view
		// Groups holds all group-combinations this view represents.
		Groups serialization.Collection
	}
)

var (
	// viewCache is used to reduce construction-time for a view.
	// The key is a combination of the node and groups requested.
	viewCache = make(map[string]*view)
)

// newViews returns a map of all views occurring in the given graph. Key is the view's name.
func newViews(g *gen.Graph) (map[string]*mergedView, error) {
	// Collect all groups ever requested together.
	gss := serialization.Collection{}
	for _, n := range g.Nodes {
		// For every operation extract the requested groups.
		for _, a := range [...]string{opCreate, opRead, opUpdate, opList} {
			// TODO: Do not return views for excluded operations.
			gs, err := groupsForOperation(n, a)
			if err != nil {
				return nil, err
			}
			if !gss.Contains(gs) {
				gss = append(gss, gs)
			}
		}
	}
	m := make(map[string]*mergedView)
	for _, n := range g.Nodes {
		for _, gs := range gss {
			// Generate the newView.
			r, err := newView(n, gs)
			if err != nil {
				return nil, err
			}
			v, err := r.Name()
			if err != nil {
				return nil, err
			}
			if mv, ok := m[v]; ok && !mv.Groups.Contains(gs) {
				mv.Groups = append(mv.Groups, gs)
			} else {
				name, _ := r.Name()
				log.Printf("name; groups: %s: %v", name, gs)
				m[v] = &mergedView{
					view:   *r,
					Groups: serialization.Collection{gs},
				}
			}
		}
	}
	return m, nil
}

// Name returns a unique name for this view.
func (v view) Name() (string, error) {
	groups := []string{}

	for _, elem := range v.Groups {
		camelized := inflect.Camelize(elem)
		if v.Node.Name != camelized {
			groups = append(groups, camelized)
		}
	}

	annotations, ok := v.Node.Annotations[elkSchemaName]
	schemaAnnotation, canCast := annotations.(map[string]interface{})
	if !ok || len(groups) == 0 || !canCast {

		// if no annotations or no groups (should be same thing) there's only one view
		return fmt.Sprintf("%sView", v.Node.Name), nil
	}

	for k, val := range schemaAnnotation {
		gs, ok := val.(serialization.Groups)
		if ok && len(gs) > 1 {
			return fmt.Sprintf("%v%sView", v.Node.Name, k), nil
		}
	}

	return fmt.Sprintf("%vWith%sView", v.Node.Name, strings.Join(groups, "And")), nil
}

// newView create a new view for the given type and serialization.Groups.
func newView(n *gen.Type, gs serialization.Groups) (*view, error) {
	var err error
	h := hashNodeAndGroups(n, gs)
	v, ok := viewCache[h]
	if !ok {
		v, err = newViewHelper(n, gs, true)
		if err != nil {
			return nil, err
		}
		viewCache[h] = v
	}
	return v, nil
}

func newViewHelper(n *gen.Type, gs serialization.Groups, loadEdges bool) (*view, error) {
	v := &view{Node: n, Groups: gs}
	ok, err := fieldNeedsSerialization(n.ID, gs)
	if err != nil {
		return nil, err
	}
	if ok {
		v.Fields = append(v.Fields, n.ID)
	}
	for _, f := range n.Fields {
		ok, err := fieldNeedsSerialization(f, gs)
		if err != nil {
			return nil, err
		}
		if ok {
			v.Fields = append(v.Fields, f)
		}
	}
	for _, e := range n.Edges {
		ok, err := edgeNeedsSerialization(e, gs)
		if err != nil {
			return nil, err
		}
		if ok {
			edg := &viewEdge{Edge: e}
			if loadEdges {
				er, err := newViewHelper(e.Type, gs, false)
				if err != nil {
					return nil, err
				}
				edg.Name, err = er.Name()
				if err != nil {
					return nil, err
				}
			}
			v.Edges = append(v.Edges, edg)
		}
	}
	return v, nil
}

// fieldNeedsSerialization checks if a field is to be serialized according to its annotations and the requested groups.
func fieldNeedsSerialization(f *gen.Field, g serialization.Groups) (bool, error) {
	// If the field is sensitive, don't serialize it.
	if f.Sensitive() {
		return false, nil
	}
	// If no groups are requested or the field has no groups defined render the field.
	if f.Annotations == nil || len(g) == 0 {
		return true, nil
	}
	// Extract the Groups defined on the edge.
	gs, err := groups(f.Annotations)
	if err != nil {
		return false, err
	}
	// If no groups are given on the field default is to include it in the output.
	if len(gs) == 0 {
		return true, nil
	}
	// If there are groups given check if the groups match the requested ones.
	return g.Match(gs), nil
}

// edgeNeedsSerialization checks if an edge is to be serialized according to its annotations and the requested groups.
func edgeNeedsSerialization(e *gen.Edge, g serialization.Groups) (bool, error) {
	// If no groups are requested or the edge has no groups defined do not render the edge.
	if e.Annotations == nil || len(g) == 0 {
		return false, nil
	}
	// Extract the Groups defined on the edge.
	gs, err := groups(e.Annotations)
	if err != nil {
		return false, err
	}
	// If no groups are given on the edge default is to exclude it.
	if len(gs) == 0 {
		return false, nil
	}
	// If there are groups given check if the groups match the requested ones.
	return g.Match(gs), nil
}

// hashNodeAndGroups returns a unique Hash for a given node and groups.
func hashNodeAndGroups(n *gen.Type, gs serialization.Groups) string {
	return fmt.Sprintf("%s_%d", n.Name, gs.Hash())
}
