package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// const (
// 	EUR Unit = "EUR"
// 	USD Unit = "USD"
// 	SKR Unit = "SKR"
// )
//
// type (
// 	Unit     string
// 	Currency struct {
// 		Amount int
// 		Unit   Unit
// 	}
// )
//
// func (c Currency) Value() (driver.Value, error) {
// 	return fmt.Sprintf("%d|%s", c.Amount, c.Unit), nil
// }
//
// func (c *Currency) Scan(v interface{}) error {
// 	if v == nil {
// 		*c = Currency{}
// 		return nil
// 	}
// 	sv, err := driver.String.ConvertValue(v)
// 	if err != nil {
// 		return err
// 	}
// 	var spl []string
// 	switch v := sv.(type) {
// 	case string:
// 		spl = strings.Split(v, "|")
// 	case []byte:
// 		spl = strings.Split(string(v), "|")
//
// 	}
// 	a, err := strconv.Atoi(spl[0])
// 	if err != nil {
// 		return err
// 	}
// 	*c = Currency{
// 		Amount: a,
// 		Unit:   Unit(spl[1]),
// 	}
// 	return nil
// }

// Category holds the schema definition for the Category entity.
type Category struct {
	ent.Schema
}

// Fields of the Category.
func (Category) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		// field.Other("price", &Currency{}).
		// 	SchemaType(map[string]string{
		// 		dialect.SQLite: "varchar",
		// 	}),
	}
}

// Edges of the Category.
func (Category) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("pets", Pet.Type),
	}
}
