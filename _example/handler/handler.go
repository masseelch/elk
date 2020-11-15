// Code generated by entc, DO NOT EDIT.

package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/facebook/ent/entc/integration/ent/pet"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/liip/sheriff"
	"github.com/masseelch/go-api-skeleton/ent"
	"github.com/masseelch/render"
	"github.com/sirupsen/logrus"
)

// Shared handler.
type handler struct {
	*chi.Mux

	client    *ent.Client
	validator *validator.Validate
	logger    *logrus.Logger
}

// The OwnerHandler.
type OwnerHandler struct {
	*handler
}

// Create a new OwnerHandler
func NewOwnerHandler(c *ent.Client, v *validator.Validate, log *logrus.Logger) *OwnerHandler {
	h := &OwnerHandler{
		&handler{
			Mux:       chi.NewRouter(),
			client:    c,
			validator: v,
			logger:    log,
		},
	}

	h.Post("/", h.Create)
	h.Get("/{id:\\d+}", h.Read)
	h.Get("/{id:\\d+}", h.Update)

	h.Get("/", h.List)

	h.Get("/{id:\\d+}/pets", h.Pets)

	return h
}

// struct to bind the post body to.
type ownerCreateRequest struct {
	Name string `json:"name,omitempty" `
	Pets []int  `json:"pets" `
}

// This function creates a new Owner model and stores it in the database.
func (h OwnerHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Get the post data.
	d := ownerCreateRequest{} // todo - allow form-url-encdoded/xml/protobuf data.
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		h.logger.WithError(err).Error("error decoding json")
		render.BadRequest(w, r, "invalid json string")
		return
	}

	// Validate the data.
	if err := h.validator.Struct(d); err != nil {
		if err, ok := err.(*validator.InvalidValidationError); ok {
			h.logger.WithError(err).Error("error validating request data")
			render.InternalServerError(w, r, nil)
			return
		}

		h.logger.WithError(err).Info("validation failed")
		render.BadRequest(w, r, err)
		return
	}

	// Save the data.
	b := h.client.Owner.Create().
		SetName(d.Name).
		AddPetIDs(d.Pets...)

	// Store in database.
	e, err := b.Save(r.Context())
	if err != nil {
		h.logger.WithError(err).Error("error saving Owner")
		render.InternalServerError(w, r, nil)
		return
	}

	// Serialize the data.
	j, err := sheriff.Marshal(&sheriff.Options{Groups: []string{"owner:read"}}, e)
	if err != nil {
		h.logger.WithError(err).WithField("Owner.id", e.ID).Error("serialization error")
		render.InternalServerError(w, r, nil)
		return
	}

	h.logger.WithField("owner", e.ID).Info("owner rendered")
	render.OK(w, r, j)
}

// This function fetches the Owner model identified by a give url-parameter from
// database and returns it to the client.
func (h OwnerHandler) Read(w http.ResponseWriter, r *http.Request) {
	id, err := h.urlParamInt(w, r, "id")
	if err != nil {
		return
	}

	qb := h.client.Owner.Query().Where(owner.ID(id))

	qb.WithPets()

	e, err := qb.Only(r.Context())

	if err != nil {
		switch err.(type) {
		case *ent.NotFoundError:
			h.logger.WithError(err).WithField("Owner.id", id).Debug("job not found")
			render.NotFound(w, r, err)
			return
		case *ent.NotSingularError:
			h.logger.WithError(err).WithField("Owner.id", id).Error("duplicate entry for id")
			render.InternalServerError(w, r, nil)
			return
		default:
			h.logger.WithError(err).WithField("Owner.id", id).Error("error fetching node from db")
			render.InternalServerError(w, r, nil)
			return
		}
	}

	d, err := sheriff.Marshal(&sheriff.Options{Groups: []string{"owner:read", "pet:list"}}, e)
	if err != nil {
		h.logger.WithError(err).WithField("Owner.id", id).Error("serialization error")
		render.InternalServerError(w, r, nil)
		return
	}

	h.logger.WithField("owner", e.ID).Info("owner rendered")
	render.OK(w, r, d)
}

