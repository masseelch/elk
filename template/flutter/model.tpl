{{ define "model" }}
    {{ template "header" }}
    import 'package:json_annotation/json_annotation.dart';

    part '{{ $.Name | snake }}.g.dart';

{{/*    @JsonSerializable()*/}}
    class {{ $.Name }} {
        {{ $.Name }}();

        {{ $.ID.Type | dartType }} {{ $.ID.Name }};
        {{- range $f := $.Fields }}
            {{ $f.Type | dartType }} {{ $f.Name }};
        {{ end }}

        factory {{ $.Name }}.fromJson(Map<String, dynamic> json) => _${{ $.Name }}FromJson(json);
        Map<String, dynamic> toJson() => _${{ $.Name }}ToJson(this);
    }
{{ end }}
