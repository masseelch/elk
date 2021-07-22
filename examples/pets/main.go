package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/masseelch/elk/examples/pets/ent"
	elk "github.com/masseelch/elk/examples/pets/ent/http"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"reflect"
	"strings"
)

func main() {
	// Create the ent client.
	c, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer c.Close()
	// Run the auto migration tool.
	if err := c.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	// Create example data.
	if err := Gen(context.Background(), c); err != nil {
		log.Fatalf("failed creating data: %v", err)
	}
	// Create a zap logger to use - colored with a nice, human readable format.
	lc := zap.NewDevelopmentConfig()
	lc.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	l, err := lc.Build()
	if err != nil {
		log.Fatalf("failed creating logger: %v", err)
	}
	// Validator used by elks handlers.
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	// Create a router.
	r := chi.NewRouter()
	// Hook up our generated handlers.
	r.Route("/pets", func(r chi.Router) {
		elk.NewPetHandler(c, l, v).RegisterHandlers(r)
	})
	r.Route("/users", func(r chi.Router) {
		elk.NewUserHandler(c, l, v).RegisterHandlers(r)
	})
	r.Route("/groups", func(r chi.Router) {
		elk.NewGroupHandler(c, l, v).RegisterHandlers(r)
	})
	// Start listen to incoming requests.
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}

func Gen(ctx context.Context, client *ent.Client) error {
	hub, err := client.Group.
		Create().
		SetName("Github").
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed creating the group: %w", err)
	}
	// Create the admin of the group.
	// Unlike `Save`, `SaveX` panics if an error occurs.
	dan := client.User.
		Create().
		SetAge(29).
		SetName("Dan").
		AddManage(hub).
		SaveX(ctx)

	// Create "Ariel" and its pets.
	a8m := client.User.
		Create().
		SetAge(30).
		SetName("Ariel").
		AddGroups(hub).
		AddFriends(dan).
		SaveX(ctx)
	pedro := client.Pet.
		Create().
		SetName("Pedro").
		SetOwner(a8m).
		SaveX(ctx)
	xabi := client.Pet.
		Create().
		SetName("Xabi").
		SetOwner(a8m).
		SaveX(ctx)

	// Create "Alex" and its pets.
	alex := client.User.
		Create().
		SetAge(37).
		SetName("Alex").
		SaveX(ctx)
	coco := client.Pet.
		Create().
		SetName("Coco").
		SetOwner(alex).
		AddFriends(pedro).
		SaveX(ctx)

	fmt.Println("Pets created:", pedro, xabi, coco)
	// Output:
	// Pets created: Pet(id=1, name=Pedro) Pet(id=2, name=Xabi) Pet(id=3, name=Coco)
	return nil
}
