{{ define "client/provider" }}
    {{ template "header" -}}
    import 'package:flutter/widgets.dart';
    import 'package:intercepted_http/intercepted_http.dart' show Client;
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
            Key? key,
            required this.client,
            this.child,
        }) : super(key: key, child: child);

        final Client client;
        final Widget? child;

        @override
        Widget buildWithChild(BuildContext context, Widget? child) {
            return MultiProvider(
                providers: [
                    {{ range $n := $.Nodes -}}
                        {{- if not $n.Annotations.HandlerGen.Skip }}
                            Provider<{{ $n.Name }}Client>(
                                create: (_) => {{ $n.Name }}Client(client: client),
                            ),
                        {{ end -}}
                    {{ end -}}
                ],
                child: child,
            );
        }
    }
{{ end }}