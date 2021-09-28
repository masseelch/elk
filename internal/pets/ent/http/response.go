// Code generated by entc, DO NOT EDIT.

package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/mailru/easyjson"
	"github.com/masseelch/elk/internal/pets/ent"
	"github.com/masseelch/elk/internal/pets/ent/badge"
	"github.com/masseelch/elk/internal/pets/ent/pet"
	"github.com/masseelch/elk/internal/pets/ent/playgroup"
	"github.com/masseelch/elk/internal/pets/ent/toy"
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
	// BadgeWithPetListAndPetReadView represents the data serialized for the following serialization group combinations:
	// [pet:list pet:read]
	BadgeWithPetListAndPetReadView struct {
		ID       uint32         `json:"id,omitempty"`
		Color    badge.Color    `json:"color,omitempty"`
		Material badge.Material `json:"material,omitempty"`
	}
	BadgeWithPetListAndPetReadViews []*BadgeWithPetListAndPetReadView
)

func NewBadgeWithPetListAndPetReadView(e *ent.Badge) *BadgeWithPetListAndPetReadView {
	if e == nil {
		return nil
	}
	return &BadgeWithPetListAndPetReadView{
		ID:       e.ID,
		Color:    e.Color,
		Material: e.Material,
	}
}

func NewBadgeWithPetListAndPetReadViews(es []*ent.Badge) BadgeWithPetListAndPetReadViews {
	if len(es) == 0 {
		return nil
	}
	r := make(BadgeWithPetListAndPetReadViews, len(es))
	for i, e := range es {
		r[i] = NewBadgeWithPetListAndPetReadView(e)
	}
	return r
}

type (
	// BadgeWithPetListView represents the data serialized for the following serialization group combinations:
	// [pet:list]
	BadgeWithPetListView struct {
		ID       uint32         `json:"id,omitempty"`
		Color    badge.Color    `json:"color,omitempty"`
		Material badge.Material `json:"material,omitempty"`
	}
	BadgeWithPetListViews []*BadgeWithPetListView
)

func NewBadgeWithPetListView(e *ent.Badge) *BadgeWithPetListView {
	if e == nil {
		return nil
	}
	return &BadgeWithPetListView{
		ID:       e.ID,
		Color:    e.Color,
		Material: e.Material,
	}
}

func NewBadgeWithPetListViews(es []*ent.Badge) BadgeWithPetListViews {
	if len(es) == 0 {
		return nil
	}
	r := make(BadgeWithPetListViews, len(es))
	for i, e := range es {
		r[i] = NewBadgeWithPetListView(e)
	}
	return r
}

type (
	// BadgeWithPetReadView represents the data serialized for the following serialization group combinations:
	// [pet:read]
	BadgeWithPetReadView struct {
		ID       uint32         `json:"id,omitempty"`
		Color    badge.Color    `json:"color,omitempty"`
		Material badge.Material `json:"material,omitempty"`
	}
	BadgeWithPetReadViews []*BadgeWithPetReadView
)

func NewBadgeWithPetReadView(e *ent.Badge) *BadgeWithPetReadView {
	if e == nil {
		return nil
	}
	return &BadgeWithPetReadView{
		ID:       e.ID,
		Color:    e.Color,
		Material: e.Material,
	}
}

func NewBadgeWithPetReadViews(es []*ent.Badge) BadgeWithPetReadViews {
	if len(es) == 0 {
		return nil
	}
	r := make(BadgeWithPetReadViews, len(es))
	for i, e := range es {
		r[i] = NewBadgeWithPetReadView(e)
	}
	return r
}

