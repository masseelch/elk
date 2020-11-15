{{ define "handler/read/route" }}h.Get("/{id:\\d+}", h.Read){{ end }}

{{ define "handler/read" }}
    // This function fetches the {{ $.Name }} model identified by a give url-parameter from
    // database and returns it to the client.
    func(h {{ $.Name }}Handler) Read(w http.ResponseWriter, r *http.Request) {
        id, err := h.urlParamInt(w, r, "id")
        if err != nil {
            return
        }

        qb := h.client.{{ $.Name }}.Query().Where({{ $.Name | snake }}.ID(id))
        {{ template "read/qb" $ }}
        e, err := qb.Only(r.Context())
        {{ template "read/error-handling" $ }}

        {{ $groups := $.Annotations.HandlerGen.ReadGroups }}
        d, err := sheriff.Marshal(&sheriff.Options{Groups: []string{
            {{- if $groups }}
                {{- range $g := $groups}}"{{$g}}",{{ end -}}
            {{ else -}}
                "{{ $.Name | snake }}:read"
            {{- end -}}
        }}, e)
        if err != nil {
            h.logger.WithError(err).WithField("{{ $.Name }}.{{ $.ID.Name }}", id).Error("serialization error")
            render.InternalServerError(w, r, nil)
            return
        }

        h.logger.WithField("{{ $.Name | snake }}", e.ID).Info("{{ $.Name | snake }} rendered")
        render.OK(w, r, d)
    }
{{end}}