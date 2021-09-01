package spec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
)

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
	buf.WriteString("{")
	buf.WriteString(`"type":"object",`)
	if r := s.Fields.required(); len(r) > 0 {
		j, err := json.Marshal(r)
		if err != nil {
			return nil, err
		}
		buf.WriteString(`"required":`)
		buf.Write(j)
		buf.WriteString(",")
	}
	j, err := json.Marshal(s.Fields)
	if err != nil {
		return nil, err
	}
	buf.WriteString(fmt.Sprintf(`"properties":%s`, j))
	buf.WriteString("}")
	return buf.Bytes(), nil
}

func (u URL) MarshalJSON() ([]byte, error) {
	uu := url.URL(u)
	return json.Marshal(uu.String())
}
