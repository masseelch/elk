// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/masseelch/elk/internal/pets/ent/badge"
	"github.com/masseelch/elk/internal/pets/ent/pet"
	"github.com/masseelch/elk/internal/pets/ent/predicate"
)

// BadgeUpdate is the builder for updating Badge entities.
type BadgeUpdate struct {
	config
	hooks    []Hook
	mutation *BadgeMutation
}

// Where appends a list predicates to the BadgeUpdate builder.
func (bu *BadgeUpdate) Where(ps ...predicate.Badge) *BadgeUpdate {
	bu.mutation.Where(ps...)
	return bu
}

// SetColor sets the "color" field.
func (bu *BadgeUpdate) SetColor(b badge.Color) *BadgeUpdate {
	bu.mutation.SetColor(b)
	return bu
}

// SetMaterial sets the "material" field.
func (bu *BadgeUpdate) SetMaterial(b badge.Material) *BadgeUpdate {
	bu.mutation.SetMaterial(b)
	return bu
}

// SetWearerID sets the "wearer" edge to the Pet entity by ID.
func (bu *BadgeUpdate) SetWearerID(id int) *BadgeUpdate {
	bu.mutation.SetWearerID(id)
	return bu
}

// SetNillableWearerID sets the "wearer" edge to the Pet entity by ID if the given value is not nil.
func (bu *BadgeUpdate) SetNillableWearerID(id *int) *BadgeUpdate {
	if id != nil {
		bu = bu.SetWearerID(*id)
	}
	return bu
}

// SetWearer sets the "wearer" edge to the Pet entity.
func (bu *BadgeUpdate) SetWearer(p *Pet) *BadgeUpdate {
	return bu.SetWearerID(p.ID)
}

// Mutation returns the BadgeMutation object of the builder.
func (bu *BadgeUpdate) Mutation() *BadgeMutation {
	return bu.mutation
}

