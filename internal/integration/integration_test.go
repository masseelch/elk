package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"entgo.io/ent/dialect"
	"github.com/go-chi/chi/v5"
	"github.com/masseelch/elk/internal/integration/pets/ent"
	"github.com/masseelch/elk/internal/integration/pets/ent/enttest"
	elkhttp "github.com/masseelch/elk/internal/integration/pets/ent/http"
	"github.com/masseelch/elk/internal/integration/pets/ent/pet"
	"github.com/masseelch/render"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type test struct {
	// name of the test
	name string
	// request to send
	req *http.Request
	// expected status
	status int
	// expected body
	body []byte
	// expected outputs to logging
	logs []map[string]interface{}
	// additional test logic on response body
	fn func(t *testing.T, tt *test, b []byte)
}

func TestHttp(t *testing.T) {
	c := enttest.Open(t, dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1", enttest.WithOptions(ent.Log(t.Log)))
	defer c.Close()

	// Load test data.
	require.NoError(t, fixtures(context.Background(), c))

	// Logger
	logs := new(bytes.Buffer)
	l := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(
			zapcore.EncoderConfig{
				MessageKey:     "msg",
				LevelKey:       "level",
				NameKey:        "logger",
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.StringDurationEncoder,
			},
		),
		zapcore.AddSync(logs),
		zap.DebugLevel,
	))
	defer l.Sync()

	// Validator
	// v := validator.New()
	// v.RegisterTagNameFunc(func(fld reflect.StructField) string {
	// 	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	// 	if name == "-" {
	// 		return ""
	// 	}
	// 	return name
	// })

	// Needed to test url param fetching
	r := chi.NewRouter()

	// Register pet endpoints.
	r.Route("/pets", func(r chi.Router) {
		elkhttp.NewPetHandler(c, l).Mount(r, elkhttp.PetRoutes)
	})

	// Create the tests.
	tests := []test{
		{
			name:   "create _ invalid json",
			req:    httptest.NewRequest(http.MethodPost, "/pets", strings.NewReader("invalid")),
			status: http.StatusBadRequest,
			body:   mustEncode(t, render.NewResponse(http.StatusBadRequest, "invalid json string")),
			logs: []map[string]interface{}{
				{
					"level":   "error",
					"msg":     "error decoding json",
					"handler": "PetHandler",
					"method":  "Create",
				},
			},
		},
		{
			name:   "create _ failed validation",
			req:    httptest.NewRequest(http.MethodPost, "/pets", bytes.NewReader(mustEncode(t, map[string]interface{}{"age": 0}))),
			status: http.StatusBadRequest,
			body: mustEncode(t, render.NewResponse(http.StatusBadRequest, map[string]interface{}{
				"age":  "This value failed validation on 'gt:0'.",
				"name": "This value is required.",
			})),
			logs: []map[string]interface{}{
				{
					"level":   "info",
					"msg":     "validation failed",
					"handler": "PetHandler",
					"method":  "Create",
				},
			},
		},
		{
			name: "create _ ok",
			req: httptest.NewRequest(http.MethodPost, "/pets", bytes.NewReader(mustEncode(t, map[string]interface{}{
				"name": "my new pet",
				"age":  1,
			}))),
			status: http.StatusOK,
			fn: func(t *testing.T, tt *test, b []byte) {
				p, err := c.Pet.Query().Order(ent.Desc(pet.FieldID)).First(context.Background())
				require.NoError(t, err)
				var j map[string]interface{}
				require.NoError(t, json.Unmarshal(b, &j))
				require.EqualValues(t, p.ID, j["id"])
				require.Equal(t, p.Age, 1)
				require.Equal(t, p.Name, "my new pet")
			},
			logs: []map[string]interface{}{
				{
					"level":   "info",
					"msg":     "pet rendered",
					"handler": "PetHandler",
					"method":  "Create",
					"id":      51,
				},
			},
		},
		{
			name:   "read _ malformed id",
			req:    httptest.NewRequest(http.MethodGet, "/pets/invalid", nil),
			status: http.StatusBadRequest,
			body:   mustEncode(t, render.NewResponse(http.StatusBadRequest, "id must be an integer greater zero")),
			logs: []map[string]interface{}{
				{
					"level":   "error",
					"msg":     "error getting id from url parameter",
					"handler": "PetHandler",
					"method":  "Read",
					"id":      "invalid",
				},
			},
		},
		{
			name:   "read _ not found",
			req:    httptest.NewRequest(http.MethodGet, "/pets/10000", nil),
			status: http.StatusNotFound,
			body:   mustEncode(t, render.NewResponse(http.StatusNotFound, "pet not found")),
			logs: []map[string]interface{}{
				{
					"level":   "info",
					"msg":     "pet not found",
					"handler": "PetHandler",
					"method":  "Read",
					"id":      10000,
				},
			},
		},
		{
			name:   "read _ ok",
			req:    httptest.NewRequest(http.MethodGet, "/pets/1", nil),
			status: http.StatusOK,
			logs: []map[string]interface{}{
				{
					"level":   "info",
					"msg":     "pet rendered",
					"handler": "PetHandler",
					"method":  "Read",
					"id":      1,
				},
			},
		},
		{
			name:   "update _ malformed id",
			req:    httptest.NewRequest(http.MethodPatch, "/pets/invalid", nil),
			status: http.StatusBadRequest,
			body:   mustEncode(t, render.NewResponse(http.StatusBadRequest, "id must be an integer greater zero")),
			logs: []map[string]interface{}{
				{
					"level":   "error",
					"msg":     "error getting id from url parameter",
					"handler": "PetHandler",
					"method":  "Update",
					"id":      "invalid",
				},
			},
		},
		{
			name:   "update _ invalid json",
			req:    httptest.NewRequest(http.MethodPatch, "/pets/1", strings.NewReader("invalid")),
			status: http.StatusBadRequest,
			body:   mustEncode(t, render.NewResponse(http.StatusBadRequest, "invalid json string")),
			logs: []map[string]interface{}{
				{
					"level":   "error",
					"msg":     "error decoding json",
					"handler": "PetHandler",
					"method":  "Update",
				},
			},
		},
		{
			name:   "update _ failed validation",
			req:    httptest.NewRequest(http.MethodPatch, "/pets/1000", bytes.NewReader(mustEncode(t, map[string]interface{}{"age": 0}))),
			status: http.StatusBadRequest,
			body:   mustEncode(t, render.NewResponse(http.StatusBadRequest, map[string]interface{}{"age": "This value failed validation on 'gt:0'."})),
			logs: []map[string]interface{}{
				{
					"level":   "info",
					"msg":     "validation failed",
					"handler": "PetHandler",
					"method":  "Update",
				},
			},
		},
		{
			name: "update _ not found",
			req: httptest.NewRequest(http.MethodPatch, "/pets/1000", bytes.NewReader(mustEncode(t, map[string]interface{}{
				"age":  1,
				"name": "this is my new name",
			}))),
			status: http.StatusNotFound,
			body:   mustEncode(t, render.NewResponse(http.StatusNotFound, "pet not found")),
			logs: []map[string]interface{}{
				{
					"level":   "info",
					"msg":     "pet not found",
					"handler": "PetHandler",
					"method":  "Update",
					"id":      1000,
				},
			},
		},
		{
			name: "update _ ok",
			req: httptest.NewRequest(http.MethodPatch, "/pets/1", bytes.NewReader(mustEncode(t, map[string]interface{}{
				"age":  1000,
				"name": "this is my new name",
			}))),
			status: http.StatusOK,
			fn: func(t *testing.T, tt *test, b []byte) {
				p, err := c.Pet.Get(context.Background(), 1)
				require.NoError(t, err)
				var j map[string]interface{}
				require.NoError(t, json.Unmarshal(b, &j))
				require.EqualValues(t, p.ID, j["id"])
				require.Equal(t, p.Age, 1000)
				require.Equal(t, p.Name, "this is my new name")
			},
			logs: []map[string]interface{}{
				{
					"level":   "info",
					"msg":     "pet rendered",
					"handler": "PetHandler",
					"method":  "Update",
					"id":      1,
				},
			},
		},
		{
			name:   "delete _ malformed id",
			req:    httptest.NewRequest(http.MethodDelete, "/pets/invalid", nil),
			status: http.StatusBadRequest,
			body:   mustEncode(t, render.NewResponse(http.StatusBadRequest, "id must be an integer greater zero")),
			logs: []map[string]interface{}{
				{
					"level":   "error",
					"msg":     "error getting id from url parameter",
					"handler": "PetHandler",
					"method":  "Delete",
					"id":      "invalid",
				},
			},
		},
		{
			name:   "delete _ not found",
			req:    httptest.NewRequest(http.MethodDelete, "/pets/1000", nil),
			status: http.StatusNotFound,
			body:   mustEncode(t, render.NewResponse(http.StatusNotFound, "pet not found")),
			logs: []map[string]interface{}{
				{
					"level":   "info",
					"msg":     "pet not found",
					"handler": "PetHandler",
					"method":  "Delete",
					"id":      1000,
				},
			},
		},
		{
			name:   "delete _ ok",
			req:    httptest.NewRequest(http.MethodDelete, "/pets/50", nil),
			status: http.StatusNoContent,
			fn: func(t *testing.T, tt *test, b []byte) {
				_, err := c.Pet.Get(context.Background(), 50)
				require.EqualError(t, err, "ent: pet not found")
			},
			logs: []map[string]interface{}{
				{
					"level":   "info",
					"msg":     "pet deleted",
					"handler": "PetHandler",
					"method":  "Delete",
					"id":      50,
				},
			},
		},
		{
			name:   "list _ malformed page",
			req:    httptest.NewRequest(http.MethodGet, "/pets?page=invalid", nil),
			status: http.StatusBadRequest,
			body:   mustEncode(t, render.NewResponse(http.StatusBadRequest, "page must be an integer greater zero")),
			logs: []map[string]interface{}{
				{
					"level":   "info",
					"msg":     "error parsing query parameter 'page'",
					"handler": "PetHandler",
					"method":  "List",
					"page":    "invalid",
				},
			},
		},
		{
			name:   "list _ malformed itemsPerPage",
			req:    httptest.NewRequest(http.MethodGet, "/pets?itemsPerPage=invalid", nil),
			status: http.StatusBadRequest,
			body:   mustEncode(t, render.NewResponse(http.StatusBadRequest, "itemsPerPage must be an integer greater zero")),
			logs: []map[string]interface{}{
				{
					"level":        "info",
					"msg":          "error parsing query parameter 'itemsPerPage'",
					"handler":      "PetHandler",
					"method":       "List",
					"itemsPerPage": "invalid",
				},
			},
		},
		{
			name:   "list _ ok",
			req:    httptest.NewRequest(http.MethodGet, "/pets", nil),
			status: http.StatusOK,
			fn: func(t *testing.T, tt *test, b []byte) {
				var j []ent.Pet
				require.NoError(t, json.Unmarshal(b, &j))
				require.Len(t, j, 30)
			},
			logs: []map[string]interface{}{
				{
					"level":   "info",
					"msg":     "pets rendered",
					"handler": "PetHandler",
					"method":  "List",
					"amount":  30,
				},
			},
		},
		{
			name:   "list _ custom page and itemsPerPage ok",
			req:    httptest.NewRequest(http.MethodGet, "/pets?page=2&itemsPerPage=2", nil),
			status: http.StatusOK,
			fn: func(t *testing.T, tt *test, b []byte) {
				var j []ent.Pet
				require.NoError(t, json.Unmarshal(b, &j))
				require.Len(t, j, 2)
				require.Equal(t, 3, j[0].ID) // default order is ascending id
			},
			logs: []map[string]interface{}{
				{
					"level":   "info",
					"msg":     "pets rendered",
					"handler": "PetHandler",
					"method":  "List",
					"amount":  2,
				},
			},
		},
		{
			name:   "relation _ unique _ malformed id",
			req:    httptest.NewRequest(http.MethodGet, "/pets/invalid/owner", nil),
			status: http.StatusBadRequest,
			body:   mustEncode(t, render.NewResponse(http.StatusBadRequest, "id must be an integer greater zero")),
			logs: []map[string]interface{}{
				{
					"level":   "error",
					"msg":     "error getting id from url parameter",
					"handler": "PetHandler",
					"method":  "Owner",
					"id":      "invalid",
				},
			},
		},
		{
			name:   "relation _ unique _ not found",
			req:    httptest.NewRequest(http.MethodGet, "/pets/10000/owner", nil),
			status: http.StatusNotFound,
			body:   mustEncode(t, render.NewResponse(http.StatusNotFound, "owner not found")),
			logs: []map[string]interface{}{
				{
					"level":   "info",
					"msg":     "owner not found",
					"handler": "PetHandler",
					"method":  "Owner",
					"id":      10000,
				},
			},
		},
		{
			name:   "relation _ unique _ ok",
			req:    httptest.NewRequest(http.MethodGet, "/pets/1/owner", nil),
			status: http.StatusOK,
			fn: func(t *testing.T, tt *test, b []byte) {
				e, err := c.Pet.Query().Where(pet.ID(1)).QueryOwner().Only(context.Background())
				require.NoError(t, err)

				var a ent.Owner
				require.NoError(t, json.Unmarshal(b, &a))

				require.Equal(t, e.ID, a.ID)
				require.Equal(t, e.Name, a.Name)
				require.Equal(t, e.Age, a.Age)
			},
			logs: []map[string]interface{}{
				{
					"level":   "info",
					"msg":     "owner rendered",
					"handler": "PetHandler",
					"method":  "Owner",
				},
			},
		},
		{
			name:   "relation _ many _ malformed page",
			req:    httptest.NewRequest(http.MethodGet, "/pets/1/friends?page=invalid", nil),
			status: http.StatusBadRequest,
			body:   mustEncode(t, render.NewResponse(http.StatusBadRequest, "page must be an integer greater zero")),
			logs: []map[string]interface{}{
				{
					"level":   "info",
					"msg":     "error parsing query parameter 'page'",
					"handler": "PetHandler",
					"method":  "Friends",
					"page":    "invalid",
				},
			},
		},
		{
			name:   "relation _ many _ malformed itemsPerPage",
			req:    httptest.NewRequest(http.MethodGet, "/pets/1/friends?itemsPerPage=invalid", nil),
			status: http.StatusBadRequest,
			body:   mustEncode(t, render.NewResponse(http.StatusBadRequest, "itemsPerPage must be an integer greater zero")),
			logs: []map[string]interface{}{
				{
					"level":        "info",
					"msg":          "error parsing query parameter 'itemsPerPage'",
					"handler":      "PetHandler",
					"method":       "Friends",
					"itemsPerPage": "invalid",
				},
			},
		},
		{
			name:   "relation _ many _ ok 1",
			req:    httptest.NewRequest(http.MethodGet, "/pets/1/friends", nil),
			status: http.StatusOK,
			fn: func(t *testing.T, tt *test, b []byte) {
				a, err := c.Pet.Query().Where(pet.ID(1)).QueryFriends().All(context.Background())
				require.NoError(t, err)

				var j []*ent.Pet
				require.NoError(t, json.Unmarshal(b, &j))

				for i := range a {
					require.Equal(t, a[i].ID, j[i].ID)
				}
			},
			logs: []map[string]interface{}{
				{
					"level":   "info",
					"msg":     "pets rendered",
					"handler": "PetHandler",
					"method":  "Friends",
				},
			},
		},
		{
			name:   "relation _ many _ ok 2",
			req:    httptest.NewRequest(http.MethodGet, "/pets/2/friends", nil),
			status: http.StatusOK,
			fn: func(t *testing.T, tt *test, b []byte) {
				a, err := c.Pet.Query().Where(pet.ID(2)).QueryFriends().All(context.Background())
				require.NoError(t, err)

				var j []*ent.Pet
				require.NoError(t, json.Unmarshal(b, &j))

				for i := range a {
					require.Equal(t, a[i].ID, j[i].ID)
				}
			},
			logs: []map[string]interface{}{
				{
					"level":   "info",
					"msg":     "pets rendered",
					"handler": "PetHandler",
					"method":  "Friends",
				},
			},
		},
		{
			name:   "relation _ many _ custom page and itemsPerPage ok",
			req:    httptest.NewRequest(http.MethodGet, "/pets/4/friends?page=1&itemsPerPage=1", nil),
			status: http.StatusOK,
			fn: func(t *testing.T, tt *test, b []byte) {
				a, err := c.Pet.Query().Where(pet.ID(4)).QueryFriends().First(context.Background())
				require.NoError(t, err)

				var j []*ent.Pet
				require.NoError(t, json.Unmarshal(b, &j))
				require.Len(t, j, 1)
				require.EqualValues(t, a.ID, j[0].ID)
			},
			logs: []map[string]interface{}{
				{
					"level":   "info",
					"msg":     "pets rendered",
					"handler": "PetHandler",
					"method":  "Friends",
					"amount":  1,
				},
			},
		},
		// TODO: Test sheriff
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logs.Reset()
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, tt.req)
			rsp := rec.Result()
			// expected status code
			require.Equal(t, tt.status, rsp.StatusCode, logs)
			b, err := io.ReadAll(rsp.Body)
			require.NoError(t, err)
			// expected body
			if tt.body != nil {
				require.Equalf(t, tt.body, b, "expected: %s\nactual  : %s\nlogs    :%s", tt.body, b, logs)
			}
			// if a func to run on response is given run it.
			if tt.fn != nil {
				tt.fn(t, &tt, b)
			}
			// If logs are given check that they indeed are present in the correct order
			if tt.logs != nil {
				// Read logs line by line.
				for i, s := range bytes.Split(bytes.TrimSpace(logs.Bytes()), []byte("\n")) {
					var j map[string]interface{}
					require.NoError(t, json.Unmarshal(s, &j))
					for k, e := range tt.logs[i] {
						v, ok := j[k]
						require.True(t, ok, "log entry not existing: %s: %v", k, e)
						require.EqualValues(t, e, v)
					}
				}
			}
		})
	}
}

func mustEncode(t *testing.T, d interface{}) []byte {
	r, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("Cannot json encode data: %s", err)
	}
	return r
}
