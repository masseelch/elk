{{ define "repository" }}
    {{ template "header" -}}
    import 'package:dio/dio.dart';
    import 'package:flutter/foundation.dart';

    import '../model/{{ $.Name | snake }}.dart';

    class {{ $.Name }}Repository {
        {{ $.Name }}Repository(
            @required this.dio,
            @required this.url,
        ) : assert(dio != null), assert(url != null && url != '');

        final Dio dio;
        final String url;

        Future<{{ $.Name }}> find({{ $.ID.Type | dartType }} id) async {
            final r = await dio.get('/$url/$id');
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

            final r = await dio.get('/$url');

            if (r.data == null) {
                return [];
            }

            return (r.data as List).map((i) => {{ $.Name }}.fromJson(i)).toList();
        }

        Future<{{ $.Name }}> create({{ $.Name }} e) async {
            final r = await dio.post('/$url', data: e.toJson());
            return ({{ $.Name }}.fromJson(r.data));
        }

        Future<{{ $.Name }}> update({{ $.Name }} e) async {
            final r = await dio.patch('/$url', data: e.toJson());
            return ({{ $.Name }}.fromJson(r.data));
        }
    }
{{ end }}

{{ define "repository/provider" }}
    {{ template "header" -}}
    import 'package:dio/dio.dart';
    import 'package:flutter/widgets.dart';
    import 'package:provider/provider.dart';
    import 'package:provider/single_child_widget.dart';

    {{ range $n := $.Nodes -}}
        import '{{ $n.Name | snake }}.dart';
    {{ end -}}

    typedef PrefixFn = String Function(String prefix);

    class GeneratedRepositoryProvider extends SingleChildStatelessWidget {
        GeneratedRepositoryProvider({
            @required this.dio,
            this.prefixFn = _defaultPrefixFn,
        }) : assert(dio != null), assert(prefixFn != null);

        final Dio dio;
        final PrefixFn prefixFn;

        @override
        Widget buildWithChild(BuildContext context, Widget child) {
            return MultiProvider(
                providers: [
                    {{ range $n := $.Nodes -}}
                        Provider<{{ $n.Name }}Repository>(
                            create: (_) => {{ $n.Name }}Repository(dio: dio, url: prefixFn('{{ replace ($n.Name | snake) "_" "-" }}')),
                        ),
                    {{ end -}}
                ],
                child: child,
            );
        }
    }

    String _defaultPrefixFn(String url) => url;
{{ end }}