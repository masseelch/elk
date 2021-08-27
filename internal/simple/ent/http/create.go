// Code generated by entc, DO NOT EDIT.

package http

import (
	"net/http"

	"github.com/mailru/easyjson"
	"github.com/masseelch/elk/internal/simple/ent"
	"github.com/masseelch/elk/internal/simple/ent/category"
	"github.com/masseelch/elk/internal/simple/ent/owner"
	"github.com/masseelch/elk/internal/simple/ent/pet"
	"go.uber.org/zap"
)

// Create creates a new ent.Category and stores it in the database.
func (h CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	l := h.log.With(zap.String("method", "Create"))
	// Get the post data.
	var d CategoryCreateRequest
	if err := easyjson.UnmarshalFromReader(r.Body, &d); err != nil {
		l.Error("error decoding json", zap.Error(err))
		BadRequest(w, "invalid json string")
		return
	}
	// Save the data.
	b := h.client.Category.Create()
	if d.Name != nil {
		b.SetName(*d.Name)
	}
	if d.Pets != nil {
		b.AddPetIDs(d.Pets...)
	}
	e, err := b.Save(r.Context())
	if err != nil {
		switch {
		default:
			l.Error("could not create category", zap.Error(err))
			InternalServerError(w, nil)
		}
		return
	}
	// Reload entry.
	q := h.client.Category.Query().Where(category.ID(e.ID))
	e, err = q.Only(r.Context())
	if err != nil {
		switch {
		case ent.IsNotFound(err):
			msg := stripEntError(err)
			l.Info(msg, zap.Error(err), zap.Int("id", e.ID))
			NotFound(w, msg)
		case ent.IsNotSingular(err):
			msg := stripEntError(err)
			l.Error(msg, zap.Error(err), zap.Int("id", e.ID))
			BadRequest(w, msg)
		default:
			l.Error("could not read category", zap.Error(err), zap.Int("id", e.ID))
			InternalServerError(w, nil)
		}
		return
	}
	l.Info("category rendered", zap.Int("id", e.ID))
	easyjson.MarshalToHTTPResponseWriter(NewCategory4094953247View(e), w)
}

// Create creates a new ent.Owner and stores it in the database.
func (h OwnerHandler) Create(w http.ResponseWriter, r *http.Request) {
	l := h.log.With(zap.String("method", "Create"))
	// Get the post data.
	var d OwnerCreateRequest
	if err := easyjson.UnmarshalFromReader(r.Body, &d); err != nil {
		l.Error("error decoding json", zap.Error(err))
		BadRequest(w, "invalid json string")
		return
	}
	// Save the data.
	b := h.client.Owner.Create()
	if d.Name != nil {
		b.SetName(*d.Name)
	}
	if d.Age != nil {
		b.SetAge(*d.Age)
	}
	if d.Pets != nil {
		b.AddPetIDs(d.Pets...)
	}
	e, err := b.Save(r.Context())
	if err != nil {
		switch {
		default:
			l.Error("could not create owner", zap.Error(err))
			InternalServerError(w, nil)
		}
		return
	}
	// Reload entry.
	q := h.client.Owner.Query().Where(owner.ID(e.ID))
	e, err = q.Only(r.Context())
	if err != nil {
		switch {
		case ent.IsNotFound(err):
			msg := stripEntError(err)
			l.Info(msg, zap.Error(err), zap.Int("id", e.ID))
			NotFound(w, msg)
		case ent.IsNotSingular(err):
			msg := stripEntError(err)
			l.Error(msg, zap.Error(err), zap.Int("id", e.ID))
			BadRequest(w, msg)
		default:
			l.Error("could not read owner", zap.Error(err), zap.Int("id", e.ID))
			InternalServerError(w, nil)
		}
		return
	}
	l.Info("owner rendered", zap.Int("id", e.ID))
	easyjson.MarshalToHTTPResponseWriter(NewOwner139708381View(e), w)
}

// Create creates a new ent.Pet and stores it in the database.
func (h PetHandler) Create(w http.ResponseWriter, r *http.Request) {
	l := h.log.With(zap.String("method", "Create"))
	// Get the post data.
	var d PetCreateRequest
	if err := easyjson.UnmarshalFromReader(r.Body, &d); err != nil {
		l.Error("error decoding json", zap.Error(err))
		BadRequest(w, "invalid json string")
		return
	}
	// Save the data.
	b := h.client.Pet.Create()
	if d.Name != nil {
		b.SetName(*d.Name)
	}
	if d.Age != nil {
		b.SetAge(*d.Age)
	}
	if d.Category != nil {
		b.AddCategoryIDs(d.Category...)
	}
	if d.Owner != nil {
		b.SetOwnerID(*d.Owner)
	}
	if d.Friends != nil {
		b.AddFriendIDs(d.Friends...)
	}
	e, err := b.Save(r.Context())
	if err != nil {
		switch {
		default:
			l.Error("could not create pet", zap.Error(err))
			InternalServerError(w, nil)
		}
		return
	}
	// Reload entry.
	q := h.client.Pet.Query().Where(pet.ID(e.ID))
	e, err = q.Only(r.Context())
	if err != nil {
		switch {
		case ent.IsNotFound(err):
			msg := stripEntError(err)
			l.Info(msg, zap.Error(err), zap.Int("id", e.ID))
			NotFound(w, msg)
		case ent.IsNotSingular(err):
			msg := stripEntError(err)
			l.Error(msg, zap.Error(err), zap.Int("id", e.ID))
			BadRequest(w, msg)
		default:
			l.Error("could not read pet", zap.Error(err), zap.Int("id", e.ID))
			InternalServerError(w, nil)
		}
		return
	}
	l.Info("pet rendered", zap.Int("id", e.ID))
	easyjson.MarshalToHTTPResponseWriter(NewPet359800019View(e), w)
}