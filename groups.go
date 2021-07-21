package elk

import "strings"

// Groups are used to determine what properties to load and serialize.
type Groups []string

// Add adds a group to the Groups. If the group is already present it does nothing.
func (gs *Groups) Add(g ...string) {
	for _, g1 := range g {
		if !gs.HasGroup(g1) {
			*gs = append(*gs, g1)
		}
	}
}

// HasGroup checks if the given group is present.
func (gs Groups) HasGroup(g string) bool {
	for _, e := range gs {
		if e == g {
			return true
		}
	}

	return false
}

// Match check if at least one of the given Groups is present.
func (gs Groups) Match(other Groups) bool {
	for _, g := range other {
		if gs.HasGroup(g) {
			return true
		}
	}

	return false
}

// StructTag returns the struct tag representation of the Groups.
func (gs Groups) StructTag() string {
	return `groups:"` + strings.Join(gs, ",") + `"`
}