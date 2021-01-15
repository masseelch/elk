{{ define "client" }}
    {{ template "header" -}}
    import 'package:dio/dio.dart';
    import 'package:flutter/widgets.dart';
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
    const {{ $.Name | snake }}Url = '{{ (replace ($.Name | snake) "_" "-") | plural }}';

    {{/* The client for a model. Consumes the generated api. */}}
    class {{ $.Name }}Client {
        {{ $.Name }}Client({@required this.dio}) : assert(dio != null);

        final Dio dio;

        {{/* Find a single node by id. */}}
        Future<{{ $.Name }}> find({{ $.ID.Type | dartType }} id) async {
            final r = await dio.get('/${{ $.Name | snake }}Url/$id');
            return {{ $.Name }}.fromJson(r.data);
        }

        {{/* List multiple nodes filtered by query params. */}}
        Future<List<{{ $.Name }}>> list({
            int page,
            int itemsPerPage,
            {{- range $f := $.Fields }}
                {{- $jsonName := index (split (tagLookup $f.StructTag "json") ",") 0 }}
                    {{ if not (eq $jsonName "-") -}}
                        {{ $f.Type | dartType }} {{ $f.StructField | camel }},
                    {{ end }}
            {{ end }}
        }) async {
            final params = const {};

            if (page != null) {
                params['page'] = page;
            }

            if (itemsPerPage != null) {
                params['itemsPerPage'] = itemsPerPage;
            }

            {{ range $f := $.Fields }}
                {{- $jsonName := index (split (tagLookup $f.StructTag "json") ",") 0 }}
                    {{ if not (eq $jsonName "-") }}
                        if ({{ $f.StructField | camel }} != null) {
                            params['{{ $jsonName }}'] = {{ $f.StructField | camel }};
                        }
                    {{ end }}
            {{ end }}

            final r = await dio.get('/${{ $.Name | snake }}Url');

            if (r.data == null) {
                return [];
            }

            return (r.data as List).map((i) => {{ $.Name }}.fromJson(i)).toList();
        }

        {{/* Create a new node on the remote. */}}
        Future<{{ $.Name }}> create({{ $.Name }}CreateRequest req) async {
            final r = await dio.post('/${{ $.Name | snake }}Url', data: req.toJson());
            return ({{ $.Name }}.fromJson(r.data));
        }

        {{/* Update a node on the remote. */}}
        Future<{{ $.Name }}> update({{ $.Name }}UpdateRequest req) async {
            final r = await dio.patch('/${{ $.Name | snake }}Url/${req.{{ $.ID.Name }}}', data: req.toJson());
            return ({{ $.Name }}.fromJson(r.data));
        }

        {{/* Fetch the nodes edges. */}}
        {{ range $e := $.Edges}}
            {{ if or (not $e.Type.Annotations.HandlerGen) (not $e.Type.Annotations.HandlerGen.Skip) }}
                Future<{{ if $e.Unique }}{{ $e.Type.Name }}{{ else }}List<{{ $e.Type.Name }}>{{ end }}> {{ $e.Name | camel }}({{ $.Name }} e) async {
                    final r = await dio.get('/${{ $.Name | snake }}Url/${e.{{ $.ID.Name }}}/${{ $e.Type.Name | snake }}Url');
                    {{ if $e.Unique -}}
                        return ({{ $e.Type.Name }}.fromJson(r.data));
                    {{ else -}}
                        return (r.data as List).map((i) => {{ $e.Type.Name }}.fromJson(i)).toList();
                    {{ end -}}
                }
            {{ end }}
        {{ end }}

        {{/* Make this node acceessible by the dart provider package. */}}
        static {{ $.Name }}Client of(BuildContext context) => Provider.of<{{ $.Name }}Client>(context, listen: false);
    }


    {{/* The message used to create a new model on the remote. */}}
    {{ $dfc := dartRequestFields $.Type "SkipCreate" }}
    @JsonSerializable(createFactory: false)
    @DateUtcConverter()
    class {{ $.Name }}CreateRequest {
        {{ $.Name }}CreateRequest({
            {{ range $dfc -}}
                this.{{ .StructField | camel }},
            {{ end -}}
        });

        {{ $.Name }}CreateRequest.from{{ $.Name }}({{ $.Name }} e) :
            {{ range $i, $f := $dfc -}}
                {{ $f.StructField | camel }} = e.
                {{- if $f.IsEdge }}edges?.{{ end -}}
                {{ $f.StructField | camel }}
                {{- if $f.IsEdge }}?.
                    {{- if $f.Edge.Unique -}}
                        {{ $f.Edge.Type.ID.Name }}
                    {{- else -}}
                        map((e) => e.{{ $f.Edge.Type.ID.Name }})?.toList()
                    {{- end -}}
                {{ end -}}
                {{- if not (eq $i (dec (len $dfc))) }},{{ end }}
            {{- end -}}
        ;

        {{ range $dfc -}}
            {{ if .Converter }}{{ .Converter }}{{ end -}}
            {{ .Type }} {{ .StructField | camel }};
        {{ end }}

        Map<String, dynamic> toJson() => _${{ $.Name }}CreateRequestToJson(this);
    }

    {{/* The message used to update a model on the remote. */}}
    {{ $dfu := dartRequestFields $.Type "SkipUpdate" }}
    @JsonSerializable(createFactory: false)
    @DateUtcConverter()
    class {{ $.Name }}UpdateRequest {
        {{ $.Name }}UpdateRequest({
            this.{{ $.ID.Name }},
            {{ range $dfu -}}
                this.{{ .StructField | camel }},
            {{ end -}}
        });

        {{ $.Name }}UpdateRequest.from{{ $.Name }}({{ $.Name }} e) :
            {{ $.ID.Name }} = e.{{ $.ID.Name }}{{ if len $dfu }},{{ end }}
            {{ range $i, $f := $dfu -}}
                {{ $f.StructField | camel }} = e.
                {{- if $f.IsEdge }}edges?.{{ end -}}
                {{ $f.StructField | camel }}
                {{- if $f.IsEdge }}?.
                    {{- if $f.Edge.Unique -}}
                        {{ $f.Edge.Type.ID.Name }}
                    {{- else -}}
                        map((e) => e.{{ $f.Edge.Type.ID.Name }})?.toList()
                    {{- end -}}
                {{ end -}}
                {{- if not (eq $i (dec (len $dfu))) }},{{ end }}
            {{- end -}}
        ;

        {{ $.ID.Type | dartType }} {{ $.ID.Name }};
        {{ range $dfu -}}
            {{ if .Converter }}{{ .Converter }}{{ end -}}
            {{ .Type }} {{ .StructField | camel }};
        {{ end }}

        Map<String, dynamic> toJson() => _${{ $.Name }}UpdateRequestToJson(this);
    }
{{ end }}
