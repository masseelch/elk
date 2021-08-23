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
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestElk(t *testing.T) {
	var tts []*test
	tts = append(tts, testInvalidJson(t)...)
	tts = append(tts, testMalformedId(t)...)
	tts = append(tts, testNotFound(t)...)
	tts = append(tts, testMount()...)
	tts = append(tts, testValidation(t)...)
	tts = append(tts, testPagination(t)...)
	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			deps := setup(t)
			req, d := tt.req(t, tt, deps)
			deps.router.ServeHTTP(deps.rec, req)
			tt.check(t, tt, deps, d)
		})
	}
}

func testInvalidJson(t *testing.T) []*test {
	name := prefixNameFn("invalid json _ ")
	body := []byte("invalid")
	check := defaultCheckFn(http.StatusBadRequest, mustEncode(t, elkhttp.ErrResponse{
		Code:   http.StatusBadRequest,
		Status: http.StatusText(http.StatusBadRequest),
		Errors: "invalid json string",
	}), nil)
	return []*test{
		{
			name:  name("create"),
			req:   defaultReqFn(http.MethodPost, "/pets", body, nil),
			check: check,
		}, {
			name:  name("update"),
			req:   defaultReqFn(http.MethodPatch, "/pets/1", body, nil),
			check: check,
		},
	} // TODO: logs
}

func testMalformedId(t *testing.T) []*test {
	const path = "/pets/invalid"
	name := prefixNameFn("malformed id _ ")
	check := defaultCheckFn(http.StatusBadRequest, mustEncode(t, elkhttp.ErrResponse{
		Code:   http.StatusBadRequest,
		Status: http.StatusText(http.StatusBadRequest),
		Errors: "id must be an integer greater zero",
	}), nil)
	return []*test{
		{
			name:  name("read"),
			req:   defaultReqFn(http.MethodGet, path, nil, nil),
			check: check,
		}, {
			name:  name("update"),
			req:   defaultReqFn(http.MethodPatch, path, nil, nil),
			check: check,
		}, {
			name:  name("delete"),
			req:   defaultReqFn(http.MethodDelete, path, nil, nil),
			check: check,
		},
	} // TODO: logs
}

func testNotFound(t *testing.T) []*test {
	const path = "/toys/10000"
	name := prefixNameFn("not found _ ")
	check := defaultCheckFn(http.StatusNotFound, mustEncode(t, elkhttp.ErrResponse{
		Code:   http.StatusNotFound,
		Status: http.StatusText(http.StatusNotFound),
		Errors: "toy not found",
	}), nil)
	return []*test{
		{
			name:  name("read"),
			req:   defaultReqFn(http.MethodGet, path, nil, nil),
			check: check,
		}, {
			name:  name("update"),
			req:   defaultReqFn(http.MethodPatch, path, mustEncode(t, map[string]interface{}{}), nil),
			check: check,
		}, {
			name:  name("delete"),
			req:   defaultReqFn(http.MethodDelete, path, nil, nil),
			check: check,
		}, // TODO: sub-resources
	} // TODO: logs
}

