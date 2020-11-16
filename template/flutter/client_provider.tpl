{{ define "client/provider" }}
    {{ template "header" -}}
    import 'package:dio/dio.dart';
    import 'package:flutter/widgets.dart';
    import 'package:provider/provider.dart';
    import 'package:provider/single_child_widget.dart';

    {{/* Import every node */}}
    {{ range $n := $.Nodes -}}
        {{- if not $n.Annotations.HandlerGen.Skip }}
            import '{{ $n.Name | snake }}.dart';
        {{ end -}}
    {{ end -}}

    {{/* Provide the clients down the widget tree. */}}
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
                        {{- if not $n.Annotations.HandlerGen.Skip }}
                            Provider<{{ $n.Name }}Client>(
                                create: (_) => {{ $n.Name }}Client(dio: dio),
                            ),
                        {{ end -}}
                    {{ end -}}
                ],
                child: child,
            );
        }
    }
{{ end }}