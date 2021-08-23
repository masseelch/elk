package integration

import (
	"context"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/masseelch/elk/internal/integration/pets/ent"
	"github.com/masseelch/elk/internal/integration/pets/ent/badge"
	"github.com/masseelch/elk/internal/integration/pets/ent/pet"
	"github.com/masseelch/elk/internal/integration/pets/ent/playgroup"
	"github.com/masseelch/elk/internal/integration/pets/ent/toy"
	"math/rand"
	"time"
)

const (
	_ = iota
	petKey
	playGroupKey
	toyKey

	petCount       = 50
	playGroupCount = 5
	toyCount       = 20
)

type refs map[uint]interface{}
type fixtureFn func(ctx context.Context, refs refs, c *ent.Client) error

func fixtures(ctx context.Context, c *ent.Client) error {
	refs := make(refs)
	for _, fn := range []fixtureFn{playGroups, pets, toys} {
		if err := fn(ctx, refs, c); err != nil {
			return err
		}
	}
	return nil
}

// pet fixtures
var badgeColors = []badge.Color{badge.ColorRed, badge.ColorOrange, badge.ColorYellow, badge.ColorGreen, badge.ColorBlue, badge.ColorIndigo, badge.ColorViolet, badge.ColorPurple, badge.ColorPink, badge.ColorSilver, badge.ColorGold, badge.ColorBeige, badge.ColorBrown, badge.ColorGrey, badge.ColorBlack, badge.ColorWhite}
var badgeMaterials = []badge.Material{badge.MaterialLeather, badge.MaterialPlastic, badge.MaterialFabric}

func (r refs) pet() *ent.Pet {
	m := r[petKey].([]*ent.Pet)
	return m[rand.Intn(len(m))]
}

func pets(ctx context.Context, refs refs, c *ent.Client) error {
	var err error
	var bday time.Time
	var sex pet.Sex
	pb := make([]*ent.PetCreate, petCount)
	bb := make([]*ent.BadgeCreate, petCount)
	for i := 0; i < len(pb); i++ {
		bb[i] = c.Badge.Create().
			SetColor(badgeColors[randomdata.Number(len(badgeColors))]).
			SetMaterial(badgeMaterials[randomdata.Number(len(badgeMaterials))])
	}
	bs, err := c.Badge.CreateBulk(bb...).Save(ctx)
	if err != nil {
		return err
	}
	for i := 0; i < len(pb); i++ {
		bday, err = time.Parse(randomdata.DateOutputLayout, randomdata.FullDate())
		if err != nil {
			return err
		}
		var nns []string
		for i := 0; i < randomdata.Number(5); i++ {
			nns = append(nns, randomdata.Adjective())
		}
		if randomdata.Boolean() {
			sex = pet.SexMale
		} else {
			sex = pet.SexFemale
		}
		pb[i] = c.Pet.Create().
			SetHeight(randomdata.Number(40, 200)).
			SetWeight(randomdata.Decimal(1, 500)).
			SetCastrated(randomdata.Boolean()).
			SetName(fmt.Sprintf("%s%d", randomdata.Noun(), i)).
			SetBirthday(bday).
			SetNicknames(nns).
			SetSex(sex).
			SetBadge(bs[i])
		// playgroups
		pg := refs.playGroups(randomdata.Number(1, 3))
		if len(pg) > 0 {
			pb[i].AddPlayGroups(pg...)
		}
	}
	refs[petKey], err = c.Pet.CreateBulk(pb...).Save(ctx)
	if err != nil {
		return err
	}
	// friends
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
	// mentor
	if err := c.Pet.UpdateOneID(1).SetMentorID(2).Exec(ctx); err != nil {
		return err
	}
	// spouse
	if err := c.Pet.UpdateOneID(1).SetSpouseID(2).Exec(ctx); err != nil {
		return err
	}
	// children
	if err := c.Pet.UpdateOneID(1).AddChildIDs(2).Exec(ctx); err != nil {
		return err
	}
	return nil
}

// toy fixtures
var toyColors = []toy.Color{toy.ColorRed, toy.ColorOrange, toy.ColorYellow, toy.ColorGreen, toy.ColorBlue, toy.ColorIndigo, toy.ColorViolet, toy.ColorPurple, toy.ColorPink, toy.ColorSilver, toy.ColorGold, toy.ColorBeige, toy.ColorBrown, toy.ColorGrey, toy.ColorBlack, toy.ColorWhite}
var toyMaterials = []toy.Material{toy.MaterialLeather, toy.MaterialPlastic, toy.MaterialFabric}

func (r refs) toy() *ent.Toy {
	m := r[toyKey].([]*ent.Toy)
	return m[rand.Intn(len(m))]
}

func (r refs) toys(c int) []*ent.Toy {
	var l []*ent.Toy
	for i := 0; i <= c; i++ {
		l = append(l, r.toy())
	}
	return l
}

func toys(ctx context.Context, refs refs, c *ent.Client) error {
	var err error
	b := make([]*ent.ToyCreate, toyCount)
	for i := 0; i < len(b); i++ {
		b[i] = c.Toy.Create().
			SetTitle(randomdata.SillyName()).
			SetColor(toyColors[randomdata.Number(len(toyColors))]).
			SetMaterial(toyMaterials[randomdata.Number(len(toyMaterials))]).
			SetOwner(refs.pet())
	}
	refs[toyKey], err = c.Toy.CreateBulk(b...).Save(ctx)
	return err
}

// playGroups fixtures
var playGroupWeekdays = []playgroup.Weekday{playgroup.WeekdayMon, playgroup.WeekdayTue, playgroup.WeekdayWed, playgroup.WeekdayThu, playgroup.WeekdayFri, playgroup.WeekdaySat, playgroup.WeekdaySun}

func (r refs) playGroup() *ent.PlayGroup {
	m := r[playGroupKey].([]*ent.PlayGroup)
	return m[rand.Intn(len(m))]
}

func (r refs) playGroups(c int) []*ent.PlayGroup {
	var l []*ent.PlayGroup
	for i := 0; i <= c; i++ {
		l = append(l, r.playGroup())
	}
	return l
}

func playGroups(ctx context.Context, refs refs, c *ent.Client) error {
	var err error
	b := make([]*ent.PlayGroupCreate, playGroupCount)
	for i := 0; i < len(b); i++ {
		b[i] = c.PlayGroup.Create().
			SetTitle(randomdata.SillyName()).
			SetDescription(randomdata.Paragraph()).
			SetWeekday(playGroupWeekdays[randomdata.Number(len(playGroupWeekdays))])
	}
	refs[playGroupKey], err = c.PlayGroup.CreateBulk(b...).Save(ctx)
	return err
}
