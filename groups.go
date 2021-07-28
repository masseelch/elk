package elk

import "strings"

// groups are used to determine what properties to load and serialize.
type groups []string

// Add adds a group to the groups. If the group is already present it does nothing.
func (gs *groups) Add(g ...string) {
	for _, g1 := range g {
		if !gs.HasGroup(g1) {
			*gs = append(*gs, g1)
		}
	}
}

// HasGroup checks if the given group is present.
func (gs groups) HasGroup(g string) bool {
	for _, e := range gs {
		if e == g {
			return true
		}
	}

	return false
}

// Match check if at least one of the given groups is present.
func (gs groups) Match(other groups) bool {
	for _, g := range other {
		if gs.HasGroup(g) {
			return true
		}
	}

	return false
}

// StructTag returns the struct tag representation of the groups.
func (gs groups) StructTag() string {
	return `groups:"` + strings.Join(gs, ",") + `"`
}