// ClearWearer clears the "wearer" edge to the Pet entity.
func (bu *BadgeUpdate) ClearWearer() *BadgeUpdate {
	bu.mutation.ClearWearer()
	return bu
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (bu *BadgeUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(bu.hooks) == 0 {
		if err = bu.check(); err != nil {
			return 0, err
		}
		affected, err = bu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*BadgeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = bu.check(); err != nil {
				return 0, err
			}
			bu.mutation = mutation
			affected, err = bu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(bu.hooks) - 1; i >= 0; i-- {
			if bu.hooks[i] == nil {
				return 0, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = bu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, bu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (bu *BadgeUpdate) SaveX(ctx context.Context) int {
	affected, err := bu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (bu *BadgeUpdate) Exec(ctx context.Context) error {
	_, err := bu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (bu *BadgeUpdate) ExecX(ctx context.Context) {
	if err := bu.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (bu *BadgeUpdate) check() error {
	if v, ok := bu.mutation.Color(); ok {
		if err := badge.ColorValidator(v); err != nil {
			return &ValidationError{Name: "color", err: fmt.Errorf("ent: validator failed for field \"color\": %w", err)}
		}
	}
	if v, ok := bu.mutation.Material(); ok {
		if err := badge.MaterialValidator(v); err != nil {
			return &ValidationError{Name: "material", err: fmt.Errorf("ent: validator failed for field \"material\": %w", err)}
		}
	}
	return nil
}

func (bu *BadgeUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   badge.Table,
			Columns: badge.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: badge.FieldID,
			},
		},
	}
	if ps := bu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := bu.mutation.Color(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: badge.FieldColor,
		})
	}
	if value, ok := bu.mutation.Material(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: badge.FieldMaterial,
		})
	}
	if bu.mutation.WearerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   badge.WearerTable,
			Columns: []string{badge.WearerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: pet.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := bu.mutation.WearerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   badge.WearerTable,
			Columns: []string{badge.WearerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: pet.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, bu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{badge.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return 0, err
	}
	return n, nil
}

// BadgeUpdateOne is the builder for updating a single Badge entity.
type BadgeUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *BadgeMutation
}

// SetColor sets the "color" field.
func (buo *BadgeUpdateOne) SetColor(b badge.Color) *BadgeUpdateOne {
	buo.mutation.SetColor(b)
	return buo
}

// SetMaterial sets the "material" field.
func (buo *BadgeUpdateOne) SetMaterial(b badge.Material) *BadgeUpdateOne {
	buo.mutation.SetMaterial(b)
	return buo
}

// SetWearerID sets the "wearer" edge to the Pet entity by ID.
func (buo *BadgeUpdateOne) SetWearerID(id int) *BadgeUpdateOne {
	buo.mutation.SetWearerID(id)
	return buo
}

// SetNillableWearerID sets the "wearer" edge to the Pet entity by ID if the given value is not nil.
func (buo *BadgeUpdateOne) SetNillableWearerID(id *int) *BadgeUpdateOne {
	if id != nil {
		buo = buo.SetWearerID(*id)
	}
	return buo
}

// SetWearer sets the "wearer" edge to the Pet entity.
func (buo *BadgeUpdateOne) SetWearer(p *Pet) *BadgeUpdateOne {
	return buo.SetWearerID(p.ID)
}

// Mutation returns the BadgeMutation object of the builder.
func (buo *BadgeUpdateOne) Mutation() *BadgeMutation {
	return buo.mutation
}

// ClearWearer clears the "wearer" edge to the Pet entity.
func (buo *BadgeUpdateOne) ClearWearer() *BadgeUpdateOne {
	buo.mutation.ClearWearer()
	return buo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (buo *BadgeUpdateOne) Select(field string, fields ...string) *BadgeUpdateOne {
	buo.fields = append([]string{field}, fields...)
	return buo
}

// Save executes the query and returns the updated Badge entity.
func (buo *BadgeUpdateOne) Save(ctx context.Context) (*Badge, error) {
	var (
		err  error
		node *Badge
	)
	if len(buo.hooks) == 0 {
		if err = buo.check(); err != nil {
			return nil, err
		}
		node, err = buo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*BadgeMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = buo.check(); err != nil {
				return nil, err
			}
			buo.mutation = mutation
			node, err = buo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(buo.hooks) - 1; i >= 0; i-- {
			if buo.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = buo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, buo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (buo *BadgeUpdateOne) SaveX(ctx context.Context) *Badge {
	node, err := buo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (buo *BadgeUpdateOne) Exec(ctx context.Context) error {
	_, err := buo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (buo *BadgeUpdateOne) ExecX(ctx context.Context) {
	if err := buo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (buo *BadgeUpdateOne) check() error {
	if v, ok := buo.mutation.Color(); ok {
		if err := badge.ColorValidator(v); err != nil {
			return &ValidationError{Name: "color", err: fmt.Errorf("ent: validator failed for field \"color\": %w", err)}
		}
	}
	if v, ok := buo.mutation.Material(); ok {
		if err := badge.MaterialValidator(v); err != nil {
			return &ValidationError{Name: "material", err: fmt.Errorf("ent: validator failed for field \"material\": %w", err)}
		}
	}
	return nil
}

func (buo *BadgeUpdateOne) sqlSave(ctx context.Context) (_node *Badge, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   badge.Table,
			Columns: badge.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: badge.FieldID,
			},
		},
	}
	id, ok := buo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing Badge.ID for update")}
	}
	_spec.Node.ID.Value = id
	if fields := buo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, badge.FieldID)
		for _, f := range fields {
			if !badge.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != badge.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := buo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := buo.mutation.Color(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: badge.FieldColor,
		})
	}
	if value, ok := buo.mutation.Material(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeEnum,
			Value:  value,
			Column: badge.FieldMaterial,
		})
	}
	if buo.mutation.WearerCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   badge.WearerTable,
			Columns: []string{badge.WearerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: pet.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := buo.mutation.WearerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: true,
			Table:   badge.WearerTable,
			Columns: []string{badge.WearerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeInt,
					Column: pet.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Badge{config: buo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, buo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{badge.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{err.Error(), err}
		}
		return nil, err
	}
	return _node, nil
}