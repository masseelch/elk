{{ define "repository" }}
    {{ template "header" }}
    import 'package:dio/dio.dart';

    class {{ $.Name }}Repository {
        {{ $.Name }}Repository();
    }
{{ end }}
