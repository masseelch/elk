package spec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
)

func (spec Spec) MarshalJSON() ([]byte, error) {
	type Local Spec
	return json.Marshal(struct {
		Local
		Version string `json:"openapi"`
	}{
		Local:   Local(spec),
		Version: version,
	})
}

func (f Field) MarshalJSON() ([]byte, error) {
	type local Field
	j, err := json.Marshal(local(f))
	if err != nil {
		return nil, err
	}
	if f.Unique {
		return j, nil
	}
	return []byte(fmt.Sprintf(`{"type":"array","items":%s}`, string(j))), nil
}

func (o MediaTypeObject) MarshalJSON() ([]byte, error) {
	if o.Ref != nil {
		ref := fmt.Sprintf(`{"$ref":"#/components/schemas/%s"}`, o.Ref.Name)
		if !o.Unique {
			ref = fmt.Sprintf(`{"type":"array","items":%s}`, ref)
		}
		return []byte(fmt.Sprintf(`{"schema":%s}`, ref)), nil
	}
	type local MediaTypeObject
	return json.Marshal(local(o))
}

func (r OperationResponse) MarshalJSON() ([]byte, error) {
	if r.Ref != nil {
		return []byte(fmt.Sprintf(`{"$ref":"#/components/responses/%s"}`, r.Ref.Name)), nil
	}
	return json.Marshal(r.Response)
}

func (fs Fields) required() []string {
	var rs []string
	for n, f := range fs {
		if f.Required {
			rs = append(rs, n)
		}
	}
	sort.Strings(rs)
	return rs
}

func (s Schema) MarshalJSON() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.WriteString(`{"type":"object",`)
	// Add the required section.
	if r := s.Fields.required(); len(r) > 0 {
		j, err := json.Marshal(r)
		if err != nil {
			return nil, err
		}
		buf.WriteString(`"required":`)
		buf.Write(j)
		buf.WriteString(",")
	}
	// Properties
	var err error
	props := make(map[string]json.RawMessage, len(s.Fields)+len(s.Edges))
	for n, f := range s.Fields {
		props[n], err = json.Marshal(f)
		if err != nil {
			return nil, err
		}
	}
	for n, e := range s.Edges {
		props[n], err = json.Marshal(e)
		if err != nil {
			return nil, err
		}
	}
	fs, err := json.Marshal(props)
	if err != nil {
		return nil, err
	}
	buf.WriteString(fmt.Sprintf(`"properties":%s}`, fs))
	return buf.Bytes(), nil
}

func (e Edge) MarshalJSON() ([]byte, error) {
	if e.Ref != nil {
		ref := fmt.Sprintf(`{"$ref":"#/components/schemas/%s"}`, e.Ref.Name)
		if e.Unique {
			return []byte(ref), nil
		}
		return []byte(fmt.Sprintf(`{"type":"array","items":%s}`, ref)), nil
	}
	return json.Marshal(e.Schema)
}
