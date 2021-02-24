{{ define "handler/create/route" }}h.Post("/", h.Create){{ end }}

{{ define "handler/create" }}
    // struct to bind the post body to.
    type {{ $.Name | camel }}CreateRequest struct {
        {{/* Add all fields that are not excluded. */}}
        {{ range $f := $.Fields -}}
            {{- $a := $f.Annotations.FieldGen }}
            {{- if not (and $a $a.SkipCreate) }}
                {{ $f.StructField }} *{{ $f.Type.String }} `json:"{{ index (split (tagLookup $f.StructTag "json") ",") 0 }}"{{ if $a.CreateValidationTag }} {{ $a.CreateValidationTag }}{{ end }}`
            {{- end }}
        {{- end -}}
        {{/* Add all edges that are not excluded. */}}
        {{- range $e := $.Edges -}}
            {{- $a := $e.Annotations.FieldGen }}
            {{- if and (not $e.Type.Annotations.HandlerGen.Skip) (not (and $a $a.SkipCreate)) }}
                {{ $e.StructField }} {{ if $e.Unique }}*{{ else }}[]{{ end }}{{ $e.Type.ID.Type.String }} `json:"{{ index (split (tagLookup $e.StructTag "json") ",") 0 }}"{{ if $a.CreateValidationTag }} {{ $a.CreateValidationTag }}{{ end }}`
            {{- end -}}
        {{- end }}
    }

    // This function creates a new {{ $.Name }} model and stores it in the database.
    func(h {{ $.Name }}Handler) Create(w http.ResponseWriter, r *http.Request) {
        // Get the post data.
        d := {{ $.Name | snake }}CreateRequest{} // todo - allow form-url-encdoded/xml/protobuf data.
        if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
            h.Logger.WithError(err).Error("error decoding json")
            render.BadRequest(w, r, "invalid json string")
            return
        }

        // Validate the data.
        if err := h.Validator.Struct(d); err != nil {
            if err, ok := err.(*validator.InvalidValidationError); ok {
                h.Logger.WithError(err).Error("error validating request data")
                render.InternalServerError(w, r, nil)
                return
            }

            h.Logger.WithError(err).Info("validation failed")
            render.BadRequest(w, r, err)
            return
        }

        // Save the data.
        b := h.Client.{{ $.Name }}.Create()
        {{- range $f := $.Fields -}}
            {{- $a := $f.Annotations.FieldGen }}
            {{- if not (and $a $a.SkipCreate) }}
                if d.{{ $f.StructField }} != nil {
                    b.Set{{ $f.StructField }}(*d.{{ $f.StructField }}) {{/* todo - what about slice fields that have custom marshallers? */}}
                }
            {{- end -}}
        {{ end }}
        {{- range $e := $.Edges -}}
            {{- $a := $e.Annotations.FieldGen }}
            {{- if and (not $e.Type.Annotations.HandlerGen.Skip) (not (and $a $a.SkipCreate)) }}
                if d.{{ $e.StructField }} != nil {
                    {{- if $e.Unique }}
                        b.{{ $e.MutationSet }}(*d.{{ $e.StructField }})
                    {{- else }}
                        b.{{ $e.MutationAdd }}(d.{{ $e.StructField }}...)
                    {{- end }}
                }
            {{- end -}}
        {{ end }}

        // Store in database.
        e, err := b.Save(r.Context())
        if err != nil {
            h.Logger.WithError(err).Error("error saving {{ $.Name }}")
            render.InternalServerError(w, r, nil)
            return
        }

        // Read new entry.
        q := h.Client.{{ $.Name }}.Query().Where({{ $.Name | snake }}.ID(e.ID))
        {{- range $e := $.Edges }}
            {{ range $g := $.Annotations.HandlerGen.CreateGroups }}
                {{ range $eg := split (tagLookup $e.StructTag "groups") "," }}
                    {{ if eq $g $eg }}q.With{{ pascal $e.Name }}(){{ end }}
                {{ end }}
            {{ end }}
        {{ end }}
        e1, err := q.Only(r.Context())
        if err != nil {
            h.Logger.WithError(err).Error("error reading {{ $.Name }}")
            render.InternalServerError(w, r, nil)
            return
        }

        // Serialize the data.
        {{- $groups := $.Annotations.HandlerGen.CreateGroups }}
        j, err := sheriff.Marshal(&sheriff.Options{Groups: []string{
            {{- if $groups }}
                {{- range $g := $groups}}"{{$g}}",{{ end -}}
            {{ else -}}
                "{{ $.Name | snake }}:read"
            {{- end -}}
        }}, e1)
        if err != nil {
            h.Logger.WithError(err).WithField("{{ $.Name }}.{{ $.ID.Name }}", e.ID).Error("serialization error")
            render.InternalServerError(w, r, nil)
            return
        }

        h.Logger.WithField("{{ $.Name | snake }}", e.ID).Info("{{ $.Name | snake }} rendered")
        render.OK(w, r, j)
    }
{{ end }}