// Code generated by entc, DO NOT EDIT.

package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mailru/easyjson"
	"github.com/masseelch/elk/internal/simple/ent"
	"github.com/masseelch/elk/internal/simple/ent/category"
	collar "github.com/masseelch/elk/internal/simple/ent/collar"
	"github.com/masseelch/elk/internal/simple/ent/media"
	"github.com/masseelch/elk/internal/simple/ent/owner"
	"github.com/masseelch/elk/internal/simple/ent/pet"
	"go.uber.org/zap"
)

// Read fetches the ent.Category identified by a given url-parameter from the
// database and renders it to the client.
func (h *CategoryHandler) Read(w http.ResponseWriter, r *http.Request) {
	l := h.log.With(zap.String("method", "Read"))
	// ID is URL parameter.
	id64, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 0)
	if err != nil {
		l.Error("error getting id from url parameter", zap.String("id", chi.URLParam(r, "id")), zap.Error(err))
		BadRequest(w, "id must be an integer greater zero")
		return
	}
	id := uint64(id64)
	// Create the query to fetch the Category
	q := h.client.Category.Query().Where(category.ID(id))
	e, err := q.Only(r.Context())
	if err != nil {
		switch {
		case ent.IsNotFound(err):
			msg := stripEntError(err)
			l.Info(msg, zap.Error(err), zap.Uint64("id", id))
			NotFound(w, msg)
		case ent.IsNotSingular(err):
			msg := stripEntError(err)
			l.Error(msg, zap.Error(err), zap.Uint64("id", id))
			BadRequest(w, msg)
		default:
			l.Error("could not read category", zap.Error(err), zap.Uint64("id", id))
			InternalServerError(w, nil)
		}
		return
	}
	l.Info("category rendered", zap.Uint64("id", id))
	easyjson.MarshalToHTTPResponseWriter(NewCategory4094953247View(e), w)
}

// Read fetches the ent.Collar identified by a given url-parameter from the
// database and renders it to the client.
func (h *CollarHandler) Read(w http.ResponseWriter, r *http.Request) {
	l := h.log.With(zap.String("method", "Read"))
	// ID is URL parameter.
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		l.Error("error getting id from url parameter", zap.String("id", chi.URLParam(r, "id")), zap.Error(err))
		BadRequest(w, "id must be an integer")
		return
	}
	// Create the query to fetch the Collar
	q := h.client.Collar.Query().Where(collar.ID(id))
	e, err := q.Only(r.Context())
	if err != nil {
		switch {
		case ent.IsNotFound(err):
			msg := stripEntError(err)
			l.Info(msg, zap.Error(err), zap.Int("id", id))
			NotFound(w, msg)
		case ent.IsNotSingular(err):
			msg := stripEntError(err)
			l.Error(msg, zap.Error(err), zap.Int("id", id))
			BadRequest(w, msg)
		default:
			l.Error("could not read collar", zap.Error(err), zap.Int("id", id))
			InternalServerError(w, nil)
		}
		return
	}
	l.Info("collar rendered", zap.Int("id", id))
	easyjson.MarshalToHTTPResponseWriter(NewCollar1522160880View(e), w)
}

// Read fetches the ent.Media identified by a given url-parameter from the
// database and renders it to the client.
func (h *MediaHandler) Read(w http.ResponseWriter, r *http.Request) {
	l := h.log.With(zap.String("method", "Read"))
	// ID is URL parameter.
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		l.Error("error getting id from url parameter", zap.String("id", chi.URLParam(r, "id")), zap.Error(err))
		BadRequest(w, "id must be an integer")
		return
	}
	// Create the query to fetch the Media
	q := h.client.Media.Query().Where(media.ID(id))
	e, err := q.Only(r.Context())
	if err != nil {
		switch {
		case ent.IsNotFound(err):
			msg := stripEntError(err)
			l.Info(msg, zap.Error(err), zap.Int("id", id))
			NotFound(w, msg)
		case ent.IsNotSingular(err):
			msg := stripEntError(err)
			l.Error(msg, zap.Error(err), zap.Int("id", id))
			BadRequest(w, msg)
		default:
			l.Error("could not read media", zap.Error(err), zap.Int("id", id))
			InternalServerError(w, nil)
		}
		return
	}
	l.Info("media rendered", zap.Int("id", id))
	easyjson.MarshalToHTTPResponseWriter(NewMedia1941033838View(e), w)
}

// Read fetches the ent.Owner identified by a given url-parameter from the
// database and renders it to the client.
func (h *OwnerHandler) Read(w http.ResponseWriter, r *http.Request) {
	l := h.log.With(zap.String("method", "Read"))
	// ID is URL parameter.
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		l.Error("error getting id from url parameter", zap.String("id", chi.URLParam(r, "id")), zap.Error(err))
		BadRequest(w, "id must be a valid UUID")
		return
	}
	// Create the query to fetch the Owner
	q := h.client.Owner.Query().Where(owner.ID(uuid.UUID(id)))
	e, err := q.Only(r.Context())
	if err != nil {
		switch {
		case ent.IsNotFound(err):
			msg := stripEntError(err)
			l.Info(msg, zap.Error(err), zap.String("id", id.String()))
			NotFound(w, msg)
		case ent.IsNotSingular(err):
			msg := stripEntError(err)
			l.Error(msg, zap.Error(err), zap.String("id", id.String()))
			BadRequest(w, msg)
		default:
			l.Error("could not read owner", zap.Error(err), zap.String("id", id.String()))
			InternalServerError(w, nil)
		}
		return
	}
	l.Info("owner rendered", zap.String("id", id.String()))
	easyjson.MarshalToHTTPResponseWriter(NewOwner139708381View(e), w)
}

// Read fetches the ent.Pet identified by a given url-parameter from the
// database and renders it to the client.
func (h *PetHandler) Read(w http.ResponseWriter, r *http.Request) {
	l := h.log.With(zap.String("method", "Read"))
	// ID is URL parameter.
	var err error
	id := chi.URLParam(r, "id")
	// Create the query to fetch the Pet
	q := h.client.Pet.Query().Where(pet.ID(id))
	// Eager load edges that are required on read operation.
	q.WithOwner().WithFriends(func(q *ent.PetQuery) {
		q.WithOwner().WithFriends(func(q *ent.PetQuery) {
			q.WithOwner().WithFriends(func(q *ent.PetQuery) {
				q.WithOwner()
			})
		})
	})
	e, err := q.Only(r.Context())
	if err != nil {
		switch {
		case ent.IsNotFound(err):
			msg := stripEntError(err)
			l.Info(msg, zap.Error(err), zap.String("id", id))
			NotFound(w, msg)
		case ent.IsNotSingular(err):
			msg := stripEntError(err)
			l.Error(msg, zap.Error(err), zap.String("id", id))
			BadRequest(w, msg)
		default:
			l.Error("could not read pet", zap.Error(err), zap.String("id", id))
			InternalServerError(w, nil)
		}
		return
	}
	l.Info("pet rendered", zap.String("id", id))
	easyjson.MarshalToHTTPResponseWriter(NewPet1876743790View(e), w)
}
