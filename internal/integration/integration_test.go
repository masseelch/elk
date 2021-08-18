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
	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			_, l, r, rec := setup(t)
			req, d := tt.req(t, tt)
			r.ServeHTTP(rec, req)
			tt.check(t, tt, rec, l, d)
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
	checkRegistered := func(t *testing.T, tt *test, rec *httptest.ResponseRecorder, logs *bytes.Buffer, d interface{}) {
		require.NotEqual(t, http.StatusNotFound, rec.Result().StatusCode)
		// TODO: logs
	}
	checkNotRegistered := func(t *testing.T, tt *test, rec *httptest.ResponseRecorder, logs *bytes.Buffer, d interface{}) {
		rsp := rec.Result()
		require.Equal(t, http.StatusNotFound, rsp.StatusCode)
		b := bodyBytes(t, rsp)
		require.Equal(t, "404 page not found\n", string(b))
		// TODO: logs
	}
	checkMethodNotAllowed := func(t *testing.T, tt *test, rec *httptest.ResponseRecorder, logs *bytes.Buffer, d interface{}) {
		require.Equal(t, http.StatusMethodNotAllowed, rec.Result().StatusCode)
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

type (
	// reqFn returns the request to send and some data to pass to the checkFn.
	reqFn func(*testing.T, *test) (*http.Request, interface{})
	// checkFn is the signature of the func to run for a test.
	checkFn func(*testing.T, *test, *httptest.ResponseRecorder, *bytes.Buffer, interface{})
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
	return func(t *testing.T, tt *test, rec *httptest.ResponseRecorder, l *bytes.Buffer, d interface{}) {
		rsp := rec.Result()
		b := bodyBytes(t, rsp)
		// expected status code
		require.Equalf(t, status, rsp.StatusCode, "expected: %s\nactual  : %s\nlogs    :%s", body, b, l)
		// expected body
		if body != nil {
			require.Equalf(t, body, b, "expected: %s\nactual  : %s\nlogs    :%s", body, b, l)
		}
		// If logs are given check that they indeed are present in the correct order
		if logs != nil { // TODO: Improve log check
			// Read logs line by line.
			for i, s := range bytes.Split(bytes.TrimSpace(l.Bytes()), []byte("\n")) {
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
	return func(t *testing.T, tt *test) (*http.Request, interface{}) {
		return httptest.NewRequest(method, path, bytes.NewReader(body)), d
	}
}

func setup(t *testing.T) (*ent.Client, *bytes.Buffer, chi.Router, *httptest.ResponseRecorder) {
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
	return c, logs, r, httptest.NewRecorder()
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