// struct to bind the post body to.
type ownerUpdateRequest struct {
	Name string `json:"name,omitempty" `
	Pets []int  `json:"pets" `
}

// This function updates a given Owner model and saves the changes in the database.
func (h OwnerHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := h.urlParamInt(w, r, "id")
	if err != nil {
		return
	}

	// Get the post data.
	d := ownerUpdateRequest{} // todo - allow form-url-encoded/xml/protobuf data.
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		h.logger.WithError(err).Error("error decoding json")
		render.BadRequest(w, r, "invalid json string")
		return
	}

	// Validate the data.
	if err := h.validator.Struct(d); err != nil {
		if err, ok := err.(*validator.InvalidValidationError); ok {
			h.logger.WithError(err).Error("error validating request data")
			render.InternalServerError(w, r, nil)
			return
		}

		h.logger.WithError(err).Info("validation failed")
		render.BadRequest(w, r, err)
		return
	}

	// Save the data.
	b := h.client.Owner.UpdateOneID(id).
		SetName(d.Name).
		AddPetIDs(d.Pets...)

	// Save in database.
	e, err := b.Save(r.Context())
	if err != nil {
		h.logger.WithError(err).Error("error saving Owner")
		render.InternalServerError(w, r, nil)
		return
	}

	// Serialize the data.
	j, err := sheriff.Marshal(&sheriff.Options{Groups: []string{"owner:read"}}, e)
	if err != nil {
		h.logger.WithError(err).WithField("Owner.id", e.ID).Error("serialization error")
		render.InternalServerError(w, r, nil)
		return
	}

	h.logger.WithField("owner", e.ID).Info("owner rendered")
	render.OK(w, r, j)
}

// This function queries for Owner models. Can be filtered by query parameters.
func (h OwnerHandler) List(w http.ResponseWriter, r *http.Request) {
	q := h.client.Owner.Query()

	// Pagination
	page, itemsPerPage, err := h.paginationInfo(w, r)
	if err != nil {
		return
	}

	q = q.Limit(itemsPerPage).Offset((page - 1) * itemsPerPage)

	// Use the query parameters to filter the query. todo - nested filter?
	if f := r.URL.Query().Get("name"); f != "" {
		q = q.Where(owner.Name(f))
	}

	es, err := q.All(r.Context())
	if err != nil {
		h.logger.WithError(err).Error("error querying database") // todo - better error
		render.InternalServerError(w, r, nil)
		return
	}

	d, err := sheriff.Marshal(&sheriff.Options{Groups: []string{"owner:read", "pet:list"}}, es)
	if err != nil {
		h.logger.WithError(err).Error("serialization error")
		render.InternalServerError(w, r, nil)
		return
	}

	h.logger.WithField("amount", len(es)).Info("owner rendered")
	render.OK(w, r, d)
}

// Enable the read operation on the pets edge.
func (h *OwnerHandler) EnablePetsEndpoint() *OwnerHandler {
	h.Get("/{id:\\d+}/pets", h.Pets)
	return h
}

func (h OwnerHandler) Pets(w http.ResponseWriter, r *http.Request) {
	id, err := h.urlParamInt(w, r, "id")
	if err != nil {
		return
	}
	qb := h.client.Owner.Query().Where(owner.ID(id)).QueryPets()

	es, err := qb.All(r.Context())
	if err != nil {
		h.logger.WithError(err).Error("error querying database") // todo - better error
		render.InternalServerError(w, r, nil)
		return
	}

	d, err := sheriff.Marshal(&sheriff.Options{Groups: []string{"pet:list"}}, es)
	if err != nil {
		h.logger.WithError(err).Error("serialization error")
		render.InternalServerError(w, r, nil)
		return
	}

	h.logger.WithField("amount", len(es)).Info("pet rendered")
	render.OK(w, r, d)

}

// The PetHandler.
type PetHandler struct {
	*handler
}

