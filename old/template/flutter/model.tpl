{{ define "model" }}
    {{ template "header" -}}
    import 'dart:convert';

    import 'package:json_annotation/json_annotation.dart';

    import '../date_utc_converter.dart';

    {{ $df := dartRequestFields $.Type "" }}

    {{/* Import the custom dart types. */}}
    {{ range $.TypeMappings -}}
        import '{{ .Import }}';
        {{- if .ConverterImport }}import '{{ .ConverterImport }}';{{ end -}}
    {{ end }}

    {{/* For every edge import the generated model. */}}
    {{ range $e := $.Edges -}}
        import '../model/{{ $e.Type.Name | snake }}.dart';
    {{ end }}

    {{/* JsonSerializable puts the generated code in this file. */}}
    part '{{ $.Name | snake }}.g.dart';

    @JsonSerializable()
    @DateUtcConverter()
    class {{ $.Name }} {
        {{ $.Name }}();

        {{/* ID of the model */}}
        {{- $jsonName := index (split (tagLookup $.ID.StructTag "json") ",") 0 }}
        {{- if and (not (eq "-" $jsonName)) (not (eq $jsonName (snake $.ID.Name))) }}
            @JsonKey(name: '{{ $jsonName }}')
        {{ end -}}
        {{ $.ID.Type | dartType }}? {{ $.ID.Name }};

        {{/* The fields of the model. */}}
        {{- range $f := $.Fields -}}
            {{- $c := $df.ConverterFor $f }}
            {{- if and $f.Annotations.FieldGen.MapGoType $f.HasGoType -}}
                {{- if $c }}{{ $c }}{{ end -}}
            {{ end -}}
            {{- $jsonName := index (split (tagLookup $f.StructTag "json") ",") 0 }}
            {{- if and (not (eq "-" $jsonName)) (not (eq $jsonName (snake $f.Name))) }}
                @JsonKey(name: '{{ $jsonName }}')
            {{ end -}}
            {{ $f.Type | dartType }}? {{ $f.Name | camel }};
        {{ end }}

        {{/* The edges of the model. */}}
        {{ $.Name }}Edges? edges;

        @override
        int get hashCode => {{ $.ID.Name }}.hashCode;

        @override
        bool operator ==(Object other) => other is {{ $.Name }} && {{ $.ID.Name }} == other.{{ $.ID.Name }};

        factory {{ $.Name }}.fromJson(Map<String, dynamic> json) => _${{ $.Name }}FromJson(json);
        Map<String, dynamic> toJson() => _${{ $.Name }}ToJson(this);

        String toString() => jsonEncode(toJson());
    }

    {{/* The edges of the model. */}}
    @JsonSerializable()
    class {{ $.Name }}Edges {
        {{ $.Name }}Edges();

        {{ range $e := $.Edges -}}
            {{ if $e.Unique }}{{ $e.Type.Name }}{{ else }}List<{{ $e.Type.Name }}>{{ end }}? {{ $e.Name | camel }};
        {{ end }}

        factory {{ $.Name }}Edges.fromJson(Map<String, dynamic> json) => _${{ $.Name }}EdgesFromJson(json);
        Map<String, dynamic> toJson() => _${{ $.Name }}EdgesToJson(this);
    }
{{ end }}
