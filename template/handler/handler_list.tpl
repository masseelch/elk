{{ define "handler/list/route" }}h.Get("/", h.List){{ end }}

{{ define "handler/list/query-filter" }}
    // Use the query parameters to filter the query. todo - nested filter?
    {{- range $f := $.Fields }}
        {{- $jsonName := index (split (tagLookup $f.StructTag "json") ",") 0 }}
        if f := r.URL.Query().Get("{{ $jsonName }}"); f != "" {
            {{- if $f.HasGoType }}
                // todo
            {{else if $f.IsBool }}
                var b bool
                if f == "true" {
                    b = true
                } else if f == "false" {
                    b = false
                } else {
                    h.logger.WithError(err).WithField("{{ $jsonName }}", f).Debug("could not parse query parameter")
                    render.BadRequest(w, r, "'{{ $jsonName }}' must be 'true' or 'false'")
                    return
                }
                q.Where({{ $.Package }}.{{$f.StructField}}(b))
            {{ else if $f.IsInt }}
                i, err := strconv.Atoi(f)
                if err != nil {
                    h.logger.WithError(err).WithField("{{ $jsonName }}", f).Debug("could not parse query parameter")
                    render.BadRequest(w, r, "'{{ $jsonName }}' must be an integer")
                    return
                }
                q.Where({{ $.Package }}.{{$f.StructField}}(i))
            {{ else if $f.IsString }}
                q.Where({{ $.Package }}.{{$f.StructField}}(f))
            {{ else if $f.IsTime }}
                // todo
            {{ end -}}
        }
    {{ end }}
{{ end }}

{{ define "handler/list" }}
    // This function queries for {{ $.Name }} models. Can be filtered by query parameters.
    func(h {{ $.Name }}Handler) List(w http.ResponseWriter, r *http.Request) {
        q := h.client.{{ $.Name }}.Query()

        {{ if $do := $.Annotations.HandlerGen.DefaultListOrder }}
            if r.URL.Query().Get("order") == "" {
                q.Order(
                    {{- range $o := $do -}}
                        ent.{{ if eq ($o.Order | lower) "desc" }}Desc{{ else }}Asc{{ end }}("{{ $o.Field }}"),
                    {{- end -}}
                )
            }
        {{ end }}

        {{- $es := eagerLoadedEdges $ "ListGroups" }}
        {{ if $es }}
            // Eager load edges.
            q
            {{- range $e := $es -}}
                .With{{ pascal $e.Name }}()
            {{- end }}
        {{ end }}

        // Pagination
        page, itemsPerPage, err := h.paginationInfo(w, r)
        if err != nil {
            return
        }

        q.Limit(itemsPerPage).Offset((page - 1) * itemsPerPage)

        {{ template "handler/list/query-filter" $ }}

        es, err := q.All(r.Context())
        if err != nil {
            h.logger.WithError(err).Error("error querying database") // todo - better error
            render.InternalServerError(w, r, nil)
            return
        }

        {{ $groups := $.Annotations.HandlerGen.ReadGroups }}
        d, err := sheriff.Marshal(&sheriff.Options{Groups: []string{
            {{- if $groups }}
                {{- range $g := $groups}}"{{$g}}",{{ end -}}
            {{ else -}}
                "{{ $.Name | snake }}:list"
            {{- end -}}
        }}, es)
        if err != nil {
            h.logger.WithError(err).Error("serialization error")
            render.InternalServerError(w, r, nil)
            return
        }

        h.logger.WithField("amount", len(es)).Info("{{ $.Name | snake }} rendered")
        render.OK(w, r, d)
    }
{{end}}