// Create a new PetHandler
func NewPetHandler(c *ent.Client, v *validator.Validate, log *logrus.Logger) *PetHandler {
	h := &PetHandler{
		&handler{
			Mux:       chi.NewRouter(),
			client:    c,
			validator: v,
			logger:    log,
		},
	}

	h.Post("/", h.Create)
	h.Get("/{id:\\d+}", h.Read)
	h.Get("/{id:\\d+}", h.Update)

	h.Get("/", h.List)

	h.Get("/{id:\\d+}/owner", h.Owner)

	return h
}

// struct to bind the post body to.
type petCreateRequest struct {
	Name  string `json:"name,omitempty" `
	Owner int    `json:"" `
}

// This function creates a new Pet model and stores it in the database.
func (h PetHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Get the post data.
	d := petCreateRequest{} // todo - allow form-url-encdoded/xml/protobuf data.
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		h.logger.WithError(err).Error("error decoding json")
		render.BadRequest(w, r, "invalid json string")
		return
	}

	// Validate the data.
	if err := h.validator.Struct(d); err != nil {
		if err, ok := err.(*validator.InvalidValidationError); ok {
			h.logger.WithError(err).Error("error validating request data")
			render.InternalServerError(w, r, nil)
			return
		}

		h.logger.WithError(err).Info("validation failed")
		render.BadRequest(w, r, err)
		return
	}

	// Save the data.
	b := h.client.Pet.Create().
		SetName(d.Name).
		SetOwnerID(d.Owner)

	// Store in database.
	e, err := b.Save(r.Context())
	if err != nil {
		h.logger.WithError(err).Error("error saving Pet")
		render.InternalServerError(w, r, nil)
		return
	}

	// Serialize the data.
	j, err := sheriff.Marshal(&sheriff.Options{Groups: []string{"pet:read"}}, e)
	if err != nil {
		h.logger.WithError(err).WithField("Pet.id", e.ID).Error("serialization error")
		render.InternalServerError(w, r, nil)
		return
	}

	h.logger.WithField("pet", e.ID).Info("pet rendered")
	render.OK(w, r, j)
}

// This function fetches the Pet model identified by a give url-parameter from
// database and returns it to the client.
func (h PetHandler) Read(w http.ResponseWriter, r *http.Request) {
	id, err := h.urlParamInt(w, r, "id")
	if err != nil {
		return
	}

	qb := h.client.Pet.Query().Where(pet.ID(id))

	e, err := qb.Only(r.Context())

	if err != nil {
		switch err.(type) {
		case *ent.NotFoundError:
			h.logger.WithError(err).WithField("Pet.id", id).Debug("job not found")
			render.NotFound(w, r, err)
			return
		case *ent.NotSingularError:
			h.logger.WithError(err).WithField("Pet.id", id).Error("duplicate entry for id")
			render.InternalServerError(w, r, nil)
			return
		default:
			h.logger.WithError(err).WithField("Pet.id", id).Error("error fetching node from db")
			render.InternalServerError(w, r, nil)
			return
		}
	}

	d, err := sheriff.Marshal(&sheriff.Options{Groups: []string{"pet:read"}}, e)
	if err != nil {
		h.logger.WithError(err).WithField("Pet.id", id).Error("serialization error")
		render.InternalServerError(w, r, nil)
		return
	}

	h.logger.WithField("pet", e.ID).Info("pet rendered")
	render.OK(w, r, d)
}

// struct to bind the post body to.
type petUpdateRequest struct {
	Name  string `json:"name,omitempty" `
	Owner int    `json:"" `
}

