package elk

import (
	"embed"
	"entgo.io/ent/entc/gen"
	"fmt"
	"github.com/stoewer/go-strcase"
	"hash/fnv"
	"text/template"
)

const (
	actionCreate = "create"
	actionRead   = "read"
	actionUpdate = "update"
	actionList   = "list"
)

var (
	actions = [...]string{actionCreate, actionRead, actionUpdate, actionList}
	//go:embed template
	templateDir embed.FS
)

var (
	// Funcs contains the extra template functions used by elk.
	Funcs = template.FuncMap{
		"edgesToLoad":     edgesToLoad,
		"kebab":           strcase.KebabCase,
		"needsValidation": needsValidation,
		"responseView":    responseView,
		"responseViews":   responseViews,
		"stringSlice":     stringSlice,
		"xextend":         xextend,
	}
	// HTTPTemplate holds all templates for generating http handlers.
	HTTPTemplate = gen.MustParse(gen.NewTemplate("elk").Funcs(Funcs).ParseFS(templateDir, "template/http/*.tmpl"))
)

// needsValidation returns if a type needs validation because there is some defined on one of its fields.
func needsValidation(n *gen.Type) bool {
	for _, f := range n.Fields {
		if f.Validators > 0 {
			return true
		}
	}
	return false
}

// fieldNeedsSerialization checks if a field is to be serialized according to its annotations and the requested groups.
func fieldNeedsSerialization(f *gen.Field, g groups) (bool, error) {
	// If the field is sensitive, don't serialize it.
	if f.Sensitive() {
		return false, nil
	}
	// If no groups are requested or the field has no groups defined render the field.
	if f.Annotations == nil || len(g) == 0 {
		return true, nil
	}
	// If there are groups given check if the groups match the requested ones.
	an := Annotation{}
	if err := an.Decode(f.Annotations[an.Name()]); err != nil {
		return false, err
	}
	// If no groups are given on the field default is to include it in the output.
	if len(an.Groups) == 0 {
		return true, nil
	}
	return g.Match(an.Groups), nil
}

type (
	ResponseView struct {
		Node   *gen.Type
		Fields []*gen.Field
		Edges  []*Edge
	}
	Edge struct {
		*gen.Edge
		ViewName string
	}
)

// ViewName returns a unique name for this view.
func (v ResponseView) ViewName() (string, error) {
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

func responseViews(g *gen.Graph) (map[string]*ResponseView, error) {
	gss, err := groupCombinations(g)
	if err != nil {
		return nil, err
	}
	m := make(map[string]*ResponseView)
	for _, n := range g.Nodes {
		for _, gs := range gss {
			// Generate the view.
			r, err := responseView(n, gs)
			if err != nil {
				return nil, err
			}
			v, err := r.ViewName()
			if err != nil {
				return nil, err
			}
			m[v] = r
		}
	}
	return m, nil
}

var responseViewCache = make(map[string]*ResponseView)

func responseView(n *gen.Type, gs groups) (*ResponseView, error) {
	var err error
	h := hashNodeAndGroups(n, gs)
	r, ok := responseViewCache[h]
	if !ok {
		r, err = responseViewHelper(n, gs, true)
		if err != nil {
			return nil, err
		}
		responseViewCache[h] = r
	}
	return r, nil
}

func responseViewHelper(n *gen.Type, gs groups, loadEdges bool) (*ResponseView, error) {
	r := &ResponseView{Node: n}
	ok, err := fieldNeedsSerialization(n.ID, gs)
	if err != nil {
		return nil, err
	}
	if ok {
		r.Fields = append(r.Fields, n.ID)
	}
	for _, f := range n.Fields {
		ok, err := fieldNeedsSerialization(f, gs)
		if err != nil {
			return nil, err
		}
		if ok {
			r.Fields = append(r.Fields, f)
		}
	}
	for _, e := range n.Edges {
		ok, err := edgeNeedsSerialization(e, gs)
		if err != nil {
			return nil, err
		}
		if ok {
			edg := &Edge{Edge: e}
			if loadEdges {
				er, err := responseViewHelper(e.Type, gs, false)
				if err != nil {
					return nil, err
				}
				edg.ViewName, err = er.ViewName()
				if err != nil {
					return nil, err
				}
			}
			r.Edges = append(r.Edges, edg)
		}
	}
	return r, nil
}

func hashNodeAndGroups(n *gen.Type, gs groups) string {
	return fmt.Sprintf("%s_%d", n.Name, gs.Hash())
}

// edgeNeedsSerialization checks if an edge is to be serialized according to its annotations and the requested groups.
func edgeNeedsSerialization(e *gen.Edge, g groups) (bool, error) {
	// If no groups are requested or the edge has no groups defined do not render the edge.
	if e.Annotations == nil || len(g) == 0 {
		return false, nil
	}
	// If there are groups given check if the groups match the requested ones.
	an := Annotation{}
	if err := an.Decode(e.Annotations[an.Name()]); err != nil {
		return false, err
	}
	// If no groups are given on the edge default is to exclude it.
	if len(an.Groups) == 0 {
		return false, nil
	}
	return g.Match(an.Groups), nil
}

// groupCombinations returns all groups ever requested together.
func groupCombinations(g *gen.Graph) ([]groups, error) {
	var gss []groups
	for _, n := range g.Nodes {
		// For every action extract the requested groups.
		for _, a := range actions {
			gs, err := groupsForAction(n, a)
			if err != nil {
				return nil, err
			}
			if !contains(gss, gs) {
				gss = append(gss, gs)
			}
		}
	}
	return gss, nil
}

func contains(haystack []groups, needle groups) bool {
	for _, gs := range haystack {
		if gs.Equal(needle) {
			return true
		}
	}
	return false
}

func stringSlice(src []interface{}) []string {
	if len(src) == 0 {
		return nil
	}
	dst := make([]string, len(src))
	for i, v := range src {
		dst[i] = v.(string)
	}
	return dst
}

// graphScope wraps the Graph object with extended scope.
type edgeToLoadScope struct {
	EdgeToLoad
	Scope map[interface{}]interface{}
}

// xextend extends the parent block with a KV pairs. Stolen from entgo.io/ent/entc/gen/func.go.
//
//	{{ with $scope := xextend $ "key" "value" }}
//		{{ template "setters" $scope }}
//	{{ end}}
//
func xextend(v interface{}, kv ...interface{}) (interface{}, error) {
	scope := make(map[interface{}]interface{})
	if len(kv)%2 != 0 {
		return nil, fmt.Errorf("invalid number of parameters: %d", len(kv))
	}
	for i := 0; i < len(kv); i += 2 {
		scope[kv[i]] = kv[i+1]
	}
	switch v := v.(type) {
	case EdgeToLoad:
		return &edgeToLoadScope{EdgeToLoad: v, Scope: scope}, nil
	case *edgeToLoadScope:
		for k := range v.Scope {
			scope[k] = v.Scope[k]
		}
		return &edgeToLoadScope{EdgeToLoad: v.EdgeToLoad, Scope: scope}, nil
	default:
		return nil, fmt.Errorf("invalid type for xextend: %T", v)
	}
}
