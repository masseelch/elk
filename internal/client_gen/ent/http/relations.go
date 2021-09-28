// Code generated by entc, DO NOT EDIT.

package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mailru/easyjson"
	"github.com/masseelch/elk/internal/client_gen/ent"
	"github.com/masseelch/elk/internal/client_gen/ent/category"
	collar "github.com/masseelch/elk/internal/client_gen/ent/collar"
	"github.com/masseelch/elk/internal/client_gen/ent/owner"
	"github.com/masseelch/elk/internal/client_gen/ent/pet"
	"go.uber.org/zap"
)

// Pets fetches the ent.pets attached to the ent.Category
// identified by a given url-parameter from the database and renders it to the client.
func (h CategoryHandler) Pets(w http.ResponseWriter, r *http.Request) {
	l := h.log.With(zap.String("method", "Pets"))
	// ID is URL parameter.
	id64, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 0)
	if err != nil {
		l.Error("error getting id from url parameter", zap.String("id", chi.URLParam(r, "id")), zap.Error(err))
		BadRequest(w, "id must be an integer greater zero")
		return
	}
	id := uint64(id64)
	// Create the query to fetch the pets attached to this category
	q := h.client.Category.Query().Where(category.ID(id)).QueryPets()
	page := 1
	if d := r.URL.Query().Get("page"); d != "" {
		page, err = strconv.Atoi(d)
		if err != nil {
			l.Info("error parsing query parameter 'page'", zap.String("page", d), zap.Error(err))
			BadRequest(w, "page must be an integer greater zero")
			return
		}
	}
	itemsPerPage := 30
	if d := r.URL.Query().Get("itemsPerPage"); d != "" {
		itemsPerPage, err = strconv.Atoi(d)
		if err != nil {
			l.Info("error parsing query parameter 'itemsPerPage'", zap.String("itemsPerPage", d), zap.Error(err))
			BadRequest(w, "itemsPerPage must be an integer greater zero")
			return
		}
	}
	es, err := q.Limit(itemsPerPage).Offset((page - 1) * itemsPerPage).All(r.Context())
	if err != nil {
		l.Error("error fetching pets from db", zap.Error(err))
		InternalServerError(w, nil)
		return
	}
	l.Info("pets rendered", zap.Int("amount", len(es)))
	easyjson.MarshalToHTTPResponseWriter(NewPetViews(es), w)
}

// Pet fetches the ent.pet attached to the ent.Collar
// identified by a given url-parameter from the database and renders it to the client.
func (h CollarHandler) Pet(w http.ResponseWriter, r *http.Request) {
	l := h.log.With(zap.String("method", "Pet"))
	// ID is URL parameter.
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		l.Error("error getting id from url parameter", zap.String("id", chi.URLParam(r, "id")), zap.Error(err))
		BadRequest(w, "id must be an integer")
		return
	}
	// Create the query to fetch the pet attached to this collar
	q := h.client.Collar.Query().Where(collar.ID(id)).QueryPet()
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
			l.Info(msg, zap.Error(err), zap.Int("id", id))
			NotFound(w, msg)
		case ent.IsNotSingular(err):
			msg := stripEntError(err)
			l.Error(msg, zap.Error(err), zap.Int("id", id))
			BadRequest(w, msg)
		default:
			l.Error("could-not-read-collar", zap.Error(err), zap.Int("id", id))
			InternalServerError(w, nil)
		}
		return
	}
	l.Info("pet rendered", zap.String("id", e.ID))
	easyjson.MarshalToHTTPResponseWriter(NewPetWithOwnerAndPetOwnerView(e), w)
}

// Pets fetches the ent.pets attached to the ent.Owner
// identified by a given url-parameter from the database and renders it to the client.
func (h OwnerHandler) Pets(w http.ResponseWriter, r *http.Request) {
	l := h.log.With(zap.String("method", "Pets"))
	// ID is URL parameter.
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		l.Error("error getting id from url parameter", zap.String("id", chi.URLParam(r, "id")), zap.Error(err))
		BadRequest(w, "id must be a valid UUID")
		return
	}
	// Create the query to fetch the pets attached to this owner
	q := h.client.Owner.Query().Where(owner.ID(id)).QueryPets()
	page := 1
	if d := r.URL.Query().Get("page"); d != "" {
		page, err = strconv.Atoi(d)
		if err != nil {
			l.Info("error parsing query parameter 'page'", zap.String("page", d), zap.Error(err))
			BadRequest(w, "page must be an integer greater zero")
			return
		}
	}
	itemsPerPage := 30
	if d := r.URL.Query().Get("itemsPerPage"); d != "" {
		itemsPerPage, err = strconv.Atoi(d)
		if err != nil {
			l.Info("error parsing query parameter 'itemsPerPage'", zap.String("itemsPerPage", d), zap.Error(err))
			BadRequest(w, "itemsPerPage must be an integer greater zero")
			return
		}
	}
	es, err := q.Limit(itemsPerPage).Offset((page - 1) * itemsPerPage).All(r.Context())
	if err != nil {
		l.Error("error fetching pets from db", zap.Error(err))
		InternalServerError(w, nil)
		return
	}
	l.Info("pets rendered", zap.Int("amount", len(es)))
	easyjson.MarshalToHTTPResponseWriter(NewPetViews(es), w)
}