// This function updates a given Pet model and saves the changes in the database.
func (h PetHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := h.urlParamInt(w, r, "id")
	if err != nil {
		return
	}

	// Get the post data.
	d := petUpdateRequest{} // todo - allow form-url-encoded/xml/protobuf data.
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		h.logger.WithError(err).Error("error decoding json")
		render.BadRequest(w, r, "invalid json string")
		return
	}

	// Validate the data.
	if err := h.validator.Struct(d); err != nil {
		if err, ok := err.(*validator.InvalidValidationError); ok {
			h.logger.WithError(err).Error("error validating request data")
			render.InternalServerError(w, r, nil)
			return
		}

		h.logger.WithError(err).Info("validation failed")
		render.BadRequest(w, r, err)
		return
	}

	// Save the data.
	b := h.client.Pet.UpdateOneID(id).
		SetName(d.Name).
		SetOwnerID(d.Owner)

	// Save in database.
	e, err := b.Save(r.Context())
	if err != nil {
		h.logger.WithError(err).Error("error saving Pet")
		render.InternalServerError(w, r, nil)
		return
	}

	// Serialize the data.
	j, err := sheriff.Marshal(&sheriff.Options{Groups: []string{"pet:read"}}, e)
	if err != nil {
		h.logger.WithError(err).WithField("Pet.id", e.ID).Error("serialization error")
		render.InternalServerError(w, r, nil)
		return
	}

	h.logger.WithField("pet", e.ID).Info("pet rendered")
	render.OK(w, r, j)
}

// This function queries for Pet models. Can be filtered by query parameters.
func (h PetHandler) List(w http.ResponseWriter, r *http.Request) {
	q := h.client.Pet.Query()

	// Pagination
	page, itemsPerPage, err := h.paginationInfo(w, r)
	if err != nil {
		return
	}

	q = q.Limit(itemsPerPage).Offset((page - 1) * itemsPerPage)

	// Use the query parameters to filter the query. todo - nested filter?
	if f := r.URL.Query().Get("name"); f != "" {
		q = q.Where(pet.Name(f))
	}

	es, err := q.All(r.Context())
	if err != nil {
		h.logger.WithError(err).Error("error querying database") // todo - better error
		render.InternalServerError(w, r, nil)
		return
	}

	d, err := sheriff.Marshal(&sheriff.Options{Groups: []string{"pet:list"}}, es)
	if err != nil {
		h.logger.WithError(err).Error("serialization error")
		render.InternalServerError(w, r, nil)
		return
	}

	h.logger.WithField("amount", len(es)).Info("pet rendered")
	render.OK(w, r, d)
}

// Enable the read operation on the owner edge.
func (h *PetHandler) EnableOwnerEndpoint() *PetHandler {
	h.Get("/{id:\\d+}/owner", h.Owner)
	return h
}

func (h PetHandler) Owner(w http.ResponseWriter, r *http.Request) {
	id, err := h.urlParamInt(w, r, "id")
	if err != nil {
		return
	}
	qb := h.client.Pet.Query().Where(pet.ID(id)).QueryOwner()

	qb.WithPets()

	e, err := qb.Only(r.Context())

	if err != nil {
		switch err.(type) {
		case *ent.NotFoundError:
			h.logger.WithError(err).WithField("Owner.id", id).Debug("job not found")
			render.NotFound(w, r, err)
			return
		case *ent.NotSingularError:
			h.logger.WithError(err).WithField("Owner.id", id).Error("duplicate entry for id")
			render.InternalServerError(w, r, nil)
			return
		default:
			h.logger.WithError(err).WithField("Owner.id", id).Error("error fetching node from db")
			render.InternalServerError(w, r, nil)
			return
		}
	}

	d, err := sheriff.Marshal(&sheriff.Options{Groups: []string{"owner:read"}}, e)
	if err != nil {
		h.logger.WithError(err).WithField("Owner.id", id).Error("serialization error")
		render.InternalServerError(w, r, nil)
		return
	}

	h.logger.WithField("owner", e.ID).Info("owner rendered")
	render.OK(w, r, d)

}

// The SkipGenerationModelHandler.
type SkipGenerationModelHandler struct {
	*handler
}

// Create a new SkipGenerationModelHandler
func NewSkipGenerationModelHandler(c *ent.Client, v *validator.Validate, log *logrus.Logger) *SkipGenerationModelHandler {
	h := &SkipGenerationModelHandler{
		&handler{
			Mux:       chi.NewRouter(),
			client:    c,
			validator: v,
			logger:    log,
		},
	}

	h.Post("/", h.Create)
	h.Get("/{id:\\d+}", h.Read)
	h.Get("/{id:\\d+}", h.Update)

	h.Get("/", h.List)

	return h
}

