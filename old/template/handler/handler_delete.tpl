{{ define "handler/delete/route" }}h.Mux.Delete("/{id{{ if $.ID.IsInt }}:\\d+{{ end }}}", h.Delete){{ end }}

{{ define "handler/delete" }}
    // This function deletes the {{ $.Name }} model identified by a given url-parameter.
    func(h {{ $.Name }}Handler) Delete(w http.ResponseWriter, r *http.Request) {
        {{- if $.ID.IsInt }}
            id, err := h.URLParamInt(w, r, "id")
        {{ else }}
            id, err := h.URLParamString(w, r, "id")
        {{ end -}}
        if err != nil {
            return
        }

        {{/* cast to go-type if needed */}}
        {{- if $.ID.HasGoType }}
            _id := {{ $.ID.Type }}(id)
        {{ end }}

        if err := h.Client.{{ $.Name }}.DeleteOneID({{ if $.ID.HasGoType }}_id{{ else }}id{{ end }}).Exec(r.Context()); err != nil {
            h.Logger.WithError(err).WithField("{{ $.Name }}.{{ $.ID.Name }}", id).Error("error deleting node from db")
            render.InternalServerError(w, r, nil)
            return
        }

        h.Logger.WithField("{{ $.Name | snake }}", id).Info("{{ $.Name | snake }} deleted")
        render.NoContent(w)
    }
{{end}}