// Collar fetches the ent.collar attached to the ent.Pet
// identified by a given url-parameter from the database and renders it to the client.
func (h PetHandler) Collar(w http.ResponseWriter, r *http.Request) {
	l := h.log.With(zap.String("method", "Collar"))
	// ID is URL parameter.
	var err error
	id := chi.URLParam(r, "id")
	// Create the query to fetch the collar attached to this pet
	q := h.client.Pet.Query().Where(pet.ID(id)).QueryCollar()
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
			l.Error("could-not-read-pet", zap.Error(err), zap.String("id", id))
			InternalServerError(w, nil)
		}
		return
	}
	l.Info("collar rendered", zap.Int("id", e.ID))
	easyjson.MarshalToHTTPResponseWriter(NewCollarView(e), w)
}

// Categories fetches the ent.categories attached to the ent.Pet
// identified by a given url-parameter from the database and renders it to the client.
func (h PetHandler) Categories(w http.ResponseWriter, r *http.Request) {
	l := h.log.With(zap.String("method", "Categories"))
	// ID is URL parameter.
	var err error
	id := chi.URLParam(r, "id")
	// Create the query to fetch the categories attached to this pet
	q := h.client.Pet.Query().Where(pet.ID(id)).QueryCategories()
	page := 1
	if d := r.URL.Query().Get("page"); d != "" {
		page, err = strconv.Atoi(d)
		if err != nil {
			l.Info("error parsing query parameter 'page'", zap.String("page", d), zap.Error(err))
			BadRequest(w, "page must be an integer greater zero")
			return
		}
	}
	itemsPerPage := 30
	if d := r.URL.Query().Get("itemsPerPage"); d != "" {
		itemsPerPage, err = strconv.Atoi(d)
		if err != nil {
			l.Info("error parsing query parameter 'itemsPerPage'", zap.String("itemsPerPage", d), zap.Error(err))
			BadRequest(w, "itemsPerPage must be an integer greater zero")
			return
		}
	}
	es, err := q.Limit(itemsPerPage).Offset((page - 1) * itemsPerPage).All(r.Context())
	if err != nil {
		l.Error("error fetching categories from db", zap.Error(err))
		InternalServerError(w, nil)
		return
	}
	l.Info("categories rendered", zap.Int("amount", len(es)))
	easyjson.MarshalToHTTPResponseWriter(NewCategoryViews(es), w)
}

// Owner fetches the ent.owner attached to the ent.Pet
// identified by a given url-parameter from the database and renders it to the client.
func (h PetHandler) Owner(w http.ResponseWriter, r *http.Request) {
	l := h.log.With(zap.String("method", "Owner"))
	// ID is URL parameter.
	var err error
	id := chi.URLParam(r, "id")
	// Create the query to fetch the owner attached to this pet
	q := h.client.Pet.Query().Where(pet.ID(id)).QueryOwner()
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
			l.Error("could-not-read-pet", zap.Error(err), zap.String("id", id))
			InternalServerError(w, nil)
		}
		return
	}
	l.Info("owner rendered", zap.String("id", e.ID.String()))
	easyjson.MarshalToHTTPResponseWriter(NewOwnerView(e), w)
}

// Friends fetches the ent.friends attached to the ent.Pet
// identified by a given url-parameter from the database and renders it to the client.
func (h PetHandler) Friends(w http.ResponseWriter, r *http.Request) {
	l := h.log.With(zap.String("method", "Friends"))
	// ID is URL parameter.
	var err error
	id := chi.URLParam(r, "id")
	// Create the query to fetch the friends attached to this pet
	q := h.client.Pet.Query().Where(pet.ID(id)).QueryFriends()
	page := 1
	if d := r.URL.Query().Get("page"); d != "" {
		page, err = strconv.Atoi(d)
		if err != nil {
			l.Info("error parsing query parameter 'page'", zap.String("page", d), zap.Error(err))
			BadRequest(w, "page must be an integer greater zero")
			return
		}
	}
	itemsPerPage := 30
	if d := r.URL.Query().Get("itemsPerPage"); d != "" {
		itemsPerPage, err = strconv.Atoi(d)
		if err != nil {
			l.Info("error parsing query parameter 'itemsPerPage'", zap.String("itemsPerPage", d), zap.Error(err))
			BadRequest(w, "itemsPerPage must be an integer greater zero")
			return
		}
	}
	es, err := q.Limit(itemsPerPage).Offset((page - 1) * itemsPerPage).All(r.Context())
	if err != nil {
		l.Error("error fetching pets from db", zap.Error(err))
		InternalServerError(w, nil)
		return
	}
	l.Info("pets rendered", zap.Int("amount", len(es)))
	easyjson.MarshalToHTTPResponseWriter(NewPetViews(es), w)
}