func testMount() []*test {
	name := prefixNameFn("mount _ ")
	checkRegistered := func(t *testing.T, tt *test, deps *deps, d interface{}) {
		require.NotEqual(t, http.StatusNotFound, deps.rec.Result().StatusCode)
		// TODO: logs
	}
	checkNotRegistered := func(t *testing.T, tt *test, deps *deps, d interface{}) {
		rsp := deps.rec.Result()
		require.Equal(t, http.StatusNotFound, rsp.StatusCode)
		b := bodyBytes(t, rsp)
		require.Equal(t, "404 page not found\n", string(b))
		// TODO: logs
	}
	checkMethodNotAllowed := func(t *testing.T, tt *test, deps *deps, d interface{}) {
		require.Equal(t, http.StatusMethodNotAllowed, deps.rec.Result().StatusCode)
		// TODO: logs
	}
	return []*test{
		// all routes registered
		{
			name:  name("registered pet create"),
			req:   defaultReqFn(http.MethodGet, "/pets", nil, nil),
			check: checkRegistered,
		}, {
			name:  name("registered pet read"),
			req:   defaultReqFn(http.MethodGet, "/pets/1", nil, nil),
			check: checkRegistered,
		}, {
			name:  name("registered pet update"),
			req:   defaultReqFn(http.MethodPatch, "/pets/1", nil, nil),
			check: checkRegistered,
		}, {
			name:  name("registered pet delete"),
			req:   defaultReqFn(http.MethodDelete, "/pets/1", nil, nil),
			check: checkRegistered,
		}, {
			name:  name("registered pet list"),
			req:   defaultReqFn(http.MethodGet, "/pets", nil, nil),
			check: checkRegistered,
		}, {
			name:  name("registered pet badge read"),
			req:   defaultReqFn(http.MethodGet, "/pets/1/badge", nil, nil),
			check: checkRegistered,
		}, {
			name:  name("registered pet mentor read"),
			req:   defaultReqFn(http.MethodGet, "/pets/1/mentor", nil, nil),
			check: checkRegistered,
		}, {
			name:  name("registered pet spouse read"),
			req:   defaultReqFn(http.MethodGet, "/pets/1/spouse", nil, nil),
			check: checkRegistered,
		}, {
			name:  name("registered pet toys list"),
			req:   defaultReqFn(http.MethodGet, "/pets/1/toys", nil, nil),
			check: checkRegistered,
		}, {
			name:  name("registered pet children list"),
			req:   defaultReqFn(http.MethodGet, "/pets/1/children", nil, nil),
			check: checkRegistered,
		}, {
			name:  name("registered pet play-groups list"),
			req:   defaultReqFn(http.MethodGet, "/pets/1/play-groups", nil, nil),
			check: checkRegistered,
		}, {
			name:  name("registered pet friends list"),
			req:   defaultReqFn(http.MethodGet, "/pets/1/friends", nil, nil),
			check: checkRegistered,
		},
		// no routes registered
		{
			name:  name("not registered badge create"),
			req:   defaultReqFn(http.MethodGet, "/badges", nil, nil),
			check: checkNotRegistered,
		}, {
			name:  name("not registered badge read"),
			req:   defaultReqFn(http.MethodGet, "/badges/1", nil, nil),
			check: checkNotRegistered,
		}, {
			name:  name("not registered badge update"),
			req:   defaultReqFn(http.MethodPatch, "/badges/1", nil, nil),
			check: checkNotRegistered,
		}, {
			name:  name("not registered badge delete"),
			req:   defaultReqFn(http.MethodDelete, "/badges/1", nil, nil),
			check: checkNotRegistered,
		}, {
			name:  name("not registered badge list"),
			req:   defaultReqFn(http.MethodGet, "/badges", nil, nil),
			check: checkNotRegistered,
		}, {
			name:  name("not registered badge wearer read"),
			req:   defaultReqFn(http.MethodGet, "/badges/wearer", nil, nil),
			check: checkNotRegistered,
		},
		// some routes registered
		{
			name:  name("registered play-group read"),
			req:   defaultReqFn(http.MethodGet, "/play-groups/1", nil, nil),
			check: checkRegistered,
		}, {
			name:  name("registered play-group list"),
			req:   defaultReqFn(http.MethodGet, "/play-groups", nil, nil),
			check: checkRegistered,
		}, {
			name:  name("not registered play-group create"),
			req:   defaultReqFn(http.MethodPost, "/play-groups", nil, nil),
			check: checkMethodNotAllowed,
		}, {
			name:  name("not registered play-group update"),
			req:   defaultReqFn(http.MethodPatch, "/play-groups/1", nil, nil),
			check: checkMethodNotAllowed,
		}, {
			name:  name("not registered play-group delete"),
			req:   defaultReqFn(http.MethodDelete, "/play-groups/1", nil, nil),
			check: checkMethodNotAllowed,
		}, {
			name:  name("not registered play-group participants list"),
			req:   defaultReqFn(http.MethodGet, "/play-groups/1/participants", nil, nil),
			check: checkNotRegistered,
		},
	}
}

func testValidation(t *testing.T) []*test {
	name := prefixNameFn("validation _ ")
	return []*test{
		{
			name: name("1"),
			req:  defaultReqFn(http.MethodPost, "/pets", mustEncode(t, make(map[string]interface{})), nil),
			check: func(t *testing.T, tt *test, deps *deps, _ interface{}) {
				rsp := deps.rec.Result()
				require.Equal(t, http.StatusBadRequest, rsp.StatusCode)
				var err elkhttp.ErrResponse
				require.NoError(t, json.Unmarshal(bodyBytes(t, rsp), &err))
				require.Equal(t, map[string]interface{}{
					"badge":     "missing required edge: \"badge\"",
					"castrated": "missing required field: \"castrated\"",
					"height":    "missing required field: \"height\"",
					"name":      "missing required field: \"name\"",
					"sex":       "missing required field: \"sex\"",
				}, err.Errors)
			},
		}, {
			name: name("2"),
			req: defaultReqFn(http.MethodPost, "/pets", mustEncode(t, map[string]interface{}{
				"castrated": true,     // valid
				"sex":       "divers", // invalid - (male or female)
				"name":      "a",      // invalid - too short
				"weight":    0,        // invalid - must be positive
			}), nil),
			check: func(t *testing.T, tt *test, deps *deps, _ interface{}) {
				rsp := deps.rec.Result()
				require.Equal(t, http.StatusBadRequest, rsp.StatusCode)
				var err elkhttp.ErrResponse
				require.NoError(t, json.Unmarshal(bodyBytes(t, rsp), &err))
				require.Equal(t, map[string]interface{}{
					"badge":  "missing required edge: \"badge\"",
					"height": "missing required field: \"height\"",
					"name":   "value is less than the required length",
					"sex":    "invalid enum value for sex field: \"divers\"",
					"weight": "value out of range",
				}, err.Errors)
			},
		},
	}
}

