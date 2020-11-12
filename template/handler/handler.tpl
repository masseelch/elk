{{ $pkg := base "handler" }}
{{- with extend $ "Package" "handler" -}}
    {{ template "header" . }}
{{ end }}

import (
    "net/http"
    "strconv"

    "github.com/go-chi/chi"
    "github.com/go-playground/validator/v10"
    "github.com/liip/sheriff"
    "github.com/masseelch/render"
    "github.com/sirupsen/logrus"

    "{{ $.Config.Package }}"
)

{{ range $n := $.Nodes }}
    {{ if not $n.Annotations.HandlerGen.SkipGeneration }}

        // The {{ $n.Name }}Handler.
        type {{ $n.Name }}Handler struct {
            *chi.Mux

            client    *ent.Client
            validator *validator.Validate
            logger    *logrus.Logger
        }

        // Create a new {{ $n.Name }}Handler
        func New{{ $n.Name }}Handler(c *ent.Client, v *validator.Validate, log *logrus.Logger) *{{ $n.Name }}Handler {
            return &{{ $n.Name }}Handler{
                Mux:         chi.NewRouter(),
                client:    c,
                validator: v,
                logger:    log,
            }
        }

        // Enable all endpoints.
        func (h *{{ $n.Name }}Handler) EnableAllEndpoints() *{{ $n.Name }}Handler {
            h.EnableCreateEndpoint()
            h.EnableReadEndpoint()
            h.EnableUpdateEndpoint()
            h.EnableListEndpoint()
            {{ range $e := $n.Edges -}}
                h.Enable{{ $e.Name | pascal }}Endpoint()
            {{ end -}}
            return h
        }

        {{ template "handler/create" $n }}
        {{ template "handler/read" $n }}
        {{ template "handler/update" $n }}
    {{/*    {{ template "delete" $ }}*/}}
        {{ template "handler/list" $n }}

        {{ template "handler/subressources" $n }}
    {{ end }}
{{ end }}