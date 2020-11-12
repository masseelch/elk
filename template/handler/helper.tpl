{{ define "id-from-request-param" }}
    idp := chi.URLParam(r, "id")
    if idp == "" {
        h.logger.WithField("id", idp).Info("empty 'id' url param")
        render.BadRequest(w, r, "id cannot be ''")
        return
    }
    {{ if $.ID.HasGoType -}}
        id := {{ $.ID.Type.String }}(idp)
    {{ else if $.ID.IsString -}}
        id := idp
    {{ else if $.ID.IsInt -}}
        id, err := strconv.Atoi(idp)
        if err != nil {
            h.logger.WithField("id", idp).Info("error parsing url parameter 'id'")
            render.BadRequest(w, r, "id must be a positive integer greater zero")
            return
        }
    {{- end}}
{{ end }}

{{ define "read/qb" }}
    {{/* If one of the given handler groups is set on the edge eager join it.*/}}
    {{/* todo - nested eager loading? */}}
    {{- range $e := $.Edges }}
        {{- range $g := $.Annotations.HandlerGen.ReadGroups }}
            {{- range $eg := split (tagLookup $e.StructTag "groups") "," }}
                {{- if eq $g $eg }}qb.With{{ pascal $e.Name }}(){{ end -}}
            {{- end }}
        {{- end }}
    {{- end }}
{{ end }}

{{ define "read/error-handling" }}
    if err != nil {
        switch err.(type) {
            case *ent.NotFoundError:
                h.logger.WithError(err).WithField("{{ $.Name }}.{{ $.ID.Name }}", id).Debug("job not found")
                render.NotFound(w, r, err)
                return
            case *ent.NotSingularError:
                h.logger.WithError(err).WithField("{{ $.Name }}.{{ $.ID.Name }}", id).Error("duplicate entry for id")
                render.InternalServerError(w, r, nil)
                return
            default:
                h.logger.WithError(err).WithField("{{ $.Name }}.{{ $.ID.Name }}", id).Error("error fetching node from db")
                render.InternalServerError(w, r, nil)
                return
        }
    }
{{ end }}
