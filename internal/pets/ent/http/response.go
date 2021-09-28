// Code generated by entc, DO NOT EDIT.

package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/mailru/easyjson"
	"github.com/masseelch/elk/internal/pets/ent"
	badge "github.com/masseelch/elk/internal/pets/ent/badge"
	pet "github.com/masseelch/elk/internal/pets/ent/pet"
	playgroup "github.com/masseelch/elk/internal/pets/ent/playgroup"
	toy "github.com/masseelch/elk/internal/pets/ent/toy"
)

// Basic HTTP Error Response
type ErrResponse struct {
	Code   int         `json:"code"`             // http response status code
	Status string      `json:"status"`           // user-level status message
	Errors interface{} `json:"errors,omitempty"` // application-level error
}

func (e ErrResponse) MarshalToHTTPResponseWriter(w http.ResponseWriter) (int, error) {
	d, err := easyjson.Marshal(e)
	if err != nil {
		return 0, err
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(d)))
	w.WriteHeader(e.Code)
	return w.Write(d)
}

func BadRequest(w http.ResponseWriter, msg interface{}) (int, error) {
	return ErrResponse{
		Code:   http.StatusBadRequest,
		Status: http.StatusText(http.StatusBadRequest),
		Errors: msg,
	}.MarshalToHTTPResponseWriter(w)
}

func Conflict(w http.ResponseWriter, msg interface{}) (int, error) {
	return ErrResponse{
		Code:   http.StatusConflict,
		Status: http.StatusText(http.StatusConflict),
		Errors: msg,
	}.MarshalToHTTPResponseWriter(w)
}

func Forbidden(w http.ResponseWriter, msg interface{}) (int, error) {
	return ErrResponse{
		Code:   http.StatusForbidden,
		Status: http.StatusText(http.StatusForbidden),
		Errors: msg,
	}.MarshalToHTTPResponseWriter(w)
}

func InternalServerError(w http.ResponseWriter, msg interface{}) (int, error) {
	return ErrResponse{
		Code:   http.StatusInternalServerError,
		Status: http.StatusText(http.StatusInternalServerError),
		Errors: msg,
	}.MarshalToHTTPResponseWriter(w)
}

func NotFound(w http.ResponseWriter, msg interface{}) (int, error) {
	return ErrResponse{
		Code:   http.StatusNotFound,
		Status: http.StatusText(http.StatusNotFound),
		Errors: msg,
	}.MarshalToHTTPResponseWriter(w)
}

func Unauthorized(w http.ResponseWriter, msg interface{}) (int, error) {
	return ErrResponse{
		Code:   http.StatusUnauthorized,
		Status: http.StatusText(http.StatusUnauthorized),
		Errors: msg,
	}.MarshalToHTTPResponseWriter(w)
}

type (
	// BadgeView represents the data serialized for the following serialization group combinations:
	// []
	// [pet:list pet:read]
	// [pet:read]
	// [pet:list]
	BadgeView struct {
		ID       uint32         `json:"id,omitempty"`
		Color    badge.Color    `json:"color,omitempty"`
		Material badge.Material `json:"material,omitempty"`
	}
	BadgeViews []*BadgeView
)

func NewBadgeView(e *ent.Badge) *BadgeView {
	if e == nil {
		return nil
	}
	return &BadgeView{
		ID:       e.ID,
		Color:    e.Color,
		Material: e.Material,
	}
}

func NewBadgeViews(es []*ent.Badge) BadgeViews {
	if len(es) == 0 {
		return nil
	}
	r := make(BadgeViews, len(es))
	for i, e := range es {
		r[i] = NewBadgeView(e)
	}
	return r
}

type (
	// PetView represents the data serialized for the following serialization group combinations:
	// []
	// [pet:list pet:read]
	// [pet:read]
	// [pet:list]
	PetView struct {
		ID        int       `json:"id,omitempty"`
		Height    int       `json:"height,omitempty"`
		Weight    float64   `json:"weight,omitempty"`
		Castrated bool      `json:"castrated,omitempty"`
		Name      string    `json:"name,omitempty"`
		Birthday  time.Time `json:"birthday,omitempty"`
		Nicknames []string  `json:"nicknames,omitempty"`
		Sex       pet.Sex   `json:"sex,omitempty"`
		Chip      uuid.UUID `json:"chip,omitempty"`
	}
	PetViews []*PetView
)

func NewPetView(e *ent.Pet) *PetView {
	if e == nil {
		return nil
	}
	return &PetView{
		ID:        e.ID,
		Height:    e.Height,
		Weight:    e.Weight,
		Castrated: e.Castrated,
		Name:      e.Name,
		Birthday:  e.Birthday,
		Nicknames: e.Nicknames,
		Sex:       e.Sex,
		Chip:      e.Chip,
	}
}

func NewPetViews(es []*ent.Pet) PetViews {
	if len(es) == 0 {
		return nil
	}
	r := make(PetViews, len(es))
	for i, e := range es {
		r[i] = NewPetView(e)
	}
	return r
}

type (
	// PlayGroupView represents the data serialized for the following serialization group combinations:
	// []
	// [pet:list pet:read]
	// [pet:read]
	// [pet:list]
	PlayGroupView struct {
		ID          int               `json:"id,omitempty"`
		Title       string            `json:"title,omitempty"`
		Description string            `json:"description,omitempty"`
		Weekday     playgroup.Weekday `json:"weekday,omitempty"`
	}
	PlayGroupViews []*PlayGroupView
)

func NewPlayGroupView(e *ent.PlayGroup) *PlayGroupView {
	if e == nil {
		return nil
	}
	return &PlayGroupView{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		Weekday:     e.Weekday,
	}
}

func NewPlayGroupViews(es []*ent.PlayGroup) PlayGroupViews {
	if len(es) == 0 {
		return nil
	}
	r := make(PlayGroupViews, len(es))
	for i, e := range es {
		r[i] = NewPlayGroupView(e)
	}
	return r
}

type (
	// ToyView represents the data serialized for the following serialization group combinations:
	// []
	// [pet:list pet:read]
	// [pet:read]
	// [pet:list]
	ToyView struct {
		ID       uuid.UUID    `json:"id,omitempty"`
		Color    toy.Color    `json:"color,omitempty"`
		Material toy.Material `json:"material,omitempty"`
		Title    string       `json:"title,omitempty"`
	}
	ToyViews []*ToyView
)

func NewToyView(e *ent.Toy) *ToyView {
	if e == nil {
		return nil
	}
	return &ToyView{
		ID:       e.ID,
		Color:    e.Color,
		Material: e.Material,
		Title:    e.Title,
	}
}

func NewToyViews(es []*ent.Toy) ToyViews {
	if len(es) == 0 {
		return nil
	}
	r := make(ToyViews, len(es))
	for i, e := range es {
		r[i] = NewToyView(e)
	}
	return r
}
