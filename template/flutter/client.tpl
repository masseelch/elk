{{ define "client" }}
    {{  template "header" -}}
    import 'package:flutter/widgets.dart';
    import 'package:intercepted_http/intercepted_http.dart' show Client;
    import 'package:json_annotation/json_annotation.dart';
    import 'package:provider/provider.dart';

    import '../date_utc_converter.dart';

    {{/* Import the custom dart types. */}}
    {{ range $.TypeMappings -}}
        import '{{ .Import }}';
        {{- if .ConverterImport }}import '{{ .ConverterImport }}';{{ end -}}
    {{ end }}

    {{/* Import the node itself and all of the edges target nodes / clients. */}}
    import '../model/{{ $.Name | snake }}.dart';
    {{ range $e := $.Edges -}}
        import '../model/{{ $e.Type.Name | snake }}.dart';
        {{- if or (not $e.Type.Annotations.HandlerGen) (not $e.Type.Annotations.HandlerGen.Skip) }}
            import '../client/{{ $e.Type.Name | snake }}.dart';
        {{ end -}}
    {{ end }}

    {{/* JsonSerializable puts the generated code in this file. */}}
    part '{{ $.Name | snake }}.g.dart';

    {{/* Make the url of this node accessible to other dart files. */}}
    const {{ $.Name | lowerFirst }}Url = '{{ (replace ($.Name | plural | snake) "_" "-") }}';

    {{/* The client for a model. Consumes the generated api. */}}
    class {{ $.Name }}Client {
        {{ $.Name }}Client({required this.client});

        final Client client;

        {{/* Find a single node by id. */}}
        Future<{{ $.Name }}> find({{ $.ID.Type | dartType }} id) async {
            final r = await client.get(Uri(path: '/${{ $.Name | lowerFirst }}Url/$id'));

            return {{ $.Name }}.fromJson(r.body);
        }

        {{/* List multiple nodes filtered by query params. */}}
        Future<List<{{ $.Name }}>> list({
            int? page,
            int? itemsPerPage,
            {{- range $f := $.Fields }}
                {{- $jsonName := index (split (tagLookup $f.StructTag "json") ",") 0 }}
                    {{ if not (eq $jsonName "-") -}}
                        {{ $f.Type | dartType }}? {{ $f.Name | camel }},
                    {{ end }}
            {{ end }}
        }) async {
            final params = <String, dynamic>{};

            if (page != null) {
                params['page'] = page.toString();
            }

            if (itemsPerPage != null) {
                params['itemsPerPage'] = itemsPerPage.toString();
            }

            {{ range $f := $.Fields }}
                {{- $jsonName := index (split (tagLookup $f.StructTag "json") ",") 0 }}
                    {{ if not (eq $jsonName "-") }}
                        if ({{ $f.Name | camel }} != null) {
                            params['{{ $jsonName }}'] = {{ $f.Name | camel }};
                        }
                    {{ end }}
            {{ end }}

            final r = await client.get(Uri(
                path: '/${{ $.Name | lowerFirst }}Url',
                queryParameters: params,
            ));

            if (r.body.isEmpty) {
                return [];
            }

            return (r.body as List).map((i) => {{ $.Name }}.fromJson(i)).toList();
        }

        {{/* Create a new node on the remote. */}}
        Future<{{ $.Name }}> create({{ $.Name }}CreateRequest req) async {
            final r = await client.post(
                Uri(path: '/${{ $.Name | lowerFirst }}Url'),
                body: req.toJson(),
            );

            return ({{ $.Name }}.fromJson(r.body));
        }

        {{/* Update a node on the remote. */}}
        Future<{{ $.Name }}> update({{ $.Name }}UpdateRequest req) async {
            final r = await client.patch(
                Uri(path: '/${{ $.Name | lowerFirst }}Url/${req.{{ $.ID.Name }}}'),
                body: req.toJson(),
            );

            return ({{ $.Name }}.fromJson(r.body));
        }

        {{/* Delete a node on the remote. */}}
        Future delete({{ $.ID.Type | dartType }} id) =>
            client.delete(Uri(path: '/${{ $.Name | lowerFirst }}Url/$id'));

        {{/* Fetch the nodes edges. */}}
        {{ range $e := $.Edges}}
            {{ if or (not $e.Type.Annotations.HandlerGen) (not $e.Type.Annotations.HandlerGen.Skip) }}
                Future<{{ if $e.Unique }}{{ $e.Type.Name }}{{ else }}List<{{ $e.Type.Name }}>{{ end }}> {{ $e.Name | lowerFirst }}({{ $.Name }} e) async {
                    final r = await client.get(Uri(path: '/${{ $.Name | lowerFirst }}Url/${e.{{ $.ID.Name }}}/{{ (replace ($e.Name | snake) "_" "-") }}'));
                    {{ if $e.Unique -}}
                        return ({{ $e.Type.Name }}.fromJson(r.body));
                    {{ else -}}
                        return (r.body as List).map((i) => {{ $e.Type.Name }}.fromJson(i)).toList();
                    {{ end -}}
                }
            {{ end }}
        {{ end }}

        {{/* Make this node acceessible by the dart provider package. */}}
        static {{ $.Name }}Client of(BuildContext context) => context.read<{{ $.Name }}Client>();
    }

    {{/* The message used to create a new model on the remote. */}}
    {{ $dfc := dartRequestFields $.Type "SkipCreate" }}
    @JsonSerializable(createFactory: false, includeIfNull: false)
    @DateUtcConverter()
    class {{ $.Name }}CreateRequest {
        {{ $.Name }}CreateRequest({
            {{ range $dfc -}}
                this.{{ .Name | camel }},
            {{ end -}}
        });

        {{ $.Name }}CreateRequest.from{{ $.Name }}({{ $.Name }} e) :
            {{ range $i, $f := $dfc -}}
                {{ $f.Name | camel }} = e.
                {{- if $f.IsEdge }}edges?.{{ end -}}
                {{ $f.Name | camel }}
                {{- if $f.IsEdge }}?.
                    {{- if $f.Edge.Unique -}}
                        {{ $f.Edge.Type.ID.Name }}
                    {{- else -}}
                        map((e) => e.{{ $f.Edge.Type.ID.Name }}).toList()
                    {{- end -}}
                {{ end -}}
                {{- if not (eq $i (dec (len $dfc))) }},{{ end }}
            {{- end -}}
        ;

        {{ range $dfc -}}
            {{ if .Converter }}{{ .Converter }}{{ end -}}
            {{ .Type }} {{ .Name | camel }};
        {{ end }}

        Map<String, dynamic> toJson() => _${{ $.Name }}CreateRequestToJson(this);
    }

    {{/* The message used to update a model on the remote. */}}
    {{ $dfu := dartRequestFields $.Type "SkipUpdate" }}
    @JsonSerializable(createFactory: false, includeIfNull: false)
    @DateUtcConverter()
    class {{ $.Name }}UpdateRequest {
        {{ $.Name }}UpdateRequest({
            this.{{ $.ID.Name }},
            {{ range $dfu -}}
                {{- if not .Immutable -}}
                    this.{{ .Name | camel }},
                {{- end -}}
            {{ end -}}
        });

        {{ $.Name }}UpdateRequest.from{{ $.Name }}({{ $.Name }} e) :
            {{ $.ID.Name }} = e.{{ $.ID.Name }}{{ if len $dfu }},{{ end }}
            {{ range $i, $f := $dfu -}}
                {{- if not .Immutable -}}
                    {{ $f.Name | camel }} = e.
                    {{- if $f.IsEdge }}edges?.{{ end -}}
                    {{ $f.Name | camel }}
                    {{- if $f.IsEdge }}?.
                        {{- if $f.Edge.Unique -}}
                            {{ $f.Edge.Type.ID.Name }}
                        {{- else -}}
                            map((e) => e.{{ $f.Edge.Type.ID.Name }}).toList()
                        {{- end -}}
                    {{ end -}}
                    {{- if not (eq $i (dec (len $dfu))) }},{{ end }}
                {{- end -}}
            {{- end -}}
        ;

        {{ $.ID.Type | dartType }}? {{ $.ID.Name }};
        {{ range $dfu -}}
            {{- if not .Immutable -}}
                {{ if .Converter }}{{ .Converter }}{{ end -}}
                {{ .Type }} {{ .Name | camel }};
            {{- end -}}
        {{ end }}

        Map<String, dynamic> toJson() => _${{ $.Name }}UpdateRequestToJson(this);
    }
{{ end }}
