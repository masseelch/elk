package elk

import (
	"entgo.io/ent/entc/gen"
	"fmt"
	"github.com/masseelch/elk/serialization"
	"hash/fnv"
)

type (
	View struct {
		Node   *gen.Type
		Fields []*gen.Field
		Edges  []*ViewEdge
	}
	ViewEdge struct {
		*gen.Edge
		Name string
	}
)

var (
	// responseViewCache is used to reduce construction-time for a view.
	// The key is a combination of the node and groups requested.
	responseViewCache = make(map[string]*View)
)

// views returns a map of all views occurring in the given graph. Key is the View's name.
func views(g *gen.Graph) (map[string]*View, error) {
	gss, err := groupCombinations(g)
	if err != nil {
		return nil, err
	}
	m := make(map[string]*View)
	for _, n := range g.Nodes {
		for _, gs := range gss {
			// Generate the view.
			r, err := view(n, gs)
			if err != nil {
				return nil, err
			}
			v, err := r.Name()
			if err != nil {
				return nil, err
			}
			m[v] = r
		}
	}
	return m, nil
}

// Name returns a unique name for this view.
func (v View) Name() (string, error) {
	h := fnv.New32a()
	if _, err := h.Write([]byte(v.Node.Name)); err != nil {
		return "", err
	}
	for _, f := range v.Fields {
		if _, err := h.Write([]byte(f.Name)); err != nil {
			return "", err
		}
	}
	for _, e := range v.Edges {
		if _, err := h.Write([]byte(e.Name)); err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%s%dView", v.Node.Name, h.Sum32()), nil
}

// view create a new View for the given type and serialization.Groups.
func view(n *gen.Type, gs serialization.Groups) (*View, error) {
	var err error
	h := hashNodeAndGroups(n, gs)
	v, ok := responseViewCache[h]
	if !ok {
		v, err = viewHelper(n, gs, true)
		if err != nil {
			return nil, err
		}
		responseViewCache[h] = v
	}
	return v, nil
}

func viewHelper(n *gen.Type, gs serialization.Groups, loadEdges bool) (*View, error) {
	v := &View{Node: n}
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
			edg := &ViewEdge{Edge: e}
			if loadEdges {
				er, err := viewHelper(e.Type, gs, false)
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

// groupCombinations returns all groups ever requested together.
func groupCombinations(g *gen.Graph) (serialization.Collection, error) {
	gss := serialization.Collection{}
	for _, n := range g.Nodes {
		// For every action extract the requested groups.
		for _, a := range [...]string{actionCreate, actionRead, actionUpdate, actionList} {
			gs, err := groupsForAction(n, a)
			if err != nil {
				return nil, err
			}
			if !gss.Contains(gs) {
				gss = append(gss, gs)
			}
		}
	}
	return gss, nil
}

// hashNodeAndGroups returns a unique Hash for a given node and groups.
func hashNodeAndGroups(n *gen.Type, gs serialization.Groups) string {
	return fmt.Sprintf("%s_%d", n.Name, gs.Hash())
}
