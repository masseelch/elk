package spec

import (
	"encoding/json"
	"fmt"
)

type ParameterPlace uint

const (
	QueryParam ParameterPlace = iota
	HeaderParam
	PathParam
	CookieParam
)

func (p ParameterPlace) MarshalJSON() ([]byte, error) {
	switch p {
	case QueryParam:
		return json.Marshal("query")
	case HeaderParam:
		return json.Marshal("header")
	case PathParam:
		return json.Marshal("path")
	case CookieParam:
		return json.Marshal("cookie")
	default:
		return nil, fmt.Errorf("cannot marshal %v to json", p)
	}
}

func (p *ParameterPlace) UnmarshalJSON(d []byte) error {
	var s string
	if err := json.Unmarshal(d, &s); err != nil {
		return err
	}
	switch s {
	case "QueryParam":
		*p = QueryParam
	case "HeaderParam":
		*p = HeaderParam
	case "PathParam":
		*p = PathParam
	case "CookieParam":
		*p = CookieParam
	default:
		return fmt.Errorf("cannot unmarshal %v to ParameterPlace", p)
	}
	return nil
}
