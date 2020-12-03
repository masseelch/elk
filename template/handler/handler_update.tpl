{{ define "handler/update/route" }}h.Patch("/{id{{ if $.ID.IsInt }}:\\d+{{ end }}}", h.Update){{ end }}

{{ define "handler/update" }}
    // struct to bind the post body to.
    type {{ $.Name | camel }}UpdateRequest struct {
        {{/* Add all fields that are not excluded. */}}
        {{ range $f := $.Fields -}}
            {{- $a := $f.Annotations.FieldGen }}
            {{- if not (and $a $a.SkipUpdate) }}
                {{ $f.StructField }} {{ $f.Type.String }} `json:"{{ tagLookup $f.StructTag "json" }}"{{ if $a.UpdateValidationTag }} {{ $a.UpdateValidationTag }}{{ end }}`
            {{- end }}
        {{- end -}}
        {{/* Add all edges that are not excluded. */}}
        {{- range $e := $.Edges -}}
            {{- $a := $e.Annotations.FieldGen }}
            {{- if and (not $e.Type.Annotations.HandlerGen.Skip) (not (and $a $a.SkipUpdate)) }}
                {{ $e.StructField }} {{ if not $e.Unique }}[]{{ end }}{{ $e.Type.ID.Type.String }} `json:"{{ tagLookup $e.StructTag "json" }}"{{ if $a.UpdateValidationTag }} {{ $a.UpdateValidationTag }}{{ end }}`
            {{- end -}}
        {{- end }}
    }

    // This function updates a given {{ $.Name }} model and saves the changes in the database.
    func(h {{ $.Name }}Handler) Update(w http.ResponseWriter, r *http.Request) {
        {{- if $.ID.IsInt }}
            id, err := h.urlParamInt(w, r, "id")
        {{ else }}
            id, err := h.urlParamString(w, r, "id")
        {{ end -}}
        if err != nil {
            return
        }

        // Get the post data.
        d := {{ $.Name | snake }}UpdateRequest{} // todo - allow form-url-encoded/xml/protobuf data.
        if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
            h.logger.WithError(err).Error("error decoding json")
            render.BadRequest(w, r, "invalid json string")
            return
        }

        // Validate the data.
        if err := h.validator.Struct(d); err != nil {
            if err, ok := err.(*validator.InvalidValidationError); ok {
                h.logger.WithError(err).Error("error validating request data")
                render.InternalServerError(w, r, nil)
                return
            }

            h.logger.WithError(err).Info("validation failed")
            render.BadRequest(w, r, err)
            return
        }

        // Save the data.
        b := h.client.{{ $.Name }}.UpdateOneID(id)
        {{- range $f := $.Fields -}}
            {{- $a := $f.Annotations.FieldGen }}
            {{- if not (and $a $a.SkipUpdate) }}.
                Set{{ $f.StructField }}(d.{{ $f.StructField }})
            {{- end -}}
        {{ end }}
        {{- range $e := $.Edges -}}
            {{- $a := $e.Annotations.FieldGen }}
            {{- if and (not $e.Type.Annotations.HandlerGen.Skip) (not (and $a $a.SkipUpdate)) }}.
                {{- if $e.Unique }}
                    Set{{ $e.Type.Name }}ID(d.{{ $e.StructField }})
                {{- else }}
                    Add{{ $e.Type.Name }}IDs(d.{{ $e.StructField }}...) {{/*// todo - remove ids that are not given in the patch-data*/}}
                {{- end }}
            {{- end -}}
        {{ end }}

        // Save in database.
        e, err := b.Save(r.Context())
        if err != nil {
            h.logger.WithError(err).Error("error saving {{ $.Name }}")
            render.InternalServerError(w, r, nil)
            return
        }

        // Serialize the data.
        {{- $groups := $.Annotations.HandlerGen.UpdateGroups }}
        j, err := sheriff.Marshal(&sheriff.Options{Groups: []string{
            {{- if $groups }}
                {{- range $g := $groups}}"{{$g}}",{{ end -}}
            {{ else -}}
                "{{ $.Name | snake }}:read"
            {{- end -}}
        }}, e)
        if err != nil {
            h.logger.WithError(err).WithField("{{ $.Name }}.{{ $.ID.Name }}", e.ID).Error("serialization error")
            render.InternalServerError(w, r, nil)
            return
        }

        h.logger.WithField("{{ $.Name | snake }}", e.ID).Info("{{ $.Name | snake }} rendered")
        render.OK(w, r, j)
    }
{{ end }}