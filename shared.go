package elk

import (
	"entgo.io/ent/entc/gen"
	"fmt"
	"github.com/masseelch/elk/serialization"
)

const (
	actionCreate = "create"
	actionRead   = "read"
	actionUpdate = "update"
	actionList   = "list"
)

// groupsForAction returns the requested groups for a given type and action.
func groupsForAction(n *gen.Type, a string) (serialization.Groups, error) {
	// If there are no annotations given do not load any groups.
	ant := &SchemaAnnotation{}
	if n.Annotations == nil || n.Annotations[ant.Name()] == nil {
		return nil, nil
	}
	// Decode the types annotation and extract the groups requested for the given action.
	if err := ant.Decode(n.Annotations[ant.Name()]); err != nil {
		return nil, err
	}
	switch a {
	case actionCreate:
		return ant.CreateGroups, nil
	case actionRead:
		return ant.ReadGroups, nil
	case actionUpdate:
		return ant.UpdateGroups, nil
	case actionList:
		return ant.ListGroups, nil
	}
	return nil, fmt.Errorf("unknown action %q", a)
}

// groups returns the groups set on elk.Annotation.
func groups(a gen.Annotations) (serialization.Groups, error) {
	an := Annotation{}
	if err := an.Decode(a[an.Name()]); err != nil {
		return nil, err
	}
	return an.Groups, nil
}