type (
	// PetListAndPetReadView represents the data serialized for the following serialization group combinations:
	// [pet:list pet:read]
	PetListAndPetReadView struct {
		ID         int                                 `json:"id,omitempty"`
		Height     int                                 `json:"height,omitempty"`
		Weight     float64                             `json:"weight,omitempty"`
		Castrated  bool                                `json:"castrated,omitempty"`
		Name       string                              `json:"name,omitempty"`
		Birthday   time.Time                           `json:"birthday,omitempty"`
		Nicknames  []string                            `json:"nicknames,omitempty"`
		Sex        pet.Sex                             `json:"sex,omitempty"`
		Chip       uuid.UUID                           `json:"chip,omitempty"`
		Badge      *BadgeWithPetListAndPetReadView     `json:"badge,omitempty"`
		Protege    *PetListAndPetReadView              `json:"protege,omitempty"`
		Spouse     *PetListAndPetReadView              `json:"spouse,omitempty"`
		Toys       ToyWithPetListAndPetReadViews       `json:"toys,omitempty"`
		Parent     *PetListAndPetReadView              `json:"parent,omitempty"`
		PlayGroups PlayGroupWithPetListAndPetReadViews `json:"play_groups,omitempty"`
		Friends    PetListAndPetReadViews              `json:"friends,omitempty"`
	}
	PetListAndPetReadViews []*PetListAndPetReadView
)

func NewPetListAndPetReadView(e *ent.Pet) *PetListAndPetReadView {
	if e == nil {
		return nil
	}
	return &PetListAndPetReadView{
		ID:         e.ID,
		Height:     e.Height,
		Weight:     e.Weight,
		Castrated:  e.Castrated,
		Name:       e.Name,
		Birthday:   e.Birthday,
		Nicknames:  e.Nicknames,
		Sex:        e.Sex,
		Chip:       e.Chip,
		Badge:      NewBadgeWithPetListAndPetReadView(e.Edges.Badge),
		Protege:    NewPetListAndPetReadView(e.Edges.Protege),
		Spouse:     NewPetListAndPetReadView(e.Edges.Spouse),
		Toys:       NewToyWithPetListAndPetReadViews(e.Edges.Toys),
		Parent:     NewPetListAndPetReadView(e.Edges.Parent),
		PlayGroups: NewPlayGroupWithPetListAndPetReadViews(e.Edges.PlayGroups),
		Friends:    NewPetListAndPetReadViews(e.Edges.Friends),
	}
}

func NewPetListAndPetReadViews(es []*ent.Pet) PetListAndPetReadViews {
	if len(es) == 0 {
		return nil
	}
	r := make(PetListAndPetReadViews, len(es))
	for i, e := range es {
		r[i] = NewPetListAndPetReadView(e)
	}
	return r
}

type (
	// PetListView represents the data serialized for the following serialization group combinations:
	// [pet:list]
	PetListView struct {
		ID    int                   `json:"id,omitempty"`
		Name  string                `json:"name,omitempty"`
		Sex   pet.Sex               `json:"sex,omitempty"`
		Chip  uuid.UUID             `json:"chip,omitempty"`
		Badge *BadgeWithPetListView `json:"badge,omitempty"`
	}
	PetListViews []*PetListView
)

func NewPetListView(e *ent.Pet) *PetListView {
	if e == nil {
		return nil
	}
	return &PetListView{
		ID:    e.ID,
		Name:  e.Name,
		Sex:   e.Sex,
		Chip:  e.Chip,
		Badge: NewBadgeWithPetListView(e.Edges.Badge),
	}
}

func NewPetListViews(es []*ent.Pet) PetListViews {
	if len(es) == 0 {
		return nil
	}
	r := make(PetListViews, len(es))
	for i, e := range es {
		r[i] = NewPetListView(e)
	}
	return r
}

type (
	// PetReadView represents the data serialized for the following serialization group combinations:
	// [pet:read]
	PetReadView struct {
		ID         int                       `json:"id,omitempty"`
		Height     int                       `json:"height,omitempty"`
		Weight     float64                   `json:"weight,omitempty"`
		Castrated  bool                      `json:"castrated,omitempty"`
		Name       string                    `json:"name,omitempty"`
		Birthday   time.Time                 `json:"birthday,omitempty"`
		Nicknames  []string                  `json:"nicknames,omitempty"`
		Sex        pet.Sex                   `json:"sex,omitempty"`
		Chip       uuid.UUID                 `json:"chip,omitempty"`
		Badge      *BadgeWithPetReadView     `json:"badge,omitempty"`
		Protege    *PetReadView              `json:"protege,omitempty"`
		Spouse     *PetReadView              `json:"spouse,omitempty"`
		Toys       ToyWithPetReadViews       `json:"toys,omitempty"`
		Parent     *PetReadView              `json:"parent,omitempty"`
		PlayGroups PlayGroupWithPetReadViews `json:"play_groups,omitempty"`
		Friends    PetReadViews              `json:"friends,omitempty"`
	}
	PetReadViews []*PetReadView
)

