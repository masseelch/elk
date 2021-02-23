{{ define "apiClient" }}
    import 'package:http/http.dart';

    abstract class ApiClient implements Client {}
{{ end }}