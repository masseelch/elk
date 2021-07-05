package integration

import (
	"context"
	"entgo.io/ent/dialect"
	"github.com/masseelch/elk/internal/integration/petstore/ent"
	"github.com/masseelch/elk/internal/integration/petstore/ent/enttest"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func TestHttpRead(t *testing.T) {
	c := enttest.Open(t, dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1", enttest.WithOptions(ent.Log(t.Log)))
	defer c.Close()

	if err := fixtures(context.Background(), c); err != nil {
		t.Fatalf("Could not load fixtures: %s", err)
	}

	// TODO: Implement tests ...
}