func NewPetReadView(e *ent.Pet) *PetReadView {
	if e == nil {
		return nil
	}
	return &PetReadView{
		ID:         e.ID,
		Height:     e.Height,
		Weight:     e.Weight,
		Castrated:  e.Castrated,
		Name:       e.Name,
		Birthday:   e.Birthday,
		Nicknames:  e.Nicknames,
		Sex:        e.Sex,
		Chip:       e.Chip,
		Badge:      NewBadgeWithPetReadView(e.Edges.Badge),
		Protege:    NewPetReadView(e.Edges.Protege),
		Spouse:     NewPetReadView(e.Edges.Spouse),
		Toys:       NewToyWithPetReadViews(e.Edges.Toys),
		Parent:     NewPetReadView(e.Edges.Parent),
		PlayGroups: NewPlayGroupWithPetReadViews(e.Edges.PlayGroups),
		Friends:    NewPetReadViews(e.Edges.Friends),
	}
}

func NewPetReadViews(es []*ent.Pet) PetReadViews {
	if len(es) == 0 {
		return nil
	}
	r := make(PetReadViews, len(es))
	for i, e := range es {
		r[i] = NewPetReadView(e)
	}
	return r
}

type (
	// PetView represents the data serialized for the following serialization group combinations:
	// []
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
	// PlayGroupWithPetListAndPetReadView represents the data serialized for the following serialization group combinations:
	// [pet:list pet:read]
	PlayGroupWithPetListAndPetReadView struct {
		ID          int               `json:"id,omitempty"`
		Title       string            `json:"title,omitempty"`
		Description string            `json:"description,omitempty"`
		Weekday     playgroup.Weekday `json:"weekday,omitempty"`
	}
	PlayGroupWithPetListAndPetReadViews []*PlayGroupWithPetListAndPetReadView
)

func NewPlayGroupWithPetListAndPetReadView(e *ent.PlayGroup) *PlayGroupWithPetListAndPetReadView {
	if e == nil {
		return nil
	}
	return &PlayGroupWithPetListAndPetReadView{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		Weekday:     e.Weekday,
	}
}

func NewPlayGroupWithPetListAndPetReadViews(es []*ent.PlayGroup) PlayGroupWithPetListAndPetReadViews {
	if len(es) == 0 {
		return nil
	}
	r := make(PlayGroupWithPetListAndPetReadViews, len(es))
	for i, e := range es {
		r[i] = NewPlayGroupWithPetListAndPetReadView(e)
	}
	return r
}

type (
	// PlayGroupWithPetListView represents the data serialized for the following serialization group combinations:
	// [pet:list]
	PlayGroupWithPetListView struct {
		ID          int               `json:"id,omitempty"`
		Title       string            `json:"title,omitempty"`
		Description string            `json:"description,omitempty"`
		Weekday     playgroup.Weekday `json:"weekday,omitempty"`
	}
	PlayGroupWithPetListViews []*PlayGroupWithPetListView
)

func NewPlayGroupWithPetListView(e *ent.PlayGroup) *PlayGroupWithPetListView {
	if e == nil {
		return nil
	}
	return &PlayGroupWithPetListView{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		Weekday:     e.Weekday,
	}
}

func NewPlayGroupWithPetListViews(es []*ent.PlayGroup) PlayGroupWithPetListViews {
	if len(es) == 0 {
		return nil
	}
	r := make(PlayGroupWithPetListViews, len(es))
	for i, e := range es {
		r[i] = NewPlayGroupWithPetListView(e)
	}
	return r
}

