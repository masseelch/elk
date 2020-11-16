{{ define "handler/subresource/get/route" -}}
    {{ range $e := $.Edges -}}
        h.Get("/{id:\\d+}/{{ replace ($e.Name | snake) "_" "-" }}", h.{{ $e.Name | pascal }})
    {{ end -}}
{{ end -}}

{{ define "handler/subresource/get" }}
    {{ range $e := $.Edges }}
        {{/* Read/List operations on subressources */}}
        func(h {{ $.Name }}Handler) {{ $e.Name | pascal }}(w http.ResponseWriter, r *http.Request) {
            id, err := h.urlParamInt(w, r, "id")
                if err != nil {
                return
            }
            qb := h.client.{{ $.Name }}.Query().Where({{ $.Name | snake }}.ID(id)).Query{{ $e.Name | pascal }}()

            {{ if $e.Unique }}
                {{ template "read/qb" $e.Type }}
                e, err := qb.Only(r.Context())
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
                {{/* todo - this is a list route. Enable query filtering. */}}
                es, err := qb.All(r.Context())
                if err != nil {
                    h.logger.WithError(err).Error("error querying database") // todo - better error
                    render.InternalServerError(w, r, nil)
                    return
                }

                {{ $groups := $e.Type.Annotations.HandlerGen.ReadGroups }}
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

        {{/* Create operations on subressources */}}
    {{ end }}
{{ end }}