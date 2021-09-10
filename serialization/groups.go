package serialization

import (
	"hash/fnv"
	"sort"
	"strconv"
)

type (
	// Groups are used to determine what properties to load and serialize.
	Groups []string
	// Collection is a slice of multiple Groups.
	Collection []Groups
)

// Add adds a group to the groups. If the group is already present it does nothing.
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

// Equal reports if two Groups have the same entries.
func (gs Groups) Equal(other Groups) bool {
	if len(gs) != len(other) {
		return false
	}
	for _, g := range other {
		if !gs.HasGroup(g) {
			return false
		}
	}
	return true
}

// Hash returns a hash value for a Groups struct.
func (gs Groups) Hash() uint32 {
	sort.Strings(gs)
	h := fnv.New32a()
	for _, g := range gs {
		h.Write([]byte(g))
	}
	h.Write([]byte(strconv.Itoa(len(gs))))
	return h.Sum32()
}

// Contains checks if a collection of groups contains the given Groups.
func (c Collection) Contains(needle Groups) bool {
	for _, gs := range c {
		if gs.Equal(needle) {
			return true
		}
	}
	return false
}