type (
	// PlayGroupWithPetReadView represents the data serialized for the following serialization group combinations:
	// [pet:read]
	PlayGroupWithPetReadView struct {
		ID          int               `json:"id,omitempty"`
		Title       string            `json:"title,omitempty"`
		Description string            `json:"description,omitempty"`
		Weekday     playgroup.Weekday `json:"weekday,omitempty"`
	}
	PlayGroupWithPetReadViews []*PlayGroupWithPetReadView
)

func NewPlayGroupWithPetReadView(e *ent.PlayGroup) *PlayGroupWithPetReadView {
	if e == nil {
		return nil
	}
	return &PlayGroupWithPetReadView{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		Weekday:     e.Weekday,
	}
}

func NewPlayGroupWithPetReadViews(es []*ent.PlayGroup) PlayGroupWithPetReadViews {
	if len(es) == 0 {
		return nil
	}
	r := make(PlayGroupWithPetReadViews, len(es))
	for i, e := range es {
		r[i] = NewPlayGroupWithPetReadView(e)
	}
	return r
}

type (
	// ToyView represents the data serialized for the following serialization group combinations:
	// []
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

type (
	// ToyWithPetListAndPetReadView represents the data serialized for the following serialization group combinations:
	// [pet:list pet:read]
	ToyWithPetListAndPetReadView struct {
		ID       uuid.UUID    `json:"id,omitempty"`
		Color    toy.Color    `json:"color,omitempty"`
		Material toy.Material `json:"material,omitempty"`
		Title    string       `json:"title,omitempty"`
	}
	ToyWithPetListAndPetReadViews []*ToyWithPetListAndPetReadView
)

func NewToyWithPetListAndPetReadView(e *ent.Toy) *ToyWithPetListAndPetReadView {
	if e == nil {
		return nil
	}
	return &ToyWithPetListAndPetReadView{
		ID:       e.ID,
		Color:    e.Color,
		Material: e.Material,
		Title:    e.Title,
	}
}

func NewToyWithPetListAndPetReadViews(es []*ent.Toy) ToyWithPetListAndPetReadViews {
	if len(es) == 0 {
		return nil
	}
	r := make(ToyWithPetListAndPetReadViews, len(es))
	for i, e := range es {
		r[i] = NewToyWithPetListAndPetReadView(e)
	}
	return r
}

type (
	// ToyWithPetListView represents the data serialized for the following serialization group combinations:
	// [pet:list]
	ToyWithPetListView struct {
		ID       uuid.UUID    `json:"id,omitempty"`
		Color    toy.Color    `json:"color,omitempty"`
		Material toy.Material `json:"material,omitempty"`
		Title    string       `json:"title,omitempty"`
	}
	ToyWithPetListViews []*ToyWithPetListView
)

func NewToyWithPetListView(e *ent.Toy) *ToyWithPetListView {
	if e == nil {
		return nil
	}
	return &ToyWithPetListView{
		ID:       e.ID,
		Color:    e.Color,
		Material: e.Material,
		Title:    e.Title,
	}
}

func NewToyWithPetListViews(es []*ent.Toy) ToyWithPetListViews {
	if len(es) == 0 {
		return nil
	}
	r := make(ToyWithPetListViews, len(es))
	for i, e := range es {
		r[i] = NewToyWithPetListView(e)
	}
	return r
}

type (
	// ToyWithPetReadView represents the data serialized for the following serialization group combinations:
	// [pet:read]
	ToyWithPetReadView struct {
		ID       uuid.UUID    `json:"id,omitempty"`
		Color    toy.Color    `json:"color,omitempty"`
		Material toy.Material `json:"material,omitempty"`
		Title    string       `json:"title,omitempty"`
	}
	ToyWithPetReadViews []*ToyWithPetReadView
)

func NewToyWithPetReadView(e *ent.Toy) *ToyWithPetReadView {
	if e == nil {
		return nil
	}
	return &ToyWithPetReadView{
		ID:       e.ID,
		Color:    e.Color,
		Material: e.Material,
		Title:    e.Title,
	}
}

func NewToyWithPetReadViews(es []*ent.Toy) ToyWithPetReadViews {
	if len(es) == 0 {
		return nil
	}
	r := make(ToyWithPetReadViews, len(es))
	for i, e := range es {
		r[i] = NewToyWithPetReadView(e)
	}
	return r
}
