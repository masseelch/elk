package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"entgo.io/ent/dialect"
	"github.com/go-chi/chi"
	"github.com/masseelch/elk/internal/integration/petstore/ent"
	"github.com/masseelch/elk/internal/integration/petstore/ent/enttest"
	elkhttp "github.com/masseelch/elk/internal/integration/petstore/ent/http"
	"github.com/masseelch/render"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"net/http"
	"net/http/httptest"
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
	logs []string
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

	r := chi.NewRouter() // Needed to test url param fetching

	// Register pet endpoints.
	ph := elkhttp.NewPetHandler(c, l) // TODO: Provide a default set of routes in generated PetHandler.
	r.Get("/pets", ph.List)
	r.Get("/pets/{id}", ph.Read)

	// Create the tests.
	tests := []test{
		{
			name:   "read _ malformed id",
			req:    httptest.NewRequest(http.MethodGet, "/pets/invalid", nil),
			status: http.StatusBadRequest,
			body:   mustEncode(t, render.NewResponse(http.StatusBadRequest, "id must be an integer greater zero")),
			logs: []string{
				"{\"level\":\"error\",\"msg\":\"error getting id from url parameter\",\"handler\":\"PetHandler\",\"method\":\"Read\",\"id\":\"invalid\",\"error\":\"strconv.Atoi: parsing \\\"invalid\\\": invalid syntax\"}",
			},
		},
		{
			name:   "read _ not found",
			req:    httptest.NewRequest(http.MethodGet, "/pets/10000", nil),
			status: http.StatusNotFound,
			body:   mustEncode(t, render.NewResponse(http.StatusNotFound, "pet not found")),
		},
		{
			name:   "read _ ok",
			req:    httptest.NewRequest(http.MethodGet, "/pets/1", nil),
			status: http.StatusOK,
		},
		{
			name:   "list _ malformed page",
			req:    httptest.NewRequest(http.MethodGet, "/pets?page=invalid", nil),
			status: http.StatusBadRequest,
			body:   mustEncode(t, render.NewResponse(http.StatusBadRequest, "page must be an integer greater zero")),
		},
		{
			name:   "list _ malformed itemsPerPage",
			req:    httptest.NewRequest(http.MethodGet, "/pets?itemsPerPage=invalid", nil),
			status: http.StatusBadRequest,
			body:   mustEncode(t, render.NewResponse(http.StatusBadRequest, "itemsPerPage must be an integer greater zero")),
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logs.Reset()
			rec := httptest.NewRecorder()
			r.ServeHTTP(rec, tt.req)
			rsp := rec.Result()
			// expected status code
			require.Equal(t, tt.status, rsp.StatusCode)
			b, err := io.ReadAll(rsp.Body)
			// expected body
			if tt.body != nil {
				require.NoError(t, err)
				require.Equalf(t, tt.body, b, "expected: %s\nactual  : %s", tt.body, b)
			}
			// if a func to run on response is given run it.
			if tt.fn != nil {
				tt.fn(t, &tt, b)
			}
			// If logs are given check that they indeed are present in the correct order
			if tt.logs != nil {
				var l []map[string]string
				require.NoError(t, json.Unmarshal(logs.Bytes(), &l))
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
