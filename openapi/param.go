package openapi

import (
	"encoding/json"
	"fmt"
)

type ParameterPlace uint

const (
	query ParameterPlace = iota
	header
	path
	cookie
)

func (p ParameterPlace) MarshalJSON() ([]byte, error) {
	switch p {
	case query:
		return json.Marshal("query")
	case header:
		return json.Marshal("header")
	case path:
		return json.Marshal("path")
	case cookie:
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
	case "query":
		*p = query
	case "header":
		*p = header
	case "path":
		*p = path
	case "cookie":
		*p = cookie
	default:
		return fmt.Errorf("cannot unmarshal %v to ParameterPlace", p)
	}
	return nil
}
