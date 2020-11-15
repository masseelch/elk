{{ define "client" }}
    {{ template "header" -}}
    import 'package:dio/dio.dart';
    import 'package:flutter/widgets.dart';
    import 'package:provider/provider.dart';

    {{/* Import the node itself and all of the edges target nodes / clients. */}}
    import '../model/{{ $.Name | snake }}.dart';
    {{ range $e := $.Edges -}}
        import '../model/{{ $e.Type.Name | snake }}.dart';
        import '../client/{{ $e.Type.Name | snake }}.dart';
    {{ end }}

    {{/* Make the url of this node accessible to other dart files. */}}
    const {{ $.Name | snake }}Url = '{{ (replace ($.Name | snake) "_" "-") | plural }}';

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
                {{ $f.Type | dartType }} {{ $jsonName }},
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
                if ({{ $jsonName }} != null) {
                    params['{{ $jsonName }}'] = {{ $jsonName }};
                }
            {{ end }}

            final r = await dio.get('/${{ $.Name | snake }}Url');

            if (r.data == null) {
                return [];
            }

            return (r.data as List).map((i) => {{ $.Name }}.fromJson(i)).toList();
        }

        {{/* Create a new node on the remote. */}}
        Future<{{ $.Name }}> create({{ $.Name }} e) async {
            final r = await dio.post('/${{ $.Name | snake }}Url', data: e.toJson());
            return ({{ $.Name }}.fromJson(r.data));
        }

        {{/* Update a node on the remote. */}}
        Future<{{ $.Name }}> update({{ $.Name }} e) async {
            final r = await dio.patch('/${{ $.Name | snake }}Url', data: e.toJson());
            return ({{ $.Name }}.fromJson(r.data));
        }

        {{/* Fetch the nodes edges. */}}
        {{ range $e := $.Edges}}
            Future<{{ if $e.Unique }}{{ $e.Type.Name }}{{ else }}List<{{ $e.Type.Name }}>{{ end }}> {{ $e.Name | camel }}({{ $.Name }} e) async {
                final r = await dio.get('/${{ $.Name | snake }}Url/${e.{{ $.ID.Name }}}/${{ $e.Type.Name | snake }}Url');
                {{ if $e.Unique -}}
                    return ({{ $e.Type.Name }}.fromJson(r.data));
                {{ else -}}
                    return (r.data as List).map((i) => {{ $e.Type.Name }}.fromJson(i)).toList();
                {{ end -}}
            }
        {{ end }}

        {{/* Make this node acceessible by the dart provider package. */}}
        static {{ $.Name }}Client of(BuildContext context) => Provider.of<{{ $.Name }}Client>(context, listen: false);
    }
{{ end }}

