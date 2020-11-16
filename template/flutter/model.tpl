{{ define "model" }}
    {{ template "header" -}}
    import 'package:json_annotation/json_annotation.dart';

    {{/* For every edge import the generated model. */}}
    {{ range $e := $.Edges }}
        import '../model/{{ $e.Type.Name | snake }}.dart';
    {{ end }}

    {{/* JsonSerializable puts the generated code in this file. */}}
    part '{{ $.Name | snake }}.g.dart';

    @JsonSerializable()
    class {{ $.Name }} {
        {{ $.Name }}();

        {{/* The fields of the model. */}}
        {{ $.ID.Type | dartType }} {{ $.ID.Name }};
        {{- range $f := $.Fields -}}
            {{ $f.Type | dartType }} {{ $f.Name }};
        {{ end }}

        {{/* The edges of the model. */}}
        {{ $.Name }}Edges edges;

        factory {{ $.Name }}.fromJson(Map<String, dynamic> json) => _${{ $.Name }}FromJson(json);
        Map<String, dynamic> toJson() => _${{ $.Name }}ToJson(this);
    }

    {{/* The edges of the model. */}}
    @JsonSerializable()
    class {{ $.Name }}Edges {
        {{ $.Name }}Edges();

        {{ range $e := $.Edges }}
            {{ if $e.Unique }}{{ $e.Type.Name }}{{ else }}List<{{ $e.Type.Name }}>{{ end }} {{ $e.Name }};
        {{ end }}

        factory {{ $.Name }}Edges.fromJson(Map<String, dynamic> json) => _${{ $.Name }}EdgesFromJson(json);
        Map<String, dynamic> toJson() => _${{ $.Name }}EdgesToJson(this);
    }
{{ end }}
