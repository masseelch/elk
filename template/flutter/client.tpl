{{ define "client" }}
    {{ template "header" -}}
    {{ $url := (replace ($.Name | snake) "_" "-") | plural }}
    import 'package:dio/dio.dart';
    import 'package:flutter/widgets.dart';
    import 'package:provider/provider.dart';

    import '../model/{{ $.Name | snake }}.dart';
    {{ range $e := $.Edges -}}
        import '../model/{{ $e.Type.Name | snake }}.dart';
    {{ end }}

    class {{ $.Name }}Client {
        {{ $.Name }}Client({@required this.dio}) : assert(dio != null);

        final Dio dio;

        Future<{{ $.Name }}> find({{ $.ID.Type | dartType }} id) async {
            final r = await dio.get('/{{ $url }}/$id');
            return {{ $.Name }}.fromJson(r.data);
        }

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

            final r = await dio.get('/{{ $url }}');

            if (r.data == null) {
                return [];
            }

            return (r.data as List).map((i) => {{ $.Name }}.fromJson(i)).toList();
        }

        Future<{{ $.Name }}> create({{ $.Name }} e) async {
            final r = await dio.post('/{{ $url }}', data: e.toJson());
            return ({{ $.Name }}.fromJson(r.data));
        }

        Future<{{ $.Name }}> update({{ $.Name }} e) async {
            final r = await dio.patch('/{{ $url }}', data: e.toJson());
            return ({{ $.Name }}.fromJson(r.data));
        }

        {{ range $e := $.Edges}}
            Future<{{ if $e.Unique }}{{ $e.Type.Name }}{{ else }}List<{{ $e.Type.Name }}>{{ end }}> {{ $e.Name | camel }}({{ $.Name }} e) async {
                final r = await dio.get('/{{ $url }}/${e.{{ $.ID.Name }}}/{{ replace ($e.Name | snake) "_" "-" }}');
                {{ if $e.Unique -}}
                    return ({{ $e.Type.Name }}.fromJson(r.data));
                {{ else -}}
                    return (r.data as List).map((i) => {{ $e.Type.Name }}.fromJson(i)).toList();
                {{ end -}}
            }
        {{ end }}

        static {{ $.Name }}Client of(BuildContext context) => Provider.of<{{ $.Name }}Client>(context, listen: false);
    }
{{ end }}

{{ define "client/provider" }}
    {{ template "header" -}}
    import 'package:dio/dio.dart';
    import 'package:flutter/widgets.dart';
    import 'package:provider/provider.dart';
    import 'package:provider/single_child_widget.dart';

    {{ range $n := $.Nodes -}}
        import '{{ $n.Name | snake }}.dart';
    {{ end -}}

    class ClientProvider extends SingleChildStatelessWidget {
        ClientProvider({
            Key key,
            @required this.dio,
            this.child,
        }) : assert(dio != null), super(key: key, child: child);

        final Dio dio;
        final Widget child;

        @override
        Widget buildWithChild(BuildContext context, Widget child) {
            return MultiProvider(
                providers: [
                    {{ range $n := $.Nodes -}}
                        Provider<{{ $n.Name }}Client>(
                            create: (_) => {{ $n.Name }}Client(dio: dio),
                        ),
                    {{ end -}}
                ],
                child: child,
            );
        }
    }
{{ end }}