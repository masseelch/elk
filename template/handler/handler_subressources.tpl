{{ define "handler/subresource/get/route" -}}
    {{ range $e := $.Edges -}}
        {{- if not $e.Type.Annotations.HandlerGen.Skip }}
            h.Get("/{id{{ if $.ID.IsInt }}:\\d+{{ end }}}/{{ replace ($e.Name | snake) "_" "-" }}", h.{{ $e.Name | pascal }})
        {{ end -}}
    {{ end -}}
{{ end -}}

{{ define "handler/subresource/get" }}
    {{ range $e := $.Edges }}
        {{ if not $e.Type.Annotations.HandlerGen.Skip }}
            {{/* Read/List operations on subressources */}}
            func(h {{ $.Name }}Handler) {{ $e.Name | pascal }}(w http.ResponseWriter, r *http.Request) {
                {{- if $.ID.IsInt }}
                    id, err := h.URLParamInt(w, r, "id")
                {{ else }}
                    id, err := h.URLParamString(w, r, "id")
                {{ end -}}
                if err != nil {
                    return
                }

                q := h.Client.{{ $.Name }}.Query().Where({{ $.Name | lower }}.ID(id)).Query{{ $e.Name | pascal }}()

                {{ if $e.Unique }}
                    {{ template "read/qb" $e.Type }}
                    e, err := q.Only(r.Context())
                    {{ template "read/error-handling" $e.Type }}

                    {{ $groups := $e.Type.Annotations.HandlerGen.ReadGroups }}
                    d, err := sheriff.Marshal(&sheriff.Options{Groups: []string{
                        {{- if $groups }}
                            {{- range $g := $groups}}"{{$g}}",{{ end -}}
                        {{ else -}}
                            "{{ $e.Type.Name | snake }}"
                        {{- end -}}
                    }, IncludeEmptyTag: true},  e)
                    if err != nil {
                        h.Logger.WithError(err).WithField("{{ $e.Type.Name }}.{{ $e.Type.ID.Name }}", id).Error("serialization error")
                        render.InternalServerError(w, r, nil)
                        return
                    }

                    h.Logger.WithField("{{ $e.Type.Name | snake }}", e.ID).Info("{{ $e.Type.Name | snake }} rendered")
                    render.OK(w, r, d)
                {{ else }}
                    {{ if $do := $e.Annotations.EdgeGen.DefaultOrder }}
                        if r.URL.Query().Get("order") == "" {
                            q.Order(
                                {{- range $o := $do -}}
                                    ent.{{ if eq ($o.Order | lower) "desc" }}Desc{{ else }}Asc{{ end }}("{{ $o.Field }}"),
                                {{- end -}}
                            )
                        }
                    {{ else if $do := $e.Type.Annotations.HandlerGen.DefaultListOrder }}
                        if r.URL.Query().Get("order") == "" {
                            q.Order(
                                {{- range $o := $do -}}
                                    ent.{{ if eq ($o.Order | lower) "desc" }}Desc{{ else }}Asc{{ end }}("{{ $o.Field }}"),
                                {{- end -}}
                            )
                        }
                    {{ end }}

                    {{- $elb := eagerLoadBuilder $e.Type "ListGroups" "q" nil nil }}
                    {{- if $elb }}{{ $elb }}{{ end }}

                    // Pagination
                    page, itemsPerPage, err := h.paginationInfo(w, r)
                    if err != nil {
                        return
                    }

                    q.Limit(itemsPerPage).Offset((page - 1) * itemsPerPage)

                    {{ template "handler/list/query-filter" $e.Type }}

                    es, err := q.All(r.Context())
                    if err != nil {
                        h.Logger.WithError(err).Error("error querying database") // todo - better error
                        render.InternalServerError(w, r, nil)
                        return
                    }

                    {{ $groups := $e.Type.Annotations.HandlerGen.ListGroups }}
                    d, err := sheriff.Marshal(&sheriff.Options{Groups: []string{
                        {{- if $groups }}
                            {{- range $g := $groups}}"{{$g}}",{{ end -}}
                        {{ else -}}
                            "{{ $e.Type.Name | snake }}"
                        {{- end -}}
                    }, IncludeEmptyTag: true},  es)
                    if err != nil {
                        h.Logger.WithError(err).Error("serialization error")
                        render.InternalServerError(w, r, nil)
                        return
                    }

                    h.Logger.WithField("amount", len(es)).Info("{{ $e.Type.Name | snake }} rendered")
                    render.OK(w, r, d)
                {{ end }}
            }
        {{ end }}

        {{/* TODO: Create operations on subressources */}}
    {{ end }}
{{ end }}