// struct to bind the post body to.
type skipgenerationmodelCreateRequest struct {
	Name string `json:"name,omitempty" `
}

// This function creates a new SkipGenerationModel model and stores it in the database.
func (h SkipGenerationModelHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Get the post data.
	d := skip_generation_modelCreateRequest{} // todo - allow form-url-encdoded/xml/protobuf data.
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		h.logger.WithError(err).Error("error decoding json")
		render.BadRequest(w, r, "invalid json string")
		return
	}

	// Validate the data.
	if err := h.validator.Struct(d); err != nil {
		if err, ok := err.(*validator.InvalidValidationError); ok {
			h.logger.WithError(err).Error("error validating request data")
			render.InternalServerError(w, r, nil)
			return
		}

		h.logger.WithError(err).Info("validation failed")
		render.BadRequest(w, r, err)
		return
	}

	// Save the data.
	b := h.client.SkipGenerationModel.Create().
		SetName(d.Name)

	// Store in database.
	e, err := b.Save(r.Context())
	if err != nil {
		h.logger.WithError(err).Error("error saving SkipGenerationModel")
		render.InternalServerError(w, r, nil)
		return
	}

	// Serialize the data.
	j, err := sheriff.Marshal(&sheriff.Options{Groups: []string{"skip_generation_model:read"}}, e)
	if err != nil {
		h.logger.WithError(err).WithField("SkipGenerationModel.id", e.ID).Error("serialization error")
		render.InternalServerError(w, r, nil)
		return
	}

	h.logger.WithField("skip_generation_model", e.ID).Info("skip_generation_model rendered")
	render.OK(w, r, j)
}

// This function fetches the SkipGenerationModel model identified by a give url-parameter from
// database and returns it to the client.
func (h SkipGenerationModelHandler) Read(w http.ResponseWriter, r *http.Request) {
	id, err := h.urlParamInt(w, r, "id")
	if err != nil {
		return
	}

	qb := h.client.SkipGenerationModel.Query().Where(skip_generation_model.ID(id))

	e, err := qb.Only(r.Context())

	if err != nil {
		switch err.(type) {
		case *ent.NotFoundError:
			h.logger.WithError(err).WithField("SkipGenerationModel.id", id).Debug("job not found")
			render.NotFound(w, r, err)
			return
		case *ent.NotSingularError:
			h.logger.WithError(err).WithField("SkipGenerationModel.id", id).Error("duplicate entry for id")
			render.InternalServerError(w, r, nil)
			return
		default:
			h.logger.WithError(err).WithField("SkipGenerationModel.id", id).Error("error fetching node from db")
			render.InternalServerError(w, r, nil)
			return
		}
	}

	d, err := sheriff.Marshal(&sheriff.Options{Groups: []string{"skip_generation_model:read"}}, e)
	if err != nil {
		h.logger.WithError(err).WithField("SkipGenerationModel.id", id).Error("serialization error")
		render.InternalServerError(w, r, nil)
		return
	}

	h.logger.WithField("skip_generation_model", e.ID).Info("skip_generation_model rendered")
	render.OK(w, r, d)
}

// struct to bind the post body to.
type skipgenerationmodelUpdateRequest struct {
	Name string `json:"name,omitempty" `
}

