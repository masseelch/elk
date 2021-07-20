package integration

import (
	"context"
	"github.com/Pallinder/go-randomdata"
	"github.com/masseelch/elk/internal/integration/pets/ent"
	"math/rand"
	"time"
)

const (
	_ = iota
	categoryKey
	petKey
	ownerKey

	petCount      = 50
	ownerCount    = 10
	categoryCount = 5
)

type refs map[uint]interface{}
type fixtureFn func(ctx context.Context, refs refs, c *ent.Client) error

func fixtures(ctx context.Context, c *ent.Client) error {
	rand.Seed(time.Now().Unix())
	refs := make(refs)

	for _, fn := range []fixtureFn{owners, categories, pets} {
		if err := fn(ctx, refs, c); err != nil {
			return err
		}
	}

	return nil
}

// category fixtures
func (r refs) category() *ent.Category {
	m := r[categoryKey].([]*ent.Category)
	return m[rand.Intn(len(m))]
}
func (r refs) categories(c int) []*ent.Category {
	ret := make([]*ent.Category, c)

	for i := range ret {
		ret[i] = r.category()
	}

	return ret
}

func categories(ctx context.Context, refs refs, c *ent.Client) error {
	var err error
	b := make([]*ent.CategoryCreate, categoryCount)

	for i := 0; i < len(b); i++ {
		b[i] = c.Category.Create().SetName(randomdata.Noun())
	}

	refs[categoryKey], err = c.Category.CreateBulk(b...).Save(ctx)
	return err
}

// pet fixtures
func (r refs) pet() *ent.Pet {
	m := r[petKey].([]*ent.Pet)
	return m[rand.Intn(len(m))]
}

func pets(ctx context.Context, refs refs, c *ent.Client) error {
	var err error
	b := make([]*ent.PetCreate, petCount)

	for i := 0; i < len(b); i++ {
		b[i] = c.Pet.Create().SetName(randomdata.Noun()).SetAge(randomdata.Number(1)).SetOwner(refs.owner()).AddCategory(refs.categories(randomdata.Number(1, 3))...)
	}

	refs[petKey], err = c.Pet.CreateBulk(b...).Save(ctx)
	if err != nil {
		return err
	}

	ps := refs[petKey].([]*ent.Pet)
	for i, p := range ps {
		if i < 4 {
			continue
		}
		q := c.Pet.UpdateOne(p).AddFriends(ps[i-1], ps[i-2])
		if randomdata.Boolean() {
			q.AddFriends(ps[i-3], ps[i-4])
		}
		if err := q.Exec(ctx); err != nil {
			return err
		}
	}

	return nil
}

// owner fixtures
func (r refs) owner() *ent.Owner {
	m := r[ownerKey].([]*ent.Owner)
	return m[rand.Intn(len(m))]
}

func owners(ctx context.Context, refs refs, c *ent.Client) error {
	var err error
	b := make([]*ent.OwnerCreate, ownerCount)

	for i := 0; i < len(b); i++ {
		b[i] = c.Owner.Create().SetName(randomdata.Noun()).SetAge(randomdata.Number(1))
	}

	refs[ownerKey], err = c.Owner.CreateBulk(b...).Save(ctx)
	return err
}
