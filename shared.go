package elk

import (
	"entgo.io/ent/entc/gen"
	"errors"
	"fmt"
	"github.com/masseelch/elk/serialization"
)

const (
	opCreate = "create"
	opRead   = "read"
	opUpdate = "update"
	opDelete = "delete"
	opList   = "list"
)

// groupsForOperation returns the requested groups for a given type and operation.
func groupsForOperation(n *gen.Type, o string) (serialization.Groups, error) {
	// If there are no annotations given do not load any groups.
	ant := &SchemaAnnotation{}
	if n.Annotations == nil || n.Annotations[ant.Name()] == nil {
		return nil, nil
	}
	// Decode the types annotation and extract the groups requested for the given operation.
	if err := ant.Decode(n.Annotations[ant.Name()]); err != nil {
		return nil, err
	}
	switch o {
	case opCreate:
		return ant.CreateGroups, nil
	case opRead:
		return ant.ReadGroups, nil
	case opUpdate:
		return ant.UpdateGroups, nil
	case opList:
		return ant.ListGroups, nil
	}
	return nil, fmt.Errorf("unknown operation %q", o)
}

// groups returns the groups set on elk.Annotation.
func groups(a gen.Annotations) (serialization.Groups, error) {
	an := Annotation{}
	if err := an.Decode(a[an.Name()]); err != nil {
		return nil, err
	}
	return an.Groups, nil
}

// config loads the elk extension Config out of gen.Config struct.
func config(cfg *gen.Config) (*Config, error) {
	c := &Config{}
	if cfg == nil || cfg.Annotations == nil || cfg.Annotations[c.Name()] == nil {
		return nil, errors.New("elk extension config not found")
	}
	return c, c.Decode(cfg.Annotations[c.Name()])
}
