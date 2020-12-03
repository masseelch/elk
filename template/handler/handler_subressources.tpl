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
                    id, err := h.urlParamInt(w, r, "id")
                {{ else }}
                    id, err := h.urlParamString(w, r, "id")
                {{ end -}}
                if err != nil {
                    return
                }

                q := h.client.{{ $.Name }}.Query().Where({{ $.Name | snake }}.ID(id)).Query{{ $e.Name | pascal }}()

                {{ if $e.Unique }}
                    {{ template "read/qb" $e.Type }}
                    e, err := q.Only(r.Context())
                    {{ template "read/error-handling" $e.Type }}

                    {{ $groups := $e.Type.Annotations.HandlerGen.ReadGroups }}
                    d, err := sheriff.Marshal(&sheriff.Options{Groups: []string{
                        {{- if $groups }}
                            {{- range $g := $groups}}"{{$g}}",{{ end -}}
                        {{ else -}}
                            "{{ $e.Type.Name | snake }}:read"
                        {{- end -}}
                    }}, e)
                    if err != nil {
                        h.logger.WithError(err).WithField("{{ $e.Type.Name }}.{{ $e.Type.ID.Name }}", id).Error("serialization error")
                        render.InternalServerError(w, r, nil)
                        return
                    }

                    h.logger.WithField("{{ $e.Type.Name | snake }}", e.ID).Info("{{ $e.Type.Name | snake }} rendered")
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

                    {{- $es := eagerLoadedEdges $e.Type "ListGroups" }}
                    {{ if $es }}
                        // Eager load edges.
                        q
                        {{- range $e := $es -}}
                            {{ if not (eq $e.Type $) }}.With{{ pascal $e.Name }}(
                                {{- if $do := $e.Type.Annotations.HandlerGen.DefaultListOrder -}}
                                    func(q *{{ $.Config.Package | base }}.{{ $e.Type.Name }}Query) {
                                        q.Order(
                                            {{- range $o := $do -}}
                                                ent.{{ if eq ($o.Order | lower) "desc" }}Desc{{ else }}Asc{{ end }}("{{ $o.Field }}"),
                                            {{- end -}}
                                        )
                                    }
                                {{- end -}}
                            ){{ end }}
                        {{- end }}
                    {{ end }}

                    // Pagination
                    page, itemsPerPage, err := h.paginationInfo(w, r)
                    if err != nil {
                        return
                    }

                    q.Limit(itemsPerPage).Offset((page - 1) * itemsPerPage)

                    {{ template "handler/list/query-filter" $e.Type }}

                    es, err := q.All(r.Context())
                    if err != nil {
                        h.logger.WithError(err).Error("error querying database") // todo - better error
                        render.InternalServerError(w, r, nil)
                        return
                    }

                    {{ $groups := $e.Type.Annotations.HandlerGen.ListGroups }}
                    d, err := sheriff.Marshal(&sheriff.Options{Groups: []string{
                        {{- if $groups }}
                            {{- range $g := $groups}}"{{$g}}",{{ end -}}
                        {{ else -}}
                            "{{ $e.Type.Name | snake }}:list"
                        {{- end -}}
                    }}, es)
                    if err != nil {
                        h.logger.WithError(err).Error("serialization error")
                        render.InternalServerError(w, r, nil)
                        return
                    }

                    h.logger.WithField("amount", len(es)).Info("{{ $e.Type.Name | snake }} rendered")
                    render.OK(w, r, d)
                {{ end }}
            }
        {{ end }}

        {{/* Create operations on subressources */}}
    {{ end }}
{{ end }}