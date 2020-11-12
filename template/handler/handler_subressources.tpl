{{ define "handler/subressources" }}
    {{ range $e := $.Edges }}
        {{/* Read/List operations on subressources */}}
        // Enable the read operation on the {{ $e.Name }} edge.
        func (h *{{ $.Name }}Handler) Enable{{ $e.Name | pascal }}Endpoint() *{{ $.Name }}Handler {
            h.Get("/{id:\\d+}/{{ replace ($e.Name | snake) "_" "-" }}", h.{{ $e.Name | pascal }})
            return h
        }

        func(h {{ $.Name }}Handler) {{ $e.Name | pascal }}(w http.ResponseWriter, r *http.Request) {
            {{- template "id-from-request-param" $ }}
            qb := h.client.{{ $.Name }}.Query().Where({{ $.Name | snake }}.ID(id)).Query{{ $e.Name | pascal }}()

            {{ if $e.Unique }}
                {{ template "read/qb" $e.Type }}
                e, err := qb.Only(r.Context())
                {{ template "read/error-handling" $e.Type }}

                {{ $groups := $e.Annotations.HandlerGen.ReadGroups }}
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