func testPagination(t *testing.T) []*test {
	name := prefixNameFn("pagination _ ")
	return []*test{
		{
			name: name("malformed page"),
			req:  defaultReqFn(http.MethodGet, "/pets?page=invalid", nil, nil),
			check: defaultCheckFn(http.StatusBadRequest, mustEncode(t, elkhttp.ErrResponse{
				Code:   http.StatusBadRequest,
				Status: http.StatusText(http.StatusBadRequest),
				Errors: "page must be an integer greater zero",
			}), nil), // TODO: logs
		}, {
			name: name("malformed itemsPerPage"),
			req:  defaultReqFn(http.MethodGet, "/pets?itemsPerPage=invalid", nil, nil),
			check: defaultCheckFn(http.StatusBadRequest, mustEncode(t, elkhttp.ErrResponse{
				Code:   http.StatusBadRequest,
				Status: http.StatusText(http.StatusBadRequest),
				Errors: "itemsPerPage must be an integer greater zero",
			}), nil), // TODO: logs
		}, {
			name: name("ok"),
			req: func(t *testing.T, tt *test, deps *deps) (*http.Request, interface{}) {
				// Get a list of the first 30 (default page size) pets
				ps, err := deps.client.Pet.Query().Limit(30).All(context.Background())
				require.NoError(t, err)
				return httptest.NewRequest(http.MethodGet, "/pets", nil), ps
			},
			check: func(t *testing.T, tt *test, deps *deps, ps interface{}) {
				rsp := deps.rec.Result()
				require.Equal(t, http.StatusOK, rsp.StatusCode)
				// List must have 30 items and must match that one saved.
				var v []*ent.Pet
				require.NoError(t, json.Unmarshal(bodyBytes(t, rsp), &v))
				require.Len(t, v, 30)
			}, // TODO: logs
		},
	}
}

type (
	// reqFn returns the request to send and some data to pass to the checkFn.
	reqFn func(*testing.T, *test, *deps) (*http.Request, interface{})
	// checkFn is the signature of the func to run for a test.
	checkFn func(*testing.T, *test, *deps, interface{})
	// test describes one test to execute. Every test has its own database.
	test struct {
		// name of the test
		name  string
		req   reqFn
		check checkFn
	}
)

// defaultCheckFn returns a postFn to simply test a responses status-code, body-bytes and logs.
func defaultCheckFn(status int, body []byte, logs []map[string]interface{}) checkFn {
	return func(t *testing.T, tt *test, deps *deps, d interface{}) {
		rsp := deps.rec.Result()
		b := bodyBytes(t, rsp)
		// expected status code
		require.Equalf(t, status, rsp.StatusCode, "expected: %s\nactual  : %s\nlogs    :%s", body, b, deps.logs)
		// expected body
		if body != nil {
			require.Equalf(t, body, b, "expected: %s\nactual  : %s\nlogs    :%s", body, b, deps.logs)
		}
		// If logs are given check that they indeed are present in the correct order
		if logs != nil { // TODO: Improve log check
			// Read logs line by line.
			for i, s := range bytes.Split(bytes.TrimSpace(deps.logs.Bytes()), []byte("\n")) {
				var j map[string]interface{}
				require.NoError(t, json.Unmarshal(s, &j))
				for k, e := range logs[i] {
					v, ok := j[k]
					require.True(t, ok, "log entry not existing: %s: %v", k, e)
					require.EqualValues(t, e, v)
				}
			}
		}
	}
}

// defaultReqFn returns a reqFn to create a very basic http.Request.
func defaultReqFn(method, path string, body []byte, d interface{}) reqFn {
	return func(t *testing.T, tt *test, _ *deps) (*http.Request, interface{}) {
		return httptest.NewRequest(method, path, bytes.NewReader(body)), d
	}
}

type deps struct {
	client *ent.Client
	logs   *bytes.Buffer
	router chi.Router
	rec    *httptest.ResponseRecorder
}

func setup(t *testing.T) *deps {
	// ent client
	c := enttest.Open(t, dialect.SQLite, ":memory:?_fk=1", enttest.WithOptions(ent.Log(t.Log)))
	require.NoError(t, fixtures(context.Background(), c))
	logs := new(bytes.Buffer)
	// logging
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
	// handlers
	r := chi.NewRouter()
	r.Route("/pets", func(r chi.Router) {
		elkhttp.NewPetHandler(c, l).Mount(r, elkhttp.PetRoutes)
	})
	r.Route("/toys", func(r chi.Router) {
		elkhttp.NewToyHandler(c, l).Mount(r, elkhttp.ToyRoutes)
	})
	r.Route("/play-groups", func(r chi.Router) {
		elkhttp.NewToyHandler(c, l).Mount(r, elkhttp.PlayGroupList|elkhttp.PlayGroupRead)
	})
	return &deps{c, logs, r, httptest.NewRecorder()}
}

func mustEncode(t *testing.T, d interface{}) []byte {
	r, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("Cannot json encode data: %s", err)
	}
	return r
}

func prefixNameFn(s string) func(string) string {
	return func(s2 string) string {
		return s + s2
	}
}

func bodyBytes(t *testing.T, rsp *http.Response) []byte {
	b, err := io.ReadAll(rsp.Body)
	require.NoError(t, err)
	return b
}
