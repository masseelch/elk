{{ define "handler/read/route" }}h.Get("/{id{{ if $.ID.IsInt }}:\\d+{{ end }}}", h.Read){{ end }}

{{ define "handler/read" }}
    // This function fetches the {{ $.Name }} model identified by a give url-parameter from
    // database and returns it to the client.
    func(h {{ $.Name }}Handler) Read(w http.ResponseWriter, r *http.Request) {
        {{- if $.ID.IsInt }}
            id, err := h.urlParamInt(w, r, "id")
        {{ else }}
            id, err := h.urlParamString(w, r, "id")
        {{ end -}}
        if err != nil {
            return
        }

        {{/* cast to go-type if needed */}}
        {{- if $.ID.HasGoType }}
            _id := {{ $.ID.Type }}(id)
        {{ end }}

        q := h.Client.{{ $.Name }}.Query().Where({{ $.Name | lower }}.ID({{ if $.ID.HasGoType }}_id{{ else }}id{{ end }}))
        {{ template "read/qb" $ }}
        e, err := q.Only(r.Context())
        {{ template "read/error-handling" $ }}

        {{ $groups := $.Annotations.HandlerGen.ReadGroups }}
        d, err := sheriff.Marshal(&sheriff.Options{Groups: []string{
            {{- if $groups }}
                {{- range $g := $groups}}"{{$g}}",{{ end -}}
            {{ else -}}
                "{{ $.Name | snake }}"
            {{- end -}}
        }}, e)
        if err != nil {
            h.Logger.WithError(err).WithField("{{ $.Name }}.{{ $.ID.Name }}", id).Error("serialization error")
            render.InternalServerError(w, r, nil)
            return
        }

        h.Logger.WithField("{{ $.Name | snake }}", e.ID).Info("{{ $.Name | snake }} rendered")
        render.OK(w, r, d)
    }
{{end}}