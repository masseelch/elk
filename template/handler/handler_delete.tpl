{{ define "handler/delete/route" }}h.Mux.Delete("/{id{{ if $.ID.IsInt }}:\\d+{{ end }}}", h.Delete){{ end }}

{{ define "handler/delete" }}
    // This function deletes the {{ $.Name }} model identified by a given url-parameter.
    func(h {{ $.Name }}Handler) Delete(w http.ResponseWriter, r *http.Request) {
        {{- if $.ID.IsInt }}
            id, err := h.urlParamInt(w, r, "id")
        {{ else }}
            id, err := h.urlParamString(w, r, "id")
        {{ end -}}
        if err != nil {
            return
        }

        if err := h.client.{{ $.Name }}.DeleteOneID(id).Exec(r.Context()); err != nil {
            h.logger.WithError(err).WithField("{{ $.Name }}.{{ $.ID.Name }}", id).Error("error deleting node from db")
            render.InternalServerError(w, r, nil)
            return
        }

        h.logger.WithField("{{ $.Name | snake }}", id).Info("{{ $.Name | snake }} deleted")
        render.NoContent(w, r)
    }
{{end}}