package spec

import (
	"encoding/json"
	"fmt"
)

type ParameterPlace uint

const (
	InQuery ParameterPlace = iota
	InHeader
	InPath
	InCookie
)

func (p ParameterPlace) MarshalJSON() ([]byte, error) {
	switch p {
	case InQuery:
		return json.Marshal("query")
	case InHeader:
		return json.Marshal("header")
	case InPath:
		return json.Marshal("path")
	case InCookie:
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
	case "InQuery":
		*p = InQuery
	case "InHeader":
		*p = InHeader
	case "InPath":
		*p = InPath
	case "InCookie":
		*p = InCookie
	default:
		return fmt.Errorf("cannot unmarshal %v to ParameterPlace", p)
	}
	return nil
}
