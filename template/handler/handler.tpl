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

{{ if not $.Annotations.HandlerGen.SkipGeneration }}

    // The {{ $.Name }}Handler.
    type {{ $.Name }}Handler struct {
        r *chi.Mux

        client    *ent.Client
        validator *validator.Validate
        logger    *logrus.Logger
    }

    // Create a new {{ $.Name }}Handler
    func New{{ $.Name }}Handler(c *ent.Client, v *validator.Validate, log *logrus.Logger) *{{ $.Name }}Handler {
        return &{{ $.Name }}Handler{
            r:         chi.NewRouter(),
            client:    c,
            validator: v,
            logger:    log,
        }
    }

    // Implement the net/http Handler interface.
    func (h {{ $.Name }}Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        h.r.ServeHTTP(w, r)
    }

    // Enable all endpoints.
    func (h *{{ $.Name }}Handler) EnableAllEndpoints() *{{ $.Name }}Handler {
        h.EnableCreateEndpoint()
        h.EnableReadEndpoint()
        h.EnableUpdateEndpoint()
        h.EnableListEndpoint()
        return h
    }

    // Enable the create operation.
    func (h *{{ $.Name }}Handler) EnableCreateEndpoint() *{{ $.Name }}Handler {
        h.r.Post("/", h.Create)
        return h
    }

    // Enable the read operation.
    func (h *{{ $.Name }}Handler) EnableReadEndpoint() *{{ $.Name }}Handler {
        h.r.Get("/{id:\\d+}", h.Read)
        return h
    }

    // Enable the update operation.
    func (h *{{ $.Name }}Handler) EnableUpdateEndpoint() *{{ $.Name }}Handler {
        h.r.Get("/{id:\\d+}", h.Update)
        return h
    }

    // Enable the list operation.
    func (h *{{ $.Name }}Handler) EnableListEndpoint() *{{ $.Name }}Handler {
        h.r.Get("/", h.List)
        return h
    }

    {{ template "create" $ }}
    {{ template "read" $ }}
    {{ template "update" $ }}
{{/*    {{ template "delete" $ }}*/}}
    {{ template "list" $ }}

{{ end }}