// This function updates a given SkipGenerationModel model and saves the changes in the database.
func (h SkipGenerationModelHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := h.urlParamInt(w, r, "id")
	if err != nil {
		return
	}

	// Get the post data.
	d := skip_generation_modelUpdateRequest{} // todo - allow form-url-encoded/xml/protobuf data.
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		h.logger.WithError(err).Error("error decoding json")
		render.BadRequest(w, r, "invalid json string")
		return
	}

	// Validate the data.
	if err := h.validator.Struct(d); err != nil {
		if err, ok := err.(*validator.InvalidValidationError); ok {
			h.logger.WithError(err).Error("error validating request data")
			render.InternalServerError(w, r, nil)
			return
		}

		h.logger.WithError(err).Info("validation failed")
		render.BadRequest(w, r, err)
		return
	}

	// Save the data.
	b := h.client.SkipGenerationModel.UpdateOneID(id).
		SetName(d.Name)

	// Save in database.
	e, err := b.Save(r.Context())
	if err != nil {
		h.logger.WithError(err).Error("error saving SkipGenerationModel")
		render.InternalServerError(w, r, nil)
		return
	}

	// Serialize the data.
	j, err := sheriff.Marshal(&sheriff.Options{Groups: []string{"skip_generation_model:read"}}, e)
	if err != nil {
		h.logger.WithError(err).WithField("SkipGenerationModel.id", e.ID).Error("serialization error")
		render.InternalServerError(w, r, nil)
		return
	}

	h.logger.WithField("skip_generation_model", e.ID).Info("skip_generation_model rendered")
	render.OK(w, r, j)
}

// This function queries for SkipGenerationModel models. Can be filtered by query parameters.
func (h SkipGenerationModelHandler) List(w http.ResponseWriter, r *http.Request) {
	q := h.client.SkipGenerationModel.Query()

	// Pagination
	page, itemsPerPage, err := h.paginationInfo(w, r)
	if err != nil {
		return
	}

	q = q.Limit(itemsPerPage).Offset((page - 1) * itemsPerPage)

	// Use the query parameters to filter the query. todo - nested filter?
	if f := r.URL.Query().Get("name"); f != "" {
		q = q.Where(skipgenerationmodel.Name(f))
	}

	es, err := q.All(r.Context())
	if err != nil {
		h.logger.WithError(err).Error("error querying database") // todo - better error
		render.InternalServerError(w, r, nil)
		return
	}

	d, err := sheriff.Marshal(&sheriff.Options{Groups: []string{"skip_generation_model:list"}}, es)
	if err != nil {
		h.logger.WithError(err).Error("serialization error")
		render.InternalServerError(w, r, nil)
		return
	}

	h.logger.WithField("amount", len(es)).Info("skip_generation_model rendered")
	render.OK(w, r, d)
}

func (h handler) urlParamInt(w http.ResponseWriter, r *http.Request, param string) (id int, err error) {
	p := chi.URLParam(r, param)
	if p == "" {
		h.logger.WithField("param", param).Info("empty url param")
		render.BadRequest(w, r, param+" cannot be ''")
		return
	}

	id, err = strconv.Atoi(p)
	if err != nil {
		h.logger.WithField(param, p).Info("error parsing url parameter")
		render.BadRequest(w, r, param+" must be a positive integer greater zero")
		return
	}

	return
}

func (h handler) urlParamTime(w http.ResponseWriter, r *http.Request, param string) (date time.Time, err error) {
	p := chi.URLParam(r, param)
	if p == "" {
		h.logger.WithField("param", param).Info("empty url param")
		render.BadRequest(w, r, param+" cannot be ''")
		return
	}

	date, err = time.Parse("2006-01-02", p)
	if err != nil {
		h.logger.WithField(param, p).Info("error parsing url parameter")
		render.BadRequest(w, r, param+" must be a valid date in yyyy-mm-dd format")
		return
	}

	return
}

func (h handler) paginationInfo(w http.ResponseWriter, r *http.Request) (page int, itemsPerPage int, err error) {
	page = 1
	itemsPerPage = 30

	if d := r.URL.Query().Get("itemsPerPage"); d != "" {
		itemsPerPage, err = strconv.Atoi(d)
		if err != nil {
			h.logger.WithField("itemsPerPage", d).Info("error parsing query parameter 'itemsPerPage'")
			render.BadRequest(w, r, "itemsPerPage must be a positive integer greater zero")
			return
		}
	}

	if d := r.URL.Query().Get("page"); d != "" {
		page, err = strconv.Atoi(d)
		if err != nil {
			h.logger.WithField("page", d).Info("error parsing query parameter 'page'")
			render.BadRequest(w, r, "page must be a positive integer greater zero")
			return
		}
	}

	return
}
