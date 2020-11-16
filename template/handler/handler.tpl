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

// Shared handler.
type handler struct {
    *chi.Mux

    client    *ent.Client
    validator *validator.Validate
    logger    *logrus.Logger
}

{{ range $n := $.Nodes }}
    {{ if not $n.Annotations.HandlerGen.Skip }}

        // The {{ $n.Name }}Handler.
        type {{ $n.Name }}Handler struct {
            *handler
        }

        // Create a new {{ $n.Name }}Handler
        func New{{ $n.Name }}Handler(c *ent.Client, v *validator.Validate, log *logrus.Logger) *{{ $n.Name }}Handler {
            h := &{{ $n.Name }}Handler{
                &handler{
                    Mux:         chi.NewRouter(),
                    client:    c,
                    validator: v,
                    logger:    log,
                },
            }

            {{ if not $n.Annotations.HandlerGen.SkipCreate }}{{ template "handler/create/route" }}{{ end }}
            {{ if not $n.Annotations.HandlerGen.SkipRead }}{{ template "handler/read/route" }}{{ end }}
            {{ if not $n.Annotations.HandlerGen.SkipUpdate }}{{ template "handler/update/route" }}{{ end }}
{{/*            {{ if not $n.Annotations.HandlerGen.SkipDelete }}{{ template "handler/delete/route" }}{{ end }}*/}}
            {{ if not $n.Annotations.HandlerGen.SkipList }}{{ template "handler/list/route" }}{{ end }}

            {{/* todo - skip resources */}}
            {{ template "handler/subresource/get/route" $n }}

            return h
        }

        {{ if not $n.Annotations.HandlerGen.SkipCreate }}{{ template "handler/create" $n }}{{ end }}
        {{ if not $n.Annotations.HandlerGen.SkipRead }}{{ template "handler/read" $n }}{{ end }}
        {{ if not $n.Annotations.HandlerGen.SkipUpdate }}{{ template "handler/update" $n }}{{ end }}
{{/*        {{ if not $n.Annotations.HandlerGen.SkipDelete }}{{ template "handler/delete" }}{{ end }}*/}}
        {{ if not $n.Annotations.HandlerGen.SkipList }}{{ template "handler/list" $n }}{{ end }}

        {{/* todo - skip resources */}}
        {{ template "handler/subresource/get" $n }}
    {{ end }}
{{ end }}

{{/* Some helpers */}}
func (h handler) urlParamInt(w http.ResponseWriter, r *http.Request, param string) (id int, err error) {
    p := chi.URLParam(r, param)
    if p == "" {
        h.logger.WithField("param", param).Info("empty url param")
        render.BadRequest(w, r, param + " cannot be ''")
        return
    }

    id, err = strconv.Atoi(p)
    if err != nil {
        h.logger.WithField(param, p).Info("error parsing url parameter")
        render.BadRequest(w, r, param + " must be a positive integer greater zero")
        return
    }

    return
}

func (h handler) urlParamTime(w http.ResponseWriter, r *http.Request, param string) (date time.Time, err error) {
    p := chi.URLParam(r, param)
    if p == "" {
        h.logger.WithField("param", param).Info("empty url param")
        render.BadRequest(w, r, param + " cannot be ''")
        return
    }

    date, err = time.Parse("2006-01-02", p)
    if err != nil {
        h.logger.WithField(param, p).Info("error parsing url parameter")
        render.BadRequest(w, r, param + " must be a valid date in yyyy-mm-dd format")
        return
    }

    return
}

func (h handler) paginationInfo(w http.ResponseWriter, r *http.Request) (page int, itemsPerPage int, err error) {
    page = 1
    itemsPerPage = 30

    if d := r.URL.Query().Get("itemsPerPage"); d != "" {
        itemsPerPage, err = strconv.Atoi(d)
        if err != nil {
            h.logger.WithField("itemsPerPage", d).Info("error parsing query parameter 'itemsPerPage'")
            render.BadRequest(w, r, "itemsPerPage must be a positive integer greater zero")
            return
        }
    }

    if d := r.URL.Query().Get("page"); d != "" {
        page, err = strconv.Atoi(d)
        if err != nil {
            h.logger.WithField("page", d).Info("error parsing query parameter 'page'")
            render.BadRequest(w, r, "page must be a positive integer greater zero")
            return
        }
    }

    return
}

